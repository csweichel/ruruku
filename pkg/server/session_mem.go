package server

import (
	"context"
	"fmt"
	api "github.com/32leaves/ruruku/pkg/server/api/v1"
    "github.com/32leaves/ruruku/pkg/types"
    "github.com/satori/go.uuid"
    log "github.com/sirupsen/logrus"
)


type memoryBackedSessionStore struct {
    Sessions map[string]*memoryBackedSession
}

type memoryBackedSession struct {
    Name string
    PlanID string
    Open bool
    Status map[string]*memoryBackedStatus
    Participants map[string]*types.Participant
}

type memoryBackedStatus struct {
    Case types.Testcase
    Claims map[string]*types.Participant
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
    if _, exists := s.Sessions[req.Name]; exists {
        return nil, fmt.Errorf("Session '%s' exists already", req.Name)
    }

    id := uuid.Must(uuid.NewV4()).String()
    planID := id
    cases := make(map[string]*memoryBackedStatus)
    if req.Plan != nil {
        planID = req.Plan.Id

        for _, cse := range req.Plan.Case {
            if _, exists := cases[cse.Id]; exists {
                return nil, fmt.Errorf("Testcase '%s' exists already", cse.Id)
            }

            cases[cse.Id] = &memoryBackedStatus{
                Case: cse.Convert(),
                Claims: make(map[string]*types.Participant),
                Results: make(map[string]*types.TestcaseRunResult),
            }
        }
    }

    session := &memoryBackedSession{
        Name: req.Name,
        PlanID: planID,
        Open: true,
        Status: cases,
        Participants: make(map[string]*types.Participant),
    }
    s.Sessions[req.Name] = session

    log.WithField("id", id).WithField("name", req.Name).Info("Starting session")

	return &api.StartSessionResponse{ Id: id }, nil
}

func (s *memoryBackedSessionStore) Close(ctx context.Context, req *api.CloseSessionRequest) (*api.CloseSessionResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSessionStore) List(*api.ListSessionsRequest, api.SessionService_ListServer) error {
	return fmt.Errorf("Not implemented")
}

func (s *memoryBackedSessionStore) Register(ctx context.Context, req *api.RegistrationRequest) (*api.RegistrationResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSessionStore) Claim(ctx context.Context, req *api.ClaimRequest) (*api.ClaimResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSessionStore) Contribute(ctx context.Context, req *api.ContributionRequest) (*api.ContributionResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSessionStore) Status(ctx context.Context, req *api.SessionStatusRequest) (*api.SessionStatusResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSessionStore) Updates(req *api.SessionUpdatesRequest, update api.SessionService_UpdatesServer) error {
	err := update.Send(&api.SessionUpdateResponse{
		Id: fmt.Sprintf("foobar"),
	})
	if err != nil {
		return err
	}

	return nil
}
