package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	var testCases = []struct {
		name    string
		project string
		out     string
		expErr  error
	}{
		{name: "success", project: "./testdata/tool/", out: "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\n", expErr: nil},
		{name: "failed", project: "./testdata/toolError/", out: "", expErr: &stepErr{step: "go build"}},
		{name: "failFormat", project: "./testdata/toolFmtErr/", out: "", expErr: &stepErr{step: "go fmt"}},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			var out bytes.Buffer
			err := run(testcase.project, &out)

			if testcase.expErr != nil {
				if err == nil {
					t.Errorf("expected error: %q. Got 'nil' instead", testcase.expErr)
					return
				}

				if !errors.Is(err, testcase.expErr) {
					t.Errorf("expected error %q. Got %q instead", testcase.expErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %q", err)
			}
			if out.String() != testcase.out {
				t.Errorf("expected output: %q. Got %q instead.", testcase.out, out.String())
			}
		})
	}
}
