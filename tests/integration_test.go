package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/takaaa220/go_template_ast_viewer/parser"
)

func TestIntegration(t *testing.T) {
	// Get all test case directories
	entries, err := os.ReadDir("cases")
	if err != nil {
		t.Fatalf("failed to read cases directory: %v", err)
	}

	var cases []string
	for _, entry := range entries {
		if entry.IsDir() {
			cases = append(cases, entry.Name())
		}
	}

	if len(cases) == 0 {
		t.Fatal("no test cases found")
	}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			// Read expected output
			expectedBytes, err := os.ReadFile(filepath.Join("cases", c, "expected.json"))
			if err != nil {
				t.Fatalf("failed to read expected output: %v", err)
			}

			var expected map[string]any
			err = json.Unmarshal(expectedBytes, &expected)
			if err != nil {
				t.Fatalf("failed to unmarshal expected output: %v", err)
			}

			// Run the parser
			templatePath := filepath.Join("cases", c, "template.tmpl")
			content, err := os.ReadFile(templatePath)
			if err != nil {
				t.Fatalf("failed to read template file: %v", err)
			}

			// Parse the template and get the actual output
			actual, err := parser.ParseTemplate(templatePath, string(content))
			if err != nil {
				t.Fatalf("failed to parse template: %v", err)
			}

			actualJSON, err := json.Marshal(actual)
			if err != nil {
				t.Fatalf("failed to marshal actual output: %v", err)
			}

			convertedActual := make(map[string]any)
			err = json.Unmarshal(actualJSON, &convertedActual)
			if err != nil {
				t.Fatalf("failed to unmarshal actual output: %v", err)
			}

			if !reflect.DeepEqual(expected, convertedActual) {
				actualJSON, err := json.MarshalIndent(actual, "", "  ")
				if err != nil {
					t.Fatalf("failed to marshal actual output: %v", err)
				}

				// write actual output to file
				err = os.WriteFile(filepath.Join("cases", c, "actual.json"), actualJSON, 0644)
				if err != nil {
					t.Fatalf("failed to write actual output: %v", err)
				}
			}
		})
	}
}
