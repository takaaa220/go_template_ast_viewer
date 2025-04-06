package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/takaaa220/go_template_ast_viewer/parser"
)

func main() {
	var (
		templateString string
		templateFile   string
	)

	flag.StringVar(&templateString, "s", "", "template string")
	flag.StringVar(&templateFile, "f", "", "template file path")
	flag.Parse()

	var (
		content string
		err     error
	)

	switch {
	case templateString != "":
		content = templateString
	case templateFile != "":
		bytes, err := os.ReadFile(templateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read template file: %v\n", err)
			os.Exit(1)
		}
		content = string(bytes)
	default:
		// Read from stdin
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read from stdin: %v\n", err)
			os.Exit(1)
		}
		content = string(bytes)
	}

	// Use "stdin" as the template name when reading from stdin or string
	templateName := "stdin"
	if templateFile != "" {
		templateName = templateFile
	}

	result, err := parser.ParseTemplate(templateName, content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse template: %v\n", err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode JSON: %v\n", err)
		os.Exit(1)
	}
}
