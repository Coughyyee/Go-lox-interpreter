package main

type ExprVisitor interface {
	VisitAssignExpr(*AssignExpr) interface{}
	VisitBinaryExpr(*BinaryExpr) interface{}
	VisitGroupingExpr(*GroupingExpr) interface{}
	VisitLiteralExpr(*LiteralExpr) interface{}
	VisitUnaryExpr(*UnaryExpr) interface{}
	VisitVariableExpr(*VariableExpr) interface{}
}

type Expr interface {
	accept(ExprVisitor) interface{}
}

type AssignExpr struct {
	name *Token
	value Expr
}

type BinaryExpr struct {
	left Expr
	operator *Token
	right Expr
}

type GroupingExpr struct {
	expression Expr
}

type LiteralExpr struct {
	value interface{}
}

type UnaryExpr struct {
	operator *Token
	right Expr
}

type VariableExpr struct {
	name *Token
}

func (a *AssignExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(a)
}

func (b *BinaryExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

func (g *GroupingExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

func (l *LiteralExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

func (u *UnaryExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

func (v *VariableExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

