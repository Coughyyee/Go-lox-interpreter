// Package main implements a Lox language interpreter
package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// Interpreter implements the execution engine for the Lox language.
// It evaluates expressions and executes statements in the AST.
type Interpreter struct {
	globals     *Environment
	environment *Environment
}

// NewInterpreter creates a new Interpreter instance.
func NewInterpreter() *Interpreter {
	globals := NewEnvironment()
	globals.define("clock", NewClock())
	return &Interpreter{
		globals:     globals,
		environment: globals,
	}
}

type Clock struct{}

func NewClock() *Clock {
	return &Clock{}
}

func (*Clock) arity() int {
	return 0
}

func (*Clock) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return float64(time.Now().UnixNano()) / 1e9
}

func (*Clock) String() string {
	return "<native fn>"
}

// Interpret interprets a list of statements.
// This is the main entry point for program execution.
func (i *Interpreter) Interpret(statements []Stmt) interface{} {
	var result interface{}
	for _, statement := range statements {
		result = i.execute(statement)
	}
	return result
}

// VisitLiteralExpr evaluates a literal expression.
// Returns the literal value directly.
func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr) interface{} {
	return expr.value
}

func (i *Interpreter) VisitLogicalExpr(expr *LogicalExpr) interface{} {
	left := i.evaluate(expr.left)

	if expr.operator.tokenType == OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.right)
}

// VisitGroupingExpr evaluates a grouping expression.
// Evaluates the expression inside the parentheses.
func (i *Interpreter) VisitGroupingExpr(expr *GroupingExpr) interface{} {
	return i.evaluate(expr.expression)
}

// VisitUnaryExpr evaluates a unary expression.
// Handles negation (-) and logical not (!) operators.
func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr) interface{} {
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		i.checkNumberOperand(expr.operator, right)
		return -right.(float64)
	}

	return nil
}

// VisitBinaryExpr evaluates a binary expression.
// Handles arithmetic, comparison, and equality operators.
func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) interface{} {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.tokenType {
	case MINUS:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float64) - right.(float64)
	case PLUS:
		// number + number.
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l + r
			}
		}

		// string + string.
		if l, ok := left.(string); ok {
			if r, ok := right.(string); ok {
				return l + r
			}
		}

		// string + number.
		if l, ok := left.(string); ok {
			if r, ok := right.(float64); ok {
				return fmt.Sprintf("%v%v", l, r)
			}
		}

		// number + string.
		if l, ok := left.(float64); ok {
			if r, ok := right.(string); ok {
				return fmt.Sprintf("%v%v", l, r)
			}
		}

		log.Fatal(ReportExit(expr.operator.line, "", "Operands must be two numbers or two strings."))
	case SLASH:
		i.checkNumberOperands(expr.operator, left, right)
		// assert no division by 0.
		if left.(float64) == 0 || right.(float64) == 0 {
			log.Fatal(ReportExit(expr.operator.line, "", "Division by 0 is not allowed."))
		}
		return left.(float64) / right.(float64)
	case STAR:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float64) * right.(float64)
	case GREATER:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		i.checkNumberOperands(expr.operator, left, right)
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) VisitCallExpr(expr *CallExpr) interface{} {
	callee := i.evaluate(expr.callee)

	var arguments []interface{}
	for _, argument := range expr.arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	if _, ok := callee.(LoxCallable); !ok {
		log.Fatal(ReportExit(expr.paren.line, "", "Can't call non-callable object."))
	}

	function := callee.(LoxCallable)
	if len(arguments) != function.arity() {
		log.Fatal(ReportExit(expr.paren.line, "", fmt.Sprintf("Expected %v arguments but got %v.", function.arity(), len(arguments))))
	}
	return function.call(i, arguments)
}

// VisitVariableExpr evaluates a variable expression.
// Looks up the variable's value in the current environment.
func (i *Interpreter) VisitVariableExpr(expr *VariableExpr) interface{} {
	return i.environment.get(expr.name)
}

// VisitAssignExpr evaluates an assignment expression.
// Updates the variable's value in the current environment.
func (i *Interpreter) VisitAssignExpr(expr *AssignExpr) interface{} {
	value := i.evaluate(expr.value)
	i.environment.assign(expr.name, value)
	return value
}

// VisitExpressionStmt executes an expression statement.
func (i *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) interface{} {
	return i.evaluate(stmt.expression)
}

