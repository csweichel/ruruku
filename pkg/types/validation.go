package types

import "gopkg.in/go-playground/validator.v9"

// caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateParticipant(obj *Participant) error {
	return validate.Struct(obj)
}

func ValidateTestcase(obj *Testcase) error {
	return validate.Struct(obj)
}

func ValidateTestcaseRunResult(obj *TestcaseRunResult) error {
	return validate.Struct(obj)
}

func ValidateTestcaseStatus(obj *TestcaseStatus) error {
	return validate.Struct(obj)
}

func ValidateTestRunStatus(obj *TestRunStatus) error {
	return validate.Struct(obj)
}

func ValidateTestPlan(obj *TestPlan) error {
	return validate.Struct(obj)
}
