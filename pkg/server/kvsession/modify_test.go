package kvsession

import (
	"testing"

	"google.golang.org/grpc/codes"

	api "github.com/32leaves/ruruku/pkg/api/v1"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/golang/mock/gomock"
)

// positive:
// modifiy annotations: TestValidModifyAnnotations
// modify existing testcase: TestValidModifyTestcase
// add testcase: TestValidAddTestcase
// delete testcase: TestValidModifyRemoveTestcase

// negative:
// without authentication: not needed - checked by the reqval context mock
// without authorization: not needed - checked by the reqval context mock
// modify closed session: TestModifyClosedSession
// modify unmodifiable session: TestModifyUnmodifiableSession
// modify non-existing session: TestModifyNonExistentSession
// modify non-existing testcase: TestModifyNonExistentTestcase
// modify to invalid testcase: TestModifyToInvalidTestcase
// add testcase with existing ID: TestModifyInvalidAddTestcase
// remove non-existent testcase: TestRemoveNonExistentTestcase

func TestValidModifyAnnotations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mod := func(s *modificationTest) {
		mreq := &api.ModifySessionRequest{
			Id:           s.SessionID,
			Modification: api.Modification_UPDATE_ANNOTATIONS,
			Annotations:  map[string]string{"foo": "bar", "testing": "this"},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		if err != nil {
			t.Errorf("Modify returned an error despite valid request: %v", err)
		}
		if resp == nil {
			t.Errorf("Modify did not return a response despite valid request")
		}
	}
	val := func(s *modificationTest) {
		na := s.Status.Status.Annotations
		if v, ok := na["foo"]; !ok || v != "bar" {
			t.Errorf("Modifiy did not update annotations: expected foo=bar")
		}
		if v, ok := na["testing"]; !ok || v != "this" {
			t.Errorf("Modifiy did not update annotations: expected foo=bar")
		}
		if len(na) != 2 {
			t.Errorf("Modify did not update all annotations. expected two annotations, actual: %v", na)
		}
	}
	testModification(t, ctrl, mod, val)
}

func TestValidModifyTestcase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var casemod *api.Testcase
	mod := func(s *modificationTest) {
		casemod = s.SessionReq.Plan.Case[0]
		casemod.Description = "foobar this is a change"
		mreq := &api.ModifySessionRequest{
			Id:           s.SessionID,
			Modification: api.Modification_MODIFY_TESTCASE,
			Case:         []*api.Testcase{casemod},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		if err != nil {
			t.Errorf("Modify returned an error despite valid request: %v", err)
		}
		if resp == nil {
			t.Errorf("Modify did not return a response despite valid request")
		}
	}
	val := func(s *modificationTest) {
		var tc00 *api.TestcaseStatus
		for _, tc := range s.Status.Status.Case {
			if tc.Case.Id == "tc00" {
				tc00 = tc
				break
			}
		}

		if tc00 == nil {
			t.Errorf("Status did not return tc00")
			return
		}

		if len(tc00.Claim) != 1 {
			t.Errorf("Modify did not preserve testcase claim: %d claims expected, %d claims actual", 1, len(tc00.Claim))
		} else if tc00.Claim[0].Name != s.User {
			t.Errorf("Modify broke testcase claim: %s expected, %s actual", s.User, tc00.Claim[0].Name)
		}
		if tc00.Case.Description != casemod.Description {
			t.Errorf("Modify did not update testcase in modify_testcase operation")
		}
	}
	testModification(t, ctrl, mod, val)
}

