package kvsession

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/server/notifier"
	"github.com/32leaves/ruruku/pkg/server/request"
	"github.com/32leaves/ruruku/pkg/types"
	bolt "github.com/etcd-io/bbolt"
	log "github.com/sirupsen/logrus"
)

func NewSession(db *bolt.DB, reqvalidator request.Validator) (*kvsessionStore, error) {
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Initialize buckets to guarantee that they exist.
	tx.CreateBucketIfNotExists([]byte(bucketSessions))
	tx.CreateBucketIfNotExists([]byte(bucketTestplan))

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	notifier := make(map[string]*notifier.Notifier)
	return &kvsessionStore{DB: db, Notifier: notifier, ReqValidator: reqvalidator}, nil
}

type kvsessionStore struct {
	DB           *bolt.DB
	Notifier     map[string]*notifier.Notifier
	ReqValidator request.Validator
}

func (s *kvsessionStore) Start(ctx context.Context, req *api.StartSessionRequest) (*api.StartSessionResponse, error) {
	if _, err := s.ReqValidator.ValidUserFromRequest(ctx, types.PermissionSessionStart); err != nil {
		return nil, err
	}

	if req.Name == "" {
		return nil, fmt.Errorf("Cannot start a session with an empty name")
	}

	sid, err := toSessionID(req.Name)
	if err != nil {
		return nil, err
	}

	if exists, err := s.sessionExists(sid); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("Session '%s' exists already", req.Name)
	}

	planID := ""
	if req.Plan != nil {
		for _, cse := range req.Plan.Case {
			tcs := cse.Convert()
			if err := types.ValidateTestcase(&tcs); err != nil {
				return nil, err
			}
		}

		if err := s.storePlan(sid, req.Plan); err != nil {
			return nil, err
		}
		planID = req.Plan.Id
	}
	if err := s.storeSession(sid, req.Name, planID); err != nil {
		return nil, err
	}

	s.Notifier[sid] = notifier.NewNotifier()

	log.WithField("id", sid).WithField("name", req.Name).Info("Starting session")
	return &api.StartSessionResponse{Id: sid}, nil
}

func (s *kvsessionStore) Close(ctx context.Context, req *api.CloseSessionRequest) (*api.CloseSessionResponse, error) {
	if _, err := s.ReqValidator.ValidUserFromRequest(ctx, types.PermissionSessionClose); err != nil {
		return nil, err
	}

	exists, err := s.sessionExists(req.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Session %s does not exist", req.Id)
	}

	if err := s.closeSession(req.Id); err != nil {
		return nil, err
	}

	defer s.broadcastChange(req.Id)
	return &api.CloseSessionResponse{}, nil
}

func (s *kvsessionStore) List(req *api.ListSessionsRequest, resp api.SessionService_ListServer) error {
	return s.listSessions(resp.Send)
}

func (s *kvsessionStore) Register(ctx context.Context, req *api.RegistrationRequest) (*api.RegistrationResponse, error) {
	user, err := s.ReqValidator.ValidUserFromRequest(ctx, types.PermissionSessionContribute)
	if err != nil {
		return nil, err
	}

	exists, err := s.isSessionOpen(req.Session)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Session %s is not open", req.Session)
	}

	if err := s.registerParticipant(req.Session, user); err != nil {
		return nil, err
	}

	defer s.broadcastChange(req.Session)
	return &api.RegistrationResponse{}, nil
}

func (s *kvsessionStore) Claim(ctx context.Context, req *api.ClaimRequest) (*api.ClaimResponse, error) {
	user, err := s.ReqValidator.ValidUserFromRequest(ctx, types.PermissionSessionContribute)
	if err != nil {
		return nil, err
	}

	if req.TestcaseID == "" {
		return nil, fmt.Errorf("Testcase does not exist in session")
	}

	// parse token
	sid := req.Session
	uid := user

	exists, err := s.isSessionOpen(sid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Session %s is closed", sid)
	}

	ok, err := s.participantInSession(sid, uid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Invalid participant %s: does not exist in session", uid)
	}

	err = s.claimTestcase(sid, req.TestcaseID, uid, req.Claim)
	if err != nil {
		return nil, err
	}

	defer s.broadcastChange(sid)
	return &api.ClaimResponse{}, nil
}

func (s *kvsessionStore) Contribute(ctx context.Context, req *api.ContributionRequest) (*api.ContributionResponse, error) {
	user, err := s.ReqValidator.ValidUserFromRequest(ctx, types.PermissionSessionContribute)
	if err != nil {
		return nil, err
	}

	if req.TestcaseID == "" {
		return nil, fmt.Errorf("Testcase does not exist in session")
	}

	// parse token
	sid := req.Session
	uid := user

	exists, err := s.isSessionOpen(sid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Session %s is closed", sid)
	}

	ok, err := s.participantInSession(sid, uid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Invalid participant token %s: does not exist in session", uid)
	}

	exists, err = s.testcaseExists(sid, req.TestcaseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Testcase %s does not exist in session", req.TestcaseID)
	}

	exists, err = s.hasClaimedTestcase(sid, req.TestcaseID, uid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Participant must claim the testcase before contributing")
	}

	result := req.Result.Convert()
	err = s.contribute(sid, kvsessionContribution{
		UserID:     uid,
		TestcaseID: req.TestcaseID,
		Result:     result,
		Comment:    req.Comment,
	})
	if err != nil {
		return nil, err
	}

	defer s.broadcastChange(sid)
	return &api.ContributionResponse{}, nil
}

func (s *kvsessionStore) Status(ctx context.Context, req *api.SessionStatusRequest) (*api.SessionStatusResponse, error) {
	if _, err := s.ReqValidator.ValidUserFromRequest(ctx, types.PermissionSessionView); err != nil {
		return nil, err
	}

	sid := req.Id
	exists, err := s.sessionExists(sid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Invalid participant token: session %s does not exist", sid)
	}

	status, err := s.getStatus(sid)
	if err != nil {
		return nil, err
	}

	return &api.SessionStatusResponse{Status: status}, nil
}

func (s *kvsessionStore) Updates(req *api.SessionUpdatesRequest, update api.SessionService_UpdatesServer) error {
	sid := req.Id
	exists, err := s.sessionExists(sid)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Session %s does not exist", sid)
	}

	notifier, exists := s.Notifier[sid]
	if !exists {
		return fmt.Errorf("Session %s does not exist", sid)
	}

	for {
		status := notifier.Listen()
		if err := update.Send(&api.SessionUpdateResponse{Id: sid, Status: &status}); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (s *kvsessionStore) broadcastChange(sid string) error {
	notifier, ok := s.Notifier[sid]
	if !ok {
		return fmt.Errorf("Session %s does not exist", sid)
	}

	status, err := s.getStatus(sid)
	if err != nil {
		return err
	}

	notifier.Update(status)
	return nil
}

func toSessionID(name string) (string, error) {
	candidate := strings.ToLower(name)
	candidate = strings.Replace(candidate, "\n", " ", -1)
	candidate = strings.Replace(candidate, "-", " ", -1)
	candidate = strings.Replace(candidate, "_", " ", -1)
	segments := strings.Split(candidate, " ")
	if len(segments) > 3 {
		segments = segments[0:2]
	}

	uid := make([]byte, 4)
	n, err := io.ReadFull(rand.Reader, uid)
	if err != nil {
		return "", err
	} else if n < len(uid) {
		return "", fmt.Errorf("Did not read enough random bytes")
	}
	segments = append(segments, hex.EncodeToString(uid))

	return strings.Join(segments, "-"), nil
}
