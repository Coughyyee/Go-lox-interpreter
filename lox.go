package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Lox struct{}

func NewLox(hadError bool) *Lox {
	return &Lox{}
}

// run is the function that calls the interpreters interpreting functionalities.
func (lox *Lox) run(source string) {
	scanner := NewScanner(source, lox)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	statements := parser.Parse()

	interpreter := NewInterpreter()
	interpreter.Interpret(statements)

	// fmt.Printf("\n%s%-15s%s %s%-50s%s %s%-50s%s\n\n",
	// 	WHITE, "TOKEN ↓", RESET,
	// 	WHITE, "LEXEME ↓", RESET,
	// 	WHITE, "LITERAL ↓", RESET)
	// for _, token := range tokens {
	// 	fmt.Println(token.toString())
	// }
}

// runFile is the function that runs when a valid file path is supplied
// into the arguments.
func (lox *Lox) runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to read file")
	}

	lox.run(string(bytes))
}

// runPrompt is the function that runs when no arguments are passed in.
// Similar to pythons prompt when running 'python<CR>'.
func (lox *Lox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Error reading input: ", err)
			continue
		}

		line = strings.TrimSuffix(line, "\n")
		lox.run(line)
	}
}
