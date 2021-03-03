package goconfig

import (
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	"github.com/jacobweinstock/registrar"

	plogr "github.com/packethost/pkg/log/logr"
	"github.com/pkg/errors"
)

// Parser struct
type Parser struct {
	Logger        logr.Logger
	Prefix        string
	File          string
	Usage         func()
	FlagInterface FlagParser
	EnvInterface  EnvParser
	FileInterface *registrar.Registry
}

// NewParser parser struct
func NewParser(opts ...Option) *Parser {
	log, _, _ := plogr.NewPacketLogr()
	c := &Parser{
		Logger:        log,
		File:          "config.yaml",
		FlagInterface: new(gflags),
		EnvInterface:  new(envConfig),
		FileInterface: registrar.NewRegistry(registrar.WithLogger(log)),
	}
	for _, opt := range opts {
		opt(c)
	}
	// len of 0 means that no Registry with any registered drivers was passed in.
	if len(c.FileInterface.Drivers) == 0 {
		c.registerFileInterfaces()
	}
	return c
}

// Parse config file, env vars and cli flags
func (c *Parser) Parse(d interface{}) error {

	// parse env then cli flags to get any config file path
	err := ParseEnv(c.Logger, c.Prefix, d, c.EnvInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing env vars")
	}
	//var f *gflags
	err = ParseFlags(c.Logger, d, c.FlagInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing cli flags")
	}

	// Read the config file
	filename := getConfigValue(d)
	if filename == "" {
		filename = c.File
	}

	err = ParseFileFromInterfaces(c.Logger, filename, d, c.FileInterface.GetDriverInterfaces())
	if err != nil {
		c.Logger.V(0).Info("problem parsing file", "file", filename, "error", err.Error())
	}

	// Overwrite config with environment variables
	err = ParseEnv(c.Logger, c.Prefix, d, c.EnvInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing env vars")
	}

	// Overwrite config with command line args
	err = ParseFlags(c.Logger, d, c.FlagInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing cli flags")
	}

	return nil
}

func (c *Parser) registerFileInterfaces() {
	// register yaml file parser implementation
	c.FileInterface.Register(fileInterfaceNameYaml, fileInterfaceProtocolYaml, fileInterfaceFeaturesYaml, nil, new(yamlParser))
	// register json file parser implementation
	c.FileInterface.Register(fileInterfaceNameJSON, fileInterfaceProtocolJSON, fileInterfaceFeaturesJSON, nil, new(jsonParser))
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
