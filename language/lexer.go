package language

import (
	"fmt"
	"regexp"
	"strings"
)

// TokenType represents the type of a token
type TokenType string

const (
	TokenField     TokenType = "FIELD"
	TokenCondition TokenType = "CONDITION"
	TokenOperator  TokenType = "OPERATOR"
	TokenValue     TokenType = "VALUE"
	TokenLogical   TokenType = "LOGICAL"
	TokenSeparator TokenType = "SEPARATOR"
	TokenTransform TokenType = "TRANSFORM"
	TokenInvalid   TokenType = "INVALID"
)

// Token represents a single token
type Token struct {
	Type  TokenType
	Value string
}

// Lexer for parsing rules
type Lexer struct {
	input string
	pos   int
}

// NewLexer initializes a lexer with the input string
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: strings.TrimSpace(input),
		pos:   0,
	}
}

// Tokenize splits the input into tokens
func (l *Lexer) Tokenize(input string) ([]Token, error) {
	var tokens []Token
	pos := 0
	patterns := map[TokenType]*regexp.Regexp{
		TokenField:     regexp.MustCompile(`^FIELD\("([^"]+)"\)`),                    // Match FIELD("field_name")
		TokenCondition: regexp.MustCompile(`^(TYPE|RANGE|MATCHES|IN|REQUIRED)`),      // Custom conditions
		TokenValue:     regexp.MustCompile(`^"([^"]*)"|'([^']*)'|[\d\.]+|\([^)]*\)`), // Match strings, numbers, lists
		TokenLogical:   regexp.MustCompile(`^(AND|OR|NOT)`),                          // Logical operators
		TokenSeparator: regexp.MustCompile(`^,`),                                     // Separators
	}

	for pos < len(input) {
		input = strings.TrimSpace(input[pos:])
		pos = 0

		matched := false
		for tokenType, pattern := range patterns {
			if loc := pattern.FindStringIndex(input); loc != nil && loc[0] == 0 {
				value := input[loc[0]:loc[1]]
				tokens = append(tokens, Token{Type: tokenType, Value: value})
				pos += len(value)
				matched = true
				break
			}
		}

		if !matched {
			return nil, fmt.Errorf("unexpected token at: %s", input)
		}
	}

	return tokens, nil
}
