package kvuser

import (
    "testing"
    api "github.com/32leaves/ruruku/pkg/api/v1"
    "github.com/32leaves/ruruku/pkg/types"
    "context"
    "google.golang.org/grpc/codes"
)

// positive:
// root lists users: TestListAsRoot
// user with user.list permission lists users: TestListAsUser

// negative:
// list without authentication: TestListNoAuthorization
// list without authorization (testuser with missing user.list permission): TestListNotAuthorized

func TestListAsRoot(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.GetUserToken("root")
	if err != nil {
		t.Errorf("Cannot get root token: %v", err)
		return
	}

	if _, err := srv.newTestUserWithPermission(types.PermissionSessionClose); err != nil {
		t.Errorf("Cannot add test user")
		return
	}

	resp, err := srv.List(newAuthorizedContext(tkn), &api.ListUsersRequest{})
    if err != nil {
        t.Errorf("List returned an error despite valid request: %v", err)
    }
    if resp == nil {
        t.Errorf("List did not return a response despite valid request")
        return
    }

    users := make(map[string]*api.User)
    for _, usr := range resp.User {
        users[usr.Name] = usr
    }

    if usr, ok := users["root"]; !ok {
        t.Error("List did not return root user")
    } else {
        if usr.Email != "" {
            t.Errorf("List returned an email for a root: %v", usr.Email)
        }
        if len(usr.Permission) != len(types.AllPermissions) {
            t.Errorf("List returned wrong number pf permissions for root: expected %d, actual %d", len(types.AllPermissions), len(usr.Permission))
        }
    }

    if usr, ok := users[testuserName]; !ok {
        t.Error("List did not return testuser")
    } else {
        if usr.Email != testuserEmail {
            t.Errorf("List returned wrong email for testuser")
        }
        if len(usr.Permission) != 1 {
            t.Errorf("List returned wrong number of permissions for testuser: %d instead of 1", len(usr.Permission))
        } else if usr.Permission[0] != api.Permission_SESSION_CLOSE {
            t.Errorf("List returned wrong permission for testuser")
        }
    }
}

func TestListAsUser(t *testing.T) {
	srv := newTestUserService()

    tkn, err := srv.newTestUserWithPermission(types.PermissionUserList)
	if err != nil {
		t.Errorf("Cannot add test user")
		return
	}
    if err := srv.addPermissions(testuserName, []types.Permission{types.PermissionUserAdd}); err != nil {
        t.Errorf("Cannot add permission to testuser: %v", err)
        return
    }

	resp, err := srv.List(newAuthorizedContext(tkn), &api.ListUsersRequest{})
    if err != nil {
        t.Errorf("List returned an error despite valid request: %v", err)
    }
    if resp == nil {
        t.Errorf("List did not return a response despite valid request")
        return
    }

    users := make(map[string]*api.User)
    for _, usr := range resp.User {
        users[usr.Name] = usr
    }

    if _, ok := users["root"]; !ok {
        t.Error("List did not return root user")
    }
    if usr, ok := users[testuserName]; !ok {
        t.Error("List did not return testuser")
    } else {
        if usr.Email != testuserEmail {
            t.Errorf("List returned wrong email for testuser")
        }
        if len(usr.Permission) != 2 {
            t.Errorf("List returned wrong number of permissions for testuser: %d instead of 2", len(usr.Permission))
        } else {
            if usr.Permission[0] != api.Permission_USER_ADD && usr.Permission[1] != api.Permission_USER_ADD {
                t.Errorf("List returned wrong permission for testuser")
            }
            if usr.Permission[0] != api.Permission_USER_LIST && usr.Permission[1] != api.Permission_USER_LIST {
                t.Errorf("List returned wrong permission for testuser")
            }
        }
    }
}

func TestListNoAuthorization(t *testing.T) {
	srv := newTestUserService()

	resp, err := srv.List(context.Background(), &api.ListUsersRequest{})
	testNegativeResponse(t, "List", codes.Unauthenticated, resp == nil, err)
}

func TestListNotAuthorized(t *testing.T) {
	srv := newTestUserService()

	tkn, err := srv.newTestUserWithPermission(types.PermissionNone)
	if err != nil {
		t.Errorf("Cannot get test user: %v", err)
		return
	}

	resp, err := srv.List(newAuthorizedContext(tkn), &api.ListUsersRequest{})
	testNegativeResponse(t, "List", codes.PermissionDenied, resp == nil, err)
}
