package kvuser

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"google.golang.org/grpc/codes"
	"testing"
)

// positive:
// authenticate valid non-root user: TestAuthenticate

// negative:
// authenticate with wrong password: TestAuthenticateWrongPassword
// authenticate with non-existent user: TestAuthenticateNonExistentUser
// authenticate with root user: TestAuthenticateRootUser

func TestAuthenticate(t *testing.T) {
	srv := newTestUserService()

	if _, err := srv.newTestUserWithPermission(types.PermissionNone); err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	resp, err := srv.AuthenticateCredentials(context.Background(), &api.AuthenticationRequest{
		Username: testuserName,
		Password: testuserPassword,
	})
	if err != nil {
		t.Errorf("AuthenticateCredentials returned an error despite valid request: %v", err)
		return
	}
	if resp == nil {
		t.Errorf("AuthenticateCredentials did not return a response despite valid request")
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

func TestAuthenticateWrongPassword(t *testing.T) {
	srv := newTestUserService()

	if _, err := srv.newTestUserWithPermission(types.PermissionNone); err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	resp, err := srv.AuthenticateCredentials(context.Background(), &api.AuthenticationRequest{
		Username: testuserName,
		Password: testuserPassword + "isWrong",
	})
	testNegativeResponse(t, "AuthenticateCredentials", codes.Unauthenticated, resp, err)
}

func TestAuthenticateNonExistentUser(t *testing.T) {
	srv := newTestUserService()

	if _, err := srv.newTestUserWithPermission(types.PermissionNone); err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	resp, err := srv.AuthenticateCredentials(context.Background(), &api.AuthenticationRequest{
		Username: testuserName + "isWrong",
		Password: testuserPassword,
	})
	testNegativeResponse(t, "AuthenticateCredentials", codes.Unauthenticated, resp, err)
}

func TestAuthenticateRootUser(t *testing.T) {
	srv := newTestUserService()

	resp, err := srv.AuthenticateCredentials(context.Background(), &api.AuthenticationRequest{
		Username: "root",
		Password: "",
	})
	testNegativeResponse(t, "AuthenticateCredentials", codes.Unauthenticated, resp, err)
}
