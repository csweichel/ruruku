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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	n := make(map[string]*notifier.Notifier)
	r := &kvsessionStore{DB: db, Notifier: n, ReqValidator: reqvalidator}
	r.listSessions(func(s *api.ListSessionsResponse) error {
		r.Notifier[s.Id] = notifier.NewNotifier()
		return nil
	})
	return r, nil
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
		plan := req.Plan.Convert()
		if err := plan.Validate(); err != nil {
			return nil, err
		}

		for _, cse := range req.Plan.Case {
			tcs := cse.Convert()
			if err := tcs.Validate(); err != nil {
				return nil, err
			}
		}

		if err := s.storePlan(sid, req.Plan); err != nil {
			return nil, err
		}
		planID = req.Plan.Id
	}
	if err := s.storeSession(sid, req.Name, planID, req.Modifiable, req.Annotations); err != nil {
		return nil, err
	}

	s.Notifier[sid] = notifier.NewNotifier()

	log.WithField("id", sid).WithField("name", req.Name).Info("Starting session")
	return &api.StartSessionResponse{Id: sid}, nil
}

func (s *kvsessionStore) Modify(ctx context.Context, req *api.ModifySessionRequest) (*api.ModifySessionResponse, error) {
	if _, err := s.ReqValidator.ValidUserFromRequest(ctx, types.PermissionSessionModify); err != nil {
		return nil, err
	}

	exists, err := s.isSessionOpen(req.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Errorf(codes.FailedPrecondition, "Session %s is not open", req.Id)
	}

	modifiable, err := s.isSessionModifiable(req.Id)
	if err != nil {
		return nil, err
	}
	if !modifiable {
		return nil, status.Errorf(codes.PermissionDenied, "Session %s is not modifiable", req.Id)
	}

	if req.Modification == api.Modification_ADD_TESTCASE || req.Modification == api.Modification_MODIFY_TESTCASE {
		if req.Case == nil {
			return nil, status.Errorf(codes.InvalidArgument, "Testcase must not be empty")
		}
		for _, tc := range req.Case {
			tcc := tc.Convert()
			if err := tcc.Validate(); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "%v", err)
			}
		}
	}

	switch req.Modification {
	case api.Modification_ADD_TESTCASE:
		err = s.addOrUpdateTestcase(req.Id, req.Case, false)
	case api.Modification_MODIFY_TESTCASE:
		err = s.addOrUpdateTestcase(req.Id, req.Case, true)
	case api.Modification_REMOVE_TESTCASE:
		if req.Case == nil {
			return nil, status.Errorf(codes.InvalidArgument, "Testcase must not be empty")
		}
		err = s.removeTestcase(req.Id, req.Case)
	case api.Modification_UPDATE_ANNOTATIONS:
		if req.Annotations == nil {
			return nil, status.Errorf(codes.InvalidArgument, "Annotations must not be nil")
		}
		err = s.modifySession(req.Id, func(session *SessionMetadata) {
			session.Annotations = req.Annotations
		})
	}
	if err != nil {
		return nil, err
	}

	if err := s.broadcastChange(req.Id); err != nil {
		return nil, err
	}
	return &api.ModifySessionResponse{}, nil
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
		return nil, status.Errorf(codes.NotFound, "Session %s does not exist", req.Id)
	}

	err = s.modifySession(req.Id, func(session *SessionMetadata) {
		session.Open = false
	})
	if err != nil {
		return nil, err
	}

	defer s.broadcastChange(req.Id)
	return &api.CloseSessionResponse{}, nil
}

func (s *kvsessionStore) List(req *api.ListSessionsRequest, resp api.SessionService_ListServer) error {
	if _, err := s.ReqValidator.ValidUserFromRequest(resp.Context(), types.PermissionSessionView); err != nil {
		return err
	}

	err := s.listSessions(resp.Send)
	if err != nil {
		log.WithError(err).Error("Error while listing sessions")
		return err
	}

	return nil
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
		return nil, status.Errorf(codes.NotFound, "Session %s does not exist", req.Session)
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
		return nil, status.Errorf(codes.NotFound, "Testcase does not exist in session")
	}

	// parse token
	sid := req.Session
	uid := user

	exists, err := s.isSessionOpen(sid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Errorf(codes.FailedPrecondition, "Session %s is closed", sid)
	}

	ok, err := s.participantInSession(sid, uid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Invalid participant %s: does not exist in session", uid)
	}

	err = s.claimTestcase(sid, req.TestcaseID, uid, req.Claim)
	if err != nil {
		return nil, err
	}

	log.WithField("session", sid).
		WithField("testcase", req.TestcaseID).
		WithField("user", uid).
		WithField("claim", req.Claim).
		Debug("Testcase claim recorded")

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
		return nil, status.Errorf(codes.FailedPrecondition, "Session %s is closed", sid)
	}

	ok, err := s.participantInSession(sid, uid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Invalid participant %s: does not exist in session", uid)
	}

	exists, err = s.testcaseExists(sid, req.TestcaseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Testcase %s does not exist in session", req.TestcaseID)
	}

	exists, err = s.hasClaimedTestcase(sid, req.TestcaseID, uid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Errorf(codes.FailedPrecondition, "Participant must claim the testcase before contributing")
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
		return nil, status.Errorf(codes.NotFound, "Session %s does not exist", sid)
	}

	status, err := s.getStatus(sid)
	if err != nil {
		return nil, err
	}

	return &api.SessionStatusResponse{Status: status}, nil
}

func (s *kvsessionStore) Updates(req *api.SessionUpdatesRequest, update api.SessionService_UpdatesServer) error {
	if _, err := s.ReqValidator.ValidUserFromRequest(update.Context(), types.PermissionSessionView); err != nil {
		return err
	}

	sid := req.Id
	exists, err := s.sessionExists(sid)
	if err != nil {
		return err
	}
	if !exists {
		return status.Errorf(codes.NotFound, "Session %s does not exist", sid)
	}

	notifier, exists := s.Notifier[sid]
	if !exists {
		return status.Errorf(codes.Internal, "Session notifier %s does not exist", sid)
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
		return status.Errorf(codes.NotFound, "Session %s does not exist", sid)
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
		return "", status.Errorf(codes.Internal, "Did not read enough random bytes")
	}
	segments = append(segments, hex.EncodeToString(uid))

	return strings.Join(segments, "-"), nil
}
