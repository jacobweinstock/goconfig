package goconfig

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/octago/sflags/gen/gflag"
)

// FlagParser interface for environment variables
type FlagParser interface {
	// Parse environment variables and update the config
	// interface with the values
	Parse(log logr.Logger, config interface{}) error
}

// ParseFlags cli flags
func ParseFlags(log logr.Logger, config interface{}, p FlagParser) error {
	return p.Parse(log, config)
}

type gflags struct{}

func (g *gflags) Parse(log logr.Logger, cfg interface{}) error {
	fset, err := gflag.Parse(cfg)
	if err != nil {
		return err
	}

	if !fset.Parsed() {
		if err := fset.Parse(os.Args[1:]); err != nil {
			return err
		}
	}

	return nil
}
