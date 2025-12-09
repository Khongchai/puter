package parser

import ast "puter/ast"

type InfixParselet interface {
	Parse(parser *Parser, left ast.Expression, token *ast.Token) ast.Expression
	Precedence() int
}

type BinaryOperatorParselet struct {
	precedence int
	isRight    bool
}

func NewbinaryOperatorParselet(precedence int, isRight bool) *BinaryOperatorParselet {
	return &BinaryOperatorParselet{
		precedence,
		isRight,
	}
}

func (b *BinaryOperatorParselet) Precedence() int {
	return b.precedence
}

func (b *BinaryOperatorParselet) Parse(parser *Parser, left ast.Expression, token *ast.Token) ast.Expression {
	if parser == nil {
		panic("Parser nil, can't continue")
	}
	if left == nil {
		panic("Left is nil. Can't continue! Infix requires left to be present")
	}

	p := b.precedence
	if b.isRight {
		p -= 1
	}
	right := parser.ParseExpression(p)

	return &ast.OperatorExpression{
		Left:     left,
		Operator: token,
		Right:    right,
	}
}

type AssignParselet struct {
}

func NewAsssignParselet() *AssignParselet {
	return &AssignParselet{}
}

func (p *AssignParselet) Parse(parser *Parser, left ast.Expression, token *ast.Token) ast.Expression {
	right := parser.ParseExpression(PrecAssignment - 1)

	if _, ok := left.(*ast.IdentExpression); !ok {
		panic("Left side of assign parselet not an ident expression.")
	}

	return &ast.AssignExpression{
		TokenValue: token,
		Right:      right,
	}
}

func (p *AssignParselet) Precedence() int {
	return PrecAssignment
}
