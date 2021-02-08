package env

import (
	"github.com/go-logr/logr"
)

// Parser interface for environment variables
type Parser interface {
	// Parse environment variables and update the config
	// interface with the values
	Parse(log logr.Logger, prefix string, config interface{}) error
}

// ParseEnv environment
func ParseEnv(log logr.Logger, prefix string, config interface{}, p Parser) error {
	return p.Parse(log, prefix, config)
}
