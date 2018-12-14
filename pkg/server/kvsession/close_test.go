package kvsession

import (
	"context"
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/golang/mock/gomock"
)

func TestValidClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	sreq := validStartSessionRequest()
	sresp, err := s.Start(reqval.GetContext("user", types.PermissionSessionStart), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	user := "user"
	rreq0 := validRegistrationRequest(sresp.Id)
	if _, err := s.Register(reqval.GetContext(user, types.PermissionSessionContribute), rreq0); err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	resp, err := s.Close(reqval.GetContext(user, types.PermissionSessionClose), &api.CloseSessionRequest{Id: sresp.Id})
	if err != nil {
		t.Errorf("Close returned an error despite sending a valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Close returned nil despite sending a valid request")
	}

	statusResp, err := s.Status(reqval.GetContext(user, types.PermissionSessionView), validStatusRequest(sresp.Id))
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	if statusResp.Status.Open {
		t.Errorf("Close did not close session")
	}

	_, err = s.Claim(reqval.GetContext(user, types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    sresp.Id,
		Claim:      true,
		TestcaseID: sreq.Plan.Case[0].Id,
	})
	if err == nil {
		t.Errorf("Claim did not return an error despite claiming a testcase on a closed session")
	}

	user01 := "user01"
	rreq1 := &api.RegistrationRequest{Session: sresp.Id}
	_, err = s.Register(reqval.GetContext(user01, types.PermissionSessionContribute), rreq1)
	if err == nil {
		t.Errorf("Register was able to join a closed session")
	}
}

func TestInvalidClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	reqval.EXPECT().ValidUserFromRequest(gomock.Any(), types.PermissionSessionClose).Return("user", nil)
	resp, err := s.Close(context.Background(), &api.CloseSessionRequest{Id: "does-not-exist"})
	if err == nil {
		t.Errorf("Close returned did not return an error despite sending an invalid request")
	}
	if resp != nil {
		t.Errorf("Close returned a response despite sending a invalid request")
	}

	reqval.EXPECT().ValidUserFromRequest(gomock.Any(), types.PermissionSessionClose).Return("user", nil)
	resp, err = s.Close(context.Background(), &api.CloseSessionRequest{Id: ""})
	if err == nil {
		t.Errorf("Close returned did not return an error despite sending an invalid request")
	}
	if resp != nil {
		t.Errorf("Close returned a response despite sending a invalid request")
	}
}
