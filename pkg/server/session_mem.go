package server

import (
	"context"
	"fmt"
	api "github.com/32leaves/ruruku/pkg/server/api/v1"
)

type memoryBackedSession struct {

}

func NewMemoryBackedSession() *memoryBackedSession {
	return &memoryBackedSession{}
}

func (s *memoryBackedSession) Version(ctx context.Context, req *api.VersionRequest) (*api.VersionResponse, error) {
	return &api.VersionResponse{
		Version:     "implement_me",
		ReleaseName: "bloated octopus",
	}, nil
}

func (s *memoryBackedSession) Start(ctx context.Context, req *api.StartSessionRequest) (*api.StartSessionResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSession) Close(ctx context.Context, req *api.CloseSessionRequest) (*api.CloseSessionResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSession) List(*api.ListSessionsRequest, api.SessionService_ListServer) error {
	return fmt.Errorf("Not implemented")
}

func (s *memoryBackedSession) Register(ctx context.Context, req *api.RegistrationRequest) (*api.RegistrationResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSession) Claim(ctx context.Context, req *api.ClaimRequest) (*api.ClaimResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSession) Contribute(ctx context.Context, req *api.ContributionRequest) (*api.ContributionResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSession) Status(ctx context.Context, req *api.SessionStatusRequest) (*api.SessionStatusResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (s *memoryBackedSession) Updates(req *api.SessionUpdatesRequest, update api.SessionService_UpdatesServer) error {
	err := update.Send(&api.SessionUpdateResponse{
		Id: fmt.Sprintf("foobar"),
	})
	if err != nil {
		return err
	}

	return nil
}
