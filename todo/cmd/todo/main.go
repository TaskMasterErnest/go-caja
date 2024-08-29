package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-caja/todo"
)

const todoFileName = ".todo.json"

func main() {
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
		for _, task := range *l {
			if !task.Done {
				fmt.Println(task.Task)
			}
		}
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
