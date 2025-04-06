# Go Template AST Viewer

A tool to visualize Go's template AST (Abstract Syntax Tree) in JSON format. This tool helps you understand how Go's template engine parses and interprets your templates.

## Features

- Parses Go templates and outputs their AST structure in JSON format
- Supports all Go template syntax including:
  - Actions (`{{.Name}}`)
  - Conditionals (`{{if}}`, `{{else}}`)
  - Ranges (`{{range}}`)
  - With blocks (`{{with}}`)
- Provides a clean and readable JSON representation of the template structure

## Installation

```bash
go install github.com/takaaa220/go_template_ast_viewer@latest
```

## Usage

```bash
go_template_ast_viewer <template_file>
```

The tool will output the AST structure in JSON format to stdout.

### Example

Given a template file `example.tmpl`:

```
{{with .User}}
  Name: {{.Name}}
  {{if .IsAdmin}}
    Role: Admin
  {{else}}
    Role: User
  {{end}}
{{end}}
```

Running the tool:

```bash
go_template_ast_viewer example.tmpl
```

Will output:

```json
{
  "type": "*parse.ListNode",
  "children": [
    {
      "type": "*parse.WithNode",
      "children": [
        {
          "type": "body",
          "children": [
            {
              "type": "*parse.TextNode",
              "value": "\n  Name: "
            },
            {
              "type": "*parse.ActionNode",
              "details": {
                "pipe": ".Name"
              }
            },
            {
              "type": "*parse.IfNode",
              "children": [
                {
                  "type": "then",
                  "children": [
                    {
                      "type": "*parse.TextNode",
                      "value": "\n    Role: Admin\n  "
                    }
                  ]
                },
                {
                  "type": "else",
                  "children": [
                    {
                      "type": "*parse.TextNode",
                      "value": "\n    Role: User\n  "
                    }
                  ]
                }
              ],
              "details": {
                "condition": ".IsAdmin"
              }
            }
          ]
        }
      ],
      "details": {
        "with": ".User"
      }
    }
  ]
}
```

## Node Types

The AST is represented using the following node types:

- `*parse.ListNode`: A list of nodes
- `*parse.TextNode`: Plain text content
- `*parse.ActionNode`: Template actions (e.g., `{{.Name}}`)
- `*parse.IfNode`: Conditional blocks
- `*parse.RangeNode`: Range/loop blocks
- `*parse.WithNode`: With blocks for dot rewriting
