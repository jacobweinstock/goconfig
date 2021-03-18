# GOCONFIG

[![Go Report Card](https://goreportcard.com/badge/github.com/jacobweinstock/goconfig)](https://goreportcard.com/report/github.com/jacobweinstock/goconfig)
[![tests](https://github.com/jacobweinstock/goconfig/actions/workflows/ci.yaml/badge.svg)](https://github.com/jacobweinstock/goconfig/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jacobweinstock/goconfig.svg)](https://pkg.go.dev/github.com/jacobweinstock/goconfig)

`goconfig` is a simple but configurable library for adding configuration file -> environment variables -> cli flag parsing to your app.

## Usage

`go run examples/main.go -config examples/example.yaml`

```go
package main

import (
    "fmt"

    "github.com/jacobweinstock/goconfig"
)

// config fields must be exported (uppercase)
// CLI flags by default split camelCase field names with dashes
// e.x. `KeyOne` would be a cli flag of `-key-one`
// To modify this, add a struct tag
// KeyOne string `flag:"keyone"` will give you a cli flag of `-keyone`
type config struct {
    KeyOne  string `flag:"keyone"`
    Workers int
    Config  string
}

func main() {
    // set any default values
    cfg := config{Workers: 4}
    // create a config parser with any options; see option.go for all options
    config := goconfig.NewParser(
        goconfig.WithPrefix("TEST"),
        goconfig.WithFile("example.yaml"),
    )
    // run the config/env/flag parser
    err := config.Parse(&cfg)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("%+v\n", cfg)
}
```

## Customization

`goconfig` is fully customizable. If you prefer to parse config files, environment variables or cli flags in a different way you can pass in your own implementations.

- `WithFileInterface(myFileImplementation)`
- `WithEnvInterface(myEnvImplementation)`
- `WithFlagInterface(myFlagImplementation)`

The interface definitions:

```go
// EnvParser interface for environment variables
type EnvParser interface {
    // Parse environment variables and update the config
    // interface with the values
    Parse(log logr.Logger, prefix string, config interface{}) error
}

// FlagParser interface for environment variables
type FlagParser interface {
    // Parse environment variables and update the config
    // interface with the values
    Parse(log logr.Logger, config interface{}) error
}
```

```go
// WithFileInterface() takes a `*registrar.Registry` so that
// multiple file parsers can be registered and used
r := registrar.NewRegistry()
// example, register a yaml file parser you write
r.Register(
    fileInterfaceNameYaml,
    fileInterfaceProtocolYaml,
    fileInterfaceFeaturesYaml,
    nil,
    new(yamlParser),
)


```
