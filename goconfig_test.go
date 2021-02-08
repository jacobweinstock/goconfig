package goconfig

import (
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	t.Skip("WIP")
	type config struct {
		ConsumerToken string
		APIKey        string
		Facility      string
		HardwareIDs   []string
		Workers       int
	}
	os.Setenv("TEST_WORKERS", "6")
	os.Setenv("TEST_HARDWAREIDS", "6,5,7")
	var cfg config
	err := Parse(&cfg, WithPrefix("TEST"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%+v\n", cfg))

	t.Fatal()
}
