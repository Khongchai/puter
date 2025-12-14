package parser

import (
	"fmt"
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
	expression := parser.parseExpression(0)
	consumed := parser.Consume()
	if consumed.Type != ast.RPAREN {
		panic("Closing is not right paren. Can't continue!")
	}
	return expression
}

type BooleanParselet struct {
}

func NewBooleanParselet() *BooleanParselet {
	return &BooleanParselet{}
}

func (p *BooleanParselet) Parse(parser *Parser, token *ast.Token) ast.Expression {
	if token.Literal != "true" && token.Literal != "false" {
		panic(fmt.Sprintf("Invalid boolean value. Expect true or false, got: %s", token.Literal))
	}

	return &ast.BooleanExpression{
		ActualValue: token.Literal == "true",
		TokenValue:  token,
	}
}

type NumberParselet struct {
}

func NewNumberParselet() *NumberParselet {
	return &NumberParselet{}
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

// Generic prefix operator parselets for stuff like +, -, /, *
type PrefixOperatorParselet struct {
	precedence int
}

func NewPrefixOperatorParselet(precedence int) *PrefixOperatorParselet {
	return &PrefixOperatorParselet{
		precedence,
	}
}

func (p *PrefixOperatorParselet) Parse(parser *Parser, token *ast.Token) ast.Expression {
	right := parser.parseExpression(p.precedence)
	return &ast.PrefixExpression{
		Right:      right,
		TokenValue: token,
	}
}
