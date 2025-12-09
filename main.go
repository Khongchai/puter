package main

import (
	"fmt"
	p "puter/parser"
)

func main() {
	parser := &p.Parser{
		Text: "3 * 2 + 1",
	}
	expression := parser.ParseExpression(0)
	fmt.Println((*expression).String())

}
