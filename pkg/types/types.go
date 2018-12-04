package types

type Participant struct {
	Name string `validate:"required"`
}

type Testcase struct {
	// ID of the testcase. Must be unique within the test suite.
	ID string `validate:"required" gorm:"primary_key"`
	// Name is the short description of the testcase
	Name string `validate:"required"`
	// Groups helps organize testcases
	Group string `validate:"required"`
	// Description is a long description
	Description string
	// Steps lists the individual steps a tester should perform
	Steps string
	// If true this testcase must pass for the suite to pass
	MustPass bool
	// MinTesterCount is the number of participants who need to run this test
	MinTesterCount uint32
}

type TestRunState string

const (
	Passed    TestRunState = "passed"
	Undecided TestRunState = "undecided"
	Failed    TestRunState = "Failed"
)

func WorseState(a TestRunState, b TestRunState) TestRunState {
	if a == Passed {
		return b
	}
	if a == Undecided && b == Failed {
		return b
	}
	return a
}

type TestcaseRunResult struct {
	// Participant who contributed this result
	Participant Participant
	// State denotes the success of a testcase
	State TestRunState `validate:"required"`
	// Comment is free text entered by the participant
	Comment string
}

type TestcaseStatus struct {
	// The testcase this run concerns
	Case Testcase `validate:"required"`
	// Claims mark testers who want to run a testcase
	Claim []Participant
	// Runs are completed testcase executions
	Result []TestcaseRunResult
	// State is the overall testcase success state
	State TestRunState `validate:"required"`
}

type TestRunStatus struct {
	// ID is the globally unique ID of this test run
	ID string `validate:"required"`
	// Name is a short description of this run
	Name string `validate:"required"`
	// Plan ID is the ID of the testplan being executed
	PlanID string `validate:"required"`
	// Status lists the status for each testcase of the plan
	Status []TestcaseStatus
	// State is the overall test run state
	State TestRunState `validate:"required"`
}

type TestPlan struct {
	// ID is the globally unique ID of this testplan
	ID string `validate:"required"`
	// Name is the short description of the testplan
	Name string `validate:"required"`
	// Description is a long description
	Description string `validate:"required"`
	// Case lists the testcases of this plan
	Case []Testcase
}