func (i *Interpreter) VisitIfStmt(stmt *IfStmt) interface{} {
	if i.isTruthy(i.evaluate(stmt.condition)) {
		return i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		return i.execute(stmt.elseBranch)
	}

	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *FunctionStmt) interface{} {
	function := NewLoxFunction(stmt, i.environment)
	i.environment.define(stmt.name.lexeme, function)
	return nil
}

// VisitPrintStmt executes a print statement.
// Evaluates the expression and prints its value.
func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt) interface{} {
	var token *Token
	// check if its a variable expression.
	if v, ok := stmt.expression.(*VariableExpr); ok {
		token = v.name
	}
	value := i.evaluate(stmt.expression)
	fmt.Println(stringify(token, value))
	return nil
}

// ReturnError is used to handle return statements
type ReturnError struct {
	value interface{}
}

func (e *ReturnError) Error() string {
	return "Return statement error"
}

func (i *Interpreter) VisitReturnStmt(stmt *ReturnStmt) interface{} {
	var value interface{}
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	}
	return &ReturnError{value: value}
}

// VisitVarStmt executes a variable declaration statement.
// Defines a new variable in the current environment.
func (i *Interpreter) VisitVarStmt(stmt *VarStmt) interface{} {
	var value interface{}
	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}

	i.environment.define(stmt.name.lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) interface{} {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*BreakError); !ok {
				panic(r) // re-panic if it's not a break
			}
		}
	}()

	var result interface{}
	for i.isTruthy(i.evaluate(stmt.condition)) {
		result = i.execute(stmt.body)
	}
	return result
}

// VisitBlockStmt executes a block statement.
// Creates a new environment for the block's scope.
func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) interface{} {
	return i.executeBlock(stmt.statements, NewEnclosingEnvironment(i.environment))
}

func (i *Interpreter) VisitBreakStmt(stmt *BreakStmt) interface{} {
	panic(&BreakError{})
}

// BreakError is used to handle break statements
type BreakError struct{}

func (e *BreakError) Error() string {
	return "Break statement"
}

// execute executes a statement.
func (i *Interpreter) execute(stmt Stmt) interface{} {
	return stmt.accept(i)
}

// executeBlock executes a block of statements.
// Creates a new environment for the block's scope.
func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) interface{} {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()

	i.environment = environment
	var result interface{}
	for _, statement := range statements {
		result = i.execute(statement)
		if _, ok := result.(*ReturnError); ok {
			return result
		}
	}
	return result
}

// evaluate evaluates an expression and returns its value.
func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(i)
}

// isTruthy determines if a value is considered true in Lox.
// nil and false are falsey, everything else is truthy.
func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if v, ok := object.(bool); ok {
		return v
	}
	return true
}

// isEqual determines if two values are equal.
// Uses the == operator for comparison.
func (i *Interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle number comparisons
	if aNum, ok := a.(float64); ok {
		if bNum, ok := b.(float64); ok {
			return aNum == bNum
		}
		return false
	}

	// Handle string comparisons
	if aStr, ok := a.(string); ok {
		if bStr, ok := b.(string); ok {
			return aStr == bStr
		}
		return false
	}

	// Handle boolean comparisons
	if aBool, ok := a.(bool); ok {
		if bBool, ok := b.(bool); ok {
			return aBool == bBool
		}
		return false
	}

	return a == b
}

// checkNumberOperand verifies that an operand is a number.
// Throws a runtime error if the operand is not a number.
func (i *Interpreter) checkNumberOperand(operator *Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}
	log.Fatal(ReportExit(operator.line, "", "Operand must be a number."))
}

// checkNumberOperands verifies that both operands are numbers.
// Throws a runtime error if either operand is not a number.
func (i *Interpreter) checkNumberOperands(operator *Token, left, right interface{}) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	log.Fatal(ReportExit(operator.line, "", "Operands must be numbers."))
}

// stringify converts a value to a string representation.
// Handles nil, numbers, and strings.
func stringify(token *Token, object interface{}) string {
	if object == nil {
		log.Fatal(ReportExit(token.line, "", fmt.Sprintf("Variable %v'%v'%v is undefined.", YELLOW, token.lexeme, RESET)))
	}

	if v, ok := object.(float64); ok {
		text := fmt.Sprintf("%f", v)
		// Trim ending if returned value number from expression isnt a float.
		if strings.HasSuffix(text, ".000000") {
			text = text[:len(text)-7]
		}
		return text
	}

	return fmt.Sprintf("%v", object)
}
