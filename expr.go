package main

type ExprVisitor interface {
	VisitAssignExpr(*AssignExpr) interface{}
	VisitBinaryExpr(*BinaryExpr) interface{}
	VisitCallExpr(*CallExpr) interface{}
	VisitGroupingExpr(*GroupingExpr) interface{}
	VisitLiteralExpr(*LiteralExpr) interface{}
	VisitLogicalExpr(*LogicalExpr) interface{}
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

type CallExpr struct {
	callee Expr
	paren *Token
	arguments []Expr
}

type GroupingExpr struct {
	expression Expr
}

type LiteralExpr struct {
	value interface{}
}

type LogicalExpr struct {
	left Expr
	operator *Token
	right Expr
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

func (c *CallExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitCallExpr(c)
}

func (g *GroupingExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

func (l *LiteralExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

func (l *LogicalExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicalExpr(l)
}

func (u *UnaryExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

func (v *VariableExpr) accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

