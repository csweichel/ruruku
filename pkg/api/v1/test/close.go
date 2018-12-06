package test

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"testing"
)

func RunTestValidClose(t *testing.T, s api.SessionServiceServer) {
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

	resp, err := s.Close(context.Background(), &api.CloseSessionRequest{Id: sresp.Id})
	if err != nil {
		t.Errorf("Close returned an error despite sending a valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Close returned nil despite sending a valid request")
	}

	statusResp, err := s.Status(context.Background(), validStatusRequest(sresp.Id))
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	if statusResp.Status.Open {
		t.Errorf("Close did not close session")
	}

	_, err = s.Claim(context.Background(), &api.ClaimRequest{
		Claim:            true,
		ParticipantToken: rresp0.Token,
		TestcaseID:       sreq.Plan.Case[0].Id,
	})
	if err == nil {
		t.Errorf("Claim did not return an error despite claiming a testcase on a closed session")
	}

	rreq1 := &api.RegistrationRequest{Name: "ZZlast-tester", SessionID: sresp.Id}
	_, err = s.Register(context.Background(), rreq1)
	if err == nil {
		t.Errorf("Register was able to join a closed session")
	}
}

func RunTestInvalidClose(t *testing.T, s api.SessionServiceServer) {
	resp, err := s.Close(context.Background(), &api.CloseSessionRequest{Id: "does-not-exist"})
	if err == nil {
		t.Errorf("Close returned did not return an error despite sending an invalid request")
	}
	if resp != nil {
		t.Errorf("Close returned a response despite sending a invalid request")
	}

	resp, err = s.Close(context.Background(), &api.CloseSessionRequest{Id: ""})
	if err == nil {
		t.Errorf("Close returned did not return an error despite sending an invalid request")
	}
	if resp != nil {
		t.Errorf("Close returned a response despite sending a invalid request")
	}
}
