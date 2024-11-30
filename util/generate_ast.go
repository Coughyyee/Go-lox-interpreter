package main

import (
	"fmt"
	"os"
	"strings"
)

type TreeType struct {
	baseClassName string
	className     string
	fields        []string
}

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate-ast <output dir>")
		os.Exit(64)
	}

	outputDir := args[1]

	defineAst(outputDir, "Expr", []string{
		"Assign : *Token name, Expr value",
		"Binary : Expr left, *Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal : interface{} value",
		"Unary : *Token operator, Expr right",
		"Variable : *Token name",
	})

	defineAst(outputDir, "Stmt", []string{
		"Block : []Stmt statements",
		"Expression : Expr expression",
		"Print : Expr expression",
		"Var : *Token name, Expr initializer",
	})
}

func defineAst(outputDir string, baseName string, types []string) error {
	path := fmt.Sprintf("%s/%s.go", outputDir, strings.ToLower(baseName))
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	var treeTypes []TreeType

	file.Write([]byte("package main\n\n"))

	// visitor interface
	file.Write([]byte(fmt.Sprintf("type %sVisitor interface {\n", baseName)))
	for _, t := range types {
		split := strings.Split(t, ":") // baseClassName : Args
		baseClassName := strings.TrimRight(split[0], " ")
		className := fmt.Sprintf("%v%v", baseClassName, baseName) // e.g Binary + Expr
		file.Write([]byte(fmt.Sprintf("\tVisit%s(*%s) interface{}\n", className, className)))
	}
	file.Write([]byte("}\n\n"))

	// data
	for _, t := range types {
		split := strings.Split(t, ":") // baseClassName : Args
		baseClassName := strings.TrimRight(split[0], " ")
		className := fmt.Sprintf("%v%v", baseClassName, baseName) // e.g Binary + Expr
		arg_split := strings.Split(split[1], ",")
		var fields []string
		for _, arg := range arg_split {
			trimed := strings.TrimLeft(arg, " ")
			f := strings.Split(trimed, " ")
			fields = append(fields, fmt.Sprintf("%s %s", f[1], f[0]))
		}
		treeTypes = append(treeTypes, TreeType{baseClassName: baseClassName, className: className, fields: fields})
	}

	// base name struct
	file.Write([]byte(fmt.Sprintf("type %s interface {\n", baseName)))
	file.Write([]byte(fmt.Sprintf("\taccept(%sVisitor) interface{}\n", baseName)))
	file.Write([]byte("}\n\n"))

	// structs
	for _, t := range treeTypes {
		file.Write([]byte(fmt.Sprintf("type %s struct {\n", t.className)))
		for _, f := range t.fields {
			file.Write([]byte(fmt.Sprintf("\t%s\n", f)))
		}
		file.Write([]byte("}\n\n"))
	}

	// func accepts
	for _, t := range treeTypes {
		implName := strings.ToLower(string(t.className[0]))
		file.Write([]byte(fmt.Sprintf("func (%s *%s) accept(visitor %sVisitor) interface{} {\n", implName, t.className, baseName)))
		file.Write([]byte(fmt.Sprintf("\treturn visitor.Visit%s(%s)\n", t.className, implName)))
		file.Write([]byte("}\n\n"))
	}

	return nil
}
