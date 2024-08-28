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
	result := count(buf, false, false)

	if result != expected {
		t.Errorf("Expected %d, got %d instead\n", expected, result)
	}
}

func TestCountLines(t *testing.T) {
	buf := bytes.NewBufferString("word1 word2\n word3\n word4 word5")
	expected := 3
	result := count(buf, true, false)
	if result != expected {
		t.Errorf("Expected %d, got %d instead\n", expected, result)
	}
}

func TestCountBytes(t *testing.T) {
	buf := bytes.NewBufferString("word1 word2 word3 word4 word5")
	expected := 29
	result := count(buf, false, true)
	if result != expected {
		t.Errorf("Expected %d, got %d instead\n", expected, result)
	}
}
