package goconfig

import (
	"github.com/go-logr/logr"
	"github.com/kelseyhightower/envconfig"
)

// envconfig struct
type envConfig struct{}

// Parse implementation of Parser interface
func (e *envConfig) Parse(log logr.Logger, prefix string, config interface{}) error {
	return envconfig.Process(prefix, config)
}
