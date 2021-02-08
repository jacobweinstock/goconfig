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
	cfg :=config{Workers: 4}
	err := goconfig.Parse(&cfg, goconfig.WithPrefix("TEST"), goconfig.WithFile("example.yaml"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", cfg)
}
