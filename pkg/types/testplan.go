package types

import (
	"strconv"

	"github.com/Knetic/govaluate"
)

func (tp *TestPlan) GetTestSet(id string) *TestSet {
	for _, k := range tp.Testset {
		if k.ID == id {
			return &k
		}
	}
	return nil
}

// SelectTestCases removes all testcases from this test plan that do not match the expression
func (tp *TestPlan) SelectTestCases(expr string) error {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return err
	}

	cases := tp.Case[:0]
	for _, tc := range tp.Case {
		matches, err := expression.Eval(tc)
		if err != nil {
			return err
		}

		if matches == true {
			cases = append(cases, tc)
		}
	}
	tp.Case = cases
	return nil
}

// Get implements govaluate.Parameters interface for the expression evaluation
func (tc Testcase) Get(name string) (interface{}, error) {
	if name == "_id" {
		return tc.ID, nil
	} else if name == "_group" {
		return tc.Group, nil
	} else if name == "_name" {
		return tc.Name, nil
	} else if val, ok := tc.Annotations[name]; ok {
		if val == "true" {
			return true, nil
		} else if val == "false" {
			return false, nil
		} else if i, err := strconv.Atoi(val); err == nil {
			return i, nil
		} else {
			return val, nil
		}
		return val, nil
	} else {
		return nil, nil
	}
}
