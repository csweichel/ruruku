package cli

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"strings"
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
		}
		nme, err := p.Run()
		if err != nil {
			return "", err
		}
		value = nme
	}

	return value, nil
}
