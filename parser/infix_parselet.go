package parser

import ast "puter/ast"

type InfixParselet interface {
	Parse(parser *Parser, token ast.Token) ast.Expression
	Precedence() int
}
