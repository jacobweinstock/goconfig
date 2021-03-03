package goconfig

import (
	"errors"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
)

type testEnvImplementation struct {
	fail bool
}

func (e *testEnvImplementation) Parse(log logr.Logger, prefix string, config interface{}) error {
	if e.fail {
		return errors.New("failed")
	}
	return nil
}

type testConfig struct {
	One   string
	Two   int
	Three []string
}

func TestParseEnv(t *testing.T) {
	conf := new(testConfig)
	err := ParseEnv(nil, "", conf, &testEnvImplementation{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnvconfigImplementation(t *testing.T) {
	expected := &testConfig{
		One:   "hello one",
		Two:   12,
		Three: []string{"one", "two", "three"},
	}
	os.Setenv("TEST_ONE", "hello one")
	os.Setenv("TEST_TWO", "12")
	os.Setenv("TEST_THREE", "one,two,three")
	defer func() {
		os.Unsetenv("ONE")
		os.Unsetenv("TWO")
		os.Unsetenv("THREE")
	}()
	conf := new(testConfig)
	t.Log(conf)
	ec := new(envConfig)
	err := ParseEnv(nil, "TEST", conf, ec)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expected, conf); diff != "" {
		t.Fatal(diff)
	}
}
