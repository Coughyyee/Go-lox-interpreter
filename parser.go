package main

import (
	"fmt"
	"log"
)

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) declaration() Stmt {
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		return &BlockStmt{
			statements: p.block(),
		}
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, fmt.Sprintf("Expect %v';'%v after value.", YELLOW, RESET))
	return &PrintStmt{
		expression: value,
	}
}

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

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, fmt.Sprintf("Expect %v';'%v after expression.", YELLOW, RESET))
	return &ExpressionStmt{
		expression: expr,
	}
}

func (p *Parser) block() []Stmt {
	var statements []Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, fmt.Sprintf("Expected %v'}'%v after block.", YELLOW, RESET))
	return statements
}

func (p *Parser) assignment() Expr {
	expr := p.equality()

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

func (p *Parser) match(types ...TokenType) bool {
	for _, ttype := range types {
		if p.check(ttype) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokenType TokenType, message string) *Token {
	if p.check(tokenType) {
		return p.advance()
	}

	log.Fatal(ReportExit(p.peek().line, "", message))
	return nil
}

func (p *Parser) check(ttype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == ttype
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

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
