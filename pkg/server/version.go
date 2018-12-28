package server

import (
	"context"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/version"
)

type versionServer struct{}

func (s *versionServer) Get(ctx context.Context, req *api.GetVersionRequest) (*api.GetVersionResponse, error) {
	return &api.GetVersionResponse{
		Tag:       version.Tag,
		Rev:       version.Rev,
		BuildDate: version.BuildDate,
	}, nil
}
