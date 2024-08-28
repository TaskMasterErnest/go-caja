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
	flag.Parse()

	// print output with the flags enabled
	fmt.Println(count(os.Stdin, *lines, *bytes))
}
