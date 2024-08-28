package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func count(r io.Reader) int {
	// read text from input
	scanner := bufio.NewScanner(r)
	// scan the words in the text given
	scanner.Split(bufio.ScanWords)
	// initialize the word counter variable
	var wc int
	// count the words, increment counter and print
	for scanner.Scan() {
		wc++
	}
	return wc
}

func main() {
	fmt.Println(count(os.Stdin))
}
