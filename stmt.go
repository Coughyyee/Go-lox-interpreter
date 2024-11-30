package main

type StmtVisitor interface {
	VisitBlockStmt(*BlockStmt) interface{}
	VisitExpressionStmt(*ExpressionStmt) interface{}
	VisitPrintStmt(*PrintStmt) interface{}
	VisitVarStmt(*VarStmt) interface{}
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

type PrintStmt struct {
	expression Expr
}

type VarStmt struct {
	name *Token
	initializer Expr
}

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

