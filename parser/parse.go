package parser

import (
	"fmt"
	"text/template/parse"
)

// Node represents a template AST node
type Node struct {
	Type     string  `json:"type"`
	Value    string  `json:"value,omitempty"`
	Children []*Node `json:"children,omitempty"`
	Details  any     `json:"details,omitempty"`
}

// ParseTemplate parses a template string and returns its AST as a JSON-compatible value
func ParseTemplate(name, content string) (*Node, error) {
	tree, err := parse.Parse(name, content, "{{", "}}", nil)
	if err != nil {
		return nil, err
	}

	if tree == nil {
		return nil, fmt.Errorf("parse.Parse returned nil tree")
	}

	root, ok := tree[name]
	if !ok {
		return nil, fmt.Errorf("template %q not found in parse tree", name)
	}

	return convertNode(root.Root), nil
}

func convertNode(node parse.Node) *Node {
	if node == nil {
		return nil
	}

	result := &Node{
		Type: getNodeType(node),
	}

	switch n := node.(type) {
	case *parse.ActionNode:
		result.Details = map[string]any{
			"pipe": n.Pipe.String(),
		}
	case *parse.TextNode:
		if len(n.Text) == 0 {
			return nil
		}
		result.Value = string(n.Text)
	case *parse.IfNode:
		result.Details = map[string]any{
			"condition": n.Pipe.String(),
		}
		result.Children = append(result.Children, &Node{
			Type:     "then",
			Children: convertNodes(n.List.Nodes),
		})
		if n.ElseList != nil {
			result.Children = append(result.Children, &Node{
				Type:     "else",
				Children: convertNodes(n.ElseList.Nodes),
			})
		}
	case *parse.RangeNode:
		result.Details = map[string]any{
			"range": n.Pipe.String(),
		}
		result.Children = append(result.Children, &Node{
			Type:     "body",
			Children: convertNodes(n.List.Nodes),
		})
		if n.ElseList != nil {
			result.Children = append(result.Children, &Node{
				Type:     "else",
				Children: convertNodes(n.ElseList.Nodes),
			})
		}
	case *parse.WithNode:
		result.Details = map[string]any{
			"with": n.Pipe.String(),
		}
		result.Children = append(result.Children, &Node{
			Type:     "body",
			Children: convertNodes(n.List.Nodes),
		})
		if n.ElseList != nil {
			result.Children = append(result.Children, &Node{
				Type:     "else",
				Children: convertNodes(n.ElseList.Nodes),
			})
		}
	case *parse.TemplateNode:
		result.Value = n.Name
	case *parse.ListNode:
		result.Children = convertNodes(n.Nodes)
	}

	return result
}

func convertNodes(nodes []parse.Node) []*Node {
	var result []*Node
	for _, n := range nodes {
		if converted := convertNode(n); converted != nil {
			result = append(result, converted)
		}
	}
	return result
}

// getNodeType returns the type name of the node without the package prefix
func getNodeType(node parse.Node) string {
	return fmt.Sprintf("%T", node)
}
