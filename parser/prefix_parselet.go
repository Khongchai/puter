package parser

import (
	ast "puter/ast"
	"strconv"
)

type PrefixParselet interface {
	Parse(parser *Parser, token ast.Token) ast.Expression
}

type NumberParselet struct {
}

func (p *NumberParselet) Parse(parser *Parser, token ast.Token) ast.Expression {
	parsed, ok := strconv.ParseFloat(token.Literal, 64)
	if ok != nil {
		panic("Error") // TODO error handling coming soon...
	}
	return &ast.NumberExpression{
		ActualValue: parsed,
		TokenValue:  token,
	}

}

type PrefixOperatorParselet struct {
}

func (p *PrefixOperatorParselet) Parse(parser *Parser, token ast.Token) ast.Expression {
	parsed := parser.Parse()
	return &ast.PrefixExpression{
		Right:      parsed,
		TokenValue: token,
	}

}
