package goconfig

import (
	"github.com/go-logr/logr"
	"github.com/kelseyhightower/envconfig"
)

// EnvParser interface for environment variables
type EnvParser interface {
	// Parse environment variables and update the config
	// interface with the values
	Parse(log logr.Logger, prefix string, config interface{}) error
}

// ParseEnv environment
func ParseEnv(log logr.Logger, prefix string, config interface{}, p EnvParser) error {
	return p.Parse(log, prefix, config)
}

// envconfig struct
type envConfig struct{}

// Parse implementation of Parser interface
func (e *envConfig) Parse(log logr.Logger, prefix string, config interface{}) error {
	return envconfig.Process(prefix, config)
}
