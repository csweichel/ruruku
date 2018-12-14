package kvsession

import (
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/golang/mock/gomock"
)

var validRegistrationRequest = func(sessionID string) *api.RegistrationRequest {
	return &api.RegistrationRequest{
		Session: sessionID,
	}
}

func TestValidRegistration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	user := "user"
	sresp, err := s.Start(reqval.GetContext(user, types.PermissionSessionStart), validStartSessionRequest())
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	sid := sresp.Id
	req := validRegistrationRequest(sid)
	resp, err := s.Register(reqval.GetContext(user, types.PermissionSessionContribute), req)
	if err != nil {
		t.Errorf("Register returned error despite valid request: %v", err)
	}
	if resp == nil {
		t.Error("Register returned nil response")
	}
}

func TestDuplicateRegistration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	user := "user"
	sresp, err := s.Start(reqval.GetContext(user, types.PermissionSessionStart), validStartSessionRequest())
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	sid := sresp.Id
	req := validRegistrationRequest(sid)
	resp00, err := s.Register(reqval.GetContext(user, types.PermissionSessionContribute), req)
	if err != nil {
		t.Errorf("Register returned error despite valid request: %v", err)
	}
	if resp00 == nil {
		t.Errorf("Register returned nil response despite valid requesst")
	}
	resp01, err := s.Register(reqval.GetContext(user, types.PermissionSessionContribute), validRegistrationRequest(sid))
	if err != nil {
		t.Errorf("Register did not return an error despite invalid request")
	}
	if resp01 == nil {
		t.Errorf("Register returned nil response despite valid requesst")
	}
}
