package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	api "github.com/32leaves/ruruku/pkg/server/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"sync"
)

type memoryBackedSessionStore struct {
	Sessions map[string]*memoryBackedSession
	Mux      sync.Mutex
}

type memoryBackedSession struct {
	Name         string
	PlanID       string
	Open         bool
	Status       map[string]*memoryBackedStatus
	Participants map[string]*types.Participant
	Mux          sync.Mutex
}

type memoryBackedStatus struct {
	Case    types.Testcase
	Claims  map[string]*types.Participant
	Results map[string]*types.TestcaseRunResult
}

func NewMemoryBackedSessionStore() *memoryBackedSessionStore {
	return &memoryBackedSessionStore{
		Sessions: make(map[string]*memoryBackedSession),
	}
}

func (s *memoryBackedSessionStore) Version(ctx context.Context, req *api.VersionRequest) (*api.VersionResponse, error) {
	return &api.VersionResponse{
		Version:     "implement_me",
		ReleaseName: "bloated octopus",
	}, nil
}

func (s *memoryBackedSessionStore) Start(ctx context.Context, req *api.StartSessionRequest) (*api.StartSessionResponse, error) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	sid, err := toSessionID(req.Name)
	if err != nil {
		return nil, err
	}
	if _, exists := s.Sessions[sid]; exists {
		return nil, fmt.Errorf("Session '%s' exists already", req.Name)
	}

	planID := ""
	cases := make(map[string]*memoryBackedStatus)
	if req.Plan != nil {
		planID = req.Plan.Id

		for _, cse := range req.Plan.Case {
			if _, exists := cases[cse.Id]; exists {
				return nil, fmt.Errorf("Testcase '%s' exists already", cse.Id)
			}

			cases[cse.Id] = &memoryBackedStatus{
				Case:    cse.Convert(),
				Claims:  make(map[string]*types.Participant),
				Results: make(map[string]*types.TestcaseRunResult),
			}
		}
	}

	session := &memoryBackedSession{
		Name:         req.Name,
		PlanID:       planID,
		Open:         true,
		Status:       cases,
		Participants: make(map[string]*types.Participant),
	}
	s.Sessions[sid] = session

	log.WithField("id", sid).WithField("name", req.Name).WithField("cases", len(cases)).Info("Starting session")

	return &api.StartSessionResponse{Id: sid}, nil
}

func (s *memoryBackedSessionStore) Close(ctx context.Context, req *api.CloseSessionRequest) (*api.CloseSessionResponse, error) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	session, exists := s.Sessions[req.Id]
	if !exists {
		return nil, fmt.Errorf("Session %s does not exist", req.Id)
	}
	session.Open = false

	log.WithField("id", req.Id).WithField("name", session.Name).Info("Closing session")

	return &api.CloseSessionResponse{}, nil
}

