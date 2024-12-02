// Package main implements a Lox language interpreter
package main

// Token represents a lexical token in the Lox language.
// It contains information about the token type, lexeme, literal value, and line number.
type Token struct {
	tokenType TokenType   // Type identifies the category of the token
	lexeme    string      // Lexeme is the actual string value from the source code
	literal   interface{} // Literal holds the actual value for literals (numbers, strings, etc.)
	line      int         // Line indicates the line number where the token appears in source
}

// NewToken returns a new Token instance.
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}
