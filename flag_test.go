package goconfig

import (
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

type testFlagImplementation struct {
	fail bool
}

func (e *testFlagImplementation) Parse(log logr.Logger, config interface{}) error {
	if e.fail {
		return errors.New("failed")
	}
	return nil
}

func TestParseFlag(t *testing.T) {
	conf := new(testConfig)
	err := ParseFlags(nil, conf, &testFlagImplementation{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGflags(t *testing.T) {
	expectedSuccess := &testConfig{
		One:   "hello one",
		Two:   12,
		Three: []string{"one", "two", "three"},
	}

	tests := map[string]struct {
		expected *testConfig
		flags    []string
		err      error
	}{
		"success":              {expected: expectedSuccess, flags: []string{"-one=hello one", "-two=12", "-three=one,two,three"}},
		"nil config interface": {err: errors.New("object cannot be nil")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := new(testConfig)
			if tc.err != nil {
				got = nil
			}
			os.Args = append([]string{"app"}, tc.flags...)
			err := ParseFlags(nil, got, new(gflags))
			if err != nil {
				if tc.err != nil {
					if diff := cmp.Diff(tc.err.Error(), err.Error()); diff != "" {
						t.Fatal(diff)
					}
				} else {
					t.Fatal(err)
				}
			}
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
