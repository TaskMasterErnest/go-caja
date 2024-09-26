package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	// test cases for run
	testCases := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{name: "RunAvg1File", col: 3, op: "avg", exp: "227.6\n", files: []string{"./testdata/example.csv"}, expErr: nil},
		{name: "RunAvgMultiFiles", col: 3, op: "avg", exp: "233.84\n", files: []string{"./testdata/example.csv", "./testdata/example2.csv"}, expErr: nil},
		{name: "RunFailRead", col: 2, op: "avg", exp: "", files: []string{"./testdata/example.csv", "./testdata/fakefile.csv"}, expErr: os.ErrNotExist},
		{name: "RunFailColumn", col: 0, op: "avg", exp: "", files: []string{"./testdata/example.csv"}, expErr: ErrInvalidColumn},
		{name: "RunFailNoFiles", col: 2, op: "avg", exp: "", files: []string{}, expErr: ErrNoFiles},
		{name: "RunFailOperation", col: 2, op: "invalid", exp: "", files: []string{"./testdata/example.csv"}, expErr: ErrInvalidOperation},
	}

	// run test execution
	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			var result bytes.Buffer
			err := run(testcase.files, testcase.op, testcase.col, &result)

			if testcase.expErr != nil {
				if err == nil {
					t.Errorf("expected error. got nil instead")
				}

				if !errors.Is(err, testcase.expErr) {
					t.Errorf("expected error %q, got %q instead", testcase.expErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %q", err)
			}

			if result.String() != testcase.exp {
				t.Errorf("expected %q, got %q instead", testcase.exp, &result)
			}
		})
	}
}

// performing some benchmark tests
func BenchmarkRun(b *testing.B) {
	filenames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}
	// reset time to ignire time used to prep for start of benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := run(filenames, "avg", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}
