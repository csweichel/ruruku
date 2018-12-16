package kvuser

import (
	"context"
	"os"
	"testing"

	"github.com/32leaves/ruruku/pkg/types"
	bolt "github.com/etcd-io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	testuserName     = "testuser"
	testuserPassword = "password"
	testuserEmail    = "email@mail.com"
)

func newTestUserService() *kvuserStore {
	if _, err := os.Stat("test.db"); err == nil {
		os.Remove("test.db")
	}

	db, err := bolt.Open("test.db", 0666, nil)
	if err != nil {
		panic(err)
	}

	s, err := NewUserStore(db)
	if err != nil {
		panic(err)
	}

	return s
}

func (s *kvuserStore) newTestUserWithPermission(permission types.Permission) (string, error) {
	if err := s.AddUser(testuserName, testuserPassword, testuserEmail); err != nil {
		return "", err
	}
	if permission != types.PermissionNone {
		if err := s.AddPermissions(testuserName, []types.Permission{permission}); err != nil {
			return "", err
		}
	}
	tkn, err := s.GetUserToken(testuserName)
	if err != nil {
		return "", err
	}
	return tkn, nil
}

func newAuthorizedContext(tkn string) context.Context {
	md := metadata.New(map[string]string{"authorization": tkn})
	return metadata.NewIncomingContext(context.Background(), md)
}

func testNegativeResponse(t *testing.T, operation string, expectedCode codes.Code, respNil bool, err error) {
	if err == nil {
		t.Errorf("%s did not return an error despite invalid request", operation)
	} else {
		stat, _ := status.FromError(err)
		if stat.Code() != expectedCode {
			t.Errorf("%s did not return an %v code, but %v", operation, expectedCode, stat.Code())
		}
	}
	if !respNil {
		t.Errorf("%s returned a response despite invalid request", operation)
	}
}
