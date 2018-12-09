package kvuser

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"google.golang.org/grpc/codes"
	"testing"
)

// positive:
// root grants permissions to other user: TestGrantAsRoot
// user with user.grant permission grants permissions to other users: TestGrantAsUser
// user with user.grant permission grants permissions themselves: TestGrantToThemselves

// negative:
// grant without authentication: TestGrantNoAuthorization
// grant without authorization (testuser with missing user.delete permission): TestGrantNotAuthorized
// grant with emoty username or on non-existent user

func newValidGrantRequest(user string) *api.GrantPermissionsRequest {
	return &api.GrantPermissionsRequest{
		Username:   user,
		Permission: []api.Permission{api.Permission_USER_ADD, api.Permission_USER_DELETE},
	}
}

func TestGrantAsRoot(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.GetUserToken("root")
	if err != nil {
		t.Errorf("Cannot get root token: %v", err)
		return
	}

	if _, err := srv.newTestUserWithPermission(types.PermissionNone); err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	checkGrantPositive(t, srv, tkn, testuserName)
}

func TestGrantAsUser(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionUserGrant)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	if err := srv.addUser("foo", "bar", "foo@bar.com"); err != nil {
		t.Errorf("Cannot add test user: %v", err)
	}

	checkGrantPositive(t, srv, tkn, "foo")
}

func TestGrantToThemselves(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionUserGrant)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	checkGrantPositive(t, srv, tkn, testuserName)
}

func checkGrantPositive(t *testing.T, srv *kvuserStore, tkn string, user string) {
	req := newValidGrantRequest(user)
	resp, err := srv.Grant(newAuthorizedContext(tkn), req)
	if err != nil {
		t.Errorf("Grant returned an error despite valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Grant did not return a response despite valid request")
	}

	for _, perm := range req.Permission {
		if ok, err := srv.hasPermission(user, perm.Convert()); err != nil {
			t.Errorf("Cannot check user permissions: %v", err)
		} else if !ok {
			t.Errorf("Grant did not add permission %v", perm)
		}
	}
}

func TestGrantNoAuthorization(t *testing.T) {
	srv := newTestUserService()

	resp, err := srv.Grant(context.Background(), newValidGrantRequest(testuserName))
	testNegativeResponse(t, "Grant", codes.Unauthenticated, resp, err)
}

func TestGrantNotAuthorized(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionNone)
	if err != nil {
		t.Errorf("Cannot get test user: %v", err)
		return
	}

	resp, err := srv.Grant(newAuthorizedContext(tkn), newValidGrantRequest(testuserName))
	testNegativeResponse(t, "Grant", codes.PermissionDenied, resp, err)
}

func TestGrantInvalidRequest(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionUserAdd)
	if err != nil {
		t.Errorf("Cannot get test user: %v", err)
		return
	}

	resp, err := srv.Grant(newAuthorizedContext(tkn), newValidGrantRequest(""))
	testNegativeResponse(t, "Grant", codes.InvalidArgument, resp, err)
	resp, err = srv.Grant(newAuthorizedContext(tkn), newValidGrantRequest("does-not-exist"))
	testNegativeResponse(t, "Grant", codes.NotFound, resp, err)
}
