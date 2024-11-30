package main

import (
	"fmt"
	"log"
)

type Environment struct {
	enclosing *Environment
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{

		enclosing: nil,
		values: make(map[string]interface{}),
	}
}

func NewEnclosingEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values: make(map[string]interface{}),
	}
}

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

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}
