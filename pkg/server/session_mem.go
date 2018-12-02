package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	api "github.com/32leaves/ruruku/pkg/api/v1"
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

	if req.Name == "" {
		return nil, fmt.Errorf("Cannot start a session with an empty name")
	}

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

			tcs := cse.Convert()
			if err := types.ValidateTestcase(&tcs); err != nil {
				return nil, err
			}

			cases[cse.Id] = &memoryBackedStatus{
				Case:    tcs,
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

	uid := ""
	for id, p := range session.Participants {
		if p.Name == req.Name {
			uid = id
			break
		}
	}

	if uid == "" {
		if id, err := uuid.NewV4(); err != nil {
			return nil, fmt.Errorf("Cannot create participant ID: %v", err)
		} else {
			uid = id.String()
		}
	}
	token := fmt.Sprintf("%s/%s", req.SessionID, uid)

	pcp := &types.Participant{Name: req.Name}
	if err := types.ValidateParticipant(pcp); err != nil {
		return nil, err
	}
	session.Participants[uid] = pcp
	log.WithField("session", req.SessionID).WithField("name", req.Name).Info("Participant joined session")

	return &api.RegistrationResponse{
		Token: token,
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
		return nil, fmt.Errorf("Invalid participant token %s: does not exist in session", uid)
	}

	tc, exists := session.Status[req.TestcaseID]
	if !exists {
		return nil, fmt.Errorf("Testcase %s does not exist in session", req.TestcaseID)
	}

	if req.Claim {
		tc.Claims[uid] = participant
		log.WithField("session", sid).WithField("participant", participant.Name).WithField("testcase", req.TestcaseID).Info("Participant claimed testcase")
	} else {
		delete(tc.Claims, uid)
		log.WithField("session", sid).WithField("participant", participant.Name).WithField("testcase", req.TestcaseID).Info("Participant unclaimed testcase")
	}

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
	participant, exists := session.Participants[uid]
	if !exists {
		return nil, fmt.Errorf("Invalid participant token: does not exist in session")
	}
	tc, exists := session.Status[req.TestcaseID]
	if !exists {
		return nil, fmt.Errorf("Testcase %s does not exist in session", req.TestcaseID)
	}
	if _, claimed := tc.Claims[uid]; !claimed {
		return nil, fmt.Errorf("Must claim the testcase before contributing")
	}

	contribution := &types.TestcaseRunResult{
		Comment:     req.Comment,
		State:       req.Result.Convert(),
		Participant: *participant,
	}
	tc.Results[uid] = contribution
	log.WithField("state", contribution.State).WithField("testcase", tc).Debug("Recording contribution")

	return &api.ContributionResponse{}, nil
}

func (s *memoryBackedSessionStore) Status(ctx context.Context, req *api.SessionStatusRequest) (*api.SessionStatusResponse, error) {
	session, exists := s.Sessions[req.Id]
	if !exists {
		return nil, fmt.Errorf("Session %s does not exist", req.Id)
	}

	session.Mux.Lock()
	status := session.getStatus(req.Id)
	session.Mux.Unlock()
	return &api.SessionStatusResponse{Status: status}, nil
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
	var state types.TestRunState
	if len(s.Status) == 0 {
		state = types.Undecided
	} else {
		state = types.Passed
	}

	tcstatus := make([]*api.TestcaseStatus, 0)
	for _, tc := range s.Status {
		claims := make([]*api.Participant, 0)
		for _, p := range tc.Claims {
			claims = append(claims, api.ConvertParticipant(p))
		}

		var tcstate types.TestRunState
		if len(tc.Results) == 0 {
			tcstate = types.Undecided
		} else {
			tcstate = types.Passed
		}

		results := make([]*api.TestcaseRunResult, 0)
		for _, r := range tc.Results {
			results = append(results, api.ConvertTestcaseRunResult(r))
			tcstate = types.WorseState(tcstate, r.State)
		}
		if tc.Case.MustPass {
			state = types.WorseState(state, tcstate)
		}

		cse := api.ConvertTestcase(&tc.Case)
		tcstatus = append(tcstatus, &api.TestcaseStatus{
			Case:   cse,
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
