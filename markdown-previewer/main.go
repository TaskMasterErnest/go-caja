package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	bluemonday "github.com/microcosm-cc/bluemonday"
	blackfriday "github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
<html lang="en">

<head>
  <meta http-equiv="content-type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Markdown Preview Tool</title>
</head>

<body>
`

	footer = `
</body>

</html>`
)

func run(filename string, out io.Writer) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// take input and transform it into html
	htmlData := parseContent(input)

	// create a temporary file to store the new content
	temp, err := os.CreateTemp("/tmp/", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()

	fmt.Fprintln(out, outName)

	return saveHTML(outName, htmlData)
}

func parseContent(input []byte) []byte {
	// parse markdown through blackfriday
	output := blackfriday.Run(input)
	// run the output through bluemonday
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// create a buffer of bytes to write to a file
	var buffer bytes.Buffer
	// write html to bytes buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outFName string, data []byte) error {
	return os.WriteFile(outFName, data, 0644)
}

func main() {
	// set flags
	filename := flag.String("file", "", "Markdown file to preview.")
	flag.Parse()

	// check if the filename is passed
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
