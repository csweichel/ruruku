package kvsession

//go:generate mockgen -package kvsession -destination mock_request_validator.go github.com/32leaves/ruruku/pkg/server/request Validator

import (
	"context"
	"fmt"
	"os"
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	bolt "github.com/etcd-io/bbolt"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func newTestServer(ctrl *gomock.Controller) (api.SessionServiceServer, *MockValidator) {
	log.SetLevel(log.WarnLevel)

	if _, err := os.Stat("test.db"); err == nil {
		os.Remove("test.db")
	}

	db, err := bolt.Open("test.db", 0666, nil)
	if err != nil {
		panic(err)
	}

	reqval := NewMockValidator(ctrl)

	store, err := NewSession(db, reqval)
	if err != nil {
		panic(err)
	}

	// will return the validator mock in the future
	return store, reqval
}

func (reqval *MockValidator) GetContext(user string, perm types.Permission) context.Context {
	tkn := fmt.Sprintf("token_for::%s", user)
	md := metadata.New(map[string]string{"authorization": tkn})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	reqval.EXPECT().ValidUserFromRequest(ctx, perm).Return(user, nil)

	return ctx
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