func TestValidAddTestcase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var casemod *api.Testcase
	mod := func(s *modificationTest) {
		casemod = s.SessionReq.Plan.Case[0]
		casemod.Id = "yet_another_one"
		casemod.Description = "foobar this is a change"
		mreq := &api.ModifySessionRequest{
			Id:           s.SessionID,
			Modification: api.Modification_ADD_TESTCASE,
			Case:         []*api.Testcase{casemod},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		if err != nil {
			t.Errorf("Modify returned an error despite valid request: %v", err)
		}
		if resp == nil {
			t.Errorf("Modify did not return a response despite valid request")
		}
	}
	val := func(s *modificationTest) {
		var tc00 *api.TestcaseStatus
		var tcYAT *api.TestcaseStatus
		for _, tc := range s.Status.Status.Case {
			if tc.Case.Id == "tc00" {
				tc00 = tc
			}
			if tc.Case.Id == casemod.Id {
				tcYAT = tc
			}
		}

		if tc00 == nil {
			t.Errorf("Status did not return tc00")
		} else {
			if len(tc00.Claim) != 1 {
				t.Errorf("Modify did not preserve testcase claim: %d claims expected, %d claims actual", 1, len(tc00.Claim))
			} else if tc00.Claim[0].Name != s.User {
				t.Errorf("Modify broke testcase claim: %s expected, %s actual", s.User, tc00.Claim[0].Name)
			}
		}

		if tcYAT == nil {
			t.Errorf("Modify did not add testcase")
		}
	}
	testModification(t, ctrl, mod, val)
}

func TestValidModifyRemoveTestcase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var casemod *api.Testcase
	mod := func(s *modificationTest) {
		casemod = s.SessionReq.Plan.Case[0]
		mreq := &api.ModifySessionRequest{
			Id:           s.SessionID,
			Modification: api.Modification_REMOVE_TESTCASE,
			Case:         []*api.Testcase{casemod},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		if err != nil {
			t.Errorf("Modify returned an error despite valid request: %v", err)
		}
		if resp == nil {
			t.Errorf("Modify did not return a response despite valid request")
		}
	}
	val := func(s *modificationTest) {
		var tc00 *api.TestcaseStatus
		for _, tc := range s.Status.Status.Case {
			if tc.Case.Id == casemod.Id {
				tc00 = tc
				break
			}
		}

		if tc00 != nil {
			t.Errorf("Modify did not remove the testcase")
			return
		}
	}
	testModification(t, ctrl, mod, val)
}

func TestModifyClosedSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mod := func(s *modificationTest) {
		if _, err := s.Server.Close(s.Reqval.GetContext(s.User, types.PermissionSessionClose), &api.CloseSessionRequest{Id: s.SessionID}); err != nil {
			t.Errorf("Cannot close session: %v", err)
			return
		}

		tc := s.SessionReq.Plan.Case[0]
		mreq := &api.ModifySessionRequest{
			Id:           s.SessionID,
			Modification: api.Modification_MODIFY_TESTCASE,
			Case:         []*api.Testcase{tc},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		testNegativeResponse(t, "Modify", codes.FailedPrecondition, resp == nil, err)
	}
	testModification(t, ctrl, mod, func(s *modificationTest) {})
}

func TestModifyUnmodifiableSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s, reqval := newTestServer(ctrl)

	user := "user"
	sreq := validStartSessionRequest()
	sreq.Modifiable = false
	sresp, err := s.Start(reqval.GetContext(user, types.PermissionSessionStart), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	tc := sreq.Plan.Case[0]
	mreq := &api.ModifySessionRequest{
		Id:           sresp.Id,
		Modification: api.Modification_MODIFY_TESTCASE,
		Case:         []*api.Testcase{tc},
	}
	resp, err := s.Modify(reqval.GetContext(user, types.PermissionSessionModify), mreq)
	testNegativeResponse(t, "Modify", codes.PermissionDenied, resp == nil, err)
}

func TestModifyNonExistentSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mod := func(s *modificationTest) {
		tc := s.SessionReq.Plan.Case[0]
		mreq := &api.ModifySessionRequest{
			Id:           "this-session-does-not-exist",
			Modification: api.Modification_MODIFY_TESTCASE,
			Case:         []*api.Testcase{tc},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		testNegativeResponse(t, "Modify", codes.NotFound, resp == nil, err)
	}
	testModification(t, ctrl, mod, func(s *modificationTest) {})
}

