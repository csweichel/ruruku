package kvsession

import (
	"context"
	"sort"
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/golang/mock/gomock"
)

var validStatusRequest = func(sessionID string) *api.SessionStatusRequest {
	return &api.SessionStatusRequest{
		Id: sessionID,
	}
}

func TestBasicStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	reqval.EXPECT().ValidUserFromRequest(gomock.Any(), types.PermissionSessionStart).Return("user", nil)
	sreq := validStartSessionRequest()
	sresp, err := s.Start(context.Background(), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	reqval.EXPECT().ValidUserFromRequest(gomock.Any(), types.PermissionSessionView).Return("user", nil)
	resp, err := s.Status(context.Background(), validStatusRequest(sresp.Id))
	if err != nil {
		t.Errorf("Status returned an error despite valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Status did not return a response despite valid request")
		return
	}

	status := resp.Status
	if status.Id != sresp.Id {
		t.Errorf("Status returned wrong session ID: expected %s, actual %s", sresp.Id, status.Id)
	}
	if status.Name != sreq.Name {
		t.Errorf("Status returned wrong session name: expected %s, actual %s", sreq.Name, status.Name)
	}
	if status.PlanID != sreq.Plan.Id {
		t.Errorf("Status returned wrong planID: expected %s, actual %s", sreq.Plan.Id, status.PlanID)
	}
	if status.Open != true {
		t.Errorf("Status returned wrong open flag: expected %v, actual %v", true, status.Open)
	}
	if status.State != api.TestRunState_UNDECIDED {
		t.Errorf("Session without a single run should have \"undecided\" as state, not %s", status.State.String())
	}

	sort.Slice(status.Status, func(i, j int) bool { return status.Status[i].Case.Id < status.Status[j].Case.Id })
	for i, tcs := range status.Status {
		orig := sreq.Plan.Case[i]
		if tcs.Case.Group != orig.Group {
			t.Errorf("Status returned wrong group for TC %d: expected %s, actual: %s", i, orig.Group, tcs.Case.Group)
		}
		if tcs.Case.Id != orig.Id {
			t.Errorf("Status returned wrong ID for TC %d: expected %s, actual: %s", i, orig.Id, tcs.Case.Id)
		}
		if tcs.Case.Name != orig.Name {
			t.Errorf("Status returned wrong name for TC %d: expected %s, actual: %s", i, orig.Name, tcs.Case.Name)
		}
		if tcs.Case.MinTesterCount != orig.MinTesterCount {
			t.Errorf("Status returned wrong minTesterCount for TC %d: expected %d, actual: %d", i, orig.MinTesterCount, tcs.Case.MinTesterCount)
		}
		if tcs.Case.MustPass != orig.MustPass {
			t.Errorf("Status returned wrong mustPass for TC %d: expected %v, actual: %v", i, orig.MustPass, tcs.Case.MustPass)
		}
		if tcs.Case.Steps != orig.Steps {
			t.Errorf("Status returned wrong steps for TC %d: expected %s, actual: %s", i, orig.Steps, tcs.Case.Steps)
		}
		if tcs.Case.Description != orig.Description {
			t.Errorf("Status returned wrong group for TC %d: expected %s, actual: %s", i, orig.Description, tcs.Case.Description)
		}

		if tcs.State != api.TestRunState_UNDECIDED {
			t.Errorf("Status returned wrong state for TC %d: without a single run it should have \"undecided\", not %s", i, tcs.State)
		}
		if len(tcs.Claim) != 0 {
			t.Errorf("Status returned wrong claim list for TC %d: without a claim it should be empty", i)
		}
	}
}
