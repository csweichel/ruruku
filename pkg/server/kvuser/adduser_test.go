package kvuser

import (
	"context"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"google.golang.org/grpc/codes"

	"testing"
)

// positive:
// root adds another user: TestAdd
// user with user.add permission adds other user

// negative:
// add without authentication: TestAddNoAuthorization
// add without authorization: TestAddNotAuthorized
// add duplicate user: TestAddDuplicate
// add user without name: TestAddWithInvalidFields
// add user without password: TestAddWithInvalidFields
// add user without email: TestAddWithInvalidFields

func newValidAddUserRequest(name string) *api.AddUserRequest {
	return &api.AddUserRequest{
		Username: name,
		Password: "foobar",
		Email:    "bla@bla.com",
	}
}

func TestAdd(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.GetUserToken("root")
	if err != nil {
		t.Errorf("Cannot get root token: %v", err)
		return
	}

	req := newValidAddUserRequest("test")
	if _, err := srv.Add(newAuthorizedContext(tkn), req); err != nil {
		t.Errorf("Add returned an error despite valid request: %v", err)
	}

	if _, err := srv.GetUserToken(req.Username); err != nil {
		t.Errorf("Could not get a token for the new user: %v", err)
	}
}

func TestAddNoAuthorization(t *testing.T) {
	srv := newTestUserService()

	resp, err := srv.Add(context.Background(), newValidAddUserRequest("test"))
	testNegativeResponse(t, "Add", codes.Unauthenticated, resp == nil, err)
}

func TestAddNotAuthorized(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionNone)
	if err != nil {
		t.Errorf("Cannot get test user: %v", err)
		return
	}

	resp, err := srv.Add(newAuthorizedContext(tkn), newValidAddUserRequest("test"))
	testNegativeResponse(t, "Add", codes.PermissionDenied, resp == nil, err)
}

func TestAddDuplicate(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.GetUserToken("root")
	if err != nil {
		t.Errorf("Cannot get root token: %v", err)
		return
	}

	req := newValidAddUserRequest("test")
	if _, err := srv.Add(newAuthorizedContext(tkn), req); err != nil {
		t.Errorf("Add returned an error despite valid request: %v", err)
	}

	resp, err := srv.Add(newAuthorizedContext(tkn), req)
	testNegativeResponse(t, "Add", codes.AlreadyExists, resp == nil, err)
}

func TestAddWithInvalidFields(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.GetUserToken("root")
	if err != nil {
		t.Errorf("Cannot get root token: %v", err)
		return
	}

	testReq := func(req *api.AddUserRequest, field string) {
		resp, err := srv.Add(newAuthorizedContext(tkn), req)
		testNegativeResponse(t, "Add", codes.InvalidArgument, resp == nil, err)
	}

	testReq(&api.AddUserRequest{Username: "", Password: "bla", Email: "foo@bar.com"}, "username")
	testReq(&api.AddUserRequest{Username: "foo", Password: "", Email: "foo@bar.com"}, "password")
	testReq(&api.AddUserRequest{Username: "foo", Password: "bla", Email: ""}, "email")
	testReq(&api.AddUserRequest{Username: "foo", Password: "bla", Email: "fobar"}, "email")
}
