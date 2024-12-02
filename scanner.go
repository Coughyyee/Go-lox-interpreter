// Package main implements a Lox language interpreter
package main

import (
	"log"
	"strconv"
)

// Scanner performs lexical analysis on Lox source code.
// It converts the source text into a sequence of tokens.
type Scanner struct {
	source   string    // The source code being scanned
	tokens   []*Token  // List of tokens found during scanning
	start    int       // Start position of the current lexeme
	current  int       // Current position in the source
	line     int       // Current line number being scanned
	keywords map[string]TokenType
}

// NewScanner creates a new Scanner instance for the given source code.
func NewScanner(source string, lox *Lox) *Scanner {
	keywords := map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}

	scanner := Scanner{
		source:   source,
		start:    0,
		current:  0,
		line:     1,
		keywords: keywords,
	}

	return &scanner
}

// ScanTokens scans the source code and returns a list of tokens.
// This is the main entry point for lexical analysis.
func (scanner *Scanner) ScanTokens() []*Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}

	scanner.tokens = append(scanner.tokens, NewToken(EOF, "", nil, scanner.line))
	return scanner.tokens
}

// scanToken scans a single token from the source code.
// It identifies keywords, identifiers, literals, and operators.
func (scanner *Scanner) scanToken() {
	c := scanner.advance()
	switch c {
	case '(':
		scanner.addToken(LEFT_PAREN)
	case ')':
		scanner.addToken(RIGHT_PAREN)
	case '{':
		scanner.addToken(LEFT_BRACE)
	case '}':
		scanner.addToken(RIGHT_BRACE)
	case ',':
		scanner.addToken(COMMA)
	case '.':
		scanner.addToken(DOT)
	case '-':
		scanner.addToken(MINUS)
	case '+':
		scanner.addToken(PLUS)
	case ';':
		scanner.addToken(SEMICOLON)
	case '*':
		scanner.addToken(STAR)
	case '!':
		if scanner.match('=') {
			scanner.addToken(BANG_EQUAL)
		} else {
			scanner.addToken(BANG)
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken(EQUAL_EQUAL)
		} else {
			scanner.addToken(EQUAL)
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken(LESS_EQUAL)
		} else {
			scanner.addToken(LESS)
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(GREATER_EQUAL)
		} else {
			scanner.addToken(GREATER)
		}
	case '/':
		if scanner.match('/') {
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else if scanner.match('*') {
			for (scanner.peek() != '*' && scanner.peekNext() != '/') && !scanner.isAtEnd() {
				scanner.advance()
				// INFO: !scanner.isAtEnd shouldnt be here it should chuck an error if no close?
			}
			scanner.advanceNext() // consume the final '*' & '/' tokens
		} else {
			scanner.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t': // INFO: i have no clue if the cunt does the '\n' or just skips is. add break?
	case '\n':
		scanner.line++
	case '"':
		scanner.string()
	default:
		if scanner.isDigit(c) {
			scanner.number()
		} else if scanner.isAlpha(c) {
			scanner.identifier()
		} else {
			// scanner.lox.error(scanner.line, "Unexpected character.")
			log.Fatal(ReportExit(scanner.line, "", "Unexpected character."))
		}
	}
}

// identifier handles identifier and keyword scanning.
// It processes variable names and reserved keywords.
func (scanner *Scanner) identifier() {
	for scanner.isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	text := scanner.source[scanner.start:scanner.current]
	tokenType, ok := scanner.keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}

	scanner.addToken(tokenType)
}

// number handles numeric literal scanning.
// It processes both integer and decimal numbers.
func (scanner *Scanner) number() {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	if scanner.peek() == '.' && scanner.isDigit(scanner.peekNext()) {
		scanner.advance() // consume the "."

		for scanner.isDigit(scanner.peek()) {
			scanner.advance()
		}
	}

	number, err := strconv.ParseFloat(scanner.source[scanner.start:scanner.current], 64)
	if err != nil {
		log.Fatal(ReportExit(scanner.line, "", "Failed to parse float [scanner.number()].")) //? DEV?
	}

	scanner.addTokenLiteral(NUMBER, number)
}

// string handles string literal scanning.
// It processes the characters between double quotes.
func (scanner *Scanner) string() {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line++
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		log.Fatal(ReportExit(scanner.line, "", "Unterminated string."))
	}

	scanner.advance()

	value := scanner.source[scanner.start+1 : scanner.current-1]
	scanner.addTokenLiteral(STRING, value)
}

// match checks if the next character matches the expected one.
// Returns true and advances the cursor if there's a match.
func (scanner *Scanner) match(expected byte) bool {
	if scanner.isAtEnd() {
		return false
	}
	if scanner.source[scanner.current] != expected {
		return false
	}
	scanner.current++
	return true
}

// peek returns the next character without advancing the cursor.
func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return 0
	}
	return scanner.source[scanner.current]
}

// peekNext returns the character after the next one without advancing.
func (scanner *Scanner) peekNext() byte {
	if scanner.current+1 >= len(scanner.source) {
		return 0
	}
	return scanner.source[scanner.current+1]
}

// isAlpha is the function that returns a bool based on if the character is
// an alphabetical letter.
func (scanner *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isAlphaNumeric is the function that returns bool based on if the character
// is an alphabetical letter of a numeric digit.
func (scanner *Scanner) isAlphaNumeric(c byte) bool {
	return scanner.isAlpha(c) || scanner.isDigit(c)
}

// isDigit is the function that returns a bool based on if the character is a
// numeric value.
func (scanner *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isAtEnd checks if we've reached the end of the source code.
func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

// advance returns the current character and moves the cursor forward.
func (scanner *Scanner) advance() byte {
	if scanner.current >= len(scanner.source) {
		return byte(EOF)
	}
	ch := scanner.source[scanner.current]
	scanner.current++
	return ch
}

// advanceNext returns the character two positions ahead and moves the cursor two positions forward.
func (scanner *Scanner) advanceNext() byte {
	if scanner.current >= len(scanner.source) {
		return byte(EOF)
	}
	ch := scanner.source[scanner.current+1]
	scanner.current += 2
	return ch
}

// addToken adds a new token to the token list.
// It creates a token with the current lexeme and given type.
func (scanner *Scanner) addToken(tokenType TokenType) {
	scanner.addTokenLiteral(tokenType, nil)
}

// addTokenLiteral adds a new token with a literal value to the token list.
func (scanner *Scanner) addTokenLiteral(tokenType TokenType, literal interface{}) {
	text := scanner.source[scanner.start:scanner.current]
	scanner.tokens = append(scanner.tokens, NewToken(tokenType, text, literal, scanner.line))
}
