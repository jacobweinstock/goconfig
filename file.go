package goconfig

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
)

// FileParser interface for configuration files
type FileParser interface {
	// Parse config file  and update the config
	// interface with the values
	Parse(log logr.Logger, name string, config interface{}) error
}

// ParseConfigFile parses a config file
func ParseConfigFile(log logr.Logger, name string, config interface{}, p FileParser) error {
	return p.Parse(log, name, config)
}

type yamlParser struct{}

func (y *yamlParser) Parse(log logr.Logger, filename string, config interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if notFound := os.IsNotExist(err); notFound {
		return fmt.Errorf("file not found: %v", filename)
	}

	return yaml.Unmarshal(file, config)
}
