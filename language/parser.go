package language

import (
	"errors"
)

// Node represents a node in the Abstract Syntax Tree (AST)
type Node struct {
	Type     TokenType
	Value    string
	Children []*Node
}

// Parser for validation and transformation rules
type Parser struct{}

// NewParser initializes a parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseRules parses the parameters into an AST
func (p *Parser) ParseRules(params []string) (*Node, error) {
	if len(params) < 3 {
		return nil, errors.New("insufficient parameters")
	}

	root := &Node{Type: "ROOT", Children: []*Node{}}

	for i := 0; i < len(params); i += 3 {
		if i+2 >= len(params) {
			return nil, errors.New("incomplete expression")
		}

		field := params[i]
		condition := params[i+1]
		value := params[i+2]

		node := &Node{Type: "EXPRESSION", Children: []*Node{
			{Type: "FIELD", Value: field},
			{Type: "CONDITION", Value: condition},
			{Type: "VALUE", Value: value},
		}}

		root.Children = append(root.Children, node)
	}

	return root, nil
}
