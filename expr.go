// Package main implements a Lox language interpreter
package main

// Expr is the interface that all expression types must implement.
// It defines the Visitor pattern for traversing the expression AST.
type Expr interface {
	accept(visitor ExprVisitor) interface{}
}

// ExprVisitor defines the interface for visiting different expression types.
// Each method corresponds to a specific expression type in the AST.
type ExprVisitor interface {
	VisitAssignExpr(expr *AssignExpr) interface{}
	VisitBinaryExpr(expr *BinaryExpr) interface{}
	VisitGroupingExpr(expr *GroupingExpr) interface{}
	VisitLiteralExpr(expr *LiteralExpr) interface{}
	VisitUnaryExpr(expr *UnaryExpr) interface{}
	VisitVariableExpr(expr *VariableExpr) interface{}
}

// AssignExpr represents a variable assignment expression.
// Example: x = 42
type AssignExpr struct {
	name  *Token      // The variable being assigned to
	value Expr        // The value being assigned
}

// BinaryExpr represents a binary operation expression.
// Example: a + b, x * y
type BinaryExpr struct {
	left     Expr    // Left operand
	operator *Token  // Operator token
	right    Expr    // Right operand
}

// GroupingExpr represents a parenthesized expression.
// Example: (1 + 2)
type GroupingExpr struct {
	expression Expr  // The expression being grouped
}

// LiteralExpr represents a literal value expression.
// Example: 42, "hello", true
type LiteralExpr struct {
	value interface{}  // The literal value
}

// UnaryExpr represents a unary operation expression.
// Example: !true, -42
type UnaryExpr struct {
	operator *Token  // Operator token
	right    Expr    // Operand
}

// VariableExpr represents a variable reference expression.
// Example: x, counter
type VariableExpr struct {
	name *Token  // The variable name token
}

// Visitor pattern implementation methods

func (e *AssignExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(e)
}

func (e *BinaryExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(e)
}

func (e *GroupingExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(e)
}

func (e *LiteralExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(e)
}

func (e *UnaryExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(e)
}

func (e *VariableExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(e)
}
