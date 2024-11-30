package main

import (
	"fmt"
)

const (
	LINE_UNKNOWN = -1
)

func ReportExit(line int, where, message string) string {
	if where == "" {
		return fmt.Sprintf("%v[line %v]%v Error: %v\n", RED, line, RESET, message)
	}
	return fmt.Sprintf("%v[line %v]%v Error %v: %v\n", RED, line, RESET, where, message)
}
