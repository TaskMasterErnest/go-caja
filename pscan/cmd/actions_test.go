package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/go-caja/pscan/scan"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// create temp file
	tempFile, err := os.CreateTemp("", "pscan")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	// initalize list if specified
	if initList {
		h := &scan.HostsList{}

		for _, host := range hosts {
			h.Add(host)
		}

		if err := h.Save(tempFile.Name()); err != nil {
			t.Fatal(err)
		}
	}

	// return the tempFile name and the cleanup function
	return tempFile.Name(), func() {
		os.Remove(tempFile.Name())
	}
}

func TestHostActions(t *testing.T) {
	// define hosts
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	// test cases for actions
	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			expectedOut:    "host1\nhost2\nhost3\n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeleteAction",
			args:           []string{"host3"},
			expectedOut:    "Deleted host: host3\n",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tempFile, cleanup := setup(t, hosts, testCase.initList)
			defer cleanup()

			// define var to capture action output
			var out bytes.Buffer

			// execute action and capture output
			if err := testCase.actionFunction(&out, tempFile, testCase.args); err != nil {
				t.Fatalf("expected no error, got %q instead", err)
			}

			// test actions output
			if out.String() != testCase.expectedOut {
				t.Errorf("expected output %q, got %q\n", testCase.expectedOut, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	// define the hosts for integration use
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	// setup the test using the setup function
	tempFile, cleanup := setup(t, hosts, false)
	defer cleanup()

	delHost := "host2"

	hostsEnd := []string{
		"host1",
		"host3",
	}

	// define var to capture output
	var out bytes.Buffer

	// define expected outputs for all actions
	expectedOut := ""
	for _, v := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", v)
	}
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("Deleted host: %s\n", delHost)
	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()

	// add hosts to the list
	if err := addAction(&out, tempFile, hosts); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// list hosts
	if err := listAction(&out, tempFile, nil); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// delete hosts
	if err := deleteAction(&out, tempFile, []string{delHost}); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// list hosts after delete
	if err := listAction(&out, tempFile, nil); err != nil {
		t.Fatalf("expected no error, got %q\n", err)
	}

	// test integration output
	if out.String() != expectedOut {
		t.Errorf("expected output %q, got %q\n", expectedOut, out.String())
	}
}
