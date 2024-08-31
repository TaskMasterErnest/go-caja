package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-caja/todo"
)

var todoFileName = ".todo.json"

func main() {
	// add the usage information
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed by TaskMasterErnest\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2024 (lol)\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage Information:")
		flag.PrintDefaults()
	}
	// adding flags for true CLI implementation
	add := flag.Bool("add", false, "task to be included in ToDo list.")
	list := flag.Bool("list", false, "list all tasks.")
	complete := flag.Int("complete", 0, "item to the marked as complete.")
	delete := flag.Int("delete", 0, "item to be deleted.")
	verbose := flag.Bool("verbose", false, "verbose output.")
	flag.Parse()

	// check if user defined the env var in a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Println(l)
	case *verbose:
		fmt.Println(l.StringVerbose(true))
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// save the file
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		// taking tasks from stdin
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

// define a function to use to get task from various input sources
func getTask(r io.Reader, args ...string) (string, error) {
	// check if args are passed
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	// if no args are passed, read whatever has been passed to the io.Reader interface
	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return s.Text(), nil
}
