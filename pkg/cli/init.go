package cli

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strings"
    "github.com/32leaves/ruruku/pkg/types"
    "gopkg.in/yaml.v2"
	"io/ioutil"
)

type Init struct {
	Filename       string
	NonInteractive bool
}

func (cfg *Init) checkOrAskString(value string, label string, dflt string, required bool, validation promptui.ValidateFunc) (string, error) {
	if required && value == "" {
		if cfg.NonInteractive {
			return "", fmt.Errorf("Runnig non-interactive and there's no %s set", strings.ToLower(label))
		}

		if validation == nil {
			validation = validateNotEmpty
		}

		p := promptui.Prompt{
			Label:    fmt.Sprintf("%s%s", strings.ToUpper(label[0:1]), label[1:]),
			Validate: validation,
			Default:  dflt,
            AllowEdit: true,
		}
		nme, err := p.Run()
		if err != nil {
			return "", err
		}
		value = nme
	}

	return value, nil
}

func LoadTestplan(fn string) (*types.TestPlan, error) {
    fc, err := ioutil.ReadFile(fn)
    if err != nil {
        return nil, err
    }

    r := types.TestPlan{}
    if err := yaml.Unmarshal(fc, &r); err != nil {
        return nil, err
    }

    return &r, nil
}
