package kvsession

import (
	"context"
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/golang/mock/gomock"
)

var validStartSessionRequest = func() *api.StartSessionRequest {
	return &api.StartSessionRequest{
		Name: "foo",
		Plan: &api.TestPlan{
			Id:          "tp",
			Name:        "testplan",
			Description: "testplan descriptipn",
			Case: []*api.Testcase{
				{
					Id:             "tc00",
					Group:          "grp",
					Name:           "tcname 00",
					Description:    "desc",
					MinTesterCount: 1,
					MustPass:       true,
					Steps:          "steps",
				},
				{
					Id:             "tc01",
					Group:          "grp",
					Name:           "tcname 01",
					Description:    "desc",
					MinTesterCount: 42,
					MustPass:       false,
					Steps:          "steps",
				},
			},
		},
	}
}

func TestStartValidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	reqval.EXPECT().ValidUserFromRequest(gomock.Any(), types.PermissionSessionStart).Return("user", nil)
	resp, err := s.Start(context.Background(), validStartSessionRequest())

	if err != nil {
		t.Errorf("Start returned error despite valid request: %v", err)
	}
	if resp == nil {
		t.Error("Start returned nil response")
		return
	}
	if resp.Id == "" {
		t.Error("Start returned empty ID")
	}
}

func TestStartInvalidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	reqval.EXPECT().ValidUserFromRequest(gomock.Any(), types.PermissionSessionStart).Return("user", nil)
	req := &api.StartSessionRequest{
		Name: "",
	}
	resp, err := s.Start(context.Background(), req)
	if err == nil {
		t.Errorf("Start accepted invalid session (name == \"\")")
	}
	if resp != nil {
		t.Errorf("Start returned a response despite an invalid request")
	}

	req = validStartSessionRequest()
	req.Plan.Case[0].Name = ""
	reqval.EXPECT().ValidUserFromRequest(gomock.Any(), types.PermissionSessionStart).Return("user", nil)
	resp, err = s.Start(context.Background(), req)
	if err == nil {
		t.Errorf("Start accepted invalid testcase (name == \"\")")
	}
	if resp != nil {
		t.Errorf("Start returned a response despite an invalid request")
	}

	req = validStartSessionRequest()
	req.Plan.Case[0].Id = req.Plan.Case[1].Id
	reqval.EXPECT().ValidUserFromRequest(gomock.Any(), types.PermissionSessionStart).Return("user", nil)
	resp, err = s.Start(context.Background(), req)
	if err == nil {
		t.Errorf("Start accepted invalid testcases (two cases with ID \"%s\")", req.Plan.Case[0].Id)
	}
	if resp != nil {
		t.Errorf("Start returned a response despite an invalid request")
	}
}
