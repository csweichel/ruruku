package kvsession

import (
	"sort"
	"testing"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/golang/mock/gomock"
)

func TestValidClaim(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	sreq := validStartSessionRequest()
	sresp, err := s.Start(reqval.GetContext("user", types.PermissionSessionStart), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	rreq0 := validRegistrationRequest(sresp.Id)
	if _, err := s.Register(reqval.GetContext("user00", types.PermissionSessionContribute), rreq0); err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	rreq1 := &api.RegistrationRequest{Session: sresp.Id}
	if _, err := s.Register(reqval.GetContext("user01", types.PermissionSessionContribute), rreq1); err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	resp, err := s.Claim(reqval.GetContext("user00", types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    sresp.Id,
		Claim:      true,
		TestcaseID: sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Claim return an error despite a valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Claim did not return a response despite a valid request")
	}

	resp, err = s.Claim(reqval.GetContext("user01", types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    sresp.Id,
		Claim:      true,
		TestcaseID: sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Claim return an error despite a valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Claim did not return a response despite a valid request")
	}

	statusResp, err := s.Status(reqval.GetContext("user", types.PermissionSessionView), validStatusRequest(sresp.Id))
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	cases := statusResp.Status.Case
	sort.Slice(cases, func(i, j int) bool { return cases[i].Case.Id < cases[j].Case.Id })

	if len(cases[0].Claim) == 0 {
		t.Errorf("Claim did not register test case claim")
	}
	claims := cases[0].Claim
	sort.Slice(claims, func(i, j int) bool { return claims[i].Name < claims[j].Name })
	if len(claims) < 2 {
		t.Errorf("Claim returned wrong number of claims: %d instead of 2", len(claims))
	} else {
		if claims[0].Name != "user00" {
			t.Errorf("Claim returned wrong participant claim name: expected %s, actual %s", "user00", claims[0].Name)
		}
		if claims[1].Name != "user01" {
			t.Errorf("Claim returned wrong participant claim name: expected %s, actual %s", "user01", claims[1].Name)
		}
	}

	resp, err = s.Claim(reqval.GetContext("user00", types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    sresp.Id,
		Claim:      false,
		TestcaseID: sreq.Plan.Case[0].Id,
	})
	if err != nil {
		t.Errorf("Claim return an error despite a valid request: %v", err)
	}
	if resp == nil {
		t.Errorf("Claim did not return a response despite a valid request")
	}

	statusResp, err = s.Status(reqval.GetContext("user", types.PermissionSessionView), validStatusRequest(sresp.Id))
	if err != nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	cases = statusResp.Status.Case
	sort.Slice(cases, func(i, j int) bool { return cases[i].Case.Id < cases[j].Case.Id })
	if len(cases) < 2 {
		t.Errorf("Status returned wrong number of testcases: %d instead of 2", len(cases))
	} else {
		if len(cases[0].Claim) != 1 {
			t.Errorf("Claim did not unregister previous claim")
		} else {
			if cases[0].Claim[0].Name != "user01" {
				t.Errorf("Claim unregistered the wrong claim")
			}
		}
	}
}

func TestInvalidClaim(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	sreq := validStartSessionRequest()
	sresp, err := s.Start(reqval.GetContext("user", types.PermissionSessionStart), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	resp, err := s.Claim(reqval.GetContext("user", types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    "",
		Claim:      true,
		TestcaseID: sreq.Plan.Case[0].Id,
	})
	if err == nil {
		t.Errorf("Claim did not return an error despite missing session name: %v", err)
	}
	if resp != nil {
		t.Errorf("Claim returned a response despite an invalid request")
	}

	user00 := "user00"
	rreq := validRegistrationRequest(sresp.Id)
	if _, err := s.Register(reqval.GetContext(user00, types.PermissionSessionContribute), rreq); err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}

	resp, err = s.Claim(reqval.GetContext(user00, types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    sresp.Id,
		Claim:      true,
		TestcaseID: "",
	})
	if err == nil {
		t.Errorf("Claim did not return an error despite missing testcase ID: %v", err)
	}
	if resp != nil {
		t.Errorf("Claim returned a response despite an invalid request")
	}

	resp, err = s.Claim(reqval.GetContext(user00, types.PermissionSessionContribute), &api.ClaimRequest{
		Session:    sresp.Id,
		Claim:      true,
		TestcaseID: "does-not-exist",
	})
	if err == nil {
		t.Errorf("Claim did not return an error despite non-existent testcase ID")
	}
	if resp != nil {
		t.Errorf("Claim returned a response despite an invalid request")
	}
}
