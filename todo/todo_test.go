package todo_test

import (
	"os"
	"testing"

	"github.com/go-caja/todo"
)

// test case for the Add method
func TestAdd(t *testing.T) {
	// initialize the List struct
	l := todo.List{}

	taskname := "New Task"
	l.Add(taskname)

	// check if task exists in the list
	if l[0].Task != taskname {
		t.Errorf("expected %q, got %q instead", taskname, l[0].Task)
	}
}

func TestComplete(t *testing.T) {
	l := todo.List{}
	taskname := "New Task2"
	l.Add(taskname)
	if l[0].Task != taskname {
		t.Errorf("expected %q, got %q instead.\n", taskname, l[0].Task)
	}

	// check if task is done
	if l[0].Done {
		t.Errorf("this task should not have been completed.\n")
	}

	// mark the task as done
	l.Complete(1)

	if !l[0].Done {
		t.Errorf("this task should be completed.\n")
	}
}

func TestDelete(t *testing.T) {
	l := todo.List{}

	// create a list of tasks
	tasks := []string{
		"New task 1",
		"New task 2",
		"New task 3",
	}

	// add these tasks
	for _, task := range tasks {
		l.Add(task)
	}

	// check if the tasks match
	if l[0].Task != tasks[0] {
		t.Errorf("expected %q, got %q instead", tasks[0], l[0].Task)
	}

	// delete the second task
	l.Delete(2)

	// check if list length has changed
	if len(l) != 2 {
		t.Errorf("expected list length of %d, got %d instead", 2, len(l))
	}

	// check if task ordering has occured
	if l[1].Task != tasks[2] {
		t.Errorf("expected %q, got %q instead", tasks[2], l[1].Task)
	}
}

// test the Save and Get methods
func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	taskname := "New task"
	l1.Add(taskname)

	if l1[0].Task != taskname {
		t.Errorf("expected %q, got %q instead", taskname, l1[0].Task)
	}

	// create temp file
	file, err := os.CreateTemp("", "temp")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}

	defer os.Remove(file.Name())

	// save the file
	if err := l1.Save(file.Name()); err != nil {
		t.Fatalf("error saving list to file: %s", err)
	}

	if err := l2.Get(file.Name()); err != nil {
		t.Fatalf("error getting list from file: %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("task %q should match %q task", l1[0].Task, l2[0].Task)
	}
}
