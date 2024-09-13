package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	// extension to filter out
	ext string
	// min file size
	size int64
	// list files
	list bool
	// delete files
	del bool
	// destination log writer
	writeLog io.Writer
}

func run(root string, out io.Writer, cfg config) error {
	// initalizing delLogger
	delLogger := log.New(cfg.writeLog, "DELETED FILE: ", log.LstdFlags)

	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if filterOut(path, cfg.ext, cfg.size, info) {
				return nil
			}

			// if list was explicitly called; do nothing else
			if cfg.list {
				return listFile(path, out)
			}

			// if an explicit delete is called
			if cfg.del {
				return deleteFile(path, delLogger)
			}

			// the default option here is list;
			return listFile(path, out)
		})
}

func main() {
	// command-line flags
	root := flag.String("root", ".", "Root directory to start from.")
	logFile := flag.String("log", "", "Log deletes to file.")
	// action options
	list := flag.Bool("list", false, "List files only.")
	delete := flag.Bool("del", false, "Delete files")
	// filter options
	ext := flag.String("ext", "", "File extension to filter out.")
	size := flag.Int64("size", 0, "Minimum file size.")
	flag.Parse()

	// passing in the logger
	var (
		file = os.Stdout
		err  error
	)

	if *logFile != "" {
		file, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
	}

	// create instance of config struct that can be passed to run function
	c := config{
		ext:      *ext,
		size:     *size,
		list:     *list,
		del:      *delete,
		writeLog: file,
	}

	// pass config struct to run function
	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
