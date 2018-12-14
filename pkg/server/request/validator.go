package request

import (
	"context"
	"github.com/32leaves/ruruku/pkg/types"
)

type Validator interface {
	// ValidUserFromRequest extracts authentication token from the context (using grpc/metadata)
	// and validates that token against a list of known users. Returns the name of the authenticated
	// user or a gRPC status error if the user was not found or does not have the appropriate permission.
	ValidUserFromRequest(ctx context.Context, reqperm types.Permission) (string, error)
}
