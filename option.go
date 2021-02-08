package goconfig

import "github.com/go-logr/logr"

// Option for config
type Option func(*Config)

// WithLogger sets the logger
func WithLogger(logger logr.Logger) Option {
	return func(args *Config) { args.Logger = logger }
}

// WithPrefix sets the prefix for env vars PREFIX_
func WithPrefix(prefix string) Option {
	return func(args *Config) { args.Prefix = prefix }
}

// WithFile sets the config file name
func WithFile(name string) Option {
	return func(args *Config) { args.File = name }
}

// WithFileType sets the config file format type. i.e. yaml, json, etc
func WithFileType(fileType string) Option {
	return func(args *Config) { args.FileType = fileType }
}

// WithUsage sets the usage func
func WithUsage(usage func()) Option {
	return func(args *Config) { args.Usage = usage }
}

// WithFlagInterface sets the flag parser to use
func WithFlagInterface(fi FlagParser) Option {
	return func(args *Config) { args.FlagInterface = fi }
}

// WithEnvInterface sets the env parser to use
func WithEnvInterface(ei EnvParser) Option {
	return func(args *Config) { args.EnvInterface = ei }
}

// WithFileInterface sets the file parser to use
func WithFileInterface(fi FileParser) Option {
	return func(args *Config) { args.FileInterface = fi }
}
