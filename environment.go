// Package main implements a Lox language interpreter
package main

import (
	"fmt"
	"log"
)

// Environment represents a scope in the Lox language.
// It maintains a mapping of variable names to their values.
type Environment struct {
	enclosing *Environment // Reference to the enclosing (outer) scope
	values    map[string]interface{} // Map of variable names to their values
}

// NewEnvironment creates a new Environment instance.
// Used for creating a new global scope.
func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    make(map[string]interface{}),
	}
}

// NewEnclosingEnvironment creates a new Environment with an enclosing scope.
// Used for creating block scopes that can access their parent scope.
func NewEnclosingEnvironment(enclosing *Environment) *Environment {
	env := NewEnvironment()
	env.enclosing = enclosing
	return env
}

// define defines a new variable in the current scope.
// If the variable already exists, its value is updated.
func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

// get retrieves the value of a variable.
// Searches in the current scope and then in enclosing scopes.
func (e *Environment) get(name *Token) interface{} {
	if value, ok := e.values[name.lexeme]; ok {
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	log.Fatal(ReportExit(name.line, "", fmt.Sprintf("Undefined variable %v'%v'%v.", YELLOW, name.lexeme, RESET)))
	return nil
}

// assign updates the value of an existing variable.
// Searches in the current scope and then in enclosing scopes.
func (e *Environment) assign(name *Token, value interface{}) {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}

	log.Fatal(ReportExit(name.line, "", fmt.Sprintf("Undefined variable %v'%v'%v.", YELLOW, name.lexeme, RESET)))
}
