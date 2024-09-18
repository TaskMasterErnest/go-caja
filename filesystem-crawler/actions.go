package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path string, ext []string, minSize int64, info os.FileInfo) bool {
	// check if the file is pointing to a dir OR the size is less than the miz size for filtering
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	// compare extensions
	fileExt := filepath.Ext(path)

	// loop through multiple extension flags is any
	if len(ext) > 0 {
		for _, e := range ext {
			if fileExt == e {
				return false
			}
		}
	}

	return true
}

func listFile(path string, out io.Writer) error {
	// print out the path of the current file to the specified io.Writer and return potential errors
	_, err := fmt.Fprintln(out, path)
	return err
}

func deleteFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	delLogger.Println(path)
	return nil
}

func archiveFile(destDir, root, path string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		fmt.Errorf("%s is not a directory", destDir)
	}

	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath := filepath.Join(destDir, relDir, dest)

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	zipWriter := gzip.NewWriter(out)
	zipWriter.Name = filepath.Base(path)

	if _, err := io.Copy(zipWriter, in); err != nil {
		return err
	}

	if err := zipWriter.Close(); err != nil {
		return err
	}

	return out.Close()
}
