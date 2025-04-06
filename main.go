package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/takaaa220/go_template_ast_viewer/parser"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: go_template_ast_viewer <template_file>")
		os.Exit(1)
	}

	templatePath := args[0]
	content, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read template file: %v\n", err)
		os.Exit(1)
	}

	result, err := parser.ParseTemplate(templatePath, string(content))
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
