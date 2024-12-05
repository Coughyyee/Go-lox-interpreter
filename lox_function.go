package main

type LoxFunction struct {
	declaration *FunctionStmt
	closure     *Environment
}

func NewLoxFunction(declaration *FunctionStmt, closure *Environment) *LoxFunction {
	return &LoxFunction{declaration: declaration, closure: closure}
}

func (f *LoxFunction) call(interpreter *Interpreter, arguments []interface{}) interface{} {
	environment := NewEnclosingEnvironment(f.closure)
	for i, param := range f.declaration.params {
		environment.define(param.lexeme, arguments[i])
	}

	result := interpreter.executeBlock(f.declaration.body, environment)
	if returnError, ok := result.(*ReturnError); ok {
		return returnError.value
	}
	return nil
}

func (f *LoxFunction) arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) String() string {
	return "<fn " + f.declaration.name.lexeme + ">"
}