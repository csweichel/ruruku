package test

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"sort"
	"testing"
)

func RunTestValidClaim(t *testing.T, s api.SessionServiceServer) {
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

	resp, err := s.Claim(context.Background(), &api.ClaimRequest{
		Claim:            true,
		ParticipantToken: rresp0.Token,
		TestcaseID:       sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Claim return an error despite a valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Claim did not return a response despite a valid request")
	}

	resp, err = s.Claim(context.Background(), &api.ClaimRequest{
		Claim:            true,
		ParticipantToken: rresp1.Token,
		TestcaseID:       sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Claim return an error despite a valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Claim did not return a response despite a valid request")
	}

	statusResp, err := s.Status(context.Background(), validStatusRequest(sresp.Id))
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	status := statusResp.Status.Status
	sort.Slice(status, func(i, j int) bool { return status[i].Case.Id < status[j].Case.Id })

	if len(status[0].Claim) == 0 {
		t.Errorf("Claim did not register test case claim")
	}
	claims := status[0].Claim
	sort.Slice(claims, func(i, j int) bool { return claims[i].Name > claims[j].Name })
	if claims[0].Name != rreq0.Name {
		t.Errorf("Claim returned wrong participant claim name: expected %s, actual %s", rreq0.Name, claims[0].Name)
	}
	if claims[1].Name != rreq1.Name {
		t.Errorf("Claim returned wrong participant claim name: expected %s, actual %s", rreq1.Name, claims[1].Name)
	}

	resp, err = s.Claim(context.Background(), &api.ClaimRequest{
		Claim:            false,
		ParticipantToken: rresp0.Token,
		TestcaseID:       sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Claim return an error despite a valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Claim did not return a response despite a valid request")
	}

	statusResp, err = s.Status(context.Background(), validStatusRequest(sresp.Id))
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	status = statusResp.Status.Status
	sort.Slice(status, func(i, j int) bool { return status[i].Case.Id < status[j].Case.Id })
	if len(status[0].Claim) != 1 {
		t.Errorf("Claim did not unregister previous claim")
	}
	if status[0].Claim[0].Name != rreq1.Name {
		t.Errorf("Claim unregistered the wrong claim")
	}
}

func RunTestInvalidClaim(t *testing.T, s api.SessionServiceServer) {
	sreq := validStartSessionRequest()
	sresp, err := s.Start(context.Background(), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	resp, err := s.Claim(context.Background(), &api.ClaimRequest{
		Claim:            true,
		ParticipantToken: "",
		TestcaseID:       sreq.Plan.Case[0].Id,
	})
	if err == nil {
		t.Errorf("Claim did not return an error despite missing participant token: %v", err)
	}
	if resp != nil {
		t.Errorf("Claim returned a response despite an invalid request")
	}

	rreq := validRegistrationRequest(sresp.Id)
	rresp, err := s.Register(context.Background(), rreq)
	if err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	resp, err = s.Claim(context.Background(), &api.ClaimRequest{
		Claim:            true,
		ParticipantToken: rresp.Token,
		TestcaseID:       "",
	})
	if err == nil {
		t.Errorf("Claim did not return an error despite missing testcase ID: %v", err)
	}
	if resp != nil {
		t.Errorf("Claim returned a response despite an invalid request")
	}

	resp, err = s.Claim(context.Background(), &api.ClaimRequest{
		Claim:            true,
		ParticipantToken: rresp.Token,
		TestcaseID:       "does-not-exist",
	})
	if err == nil {
		t.Errorf("Claim did not return an error despite non-existent testcase ID")
	}
	if resp != nil {
		t.Errorf("Claim returned a response despite an invalid request")
	}
}
