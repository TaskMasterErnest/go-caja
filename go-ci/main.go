package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func run(project string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project directory is required :%w", ErrValidation)
	}
	// build a string for arguments
	args := []string{"build", ".", "errors"}
	// set the default command to run
	cmd := exec.Command("go", args...)
	// set the target directory on which the command should be executed
	cmd.Dir = project

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("'go build' failed: %s", err)
	}

	_, err := fmt.Fprintln(out, "GO Build: SUCCESS")

	return err

}

func main() {
	project := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*project, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
