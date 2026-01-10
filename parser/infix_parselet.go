package parser

import (
	ast "puter/ast"
)

type InfixParselet interface {
	Parse(parser *Parser, left ast.Expression, token *ast.Token) (ast.Expression, *ast.Diagnostic)
	Precedence() int
}

type PostfixOperatorParselet struct {
	precedence int
}

func NewPostfixOperatorParselet(precedence int) *PostfixOperatorParselet {
	return &PostfixOperatorParselet{precedence: precedence}
}

func (p *PostfixOperatorParselet) Precedence() int {
	return p.precedence
}

func (b *PostfixOperatorParselet) Parse(parser *Parser, left ast.Expression, token *ast.Token) (ast.Expression, *ast.Diagnostic) {
	if parser == nil {
		panic("Parser nil, can't continue")
	}
	if left == nil {
		panic("Left is nil. Can't continue! Postfix requires left to be present")
	}
	if token == nil {
		return nil, ast.NewDiagnosticAtToken("Expected postfix here, got nothing", left.Token())
	}

	return &ast.PostfixExpression{
		Left:       left,
		TokenValue: token,
	}, nil
}

type BinaryOperatorParselet struct {
	precedence int
	isRight    bool
}

func (b *BinaryOperatorParselet) Precedence() int {
	return b.precedence
}

func NewbinaryOperatorParselet(precedence int, isRight bool) *BinaryOperatorParselet {
	return &BinaryOperatorParselet{
		precedence,
		isRight,
	}
}

func (b *BinaryOperatorParselet) Parse(parser *Parser, left ast.Expression, token *ast.Token) (ast.Expression, *ast.Diagnostic) {
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
	right, err := parser.parseExpression(p)
	if err != nil {
		return nil, err
	}

	return &ast.OperatorExpression{
		Left:     left,
		Operator: token,
		Right:    right,
	}, nil
}

type AssignParselet struct {
}

func NewAsssignParselet() *AssignParselet {
	return &AssignParselet{}
}

func (a *AssignParselet) Parse(parser *Parser, left ast.Expression, token *ast.Token) (ast.Expression, *ast.Diagnostic) {
	right, err := parser.parseExpression(PrecAssignment - 1)
	if err != nil {
		return nil, err
	}

	if _, ok := left.(*ast.IdentExpression); !ok {
		return nil, ast.NewDiagnosticAtToken("Left side of assign parselet not an ident expression", left.Token())
	}

	return &ast.AssignExpression{
		Name:  left,
		Right: right,
	}, nil
}

func (a *AssignParselet) Precedence() int {
	return PrecAssignment
}

type CallParselet struct {
}

func NewCallParselet() *CallParselet {
	return &CallParselet{}
}

func (cp *CallParselet) Parse(parser *Parser, left ast.Expression, token *ast.Token) (ast.Expression, *ast.Diagnostic) {
	var args []ast.Expression

	// if next token is right paren, consume it and forward
	if parser.Peek(0).Type == ast.RPAREN {
		parser.Consume()
		return &ast.CallExpression{FunctionNameExpression: left, Args: args}, nil
	}

	// otherwise loop and collect expressions delimited by a comma until right paren is encountered.
	for {
		exp, err := parser.parseExpression(PrecLowest)
		if err != nil {
			return nil, err
		}
		args = append(args, exp)
		peeked := parser.Peek(0)
		if peeked.Type == ast.COMMA {
			parser.Consume()
			continue
		}
		consumed := parser.Consume()
		if consumed.Type == ast.RPAREN {
			break
		}
		return nil, ast.NewDiagnosticAtToken("Missing closing paren", consumed)
	}

	return &ast.CallExpression{FunctionNameExpression: left, Args: args}, nil
}

func (c *CallParselet) Precedence() int {
	return PrecCall
}
