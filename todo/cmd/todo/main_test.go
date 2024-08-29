package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// command to build the binary
	build := exec.Command("go", "build", "-o", binName)

	// check if the command will run
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	// run the command
	fmt.Println("Running tests...")
	// run the tests on the main function
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	// get the current working dir
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// create the command path, path to binary
	cmdPath := filepath.Join(dir, binName)

	// create the first test to ensure the tool can add a new task using the t.Run
	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-task", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := task + "\n"

		if expected != string(out) {
			t.Errorf("expected %q, got %q instead\n", expected, string(out))
		}
	})
}
