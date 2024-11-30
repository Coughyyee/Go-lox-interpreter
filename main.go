package main

import (
	"log"
	"os"
)

// ANSI escape codes for colored text.
const (
	RED    = "\033[31m"
	YELLOW = "\033[33m"
	GREEN  = "\033[32m"
	WHITE  = "\033[97m"
	RESET  = "\033[0m"
)

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
