package kvuser

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"google.golang.org/grpc/codes"
	"testing"
)

// Positive:
// change password on yourself: TestChangePasswordYourself
// root change password on other user: TestChangePasswordRoot
// user with USER_CHPWD permission change password on different user: TestChangePasswordOthers

// Negative:
// no authentication: TestChangePasswordNoAuthentication
// no authorization: TestChangePasswordNoAuthorization
// change password on root: TestChangePasswordOnRoot
// change password on non-existent user: TestChangePasswordOnNonExistentUser

func TestChangePasswordYourself(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionNone)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	srv.changePasswordAndCheckPositive(t, tkn, testuserName)
}

func TestChangePasswordRoot(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.GetUserToken("root")
	if err != nil {
		t.Errorf("Cannot get root token: %v", err)
		return
	}

	_, err = srv.newTestUserWithPermission(types.PermissionNone)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	srv.changePasswordAndCheckPositive(t, tkn, testuserName)
}

func TestChangePasswordOthers(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionUserChpwd)
	if err != nil {
		t.Errorf("Cannot add test user: %v", err)
		return
	}

	if err := srv.addUser("foo", "bar", "foo@bar.com"); err != nil {
		t.Errorf("Cannot add another test user: %v", err)
	}

	srv.changePasswordAndCheckPositive(t, tkn, "foo")
}

func (srv *kvuserStore) changePasswordAndCheckPositive(t *testing.T, tkn, username string) {
	resp, err := srv.ChangePassword(newAuthorizedContext(tkn), &api.ChangePasswordRequest{
		Username:    username,
		NewPassword: "new-password",
	})
	if err != nil {
		t.Errorf("ChangePassword returned an error despite a valid request: %v", err)
	}
	if resp == nil {
		t.Error("ChangePassword did not return a response despite a valid request")
	}

	if ok, err := srv.validatePassword("foo", "new-password"); err != nil {
		t.Errorf("Error while validating password: %v", err)
	} else if !ok {
		t.Errorf("ChangePassword did not update password correctly. Could not validate.")
	}
}

func TestChangePasswordNoAuthentication(t *testing.T) {
	srv := newTestUserService()

	_, err := srv.newTestUserWithPermission(types.PermissionUserChpwd)
	if err != nil {
		t.Errorf("Cannot add test user: %v", err)
		return
	}

	resp, err := srv.ChangePassword(context.Background(), &api.ChangePasswordRequest{Username: testuserName, NewPassword: "foobar"})
	testNegativeResponse(t, "ChangePassword", codes.Unauthenticated, resp, err)
}

func TestChangePasswordNoAuthorization(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionNone)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	if err := srv.addUser("foo", "bar", "foo@bar.com"); err != nil {
		t.Errorf("Cannot create test user foo")
	}

	resp, err := srv.ChangePassword(newAuthorizedContext(tkn), &api.ChangePasswordRequest{Username: "foo", NewPassword: "foobar"})
	testNegativeResponse(t, "ChangePassword", codes.PermissionDenied, resp, err)
}

func TestChangePasswordOnRoot(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.GetUserToken("root")
	if err != nil {
		t.Errorf("Cannot get root token: %v", err)
		return
	}

	resp, err := srv.ChangePassword(newAuthorizedContext(tkn), &api.ChangePasswordRequest{Username: "root", NewPassword: "foobar"})
	testNegativeResponse(t, "ChangePassword", codes.PermissionDenied, resp, err)
}

func TestChangePasswordOnNonExistentUser(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.GetUserToken("root")
	if err != nil {
		t.Errorf("Cannot get root token: %v", err)
		return
	}

	resp, err := srv.ChangePassword(newAuthorizedContext(tkn), &api.ChangePasswordRequest{Username: "does-not-exist", NewPassword: "foobar"})
	testNegativeResponse(t, "ChangePassword", codes.NotFound, resp, err)
}
