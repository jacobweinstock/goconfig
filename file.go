package goconfig

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
)

type yamlParser struct{}

func (y *yamlParser) Parse(log logr.Logger, filename string, config interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if notFound := os.IsNotExist(err); notFound {
		return fmt.Errorf("file not found: %v", filename)
	}

	return yaml.Unmarshal(file, config)
}
