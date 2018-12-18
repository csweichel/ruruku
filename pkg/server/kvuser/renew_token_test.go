package kvuser

import (
	"context"
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"google.golang.org/grpc/codes"
)

// positive:
// renew token with valid authentication: TestRenewToken

// negative:
// renew token without valid authentication: TestRenewTokenWithInvalidAuthentication

func TestRenewToken(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionUserList)
	if err != nil {
		t.Errorf("Cannot add test user: %v", err)
		return
	}

	resp, err := srv.RenewToken(newAuthorizedContext(tkn), &api.RenewTokenRequest{})
	if err != nil {
		t.Errorf("RenewToken returned an error despite valid request: %v", err)
		return
	}
	if resp == nil {
		t.Errorf("RenewToken did not return a response despite valid request")
		return
	}

	usr, err := srv.getUserFromToken(resp.Token)
	if err != nil {
		t.Errorf("Cannot get user from token: %v", err)
	}
	if usr == "" {
		t.Errorf("Could not find user from token")
	}
	if usr != testuserName {
		t.Errorf("Token resolved to %s instead of %s", usr, testuserName)
	}
}

func TestRenewTokenWithInvalidAuthentication(t *testing.T) {
	srv := newTestUserService()

	resp, err := srv.RenewToken(context.Background(), &api.RenewTokenRequest{})
	testNegativeResponse(t, "RenewToken", codes.Unauthenticated, resp == nil, err)

	resp, err = srv.RenewToken(newAuthorizedContext("this-token-does-not-exist"), &api.RenewTokenRequest{})
	testNegativeResponse(t, "RenewToken", codes.Unauthenticated, resp == nil, err)
}
