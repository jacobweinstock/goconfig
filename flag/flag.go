package flag

import "github.com/go-logr/logr"

// Parser interface for environment variables
type Parser interface {
	// Parse environment variables and update the config
	// interface with the values
	Parse(log logr.Logger, config interface{}) error
}

// ParseFlags cli flags
func ParseFlags(log logr.Logger, config interface{}, p Parser) error {
	return p.Parse(log, config)
}
