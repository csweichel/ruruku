package storage

import (
	"github.com/32leaves/ruruku/protocol"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadSuite(fn string) (*protocol.TestSuite, error) {
	fc, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var result protocol.TestSuite
	if err := yaml.Unmarshal(fc, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
