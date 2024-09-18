package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Extensions []string

func (e *Extensions) Set(value string) error {
	*e = append(*e, value)
	return nil
}

func (e *Extensions) String() string {
	return fmt.Sprint(*e)
}

type config struct {
	// extension to filter out
	ext []string
	// min file size
	size int64
	// list files
	list bool
	// delete files
	del bool
	// destination log writer
	writeLog io.Writer
	// archive directory
	archive string
	// modification date
	modDate time.Time
}

func run(root string, out io.Writer, cfg config) error {
	// initalizing delLogger
	delLogger := log.New(cfg.writeLog, "DELETED FILE: ", log.LstdFlags)

	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if filterOut(path, cfg.ext, cfg.size, cfg.modDate, info) {
				return nil
			}

			// if list was explicitly called; do nothing else
			if cfg.list {
				return listFile(path, out)
			}

			// adding the archive function
			if cfg.archive != "" {
				if err := archiveFile(cfg.archive, root, path); err != nil {
					return err
				}
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
	archive := flag.String("archive", "", "Archive directory")
	// filter options
	var extensions Extensions
	// ext := flag.String("ext", "", "File extension to filter out.")
	flag.Var(&extensions, "ext", "File extension to filter out")
	moddedTime := flag.String("moddedAfter", "", "Modification time of file.")
	size := flag.Int64("size", 0, "Minimum file size.")
	flag.Parse()

	// passing in the logger
	var file = os.Stdout
	var err error

	if *logFile != "" {
		file, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
	}

	// check if a modification time filter was presented
	var modTime time.Time
	if *moddedTime != "" {
		if modTime, err = time.Parse("2006-01-02", *moddedTime); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// create instance of config struct that can be passed to run function
	c := config{
		ext:      extensions,
		size:     *size,
		list:     *list,
		del:      *delete,
		writeLog: file,
		archive:  *archive,
		modDate:  modTime,
	}

	// pass config struct to run function
	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