func TestModifyNonExistentTestcase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mod := func(s *modificationTest) {
		tc := s.SessionReq.Plan.Case[0]
		tc.Id = "does-not-exist"
		mreq := &api.ModifySessionRequest{
			Id:           s.SessionID,
			Modification: api.Modification_MODIFY_TESTCASE,
			Case:         []*api.Testcase{tc},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		testNegativeResponse(t, "Modify", codes.NotFound, resp == nil, err)
	}
	testModification(t, ctrl, mod, func(s *modificationTest) {})
}

func TestModifyToInvalidTestcase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mod := func(s *modificationTest) {
		tc := api.Testcase{
			Id: s.SessionReq.Plan.Case[0].Id,
		}
		mreq := &api.ModifySessionRequest{
			Id:           s.SessionID,
			Modification: api.Modification_MODIFY_TESTCASE,
			Case:         []*api.Testcase{&tc},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		testNegativeResponse(t, "Modify", codes.InvalidArgument, resp == nil, err)
	}
	testModification(t, ctrl, mod, func(s *modificationTest) {})
}

func TestModifyInvalidAddTestcase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mod := func(s *modificationTest) {
		tc := s.SessionReq.Plan.Case[0]
		mreq := &api.ModifySessionRequest{
			Id:           s.SessionID,
			Modification: api.Modification_ADD_TESTCASE,
			Case:         []*api.Testcase{tc},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		testNegativeResponse(t, "Modify", codes.AlreadyExists, resp == nil, err)
	}
	testModification(t, ctrl, mod, func(s *modificationTest) {})
}

func TestRemoveNonExistentTestcase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mod := func(s *modificationTest) {
		tc := api.Testcase{
			Id: "does-not-exist",
		}
		mreq := &api.ModifySessionRequest{
			Id:           s.SessionID,
			Modification: api.Modification_REMOVE_TESTCASE,
			Case:         []*api.Testcase{&tc},
		}
		resp, err := s.Server.Modify(s.Reqval.GetContext(s.User, types.PermissionSessionModify), mreq)
		testNegativeResponse(t, "Modify", codes.NotFound, resp == nil, err)
	}
	testModification(t, ctrl, mod, func(s *modificationTest) {})
}

type modificationTest struct {
	Server     api.SessionServiceServer
	Reqval     *MockValidator
	SessionReq *api.StartSessionRequest
	SessionID  string
	Status     *api.SessionStatusResponse
	User       string
}

func testModification(t *testing.T, ctrl *gomock.Controller, mod func(s *modificationTest), validate func(s *modificationTest)) {
	s, reqval := newTestServer(ctrl)

	user := "user"
	sreq := validStartSessionRequest()
	sresp, err := s.Start(reqval.GetContext(user, types.PermissionSessionStart), sreq)
	if err != nil {
		t.Errorf("Cannot start session: %v", err)
		return
	}

	if _, err := s.Register(reqval.GetContext(user, types.PermissionSessionContribute), validRegistrationRequest(sresp.Id)); err != nil {
		t.Errorf("Cannot join session: %v", err)
		return
	}
	if _, err := s.Claim(reqval.GetContext(user, types.PermissionSessionContribute), &api.ClaimRequest{Session: sresp.Id, TestcaseID: "tc00", Claim: true}); err != nil {
		t.Errorf("Cannot claim testcase: %v", err)
		return
	}
	status, err := s.Status(reqval.GetContext(user, types.PermissionSessionView), &api.SessionStatusRequest{Id: sresp.Id})
	if err != nil || status == nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}

	state := modificationTest{
		Server:     s,
		Reqval:     reqval,
		SessionReq: sreq,
		SessionID:  sresp.Id,
		User:       user,
		Status:     status,
	}
	mod(&state)

	status, err = s.Status(reqval.GetContext(user, types.PermissionSessionView), &api.SessionStatusRequest{Id: sresp.Id})
	if err != nil || status == nil {
		t.Errorf("Cannot get status: %v", err)
		return
	}
	state.Status = status
	validate(&state)
}
