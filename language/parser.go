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
type Parser struct {
	tokens []Token
	pos    int
}

// NewParser initializes a parser with tokens
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// ParseRules parses the tokens into an AST
func (p *Parser) ParseRules() (*Node, error) {
	root := &Node{Type: "ROOT", Children: []*Node{}}

	for p.pos < len(p.tokens) {
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		root.Children = append(root.Children, node)
	}

	return root, nil
}

func (p *Parser) parseExpression() (*Node, error) {
	// Example rule: FIELD CONDITION VALUE [LOGICAL FIELD CONDITION VALUE]
	field := p.consume(TokenField)
	if field == nil {
		return nil, errors.New("expected field")
	}

	condition := p.consume(TokenCondition)
	if condition == nil {
		return nil, fmt.Errorf("expected condition after field %s", field.Value)
	}

	value := p.consume(TokenValue)
	if value == nil {
		return nil, fmt.Errorf("expected value after condition %s", condition.Value)
	}

	node := &Node{Type: "EXPRESSION", Children: []*Node{
		{Type: field.Type, Value: field.Value},
		{Type: condition.Type, Value: condition.Value},
		{Type: value.Type, Value: value.Value},
	}}

	// Check for logical operators (AND, OR, NOT)
	logical := p.consume(TokenLogical)
	if logical != nil {
		rightExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, &Node{
			Type:     logical.Type,
			Value:    logical.Value,
			Children: []*Node{rightExpr},
		})
	}

	return node, nil
}

// consume retrieves the next token if it matches the expected type
func (p *Parser) consume(expected TokenType) *Token {
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == expected {
		token := p.tokens[p.pos]
		p.pos++
		return &token
	}
	return nil
}
