package kvsession

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/server/notifier"
	"github.com/32leaves/ruruku/pkg/types"
	bolt "github.com/etcd-io/bbolt"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

func NewSession(path string) (*kvsessionStore, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

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
	return &kvsessionStore{DB: db, Notifier: notifier}, nil
}

type kvsessionStore struct {
	DB       *bolt.DB
	Notifier map[string]*notifier.Notifier
}

func (s *kvsessionStore) Version(ctx context.Context, req *api.VersionRequest) (*api.VersionResponse, error) {
	return &api.VersionResponse{
		Version:     "implement_me",
		ReleaseName: "bloated octopus",
	}, nil
}

func (s *kvsessionStore) Start(ctx context.Context, req *api.StartSessionRequest) (*api.StartSessionResponse, error) {
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
	exists, err := s.isSessionOpen(req.SessionID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Session %s is not open", req.SessionID)
	}

	if req.Name == "" {
		return nil, fmt.Errorf("Cannot register with an empty name")
	}

	token, err := s.registerParticipant(req.SessionID, req.Name)
	if err != nil {
		return nil, err
	}

	defer s.broadcastChange(req.SessionID)
	return &api.RegistrationResponse{Token: token}, nil
}

func (s *kvsessionStore) Claim(ctx context.Context, req *api.ClaimRequest) (*api.ClaimResponse, error) {
	if req.TestcaseID == "" {
		return nil, fmt.Errorf("Testcase does not exist in session")
	}

	// parse token
	token, err := types.ParseParticipantToken(req.ParticipantToken)
	if err != nil {
		return nil, err
	}
	sid := token.SessionID
	uid := token.ParticipantID

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

	err = s.claimTestcase(sid, req.TestcaseID, uid, req.Claim)
	if err != nil {
		return nil, err
	}

	defer s.broadcastChange(sid)
	return &api.ClaimResponse{}, nil
}

func (s *kvsessionStore) Contribute(ctx context.Context, req *api.ContributionRequest) (*api.ContributionResponse, error) {
	if req.TestcaseID == "" {
		return nil, fmt.Errorf("Testcase does not exist in session")
	}

	// parse token
	token, err := types.ParseParticipantToken(req.ParticipantToken)
	if err != nil {
		return nil, err
	}
	sid := token.SessionID
	uid := token.ParticipantID

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
