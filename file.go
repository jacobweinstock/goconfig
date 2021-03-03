package goconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-logr/logr"
	"github.com/hashicorp/go-multierror"
	"github.com/jacobweinstock/registrar"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	fileInterfaceNameYaml     = "yaml"
	fileInterfaceProtocolYaml = "yaml"
	fileInterfaceNameJSON     = "json"
	fileInterfaceProtocolJSON = "json"
)

var (
	fileInterfaceFeaturesYaml = registrar.Features{"yaml", "yml"}
	fileInterfaceFeaturesJSON = registrar.Features{"json"}
)

// FileParser interface for configuration files
type FileParser interface {
	// Parse config file  and update the config
	// interface with the values
	Parse(log logr.Logger, name string, config interface{}) error
}

// ParseFile parses a file, trying all interface implementations passed in
func ParseFile(log logr.Logger, name string, config interface{}, o []FileParser) (err error) {
	var fileParsed bool
	for _, elem := range o {
		if elem != nil {
			openErr := elem.Parse(log, name, config)
			if openErr != nil {
				err = multierror.Append(err, openErr)
				continue
			}
			fileParsed = true
			break
		}
	}
	if fileParsed {
		return nil
	}
	return multierror.Append(err, errors.New("failed to parse file"))
}

// ParseFileFromInterfaces pass through to ParseFile function
func ParseFileFromInterfaces(log logr.Logger, name string, config interface{}, generic []interface{}) (err error) {
	parsers := make([]FileParser, 0)
	for _, elem := range generic {
		switch p := elem.(type) {
		case FileParser:
			parsers = append(parsers, p)
		default:
			e := fmt.Sprintf("not a FileParser implementation: %T", p)
			err = multierror.Append(err, errors.New(e))
		}
	}
	if len(parsers) == 0 {
		return multierror.Append(err, errors.New("no FileParser implementations found"))
	}
	return ParseFile(log, name, config, parsers)
}

type yamlParser struct{}

func (y *yamlParser) Parse(log logr.Logger, filename string, config interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if notFound := os.IsNotExist(err); notFound {
		return fmt.Errorf("file not found: %v", filename)
	}

	return yaml.Unmarshal(file, config)
}

type jsonParser struct{}

func (j *jsonParser) Parse(log logr.Logger, filename string, config interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if notFound := os.IsNotExist(err); notFound {
		return fmt.Errorf("file not found: %v", filename)
	}

	return json.Unmarshal(file, config)
}

type tomlParser struct{} // nolint

// Parse toml file, WIP
func (t *tomlParser) Parse(log logr.Logger, filename string, config interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if notFound := os.IsNotExist(err); notFound {
		return fmt.Errorf("file not found: %v", filename)
	}

	return toml.Unmarshal(file, config)
}
