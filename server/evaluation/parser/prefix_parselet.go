package parser

import (
	"fmt"
	ast "puter/evaluation/ast"
	"strconv"
)

type PrefixParselet interface {
	Parse(parser *Parser, token *ast.Token) (ast.Expression, *ast.Diagnostic)
}

type GroupParselet struct {
}

func NewGroupParselet() *GroupParselet {
	return &GroupParselet{}
}

func (p *GroupParselet) Parse(parser *Parser, token *ast.Token) (ast.Expression, *ast.Diagnostic) {
	expression, err := parser.parseExpression(0)
	if err != nil {
		return nil, err
	}
	consumed := parser.Consume()
	if consumed.Type != ast.RPAREN {
		diag := ast.NewDiagnosticAtToken(fmt.Sprintf("Expected right paren, got: %s", consumed.Type), consumed)
		return nil, diag
	}
	return expression, nil
}

type BooleanParselet struct {
}

func NewBooleanParselet() *BooleanParselet {
	return &BooleanParselet{}
}

func (p *BooleanParselet) Parse(parser *Parser, token *ast.Token) (ast.Expression, *ast.Diagnostic) {
	if token.Literal != "true" && token.Literal != "false" {
		return nil, ast.NewDiagnosticAtToken(
			fmt.Sprintf("Invalid boolean value. Expect true or false, got: %s", token.Literal),
			token,
		)
	}

	return &ast.BooleanExpression{
		ActualValue: token.Literal == "true",
		TokenValue:  token,
	}, nil
}

type NumberParselet struct {
}

func NewNumberParselet() *NumberParselet {
	return &NumberParselet{}
}

func (p *NumberParselet) Parse(parser *Parser, token *ast.Token) (ast.Expression, *ast.Diagnostic) {
	parsed, ok := strconv.ParseFloat(token.Literal, 64)
	if ok != nil {
		return nil, ast.NewDiagnosticAtToken(fmt.Sprintf("Invalid number: %s", token.Literal), token)
	}
	return &ast.NumberExpression{
		ActualValue: parsed,
		TokenValue:  token,
	}, nil
}

type IdentParselet struct {
}

func NewIdentParselet() *IdentParselet {
	return &IdentParselet{}
}

func (p *IdentParselet) Parse(parser *Parser, token *ast.Token) (ast.Expression, *ast.Diagnostic) {
	return &ast.IdentExpression{
		ActualValue: token.Literal,
		TokenValue:  token,
	}, nil

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

func (p *PrefixOperatorParselet) Parse(parser *Parser, token *ast.Token) (ast.Expression, *ast.Diagnostic) {
	right, err := parser.parseExpression(p.precedence)
	if err != nil {
		return nil, err
	}
	return &ast.PrefixExpression{
		Right:      right,
		TokenValue: token,
	}, nil
}
