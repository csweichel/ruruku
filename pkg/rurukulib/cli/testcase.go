package cli

import (
	"fmt"
	"github.com/32leaves/ruruku/pkg/rurukulib"
	"github.com/32leaves/ruruku/protocol"
	"github.com/manifoldco/promptui"
	"strconv"
	"strings"
)

type InitTestCase struct {
	Init
	protocol.TestCase
	MinTesterCount    int32
	MinTesterCountSet bool
}

func (cfg *InitTestCase) Complete(suite *protocol.TestSuite) error {
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
		if mtc, err := strconv.ParseInt(val, 10, 32); err != nil {
			return err
		} else {
			cfg.MinTesterCount = int32(mtc)
		}
	}
	cfg.TestCase.MinTesterCount = float64(cfg.MinTesterCount)

	return nil
}

func validateGroup(suite *protocol.TestSuite) func(string) error {
	return func(val string) error {
		if err := validateNotEmpty(val); err != nil {
			return err
		}

		if err := validateIdentifier(val); err != nil {
			return err
		}

		for _, tc := range suite.Cases {
			if strings.ToLower(tc.Group) == strings.ToLower(val) && tc.Group != val {
				return fmt.Errorf("A group named %s already exists, but case does not match (use %s instead)", tc.Group, tc.Group)
			}
		}

		return nil
	}
}

func validateID(suite *protocol.TestSuite, grp string) func(string) error {
	return func(val string) error {
		if err := validateNotEmpty(val); err != nil {
			return err
		}

		if err := suiteHasTestcase(suite, grp, val); err != nil {
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

func (cfg *InitTestCase) Run() error {
	fn, err := cfg.checkOrAskString(cfg.Filename, "filename", "", true, nil)
	if err != nil {
		return err
	}

	suite, err := rurukulib.LoadSuite(fn)
	if err != nil {
		return err
	}

	if err := cfg.Complete(suite); err != nil {
		return err
	}

	if err := suiteHasTestcase(suite, cfg.Group, cfg.ID); err != nil {
		return err
	}

	suite.Cases = append(suite.Cases, cfg.TestCase)

	return cfg.saveSuite(*suite, true)
}

func suiteHasTestcase(suite *protocol.TestSuite, group string, id string) error {
	for _, tc := range suite.Cases {
		if group == tc.Group && id == tc.ID {
			key := fmt.Sprintf("%s/%s", group, id)
			return fmt.Errorf("Testcase %s already exists", key)
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
