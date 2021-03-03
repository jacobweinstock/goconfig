package goconfig

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jacobweinstock/registrar"
	plogr "github.com/packethost/pkg/log/logr"
)

func TestNewParser(t *testing.T) {
	log, _, _ := plogr.NewPacketLogr()
	usage := func() {}
	flagI := &gflags{}
	envI := &envConfig{}
	fileI := registrar.NewRegistry(registrar.WithLogger(log))
	expected := &Parser{
		Prefix:        "TEST",
		File:          "config.yaml",
		FlagInterface: flagI,
		EnvInterface:  envI,
	}
	cfgParser := NewParser(
		WithPrefix("TEST"),
		WithLogger(log),
		WithFile("config.yaml"),
		WithUsage(usage),
		WithFlagInterface(flagI),
		WithEnvInterface(envI),
		WithFileInterface(fileI),
	)
	if diff := cmp.Diff(expected, cfgParser, cmpopts.IgnoreFields(Parser{}, "Logger", "FileInterface", "Usage")); diff != "" {
		t.Fatal(diff)
	}
}

func TestParse(t *testing.T) {
	original := os.Args
	defer func() {
		os.Args = original
	}()
	expected := &testConfig{
		One:   "hello one",
		Two:   12,
		Three: []string{"one", "two", "three"},
	}

	tests := map[string]struct {
		failEnv     bool
		failFlags   bool
		err         error
		expectedCfg *testConfig
	}{
		"success":       {expectedCfg: expected},
		"fail on env":   {expectedCfg: new(testConfig), failEnv: true, err: errors.New("error parsing env vars: failed")},
		"fail on flags": {expectedCfg: new(testConfig), failFlags: true, err: errors.New("error parsing cli flags: failed")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cfg := new(testConfig)
			os.Args = []string{""}

			var cfgParser *Parser
			if tc.failEnv {
				cfgParser = NewParser(WithEnvInterface(&testEnvImplementation{fail: true}))
			} else if tc.failFlags {
				cfgParser = NewParser(WithFlagInterface(&testFlagImplementation{fail: true}))
			} else {
				os.Setenv("ONE", "hello one")
				os.Setenv("TWO", "12")
				os.Setenv("THREE", "one,two,three")
				cfgParser = NewParser()
			}

			err := cfgParser.Parse(cfg)
			if err != nil {
				if tc.err != nil {
					if diff := cmp.Diff(tc.err.Error(), err.Error()); diff != "" {
						t.Fatal(diff)
					}
				} else {
					t.Fatal(err)
				}
			}
			if diff := cmp.Diff(tc.expectedCfg, cfg); diff != "" {
				t.Fatal(diff)
			}
			os.Unsetenv("ONE")
			os.Unsetenv("TWO")
			os.Unsetenv("THREE")
		})
	}
}
