package goconfig

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/jacobweinstock/registrar"
	"go.uber.org/zap"

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
	log := defaultLogger()
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

// Parse a config file, env vars and cli flags (override is in that order).
//
// The fields of the confStruct passed in must be exported (uppercase).
//
// CLI flags by default split camelCase field names with dashes.
// e.x. `KeyOne` would be a cli flag of `-key-one`.
// To modify this, add a struct tag.
// KeyOne string `flag:"keyone"` will give you a cli flag of `-keyone`.
func (c *Parser) Parse(confStruct interface{}) error {

	// parse env then cli flags to get any config file path
	err := ParseEnv(c.Logger, c.Prefix, confStruct, c.EnvInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing env vars")
	}
	//var f *gflags
	err = ParseFlags(c.Logger, confStruct, c.FlagInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing cli flags")
	}

	// Read the config file
	filename := getConfigValue(confStruct)
	if filename == "" {
		filename = c.File
	}

	err = ParseFileFromInterfaces(c.Logger, filename, confStruct, c.FileInterface.GetDriverInterfaces())
	if err != nil {
		c.Logger.V(0).Info("problem parsing file", "file", filename, "error", err.Error())
	}

	// Overwrite config with environment variables
	err = ParseEnv(c.Logger, c.Prefix, confStruct, c.EnvInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing env vars")
	}

	// Overwrite config with command line args
	err = ParseFlags(c.Logger, confStruct, c.FlagInterface)
	if err != nil {
		return errors.WithMessage(err, "error parsing cli flags")
	}

	return validateRequired(confStruct)
}

// validateRequired will look at the struct tags for a `valid:"required"` tag
// and make sure those fields have a value.
// If a default value is specified for a required field that will satisfy
func validateRequired(config interface{}) error {
	_, err := govalidator.ValidateStruct(config)
	return err
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

// defaultLogger is zap logr implementation
func defaultLogger() logr.Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	zapLogger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}

	return zapr.NewLogger(zapLogger)
}
