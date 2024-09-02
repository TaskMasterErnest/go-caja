package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

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

func run(filename string, out io.Writer, skipPreview bool) error {
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

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	return preview(outName)
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

// adding an auto-preview feature

func preview(filename string) error {
	cmdName := ""
	cmdParams := []string{}

	// define the executable based on the OS
	switch runtime.GOOS {
	case "linux":
		cmdName = "xdg-open"
	case "windows":
		cmdName = "cmd.exe"
		cmdParams = []string{"/C", "start"}
	case "darwin":
		cmdName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	// append filename to parameter slice
	cmdParams = append(cmdParams, filename)

	// locate the executable in the path
	cmdPath, err := exec.LookPath(cmdName)
	if err != nil {
		return err
	}

	// open the file using the default program
	return exec.Command(cmdPath, cmdParams...).Run()
}

func main() {
	// set flags
	filename := flag.String("file", "", "Markdown file to preview.")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	// check if the filename is passed
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
