package kvuser

import (
	"context"

	"github.com/32leaves/ruruku/pkg/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *kvuserStore) getUserFromRequest(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) < 1 {
		return "", status.Error(codes.Unauthenticated, "No authorization token supplied")
	}

	token := authHeader[0]
	if usr, err := s.getUserFromToken(token); err != nil {
		return "", status.Error(codes.Internal, err.Error())
	} else if usr == "" {
		return "", status.Error(codes.Unauthenticated, "Authorization token is not valid")
	} else {
		return usr, nil
	}
}

func (s *kvuserStore) ValidUserFromRequest(ctx context.Context, reqperm types.Permission) (string, error) {
	usr, err := s.getUserFromRequest(ctx)
	if err != nil {
		return "", err
	}

	if reqperm != types.PermissionNone {
		if ok, err := s.hasPermission(usr, reqperm); err != nil || !ok {
			return "", status.Errorf(codes.PermissionDenied, "%s does not have %v permission", usr, reqperm)
		}
	}

	return usr, nil
}
