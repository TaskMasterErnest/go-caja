package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

func main() {
	// passing in flags
	operation := flag.String("op", "sum", "Operation to perform")
	column := flag.Int("col", 1, "CSV column on which to execute operation")
	flag.Parse()

	if err := run(flag.Args(), *operation, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filenames []string, op string, column int, out io.Writer) error {
	// make the operation function a type of statsFunc
	var opFunc statsFunc

	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	consolidate := make([]float64, 0)

	// add channels to receive results or errors of operation
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	// files channel to receive files
	filesCh := make(chan string)

	// define waitgroups
	wg := sync.WaitGroup{}

	// loop through the files ans send them to the channel to be processed when a worker is available
	go func() {
		defer close(filesCh)
		for _, filename := range filenames {
			filesCh <- filename
		}
	}()

	// loop through files and create goroutine to process each file concurrently
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filename := range filesCh {
				file, err := os.Open(filename)
				if err != nil {
					errCh <- fmt.Errorf("cannot open file: %w", err)
					return
				}
				// parse the csv into a slice of float64 nums
				data, err := csv2float(file, column)
				if err != nil {
					errCh <- err
				}

				if err := file.Close(); err != nil {
					errCh <- err
				}

				resCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err
		}
	}
}
