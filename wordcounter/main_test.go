package main

import (
	"bytes"
	"testing"
)

// test the wordcounter
func TestCountWords(t *testing.T) {
	// create buffer with strings
	buf := bytes.NewBufferString("word1 word2 word3 word4 word5\n")
	// expected outcome
	expected := 5
	// call the count function
	result := count(buf)

	if result != expected {
		t.Errorf("Expected %d, got %d instead", expected, result)
	}
}
