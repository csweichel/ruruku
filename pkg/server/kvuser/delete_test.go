package kvuser

import (
	"context"
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"google.golang.org/grpc/codes"
)

// positive:
// root deletes other user: TestDeleteAsRoot
// user with user.delete permission deletes other user: TestDeleteAsUser
// user with user.delete permission deletes themselves: TestDeleteYourself

// negative:
// delete without authentication: TestDeleteWithoutAuthentication
// delete without authoerization (testuser with missing user.delete permission): TestDeleteWithoutAuthorization
// delete root user: TestDeleteInvalidUser
// delete with emoty username: TestDeleteInvalidUser
// delete non-existent user: TestDeleteInvalidUser

func TestDeleteAsRoot(t *testing.T) {
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

	srv.deleteUserAndCheck(t, tkn, testuserName)
}

func TestDeleteAsUser(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionUserDelete)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	const testuser = "foo"
	srv.AddUser(testuser, "bar", "foo@bar.com")

	srv.deleteUserAndCheck(t, tkn, testuser)
}

func TestDeleteYourself(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionUserDelete)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	srv.deleteUserAndCheck(t, tkn, testuserName)

	if usr, err := srv.getUserFromToken(tkn); err != nil {
		t.Errorf("Cannot check if token is still valid: %v", err)
	} else if usr != "" {
		t.Errorf("Token resolved to valid user (%s) despite user being deleted", usr)
	}
}

func TestDeleteWithoutAuthentication(t *testing.T) {
	srv := newTestUserService()

	_, err := srv.newTestUserWithPermission(types.PermissionNone)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	srv.deleteUserAndCheckNegative(t, context.Background(), testuserName, codes.Unauthenticated)
}

func TestDeleteWithoutAuthorization(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionNone)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	srv.deleteUserAndCheckNegative(t, newAuthorizedContext(tkn), testuserName, codes.PermissionDenied)
}

func TestDeleteInvalidUser(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionUserDelete)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	srv.deleteUserAndCheckNegative(t, newAuthorizedContext(tkn), "root", codes.PermissionDenied)
	srv.deleteUserAndCheckNegative(t, newAuthorizedContext(tkn), "", codes.NotFound)
	srv.deleteUserAndCheckNegative(t, newAuthorizedContext(tkn), "does-not-exist", codes.NotFound)
}

func (srv *kvuserStore) deleteUserAndCheck(t *testing.T, tkn, username string) {
	resp, err := srv.Delete(newAuthorizedContext(tkn), &api.DeleteUserRequest{
		Username: testuserName,
	})
	if err != nil {
		t.Errorf("Delete returned an error despite valid request: %v", err)
		return
	}
	if resp == nil {
		t.Errorf("Delete did not return a response despite valid request")
		return
	}

	if exists, err := srv.UserExists(testuserName); err != nil {
		t.Errorf("Cannot check if testuser still exists")
	} else if exists {
		t.Errorf("Testuser still exists despite successful call to delete")
	}
}

func (srv *kvuserStore) deleteUserAndCheckNegative(t *testing.T, ctx context.Context, username string, expectedCode codes.Code) {
	resp, err := srv.Delete(ctx, &api.DeleteUserRequest{
		Username: username,
	})
	testNegativeResponse(t, "Delete", expectedCode, resp == nil, err)
}
