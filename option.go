package goconfig

import (
	"github.com/go-logr/logr"
	"github.com/jacobweinstock/registrar"
)

// Option for config
type Option func(*Parser)

// WithLogger sets the logger
func WithLogger(logger logr.Logger) Option {
	return func(args *Parser) { args.Logger = logger }
}

// WithPrefix sets the prefix for env vars PREFIX_
func WithPrefix(prefix string) Option {
	return func(args *Parser) { args.Prefix = prefix }
}

// WithFile sets the config file name
func WithFile(name string) Option {
	return func(args *Parser) { args.File = name }
}

// WithUsage sets the usage func
func WithUsage(usage func()) Option {
	return func(args *Parser) { args.Usage = usage }
}

// WithFlagInterface sets the flag parser to use
func WithFlagInterface(fi FlagParser) Option {
	return func(args *Parser) { args.FlagInterface = fi }
}

// WithEnvInterface sets the env parser to use
func WithEnvInterface(ei EnvParser) Option {
	return func(args *Parser) { args.EnvInterface = ei }
}

// WithFileInterface sets the file parser to use
func WithFileInterface(fi *registrar.Registry) Option {
	return func(args *Parser) { args.FileInterface = fi }
}
