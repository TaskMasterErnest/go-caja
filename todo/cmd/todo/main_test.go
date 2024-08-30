package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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
		cmd := exec.Command(cmdPath, "-add", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		// connect to the stdin pipe
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		// write the string to the pipe
		io.WriteString(cmdStdIn, task2)
		// close the pipe
		cmdStdIn.Close()

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

		expected := fmt.Sprintf("\u2615  1: %s\n\u2615  2: %s\n\n", task, task2)

		if expected != string(out) {
			t.Errorf("expected %q, got %q instead\n", expected, string(out))
		}
	})

	// test the delete
	var taskNumber int = 1
	t.Run("DeleteTask", func(t *testing.T) {
		cmdDel := exec.Command(cmdPath, "-delete", strconv.Itoa(taskNumber))
		if err := cmdDel.Run(); err != nil {
			t.Error(err)
		}

		cmdList := exec.Command(cmdPath, "-list")
		out, err := cmdList.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("\u2615  1: %s\n\n", task2)

		if string(out) != expected {
			t.Errorf("expected %q, got %q instead\n", expected, string(out))
		}
	})
}
