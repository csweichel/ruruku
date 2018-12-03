package test

import (
	"context"
	api "github.com/32leaves/ruruku/pkg/api/v1"
	"testing"
)

var validRegistrationRequest = func(sessionID string) *api.RegistrationRequest {
	return &api.RegistrationRequest{
		Name:      "tester",
		SessionID: sessionID,
	}
}

func RunTestValidRegistration(t *testing.T, s api.SessionServiceServer) {
	sresp, err := s.Start(context.Background(), validStartSessionRequest())
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	sid := sresp.Id
	resp, err := s.Register(context.Background(), validRegistrationRequest(sid))
	if err != nil {
		t.Errorf("Register returned error despite valid request: %v", err)
	}
	if resp == nil {
		t.Error("Register returned nil response")
	}
	if resp.Token == "" {
		t.Error("Register returned empty user token")
	}
}

func RunTestDuplicateRegistration(t *testing.T, s api.SessionServiceServer) {
	sresp, err := s.Start(context.Background(), validStartSessionRequest())
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	sid := sresp.Id
	resp00, err := s.Register(context.Background(), validRegistrationRequest(sid))
	if err != nil {
		t.Errorf("Register returned error despite valid request: %v", err)
	}
	resp01, err := s.Register(context.Background(), validRegistrationRequest(sid))
	if err != nil {
		t.Errorf("Register did not return an error despite invalid request")
	}

	if resp00.Token != resp01.Token {
		t.Errorf("Register returned different tokens for the same user: %s != %s", resp00.Token, resp01.Token)
	}
}

func RunTestInvalidRegistration(t *testing.T, s api.SessionServiceServer) {
	_, err := s.Register(context.Background(), &api.RegistrationRequest{
		Name:      "foobar",
		SessionID: "doesNotExist",
	})
	if err == nil {
		t.Errorf("Register did not return an error despite non-existent session")
	}

	resp, err := s.Start(context.Background(), validStartSessionRequest())
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	_, err = s.Register(context.Background(), &api.RegistrationRequest{
		Name:      "",
		SessionID: resp.Id,
	})
	if err == nil {
		t.Errorf("Register did not return an error empty user name")
	}
}
