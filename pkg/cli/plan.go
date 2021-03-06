package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/32leaves/ruruku/pkg/types"
	"github.com/manifoldco/promptui"
	"github.com/technosophos/moniker"
	yaml "gopkg.in/yaml.v2"
)

type InitPlan struct {
	Init
	ID   string
	Name string
	Plan *types.TestPlan
}

func validatePlanID(val string) error {
	if err := validateNotEmpty(val); err != nil {
		return err
	}

	if err := validateIdentifier(val); err != nil {
		return err
	}

	return nil
}

func (cfg *InitPlan) Run() error {
	var r types.TestPlan
	if cfg.Plan != nil {
		r = *cfg.Plan
	} else {
		r = types.TestPlan{
			Case: make([]types.Testcase, 0),
		}
	}
	if cfg.ID != "" {
		r.ID = cfg.ID
	}
	if cfg.Name != "" {
		r.Name = cfg.Name
	}

	var err error
	r.ID, err = cfg.checkOrAskString(r.ID, "id", "", true, validatePlanID)
	if err != nil {
		return err
	}

	r.Name, err = cfg.checkOrAskString(r.Name, "name", moniker.New().Name(), true, nil)
	if err != nil {
		return err
	}

	if err := cfg.saveSuite(r, false); err != nil {
		return err
	}

	return nil
}

func validateNotEmpty(val string) error {
	if val == "" {
		return fmt.Errorf("must not be empty")
	}
	return nil
}

func (cfg *Init) saveSuite(ts types.TestPlan, forceOverwrite bool) error {
	fn := cfg.Filename
	fn, err := cfg.checkOrAskString(fn, "filename", "testplan.yaml", true, nil)
	if err != nil {
		return err
	}
	cfg.Filename = fn

	if _, err := os.Stat(fn); err == nil && !forceOverwrite {
		if !cfg.NonInteractive {
			p := promptui.Prompt{
				Label:     fmt.Sprintf("File %s exists. Overwrite?", fn),
				IsConfirm: true,
				Default:   "y",
			}
			_, err := p.Run()
			if err != nil {
				return fmt.Errorf("%s exists and must not be overwritten", fn)
			}
		} else {
			return fmt.Errorf("%s exists and must not be overwritten", fn)
		}
	}

	fc, err := yaml.Marshal(ts)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fn, fc, 0644)
}
