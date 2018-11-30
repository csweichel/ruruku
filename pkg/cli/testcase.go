package cli

import (
	"fmt"
	"github.com/32leaves/ruruku/pkg/types"
	"github.com/manifoldco/promptui"
	"os"
	"strconv"
	"strings"
)

type InitTestcase struct {
	Init
	types.Testcase
	MinTesterCountSet bool
}

func (cfg *InitTestcase) Complete(suite *types.TestPlan) error {
	var err error

	if cfg.Group, err = cfg.checkOrAskString(cfg.Group, "Group", "", true, validateGroup(suite)); err != nil {
		return err
	}

	if cfg.ID, err = cfg.checkOrAskString(cfg.ID, "ID", "", true, validateID(suite, cfg.Group)); err != nil {
		return err
	}

	if cfg.Name, err = cfg.checkOrAskString(cfg.Name, "Name", "", true, nil); err != nil {
		return err
	}

	if cfg.Description, err = cfg.checkOrAskString(cfg.Description, "Description", "", false, nil); err != nil {
		return err
	}

	if !cfg.MinTesterCountSet && !cfg.NonInteractive {
		p := promptui.Prompt{
			Label:    "minimum tester count",
			Validate: validateMinTesterCount,
			Default:  "0",
		}
		val, err := p.Run()
		if err != nil {
			return err
		}
		if mtc, err := strconv.ParseUint(val, 10, 32); err != nil {
			return err
		} else {
			cfg.MinTesterCount = uint32(mtc)
		}
	}

	return nil
}

func validateGroup(suite *types.TestPlan) func(string) error {
	return func(val string) error {
		if err := validateNotEmpty(val); err != nil {
			return err
		}

		if err := validateIdentifier(val); err != nil {
			return err
		}

		for _, tc := range suite.Case {
			if strings.ToLower(tc.Group) == strings.ToLower(val) && tc.Group != val {
				return fmt.Errorf("A group named %s already exists, but case does not match (use %s instead)", tc.Group, tc.Group)
			}
		}

		return nil
	}
}

func validateID(suite *types.TestPlan, grp string) func(string) error {
	return func(val string) error {
		if err := validateNotEmpty(val); err != nil {
			return err
		}

		if err := suiteHasTestcase(suite, val); err != nil {
			return err
		}

		if err := validateIdentifier(val); err != nil {
			return err
		}

		return nil
	}
}

func validateMinTesterCount(val string) error {
	_, err := strconv.ParseInt(val, 10, 32)
	return err
}

func (cfg *InitTestcase) Run() error {
	var err error
	cfg.Filename, err = cfg.checkOrAskString(cfg.Filename, "filename", "testplan.yaml", true, validateFileExists)
	if err != nil {
		return err
	}

	suite, err := LoadTestplan(cfg.Filename)
	if err != nil {
		return err
	}

	if err := cfg.Complete(suite); err != nil {
		return err
	}

	if err := suiteHasTestcase(suite, cfg.ID); err != nil {
		return err
	}

	suite.Case = append(suite.Case, cfg.Testcase)

	return cfg.saveSuite(*suite, true)
}

func validateFileExists(val string) error {
	_, err := os.Stat(val)
	return err
}

func suiteHasTestcase(suite *types.TestPlan, id string) error {
	for _, tc := range suite.Case {
		if id == tc.ID {
			return fmt.Errorf("Testcase %s already exists", tc.ID)
		}
	}
	return nil
}

func validateIdentifier(val string) error {
	if strings.ContainsAny(val, " ") || strings.ContainsAny(val, "/") {
		return fmt.Errorf("Identifier must not contain whitespace or slashes")
	}
	return nil
}
