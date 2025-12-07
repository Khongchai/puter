package main

import (
	"fmt"
	"puter/ast"
	s "puter/scanner"
)

// https://engineering.desmos.com/articles/pratt-parser/
func main() {
	line := 0
	for tok := range s.Scan("=!02938#*Hello", line) {
		if tok.Type == ast.EOF {
			return
		}
		fmt.Printf("%+v \n", tok.Literal)
	}
}
