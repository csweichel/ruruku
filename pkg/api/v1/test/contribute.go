package test

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"sort"
	"testing"
)

func validContributionRequest(token string, tcid string) *api.ContributionRequest {
	return &api.ContributionRequest{
		ParticipantToken: token,
		TestcaseID:       tcid,
		Comment:          "my comment",
		Result:           api.TestRunState_UNDECIDED,
	}
}

func RunTestValidContribution(t *testing.T, s api.SessionServiceServer) {
	sreq := validStartSessionRequest()
	sresp, err := s.Start(context.Background(), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	rreq0 := validRegistrationRequest(sresp.Id)
	rresp0, err := s.Register(context.Background(), rreq0)
	if err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}
	rreq1 := &api.RegistrationRequest{Name: "ZZlast-tester", SessionID: sresp.Id}
	rresp1, err := s.Register(context.Background(), rreq1)
	if err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	_, err = s.Claim(context.Background(), &api.ClaimRequest{
		Claim:            true,
		ParticipantToken: rresp0.Token,
		TestcaseID:       sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Cannot claim testcase: %v", err)
	}
	_, err = s.Claim(context.Background(), &api.ClaimRequest{
		Claim:            true,
		ParticipantToken: rresp1.Token,
		TestcaseID:       sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Cannot claim testcase: %v", err)
	}

	req0 := validContributionRequest(rresp0.Token, sreq.Plan.Case[0].Id)
	_, err = s.Contribute(context.Background(), req0)
	if err != nil {
		t.Errorf("Contribute returned an error despite a valid request: %v", err)
	}
	req1 := validContributionRequest(rresp1.Token, sreq.Plan.Case[0].Id)
	_, err = s.Contribute(context.Background(), req1)
	if err != nil {
		t.Errorf("Contribute returned an error despite a valid request: %v", err)
	}

	statusResp, err := s.Status(context.Background(), validStatusRequest(sresp.Id))
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
	sort.Slice(contributions, func(i, j int) bool { return contributions[i].Participant.Name > contributions[j].Participant.Name })

	if contributions[0].Participant.Name != rreq0.Name {
		t.Errorf("Contribute did not record participant correctly. %v expect, %v actual", rreq0.Name, contributions[0].Participant.Name)
	}
	if contributions[0].Comment != req0.Comment {
		t.Errorf("Contribute did not record comment correctly. %v expected, %v actual", req0.Comment, contributions[0].Comment)
	}
	if contributions[0].State != req0.Result {
		t.Errorf("Contribute did not record result correctly. %v expected, %v actual", req0.Result, contributions[0].State)
	}

	newcomment := "this is a new comment"
	req0.Comment = newcomment
	_, err = s.Contribute(context.Background(), req0)
	if err != nil {
		t.Errorf("Contribute returned an error despite a valid request: %v", err)
	}
	statusResp, err = s.Status(context.Background(), validStatusRequest(sresp.Id))
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
	sort.Slice(contributions, func(i, j int) bool { return contributions[i].Participant.Name > contributions[j].Participant.Name })
	if contributions[0].Comment != req0.Comment {
		t.Errorf("Contribute did not record comment correctly. %v expected, %v actual", req0.Comment, contributions[0].Comment)
	}
}

func RunTestInvalidContribution(t *testing.T, s api.SessionServiceServer) {
	sreq := validStartSessionRequest()
	sresp, err := s.Start(context.Background(), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	_, err = s.Contribute(context.Background(), validContributionRequest("does-not-exist", sreq.Plan.Case[0].Id))
	if err == nil {
		t.Errorf("Contribute did not return an error despite invalid request")
	}

	rreq := validRegistrationRequest(sresp.Id)
	rresp, err := s.Register(context.Background(), rreq)
	if err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	_, err = s.Contribute(context.Background(), validContributionRequest(rresp.Token, sreq.Plan.Case[0].Id))
	if err == nil {
		t.Errorf("Contribute did not return an error even though participant did not claim the testcase")
	}

	_, err = s.Contribute(context.Background(), validContributionRequest(rresp.Token, "does-not-exist"))
	if err == nil {
		t.Errorf("Contribute did not return an error despite contributing to non-existent testcase")
	}
}
