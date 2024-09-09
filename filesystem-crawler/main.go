package main

import (
	"flag"
	"fmt"
	"io"
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
}

func run(root string, out io.Writer, cfg config) error {
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

			// the default option here is list;
			return listFile(path, out)
		})
}

func main() {
	// command-line flags
	root := flag.String("root", ".", "root directory to start from.")
	// action options
	list := flag.Bool("list", false, "list files only.")
	// filter options
	ext := flag.String("ext", "", "file extension to filter out.")
	size := flag.Int64("size", 0, "minimum file size.")
	flag.Parse()

	// create instance of config struct that can be passed to run function
	c := config{
		ext:  *ext,
		size: *size,
		list: *list,
	}

	// pass config struct to run function
	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
