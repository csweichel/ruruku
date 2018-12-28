package v1

import "github.com/32leaves/ruruku/pkg/types"

func (s *Participant) Convert() types.Participant {
	return types.Participant{
		Name: s.Name,
	}
}

func ConvertParticipant(s *types.Participant) *Participant {
	return &Participant{
		Name: s.Name,
	}
}

func (s *Testcase) Convert() types.Testcase {
	return types.Testcase{
		ID:             s.Id,
		Name:           s.Name,
		Group:          s.Group,
		Description:    s.Description,
		Steps:          s.Steps,
		MustPass:       s.MustPass,
		MinTesterCount: s.MinTesterCount,
	}
}

func ConvertTestcase(s *types.Testcase) *Testcase {
	return &Testcase{
		Id:             s.ID,
		Name:           s.Name,
		Group:          s.Group,
		Description:    s.Description,
		Steps:          s.Steps,
		MustPass:       s.MustPass,
		MinTesterCount: s.MinTesterCount,
	}
}

func (s *TestcaseRunResult) Convert() types.TestcaseRunResult {
	return types.TestcaseRunResult{
		Participant: s.Participant.Convert(),
		State:       s.State.Convert(),
		Comment:     s.Comment,
	}
}

func ConvertTestcaseRunResult(s *types.TestcaseRunResult) *TestcaseRunResult {
	return &TestcaseRunResult{
		Participant: ConvertParticipant(&s.Participant),
		State:       ConvertTestRunState(s.State),
		Comment:     s.Comment,
	}
}

func (s *TestcaseStatus) Convert() types.TestcaseStatus {
	claims := make([]types.Participant, len(s.Claim))
	for i, c := range s.Claim {
		claims[i] = c.Convert()
	}
	results := make([]types.TestcaseRunResult, len(s.Result))
	for i, r := range s.Result {
		results[i] = r.Convert()
	}

	return types.TestcaseStatus{
		Case:   s.Case.Convert(),
		Claim:  claims,
		Result: results,
		State:  s.State.Convert(),
	}
}

func (s *TestRunStatus) Convert() types.TestRunStatus {
	cases := make([]types.TestcaseStatus, len(s.Case))
	for i, c := range s.Case {
		cases[i] = c.Convert()
	}
	return types.TestRunStatus{
		ID:          s.Id,
		Name:        s.Name,
		PlanID:      s.PlanID,
		Case:        cases,
		State:       s.State.Convert(),
		Annotations: s.Annotations,
	}
}

func (s *TestPlan) Convert() types.TestPlan {
	cases := make([]types.Testcase, len(s.Case))
	for i, c := range s.Case {
		cases[i] = c.Convert()
	}
	return types.TestPlan{
		ID:          s.Id,
		Name:        s.Name,
		Description: s.Description,
		Case:        cases,
	}
}

func ConvertTestPlan(s *types.TestPlan) *TestPlan {
	cases := make([]*Testcase, len(s.Case))
	for i, c := range s.Case {
		cases[i] = ConvertTestcase(&c)
	}
	return &TestPlan{
		Id:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Case:        cases,
	}
}

func (s TestRunState) Convert() types.TestRunState {
	if s == TestRunState_PASSED {
		return types.Passed
	}
	if s == TestRunState_UNDECIDED {
		return types.Undecided
	}
	return types.Failed
}

func ConvertTestRunState(s types.TestRunState) TestRunState {
	if s == types.Passed {
		return TestRunState_PASSED
	}
	if s == types.Undecided {
		return TestRunState_UNDECIDED
	}
	return TestRunState_FAILED
}

func (s Permission) Convert() types.Permission {
	if s == Permission_USER_ADD {
		return types.PermissionUserAdd
	}
	if s == Permission_USER_DELETE {
		return types.PermissionUserDelete
	}
	if s == Permission_USER_GRANT {
		return types.PermissionUserGrant
	}
	if s == Permission_USER_CHPWD {
		return types.PermissionUserChpwd
	}
	if s == Permission_USER_LIST {
		return types.PermissionUserList
	}
	if s == Permission_SESSION_START {
		return types.PermissionSessionStart
	}
	if s == Permission_SESSION_CLOSE {
		return types.PermissionSessionClose
	}
	if s == Permission_SESSION_VIEW {
		return types.PermissionSessionView
	}
	if s == Permission_SESSION_CONTRIBUTE {
		return types.PermissionSessionContribute
	}
	if s == Permission_SESSION_MODIFY {
		return types.PermissionSessionModify
	}
	return types.PermissionNone
}

func ConvertPermission(s types.Permission) Permission {
	if s == types.PermissionUserAdd {
		return Permission_USER_ADD
	}
	if s == types.PermissionUserDelete {
		return Permission_USER_DELETE
	}
	if s == types.PermissionUserGrant {
		return Permission_USER_GRANT
	}
	if s == types.PermissionUserChpwd {
		return Permission_USER_CHPWD
	}
	if s == types.PermissionUserList {
		return Permission_USER_LIST
	}
	if s == types.PermissionSessionStart {
		return Permission_SESSION_START
	}
	if s == types.PermissionSessionClose {
		return Permission_SESSION_CLOSE
	}
	if s == types.PermissionSessionView {
		return Permission_SESSION_VIEW
	}
	if s == types.PermissionSessionContribute {
		return Permission_SESSION_CONTRIBUTE
	}
	if s == types.PermissionSessionModify {
		return Permission_SESSION_MODIFY
	}
	return Permission_NONE
}
