package parser

import (
	ast "puter/ast"
	s "puter/scanner"
)

type Parser struct {
	Text           string
	prefixParseFns map[ast.TokenType]PrefixParselet
	infixParseFns  map[ast.TokenType]InfixParselet
}

// Only parses math expression for now.
func (p *Parser) Parse() ast.Expression {
	return p.ParseExpression(0)
}

// new line not handled yet.
func (p *Parser) ParseExpression(precedence int) ast.Expression {

	// TODO
	for token := range s.Scan(p.Text, 0) {
		switch token.Type {
		case ast.EOF:
			break
		case ast.PLUS:
			break
		case ast.MINUS:
			break
		case ast.ASSIGN:
			break
		case ast.SLASH:
			break
		case ast.NUMBER:
			break
		default:
			panic("token not supported yet!")
		}
	}
}
