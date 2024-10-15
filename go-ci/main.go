package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type executer interface {
	execute() (string, error)
}

func run(project string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project directory is required :%w", ErrValidation)
	}

	pipeline := make([]step, 3)

	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		project,
		[]string{"build", ".", "errors"},
	)

	pipeline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		project,
		[]string{"test", "-v"},
	)

	pipeline[2] = newStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS",
		project,
		[]string{"-l", "."},
	)

	for _, s := range pipeline {
		message, err := s.Execute()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	project := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*project, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
