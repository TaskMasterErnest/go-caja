package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "NoFilter", root: "testdata", cfg: config{ext: "", size: 0, list: true},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata", cfg: config{ext: ".log", size: 0, list: true},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata", cfg: config{ext: ".log", size: 10, list: true},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata", cfg: config{ext: ".log", size: 20, list: true},
			expected: ""},
		{name: "FilterExtensionNoMatch", root: "testdata", cfg: config{ext: ".gz", size: 0, list: true},
			expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if tc.expected != res {
				t.Errorf("expected %q, got %q instead\n", tc.expected, res)
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()

	temp, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	// iterate over the files map and create dummy files of a specified number
	for idx, num := range files {
		for j := 1; j <= num; j++ {
			filename := fmt.Sprintf("file%d%s", j, idx)
			path := filepath.Join(temp, filename)
			if err := os.WriteFile(path, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return temp, func() { os.RemoveAll(temp) }
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{name: "DeleteExtensionNoMatch", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 0, nNoDelete: 10,
			expected: ""},
		{name: "DeleteExtensionMatch", cfg: config{ext: ".log", del: true}, extNoDelete: "", nDelete: 10, nNoDelete: 0,
			expected: ""},
		{name: "DeleteExtensionMixed", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 5, nNoDelete: 5,
			expected: ""},
	}

	// execute the delete test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)

			//instantiate the delLogger
			testCase.cfg.writeLog = &logBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				testCase.cfg.ext:     testCase.nDelete,
				testCase.extNoDelete: testCase.nNoDelete,
			})

			defer cleanup()

			if err := run(tempDir, &buffer, testCase.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if testCase.expected != res {
				t.Errorf("expected %q, got %q instead\n", testCase.expected, res)
			}

			// read the files left after the delete operation
			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != testCase.nNoDelete {
				t.Errorf("expected %d files left, got %d instead\n", testCase.nNoDelete, len(filesLeft))
			}

			expectedLogLines := testCase.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expectedLogLines {
				t.Errorf("expected %d log lines, got %d instead", expectedLogLines, len(lines))
			}

		})
	}
}
