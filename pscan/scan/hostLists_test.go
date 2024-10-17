package scan_test

import (
	"errors"
	"os"
	"testing"

	"github.com/go-caja/pscan/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name        string
		host        string
		expectedLen int
		expectedErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// initialize list
			h := &scan.HostsList{}

			if err := h.Add("host1"); err != nil {
				t.Fatal(err)
			}

			err := h.Add(testCase.host)

			if testCase.expectedErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil instead\n")
				}

				if !errors.Is(err, testCase.expectedErr) {
					t.Errorf("expected error %q, got %q instead\n", testCase.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %q instead\n", err)
			}

			if len(h.Hosts) != testCase.expectedLen {
				t.Errorf("expected list length: %d, got %d instead\n", testCase.expectedLen, len(h.Hosts))
			}

			if h.Hosts[1] != testCase.host {
				t.Errorf("expected host name %q as index 1, got %q instead\n", testCase.host, h.Hosts[1])
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name        string
		host        string
		expectedLen int
		expectedErr error
	}{
		{"RemoveExisting", "host1", 1, nil},
		{"RemoveNotFound", "host3", 1, scan.ErrNotExists},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// initialize the list
			h := &scan.HostsList{}

			for _, host := range []string{"host1", "host2"} {
				if err := h.Add(host); err != nil {
					t.Fatal(err)
				}
			}

			err := h.Remove(testCase.host)

			if testCase.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil instead\n")
				}

				if !errors.Is(err, testCase.expectedErr) {
					t.Errorf("expected error: %q, got %q instead", testCase.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %q instead\n", err)
			}

			if len(h.Hosts) != testCase.expectedLen {
				t.Errorf("expected list length: %d, got %d instead\n", testCase.expectedLen, len(h.Hosts))
			}

			if h.Hosts[0] == testCase.host {
				t.Errorf("host name: %q should not be in the list\n", testCase.host)
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	h1 := &scan.HostsList{}
	h2 := &scan.HostsList{}

	hostName := "host1"
	h1.Add(hostName)

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("error create temp file: %s", err)
	}
	defer tempFile.Close()

	if err := h1.Save(tempFile.Name()); err != nil {
		t.Fatalf("error getting list from file: %s", err)
	}

	// loading the saved data from file
	if err := h2.Load(tempFile.Name()); err != nil {
		t.Fatalf("error getting list from file: %s", err)
	}

	if h1.Hosts[0] != h2.Hosts[0] {
		t.Errorf("host: %q should match %q host", h1.Hosts[0], h2.Hosts[0])
	}
}

func TestLoadNoFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}

	if err := os.Remove(tempFile.Name()); err != nil {
		t.Fatalf("error deleting temp file: %s", err)
	}

	h := &scan.HostsList{}

	if err := h.Load(tempFile.Name()); err != nil {
		t.Errorf("expected no error, got %q instead\n", err)
	}
}
