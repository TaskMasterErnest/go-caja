package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-caja/todo"
)

const todoFileName = ".todo.json"

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
	task := flag.String("task", "", "task to be included in ToDo list")
	list := flag.Bool("list", false, "list all tasks.")
	complete := flag.Int("complete", 0, "item to the marked as complete")
	flag.Parse()

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		// use the API defined method
		fmt.Println(l)
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
	case *task != "":
		l.Add(*task)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}
