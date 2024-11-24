package language

import (
	"errors"
	"fmt"
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

func (p *Parser) ParseRules(tokens []Token) (*Node, error) {
	if len(tokens) < 3 {
		return nil, errors.New("insufficient parameters")
	}

	root := &Node{Type: "ROOT", Children: []*Node{}}
	var currentField string

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if token.Type == "FIELD" {
			// Set the current field and continue to the next token
			currentField = token.Value
		} else if token.Type == "CONDITION" {
			// Ensure there is a following value
			if i+1 >= len(tokens) {
				return nil, errors.New("expected value after condition")
			}

			condition := token
			value := tokens[i+1] // Next token is the value

			node := &Node{Type: "EXPRESSION", Children: []*Node{
				{Type: "FIELD", Value: currentField},
				{Type: "CONDITION", Value: condition.Value},
				{Type: "VALUE", Value: value.Value},
			}}

			root.Children = append(root.Children, node)

			// Move past the value token
			i++
		} else {
			return nil, fmt.Errorf("unexpected token: %s", token.Value)
		}
	}

	return root, nil
}
