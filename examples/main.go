package main

import (
	"fmt"

	"github.com/jacobweinstock/goconfig"
)

type config struct {
	ConsumerToken string `flag:"consumertoken"`
	APIKey        string `flag:"apikey"`
	Facility      string
	HardwareIDs   []string `flag:"hardwareids"`
	Workers       int
	Config        string
}

func main() {
	cfg := config{Workers: 4}
	config := goconfig.NewParser(goconfig.WithPrefix("TEST"), goconfig.WithFile("example.yaml"))
	err := config.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", cfg)
}
