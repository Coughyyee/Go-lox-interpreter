// Package main implements a Lox language interpreter
package main

import (
	"fmt"
	"log"
)

// Parser implements a recursive descent parser for the Lox language.
// It takes a sequence of tokens and produces an abstract syntax tree (AST).
type Parser struct {
	tokens  []*Token // List of tokens to parse
	current int      // Current position in the token list
}

// NewParser creates a new Parser instance with the given tokens.
func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Parse parses the tokens and returns a slice of statements.
// This is the entry point for syntactic analysis.
func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

// expression parses an expression.
// Handles the lowest precedence level of expressions.
func (p *Parser) expression() Expr {
	return p.assignment()
}

// declaration parses a declaration statement (var, function, etc.).
func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

// statement parses a statement (expression, print, block, etc.).
func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(LEFT_BRACE) {
		return &BlockStmt{
			statements: p.block(),
		}
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect %v'('%v after %v'for'%v.", YELLOW, RESET, YELLOW, RESET))

	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, fmt.Sprintf("Expected %v';'%v after loop condition.", YELLOW, RESET))

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, fmt.Sprintf("Expected %v')'%v after for clauses.", YELLOW, RESET))
	body := p.statement()

	if increment != nil {
		body = &BlockStmt{
			statements: []Stmt{
				body,
				&ExpressionStmt{
					expression: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = &LiteralExpr{
			value: true,
		}
	}
	body = &WhileStmt{
		condition: condition,
		body: body,
	}

	if initializer != nil {
		body = &BlockStmt{
			[]Stmt{
				initializer,
				body,
			},
		}
	}

	return body
}

// ifStatement parses an if statement.
func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect %v'('%v after %v'if'%v.", YELLOW, RESET, YELLOW, RESET))
	condition := p.expression()
	p.consume(RIGHT_PAREN, fmt.Sprintf("Expect %v')'%v after if condition.", YELLOW, RESET))

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return &IfStmt{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

// printStatement parses a print statement.
func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, fmt.Sprintf("Expect %v';'%v after value.", YELLOW, RESET))
	return &PrintStmt{
		expression: value,
	}
}

// varDeclaration parses a variable declaration statement.
func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, fmt.Sprintf("Expected %v';'%v after variable declaration.", YELLOW, RESET))
	return &VarStmt{
		name:        name,
		initializer: initializer,
	}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect %v'('%v after '%v'while'%v.", YELLOW, RESET, YELLOW, RESET))
	condition := p.expression()
	p.consume(RIGHT_PAREN, fmt.Sprintf("Expect %v')'%v after condition.", YELLOW, RESET))
	body := p.statement()

	return &WhileStmt{
		condition: condition,
		body: body,
	}
}

// expressionStatement parses an expression statement.
func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, fmt.Sprintf("Expect %v';'%v after expression.", YELLOW, RESET))
	return &ExpressionStmt{
		expression: expr,
	}
}

// block parses a block of statements.
func (p *Parser) block() []Stmt {
	var statements []Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, fmt.Sprintf("Expected %v'}'%v after block.", YELLOW, RESET))
	return statements
}

// assignment parses an assignment expression.
func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		token, ok := expr.(*VariableExpr)
		if ok {
			name := token.name
			return &AssignExpr{
				name:  name,
				value: value,
			}
		}

		log.Fatal(ReportExit(p.peek().line, "", fmt.Sprintf("%v[%v]%v Invalid assignment target.", YELLOW, equals, RESET)))
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &LogicalExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = &LogicalExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

// equality parses equality expressions (==, !=).
func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

// comparison parses comparison expressions (>, >=, <, <=).
func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

// term parses addition and subtraction expressions.
func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

// factor parses multiplication and division expressions.
func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &BinaryExpr{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}

	return expr
}

// unary parses unary expressions (!expr, -expr).
func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return &UnaryExpr{
			operator: operator,
			right:    right,
		}
	}

	return p.primary()
}

// primary parses primary expressions (literals, grouping).
func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return &LiteralExpr{value: false}
	}

	if p.match(TRUE) {
		return &LiteralExpr{value: true}
	}

	if p.match(NIL) {
		return &LiteralExpr{value: nil}
	}

	if p.match(NUMBER, STRING) {
		return &LiteralExpr{
			value: p.previous().literal,
		}
	}

	if p.match(IDENTIFIER) {
		return &VariableExpr{p.previous()}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, fmt.Sprintf("Expect %v')'%v after expression.", YELLOW, RESET))
		return &GroupingExpr{expression: expr}
	}

	log.Fatal(ReportExit(p.peek().line, "", "Expected expression."))
	return nil
}

// match checks if the current token matches any of the given types.
// Returns true and advances if there's a match.
func (p *Parser) match(types ...TokenType) bool {
	for _, ttype := range types {
		if p.check(ttype) {
			p.advance()
			return true
		}
	}

	return false
}

// consume consumes the current token if it matches the expected type.
// Throws an error if it doesn't match.
func (p *Parser) consume(tokenType TokenType, message string) *Token {
	if p.check(tokenType) {
		return p.advance()
	}

	log.Fatal(ReportExit(p.peek().line, "", message))
	return nil
}

// check checks if the current token is of the expected type.
func (p *Parser) check(ttype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == ttype
}

// advance moves to the next token and returns the previous one.
func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd checks if we've reached the end of the token list.
func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

// peek returns the current token without advancing.
func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

// previous returns the previous token.
func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

// synchronize recovers from a parse error by discarding tokens
// until it reaches a likely statement boundary.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS:
		case FUN:
		case VAR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}
