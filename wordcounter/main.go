package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func count(r io.Reader, countLines bool, countBytes bool) int {
	// read text from input
	scanner := bufio.NewScanner(r)
	// check if the countLines is enabled
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}
	if countBytes {
		scanner.Split(bufio.ScanBytes)
	}
	// initialize the word counter variable
	var wc int
	// count the words, increment counter and print
	for scanner.Scan() {
		wc++
	}
	return wc
}

func main() {
	// initialize flags
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count bytes")
	filename := flag.String("f", "", "File to read from")
	flag.Parse()

	var reader io.Reader

	if *filename != "" {
		file, err := os.Open(*filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	fmt.Println(count(reader, *lines, *bytes))
}
