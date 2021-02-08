package goconfig

import (
	"reflect"
	"strings"

	"github.com/go-logr/logr"

	plogr "github.com/packethost/pkg/log/logr"
	"github.com/pkg/errors"
)

// Config struct
type Config struct {
	Logger        logr.Logger
	Prefix        string
	File          string
	FileType      string
	Usage         func()
	FlagInterface FlagParser
	EnvInterface  EnvParser
	FileInterface FileParser
}

// Parse config file, env vars and cli flags
func Parse(c interface{}, opts ...Option) error {
	log, _, _ := plogr.NewPacketLogr()
	var defaultConfig = Config{
		Logger:        log,
		File:          "config.yaml",
		FileType:      "yaml",
		Prefix:        "APP",
		FlagInterface: new(gflags),
		EnvInterface:  new(envConfig),
		FileInterface: new(yamlParser),
	}

	for _, opt := range opts {
		opt(&defaultConfig)
	}

	// parse env then cli flags to get any config file path
	err := ParseEnv(log, defaultConfig.Prefix, c, defaultConfig.EnvInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing env vars")
	}
	//var f *gflags
	err = ParseFlags(log, c, defaultConfig.FlagInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing cli flags")
	}

	// Read the config file
	filename := getConfigValue(c)
	if filename == "" {
		filename = defaultConfig.File
	}
	err = ParseConfigFile(log, filename, c, defaultConfig.FileInterface)
	if err != nil {
		log.V(1).Info("file not found", "file", filename)
	}

	// Overwrite config with environment variables
	err = ParseEnv(log, defaultConfig.Prefix, c, defaultConfig.EnvInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing env vars")
	}

	// Overwrite config with command line args
	err = ParseFlags(log, c, defaultConfig.FlagInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing cli flags")
	}

	return nil
}

func getConfigValue(config interface{}) string {
	val := reflect.ValueOf(config).Elem()
	var name string
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if strings.ToLower(typeField.Name) == "config" {
			name = valueField.Interface().(string)
			break
		}
	}
	return name
}
