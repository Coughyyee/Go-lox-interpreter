package main

type StmtVisitor interface {
	VisitBlockStmt(*BlockStmt) interface{}
	VisitExpressionStmt(*ExpressionStmt) interface{}
	VisitIfStmt(*IfStmt) interface{}
	VisitPrintStmt(*PrintStmt) interface{}
	VisitVarStmt(*VarStmt) interface{}
	VisitWhileStmt(*WhileStmt) interface{}
	VisitBreakStmt(*BreakStmt) interface{}
}

type Stmt interface {
	accept(StmtVisitor) interface{}
}

type BlockStmt struct {
	statements []Stmt
}

type ExpressionStmt struct {
	expression Expr
}

type IfStmt struct {
	condition Expr
	thenBranch Stmt
	elseBranch Stmt
}

type PrintStmt struct {
	expression Expr
}

type VarStmt struct {
	name *Token
	initializer Expr
}

type WhileStmt struct {
	condition Expr
	body Stmt
}

type BreakStmt struct {
}

func (b *BlockStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBlockStmt(b)
}

func (e *ExpressionStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

func (i *IfStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitIfStmt(i)
}

func (p *PrintStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(p)
}

func (v *VarStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(v)
}

func (w *WhileStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitWhileStmt(w)
}

func (b *BreakStmt) accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBreakStmt(b)
}

