// Package main implements a Lox language interpreter
package main

// Stmt is the interface that all statement types must implement.
// It defines the Visitor pattern for traversing the statement AST.
type Stmt interface {
	accept(visitor StmtVisitor) interface{}
}

// StmtVisitor defines the interface for visiting different statement types.
// Each method corresponds to a specific statement type in the AST.
type StmtVisitor interface {
	VisitBlockStmt(stmt *BlockStmt) interface{}
	VisitExpressionStmt(stmt *ExpressionStmt) interface{}
	VisitPrintStmt(stmt *PrintStmt) interface{}
	VisitVarStmt(stmt *VarStmt) interface{}
}

// BlockStmt represents a block of statements.
// Example: { stmt1; stmt2; }
type BlockStmt struct {
	statements []Stmt  // List of statements in the block
}

// ExpressionStmt represents an expression statement.
// Example: print "hello";
type ExpressionStmt struct {
	expression Expr  // The expression to evaluate
}

// PrintStmt represents a print statement.
// Example: print "hello";
type PrintStmt struct {
	expression Expr  // The expression to print
}

// VarStmt represents a variable declaration statement.
// Example: var x = 42;
type VarStmt struct {
	name        *Token  // The variable name token
	initializer Expr    // The initial value expression (may be nil)
}

// Visitor pattern implementation methods

func (b *BlockStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBlockStmt(b)
}

func (e *ExpressionStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

func (p *PrintStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(p)
}

func (v *VarStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(v)
}
