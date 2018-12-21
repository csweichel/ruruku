package types

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

// caches struct info
var validTestRunState = validation.In(Passed, Undecided, Failed)

func (obj *Participant) Validate() error {
	err := validation.ValidateStruct(obj,
		validation.Field(&obj.Name, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("Participant: %v", err)
	}
	return nil
}

func (obj *Testcase) Validate() error {
	err := validation.ValidateStruct(obj,
		validation.Field(&obj.ID, validation.Required, validation.NewStringRule(containsNoSlashes, "must not contain /")),
		validation.Field(&obj.Group, validation.Required),
		validation.Field(&obj.Name, validation.Required, validation.Length(5, 300)),
	)
	if err != nil {
		return fmt.Errorf("Testcase: %v", err)
	}
	return nil
}

func (obj *TestcaseRunResult) Validate() error {
	err := validation.ValidateStruct(obj,
		validation.Field(&obj.State, validation.Required, validTestRunState),
	)
	if err != nil {
		return fmt.Errorf("TestcaseRunResult: %v", err)
	}
	return nil
}

func (obj *TestcaseStatus) Validate() error {
	err := validation.ValidateStruct(obj,
		validation.Field(&obj.Case, validation.Required),
		validation.Field(&obj.State, validation.Required, validTestRunState),
	)
	if err != nil {
		return fmt.Errorf("TestcaseStatus: %v", err)
	}
	return nil
}

func (obj *TestRunStatus) Validate() error {
	err := validation.ValidateStruct(obj,
		validation.Field(&obj.ID, validation.Required, validation.NewStringRule(containsNoSlashes, "must not contain /")),
		validation.Field(&obj.Name, validation.Required, validation.Length(5, 300)),
		validation.Field(&obj.PlanID, validation.Required),
		validation.Field(&obj.State, validation.Required, validTestRunState),
	)
	if err != nil {
		return fmt.Errorf("TestRunStatus: %v", err)
	}
	return nil
}

func (obj *TestPlan) Validate() error {
	err := validation.ValidateStruct(obj,
		validation.Field(&obj.ID, validation.Required, validation.NewStringRule(containsNoSlashes, "must not contain /")),
		validation.Field(&obj.Name, validation.Required, validation.Length(5, 300)),
	)
	if err != nil {
		return fmt.Errorf("TestPlan: %v", err)
	}
	return nil
}

func containsNoSlashes(val string) bool {
	return !strings.Contains(val, "/")
}
