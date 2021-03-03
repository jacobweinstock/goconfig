package goconfig

import (
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-multierror"
)

type testFileImplementation struct {
	shouldErr bool
}

func (e *testFileImplementation) Parse(log logr.Logger, name string, config interface{}) error {
	if e.shouldErr {
		return errors.New("failed")
	}
	return nil
}

func TestParseFile(t *testing.T) {
	tests := map[string]struct {
		shouldErr bool
		err       error
	}{
		"success": {},
		"fail":    {shouldErr: true, err: &multierror.Error{Errors: []error{errors.New("failed"), errors.New("failed to parse file")}}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			conf := new(testConfig)
			err := ParseFile(nil, "", conf, []FileParser{&testFileImplementation{shouldErr: tc.shouldErr}})
			if err != nil {
				if tc.shouldErr {
					if diff := cmp.Diff(tc.err.Error(), err.Error()); diff != "" {
						t.Fatal(diff)
					}
				} else {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestParseFileFromInterfaces(t *testing.T) {
	tests := map[string]struct {
		shouldErr         bool
		err               error
		badImplementation bool
	}{
		"success":            {},
		"bad implementation": {shouldErr: true, err: &multierror.Error{Errors: []error{errors.New("not a FileParser implementation: *struct {}"), errors.New("no FileParser implementations found")}}, badImplementation: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			conf := new(testConfig)
			var generic []interface{}
			if tc.badImplementation {
				badImplementation := struct{}{}
				generic = []interface{}{&badImplementation}
			} else {
				testImplementation := &testFileImplementation{}
				generic = []interface{}{testImplementation}
			}
			err := ParseFileFromInterfaces(nil, "", conf, generic)
			if err != nil {
				if tc.shouldErr {
					if diff := cmp.Diff(tc.err.Error(), err.Error()); diff != "" {
						t.Fatal(diff)
					}
				} else {
					t.Fatal(err)
				}
			}
		})
	}
}
