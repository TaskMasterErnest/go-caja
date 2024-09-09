package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	// check if the file is pointing to a dir OR the size is less than the miz size for filtering
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	// compare extensions
	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func listFile(path string, out io.Writer) error {
	// print out the path of the current file to the specified io.Writer and return potential errors
	_, err := fmt.Fprintln(out, path)
	return err
}

func deleteFile(path string) error {
	// remove the file path
	return os.Remove(path)
}
