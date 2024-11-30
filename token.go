package main

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

// toString is the function that returns a string with metadata from the Token struct.
func (token *Token) toString() string {
	return fmt.Sprintf("%s%-15v%s %s%-50v%s %s%-50v%s",
		RED, token.tokenType.toString(), RESET,
		YELLOW, token.lexeme, RESET,
		GREEN, token.literal, RESET)
}
