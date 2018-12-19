package kvsession

//go:generate mockgen -package kvsession -destination mock_sessionService_listserver.go github.com/32leaves/ruruku/pkg/api/v1 SessionService_ListServer

import (
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
		Annotations: map[string]string{
			"commit":  "91816720613ad4ecbd0198246aa6f44c264c129f",
			"version": "v0.3.0",
		},
	}
}

func TestStartValidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	user := "user"
	req := validStartSessionRequest()
	resp, err := s.Start(reqval.GetContext(user, types.PermissionSessionStart), req)
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

	listsrv := NewMockSessionService_ListServer(ctrl)
	listsrv.EXPECT().Context().Return(reqval.GetContext(user, types.PermissionSessionView))
	listsrv.EXPECT().Send(&api.ListSessionsResponse{Id: resp.Id, Name: req.Name, IsOpen: true}).Return(nil)
	err = s.List(&api.ListSessionsRequest{}, listsrv)
	if err != nil {
		t.Errorf("List returned an error despite valie request: %v", err)
	}
}

func TestStartInvalidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	user := "user"
	req := &api.StartSessionRequest{
		Name: "",
	}
	resp, err := s.Start(reqval.GetContext(user, types.PermissionSessionStart), req)
	if err == nil {
		t.Errorf("Start accepted invalid session (name == \"\")")
	}
	if resp != nil {
		t.Errorf("Start returned a response despite an invalid request")
	}

	req = validStartSessionRequest()
	req.Plan.Case[0].Name = ""
	resp, err = s.Start(reqval.GetContext(user, types.PermissionSessionStart), req)
	if err == nil {
		t.Errorf("Start accepted invalid testcase (name == \"\")")
	}
	if resp != nil {
		t.Errorf("Start returned a response despite an invalid request")
	}

	req = validStartSessionRequest()
	req.Plan.Case[0].Id = req.Plan.Case[1].Id
	resp, err = s.Start(reqval.GetContext(user, types.PermissionSessionStart), req)
	if err == nil {
		t.Errorf("Start accepted invalid testcases (two cases with ID \"%s\")", req.Plan.Case[0].Id)
	}
	if resp != nil {
		t.Errorf("Start returned a response despite an invalid request")
	}
}
