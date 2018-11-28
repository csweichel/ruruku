package cli

import (
	"fmt"
	"github.com/32leaves/ruruku/protocol"
	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type InitSuite struct {
	Init
	Name string
}

func (cfg *InitSuite) Run() error {
	r := protocol.TestSuite{
		Name:  cfg.Name,
		Cases: make([]protocol.TestCase, 0),
	}

	val, err := cfg.checkOrAskString(r.Name, "name", "", true, nil)
	if err != nil {
		return err
	}
	r.Name = val

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

func (cfg *Init) saveSuite(ts protocol.TestSuite, forceOverwrite bool) error {
	fn := cfg.Filename
	fn, err := cfg.checkOrAskString(fn, "filename", "testsuite.yaml", true, nil)
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
