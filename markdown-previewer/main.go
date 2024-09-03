package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	bluemonday "github.com/microcosm-cc/bluemonday"
	blackfriday "github.com/russross/blackfriday/v2"
)

// define html content struct used to define a template
type content struct {
	Title string
	Body  template.HTML
}

const (
	defaultTemplate = `<!DOCTYPE html>
<html lang="en">

<head>
  <meta http-equiv="content-type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{ .Title }}</title>
</head>

<body>
{{ .Body }}
</body>

</html>`
)

func run(filename string, tFname string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// take input and transform it into html
	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err
	}

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

	defer os.Remove(outName)

	return preview(outName)
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	// parse markdown through blackfriday
	output := blackfriday.Run(input)
	// run the output through bluemonday
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// parse the contents of the default template const into a new template
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// if the user provides alternate template file, replace default template
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	// instantiate content type
	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(body),
	}

	// create a buffer of bytes to write to a file
	var buffer bytes.Buffer

	// execute the template with the content type
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
	err = exec.Command(cmdPath, cmdParams...).Run()

	// add some time for browser the file before deletion
	// a temporary measure
	time.Sleep(5 * time.Second)
	return err
}

func main() {
	// set flags
	filename := flag.String("file", "", "Markdown file to preview.")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()

	// check if the filename is passed
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
