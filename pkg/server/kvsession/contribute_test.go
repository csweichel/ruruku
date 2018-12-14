package kvsession

import (
	"context"
	"sort"
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/golang/mock/gomock"
)

func validContributionRequest(session string, tcid string) *api.ContributionRequest {
	return &api.ContributionRequest{
		Session:    session,
		TestcaseID: tcid,
		Comment:    "my comment",
		Result:     api.TestRunState_UNDECIDED,
	}
}

func TestValidContribution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	user := "user"
	sreq := validStartSessionRequest()
	sresp, err := s.Start(reqval.GetContext(user, types.PermissionSessionStart), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	user00 := "user00"
	rreq0 := validRegistrationRequest(sresp.Id)
	if _, err := s.Register(reqval.GetContext(user00, types.PermissionSessionContribute), rreq0); err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	user01 := "user01"
	rreq1 := validRegistrationRequest(sresp.Id)
	if _, err := s.Register(reqval.GetContext(user01, types.PermissionSessionContribute), rreq1); err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	_, err = s.Claim(reqval.GetContext(user00, types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    sresp.Id,
		Claim:      true,
		TestcaseID: sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Cannot claim testcase: %v", err)
	}
	_, err = s.Claim(reqval.GetContext(user01, types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    sresp.Id,
		Claim:      true,
		TestcaseID: sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Cannot claim testcase: %v", err)
	}

	req0 := validContributionRequest(sresp.Id, sreq.Plan.Case[0].Id)
	_, err = s.Contribute(reqval.GetContext(user00, types.PermissionSessionContribute), req0)
	if err != nil {
		t.Errorf("Contribute returned an error despite a valid request: %v", err)
	}
	req1 := validContributionRequest(sresp.Id, sreq.Plan.Case[0].Id)
	_, err = s.Contribute(reqval.GetContext(user01, types.PermissionSessionContribute), req1)
	if err != nil {
		t.Errorf("Contribute returned an error despite a valid request: %v", err)
	}

	statusResp, err := s.Status(reqval.GetContext(user, types.PermissionSessionView), validStatusRequest(sresp.Id))
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	status := statusResp.Status.Status
	sort.Slice(status, func(i, j int) bool { return status[i].Case.Id < status[j].Case.Id })
	contributions := status[0].Result
	if len(contributions) == 0 {
		t.Errorf("Contribute did not record contribution")
		return
	}
	sort.Slice(contributions, func(i, j int) bool { return contributions[i].Participant.Name < contributions[j].Participant.Name })

	if contributions[0].Participant.Name != user00 {
		t.Errorf("Contribute did not record participant correctly. %v expect, %v actual", user00, contributions[0].Participant.Name)
	}
	if contributions[0].Comment != req0.Comment {
		t.Errorf("Contribute did not record comment correctly. %v expected, %v actual", req0.Comment, contributions[0].Comment)
	}
	if contributions[0].State != req0.Result {
		t.Errorf("Contribute did not record result correctly. %v expected, %v actual", req0.Result, contributions[0].State)
	}

	newcomment := "this is a new comment"
	req0.Comment = newcomment
	_, err = s.Contribute(reqval.GetContext(user00, types.PermissionSessionContribute), req0)
	if err != nil {
		t.Errorf("Contribute returned an error despite a valid request: %v", err)
	}
	statusResp, err = s.Status(reqval.GetContext(user, types.PermissionSessionView), validStatusRequest(sresp.Id))
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	status = statusResp.Status.Status
	sort.Slice(status, func(i, j int) bool { return status[i].Case.Id < status[j].Case.Id })
	contributions = status[0].Result
	if len(contributions) == 0 {
		t.Errorf("Contribute did not record contribution")
		return
	}
	sort.Slice(contributions, func(i, j int) bool { return contributions[i].Participant.Name < contributions[j].Participant.Name })
	if contributions[0].Comment != req0.Comment {
		t.Errorf("Contribute did not record comment correctly. %v expected, %v actual", req0.Comment, contributions[0].Comment)
	}
}

func TestInvalidContribution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	user := "user"
	sreq := validStartSessionRequest()
	sresp, err := s.Start(reqval.GetContext(user, types.PermissionSessionStart), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	reqval.EXPECT().ValidUserFromRequest(gomock.Any(), types.PermissionSessionContribute).Return("user", nil)
	_, err = s.Contribute(context.Background(), validContributionRequest("does-not-exist", sreq.Plan.Case[0].Id))
	if err == nil {
		t.Errorf("Contribute did not return an error despite invalid request")
	}

	rreq := validRegistrationRequest(sresp.Id)
	if _, err := s.Register(reqval.GetContext(user, types.PermissionSessionContribute), rreq); err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	_, err = s.Contribute(reqval.GetContext(user, types.PermissionSessionContribute), validContributionRequest(sresp.Id, sreq.Plan.Case[0].Id))
	if err == nil {
		t.Errorf("Contribute did not return an error even though participant did not claim the testcase")
	}

	_, err = s.Contribute(reqval.GetContext(user, types.PermissionSessionContribute), validContributionRequest(sresp.Id, "does-not-exist"))
	if err == nil {
		t.Errorf("Contribute did not return an error despite contributing to non-existent testcase")
	}
}
