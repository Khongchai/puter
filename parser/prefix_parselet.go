package parser

import (
	ast "puter/ast"
	"strconv"
)

type PrefixParselet interface {
	Parse(parser *Parser, token *ast.Token) ast.Expression
}

type GroupParselet struct {
}

func NewGroupParselet() *GroupParselet {
	return &GroupParselet{}
}

func (p *GroupParselet) Parse(parser *Parser, token *ast.Token) ast.Expression {
	expression := parser.ParseExpression(0)
	consumed := parser.Consume()
	if consumed.Type != ast.RPAREN {
		panic("Closing is not right paren. Can't continue!")
	}
	return expression
}

type NumberParselet struct {
}

func (p *NumberParselet) Parse(parser *Parser, token *ast.Token) ast.Expression {
	parsed, ok := strconv.ParseFloat(token.Literal, 64)
	if ok != nil {
		panic("Error") // TODO error handling coming soon...
	}
	return &ast.NumberExpression{
		ActualValue: parsed,
		TokenValue:  token,
	}
}

type IdentParselet struct {
}

func NewIdentParselet() *IdentParselet {
	return &IdentParselet{}
}

func (p *IdentParselet) Parse(parser *Parser, token *ast.Token) ast.Expression {
	return &ast.IdentExpression{
		ActualValue: token.Literal,
		TokenValue:  token,
	}

}

type PrefixOperatorParselet struct {
	precedence int
}

func (p *PrefixOperatorParselet) Parse(parser *Parser, token *ast.Token) ast.Expression {
	right := parser.ParseExpression(p.precedence)
	return &ast.PrefixExpression{
		Right:      right,
		TokenValue: token,
	}
}
