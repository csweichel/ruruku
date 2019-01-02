package kvsession

import (
	"sort"
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/golang/mock/gomock"
)

var (
	gh70plan = types.TestPlan{
		ID:   "abc",
		Name: "This is a test 1",
		Case: []types.Testcase{
			{
				ID:    "1",
				Name:  "test1",
				Group: "one",
			},
			{
				ID:    "10",
				Name:  "test1",
				Group: "one",
			},
			{
				ID:    "100",
				Name:  "test1",
				Group: "one",
			},
		},
	}
)

func TestGH70(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	user := "user"
	sreq := &api.StartSessionRequest{
		Name: "foo",
		Plan: api.ConvertTestPlan(&gh70plan),
	}
	sresp, err := s.Start(reqval.GetContext(user, types.PermissionSessionStart), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}
	if sresp == nil {
		t.Errorf("Start returned nil response")
		return
	}

	if _, err := s.Register(reqval.GetContext(user, types.PermissionSessionContribute), &api.RegistrationRequest{Session: sresp.Id}); err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	_, err = s.Claim(reqval.GetContext(user, types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    sresp.Id,
		Claim:      true,
		TestcaseID: "100",
	})
	if err != nil {
		t.Errorf("Cannot claim testcase: %v", err)
	}

	statusResp, err := s.Status(reqval.GetContext(user, types.PermissionSessionView), &api.SessionStatusRequest{Id: sresp.Id})
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	status := statusResp.Status
	sort.Slice(status.Case, func(i, j int) bool { return status.Case[i].Case.Id < status.Case[j].Case.Id })
	if len(status.Case[0].Claim) > 0 {
		t.Errorf("Testcase is claimed but should not be: %v", status.Case[0])
	}
	if len(status.Case[1].Claim) > 0 {
		t.Errorf("Testcase is claimed but should not be: %v", status.Case[1])
	}
	if len(status.Case[2].Claim) == 0 {
		t.Errorf("Testcase is not claimed but should be: %v", status.Case[2])
	}

	_, err = s.Contribute(reqval.GetContext(user, types.PermissionSessionContribute), &api.ContributionRequest{
		Session:    sresp.Id,
		TestcaseID: "100",
		Result:     api.TestRunState_PASSED,
	})
	if err != nil {
		t.Errorf("Cannot contribute: %v", err)
		return
	}

	statusResp, err = s.Status(reqval.GetContext(user, types.PermissionSessionView), &api.SessionStatusRequest{Id: sresp.Id})
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	status = statusResp.Status
	sort.Slice(status.Case, func(i, j int) bool { return status.Case[i].Case.Id < status.Case[j].Case.Id })
	if len(status.Case[0].Result) > 0 {
		t.Errorf("Testcase has a contribution but should not: %v", status.Case[0])
	}
	if len(status.Case[1].Result) > 0 {
		t.Errorf("Testcase has a contribution but should not: %v", status.Case[1])
	}
	if len(status.Case[2].Result) == 0 {
		t.Errorf("Testcase does not have a contribution but should have: %v", status.Case[2])
	}
}
