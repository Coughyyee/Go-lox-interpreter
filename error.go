// Package main implements a Lox language interpreter
package main

import (
	"fmt"
)

// Terminal colors for error reporting
const (
	RED    = "\033[31m"
	YELLOW = "\033[33m"
	RESET  = "\033[0m"
	LINE_UNKNOWN = -1
)

// Report generates an error message with line number and location information.
// Used for reporting syntax and runtime errors.
// Parameters:
//   - line: The line number where the error occurred
//   - where: Additional location information (e.g., token or expression)
//   - message: The error message describing the problem
func Report(line int, where string, message string) string {
	if where == "" {
		return fmt.Sprintf("%v[line %v]%v Error: %v\n", RED, line, RESET, message)
	}
	return fmt.Sprintf("%v[line %v]%v Error %v: %v\n", RED, line, RESET, where, message)
}

// ReportExit generates an error message and formats it for display before exit.
// Used for fatal errors that should terminate the program.
// Parameters:
//   - line: The line number where the error occurred
//   - where: Additional location information
//   - message: The error message
func ReportExit(line int, where string, message string) string {
	return Report(line, where, message)
}
