package main

import (
	"log"
	"strconv"
)

type Scanner struct {
	source   string
	tokens   []*Token
	start    int
	current  int
	line     int
	keywords map[string]TokenType
}

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

// scanTokens is the function that scans the tokens until the end of the file.
// Each token is appended to the scanner.tokens array in the struct.
func (scanner *Scanner) scanTokens() []*Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}

	scanner.tokens = append(scanner.tokens, NewToken(EOF, "", nil, scanner.line))
	return scanner.tokens
}

// scanToken is the function that scans indivisual tokens and adds the token
// to the scanner.tokens arrray in the struct.
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

// identifier is the function that manages the type of token if a number
// is scanned.
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

// number is the function that manages numbers scanned.
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

// string is the function that manages strings scanned.
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

// match is the function that returns a bool based on if the current character
// is the same as the one passed into the parameter.
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

// peek is the function that returns the character that is one character ahead
// of the current character being scanned.
func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return 0
	}
	return scanner.source[scanner.current]
}

// peekNext is the function that returns the character two characters ahead
// of the current character being scanned.
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

// isAtEnd is the function that returns bool based on if the scanner is at the
// end of the source.
func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

// advance is the function that moves the scanners position one ahead and returns
// the new current character.
func (scanner *Scanner) advance() byte {
	if scanner.current >= len(scanner.source) {
		return byte(EOF)
	}
	ch := scanner.source[scanner.current]
	scanner.current++
	return ch
}

// advance is the function that moves the scanners position two ahead and returns
// the new current character.
func (scanner *Scanner) advanceNext() byte {
	if scanner.current >= len(scanner.source) {
		return byte(EOF)
	}
	ch := scanner.source[scanner.current+1]
	scanner.current += 2
	return ch
}

// addToken is the function that adds a token without any literal for the token.
func (scanner *Scanner) addToken(tokenType TokenType) {
	scanner.addTokenLiteral(tokenType, nil)
}

// addTokenLiteral is the function that adds a token and a literal for the token.
func (scanner *Scanner) addTokenLiteral(tokenType TokenType, literal interface{}) {
	text := scanner.source[scanner.start:scanner.current]
	scanner.tokens = append(scanner.tokens, NewToken(tokenType, text, literal, scanner.line))
}
