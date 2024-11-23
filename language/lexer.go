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
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	patterns := map[TokenType]*regexp.Regexp{
		TokenField:     regexp.MustCompile(`^[a-zA-Z0-9_\.]+`),         // Matches field names
		TokenCondition: regexp.MustCompile(`^(==|!=|>=|<=|>|<)`),       // Matches conditions
		TokenOperator:  regexp.MustCompile(`^(->|=>)`),                 // Matches transformation operators
		TokenValue:     regexp.MustCompile(`^"([^"]*)"|'([^']*)'|\d+`), // Matches strings or numbers
		TokenLogical:   regexp.MustCompile(`^(AND|OR|NOT)`),            // Matches logical operators
		TokenSeparator: regexp.MustCompile(`^,`),                       // Matches separators
		TokenTransform: regexp.MustCompile(`^TRANSFORM\(`),             // Matches the transform keyword
	}

	for l.pos < len(l.input) {
		// Skip whitespace
		l.input = strings.TrimSpace(l.input[l.pos:])
		l.pos = 0

		matched := false
		for tokenType, pattern := range patterns {
			loc := pattern.FindStringIndex(l.input)
			if loc != nil && loc[0] == 0 {
				value := l.input[loc[0]:loc[1]]
				tokens = append(tokens, Token{Type: tokenType, Value: value})
				l.pos += len(value)
				matched = true
				break
			}
		}

		if !matched {
			return nil, fmt.Errorf("unexpected token at: %s", l.input)
		}
	}

	return tokens, nil
}
