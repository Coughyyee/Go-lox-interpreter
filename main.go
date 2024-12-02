// Package main implements a Lox language interpreter.
// This is the main entry point for the interpreter, handling file execution
// and interactive REPL mode.
package main

import (
	"log"
	"os"
)

// main is the entry point of the Lox interpreter.
// It supports two modes of operation:
// 1. File execution: jlox [script]
// 2. Interactive REPL: jlox
func main() {
	// log.SetFlags(0) // Removes the date before any log.Fatal().
	args := os.Args
	lox := NewLox(false)
	if len(args) > 2 {
		log.Fatal("Usage: jlox [script]")
		os.Exit(64)
	} else if len(args) == 2 {
		lox.runFile(args[1])
	} else {
		lox.runPrompt()
	}
}