func (s *memoryBackedSessionStore) List(req *api.ListSessionsRequest, resp api.SessionService_ListServer) error {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	for id, session := range s.Sessions {
		err := resp.Send(&api.ListSessionsResponse{
			Id:     id,
			Name:   session.Name,
			IsOpen: session.Open,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *memoryBackedSessionStore) Register(ctx context.Context, req *api.RegistrationRequest) (*api.RegistrationResponse, error) {
	session, err := s.getOpenSession(req.SessionID)
	if err != nil {
		return nil, err
	}

	session.Mux.Lock()
	defer session.Mux.Unlock()
	for id, p := range session.Participants {
		if p.Name == req.Name {
			return &api.RegistrationResponse{
				Token:  id,
				Status: session.getStatus(req.SessionID),
			}, nil
		}
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("Cannot create participant ID: %v", err)
	}
	uid := fmt.Sprintf("%s/%s", req.SessionID, id.String())

	session.Participants[uid] = &types.Participant{Name: req.Name}

	return &api.RegistrationResponse{
		Token:  uid,
		Status: session.getStatus(req.SessionID),
	}, nil
}

func (s *memoryBackedSessionStore) Claim(ctx context.Context, req *api.ClaimRequest) (*api.ClaimResponse, error) {
	sid, uid, err := parseParticipantToken(req.ParticipantToken)
	if err != nil {
		return nil, err
	}

	session, err := s.getOpenSession(sid)
	if err != nil {
		return nil, err
	}

	session.Mux.Lock()
	defer session.Mux.Unlock()
	participant, exists := session.Participants[uid]
	if !exists {
		return nil, fmt.Errorf("Invalid participant token: does not exist in session")
	}
	tc, exists := session.Status[req.TestcaseID]
	if !exists {
		return nil, fmt.Errorf("Testcase %s does not exist in session", req.TestcaseID)
	}

	tc.Claims[uid] = participant

	return &api.ClaimResponse{}, nil
}

func (s *memoryBackedSessionStore) Contribute(ctx context.Context, req *api.ContributionRequest) (*api.ContributionResponse, error) {
	sid, uid, err := parseParticipantToken(req.ParticipantToken)
	if err != nil {
		return nil, err
	}

	session, err := s.getOpenSession(sid)
	if err != nil {
		return nil, err
	}

	session.Mux.Lock()
	defer session.Mux.Unlock()
	if _, exists := session.Participants[uid]; !exists {
		return nil, fmt.Errorf("Invalid participant token: does not exist in session")
	}
	tc, exists := session.Status[req.TestcaseID]
	if !exists {
		return nil, fmt.Errorf("Testcase %s does not exist in session", req.TestcaseID)
	}
	if _, claimed := tc.Claims[uid]; !claimed {
		return nil, fmt.Errorf("Must claim the testcase before contributing")
	}

	tc.Results[uid] = &types.TestcaseRunResult{
		State:   types.TestRunState(req.Result),
		Comment: req.Comment,
	}

	return &api.ContributionResponse{}, nil
}

func (s *memoryBackedSessionStore) Status(ctx context.Context, req *api.SessionStatusRequest) (*api.SessionStatusResponse, error) {
	session, exists := s.Sessions[req.Id]
	if !exists {
		return nil, fmt.Errorf("Session %s does not exist", req.Id)
	}

	session.Mux.Lock()
	defer session.Mux.Unlock()
	return &api.SessionStatusResponse{Status: session.getStatus(req.Id)}, nil
}

func (s *memoryBackedSessionStore) Updates(req *api.SessionUpdatesRequest, update api.SessionService_UpdatesServer) error {
	return fmt.Errorf("Not implemented")
}

// assumes you hold the store lock
func (s *memoryBackedSessionStore) getOpenSession(id string) (*memoryBackedSession, error) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	session, exists := s.Sessions[id]
	if !exists {
		return nil, fmt.Errorf("Session %s does not exist", id)
	}
	if !session.Open {
		return nil, fmt.Errorf("Session %s is closed", id)
	}
	return session, nil
}

func (s *memoryBackedSession) getStatus(id string) *api.TestRunStatus {
	state := types.Passed
	tcstatus := make([]*api.TestcaseStatus, len(s.Status))
	for _, tc := range s.Status {
		claims := make([]*api.Participant, 0)
		for _, p := range tc.Claims {
			claims = append(claims, api.ConvertParticipant(p))
		}

		results := make([]*api.TestcaseRunResult, 0)
		tcstate := types.Passed
		for _, r := range tc.Results {
			results = append(results, api.ConvertTestcaseRunResult(r))
			tcstate = types.WorseState(tcstate, r.State)
		}
		if tc.Case.MustPass {
			state = types.WorseState(state, tcstate)
		}

		tcstatus = append(tcstatus, &api.TestcaseStatus{
			Case:   api.ConvertTestcase(&tc.Case),
			Claim:  claims,
			Result: results,
			State:  api.ConvertTestRunState(tcstate),
		})
	}

	return &api.TestRunStatus{
		Id:     id,
		Name:   s.Name,
		PlanID: s.PlanID,
		Status: tcstatus,
		State:  api.ConvertTestRunState(state),
	}
}

func parseParticipantToken(token string) (string, string, error) {
	if token == "" {
		return "", "", fmt.Errorf("Token must not be empty")
	}

	parts := strings.Split(token, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Token has an invalid format")
	}

	// session id, participant id
	return parts[0], parts[1], nil
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
