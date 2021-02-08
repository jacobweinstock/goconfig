package fileparse

import (
	"github.com/go-logr/logr"
)

// Parser interface for configuration files
type Parser interface {
	// Parse config file  and update the config
	// interface with the values
	Parse(log logr.Logger, name string, config interface{}) error
}

// ParseConfigFile parses a config file
func ParseConfigFile(log logr.Logger, name string, config interface{}, p Parser) error {
	return p.Parse(log, name, config)
}
