package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

// format the List output
func (l *List) String() string {
	formattedOutput := ""

	for idx, task := range *l {
		prefix := "\u2615  "
		if task.Done {
			prefix = "\u2705  " // alt codes
		}
		// adjust item number in formatted output
		formattedOutput += fmt.Sprintf("%s%d: %s\n", prefix, idx+1, task.Task)
		// verboseFormattedOutput += fmt.Sprintf("%s%d: %s\t\t[created: %s | completed: %s]\n", prefix, idx+1, task.Task, task.CreatedAt.Format(time.ANSIC), task.CompletedAt.Format(time.ANSIC))
	}

	return formattedOutput
}

// create a new todo item and add it to list
func (l *List) Add(task string) {
	//
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	// append task to list
	*l = append(*l, t)
}

// define the Complete method to mark a todo item as completed
func (l *List) Complete(i int) error {
	list := *l
	// check of the i number is valid/exists in the list of tasks
	if i <= 0 || i > len(list) {
		return fmt.Errorf("item %d does not exist", i)
	}
	// adjust the index for zero(0) based indexing
	list[i-1].Done = true
	list[i-1].CompletedAt = time.Now()

	return nil
}

// define the Delete method to delete an item task from list
func (l *List) Delete(i int) error {
	list := *l
	if i <= 0 || i > len(list) {
		return fmt.Errorf("item %d does not exist", i)
	}
	// adjust index to clean out deleted item
	*l = append(list[:i-1], list[i:]...)

	return nil
}

// define the Save method; converts data to JSON and writes it to a file
func (l *List) Save(filename string) error {
	jsonData, err := json.MarshalIndent(l, "", "	")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// define the Get method; opens the file, decodes json into List struct
func (l *List) Get(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, l)

	return err
}
