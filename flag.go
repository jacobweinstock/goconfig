package goconfig

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/octago/sflags/gen/gflag"
)

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
