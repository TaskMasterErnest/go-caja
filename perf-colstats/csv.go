package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// adding a generic function
type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	sum := 0.0

	for _, v := range data {
		sum += v
	}

	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func csv2float(r io.Reader, column int) ([]float64, error) {
	reader := csv.NewReader(r)
	// reduce mem allocation by using the same slice for each read operation
	reader.ReuseRecord = true
	// adjusting for a zero-based index
	column--

	var data []float64
	
	for i := 0; ; i++ {
		// read each record from each row in the file
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read data from file: %w", err)
		}
		if i == 0 {
			continue
		}
		// checking number of columns in csv file
		if len(row) <= column {
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}
		// convert data into a float
		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}
		// append converted float to data slice
		data = append(data, v)
	}

	return data, nil
}
