package parser

import (
	"fmt"
	ast "puter/ast"
)

type InfixParselet interface {
	Parse(parser *Parser, left ast.Expression, token *ast.Token) ast.Expression
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
		Name:  left,
		Right: right,
	}
}

type InParselet struct {
}

func NewInParselet() *InParselet {
	return &InParselet{}
}

func (cp *InParselet) Parse(parser *Parser, left ast.Expression, token *ast.Token) ast.Expression {
	// if right is not a static ident, like this: (...) in usd, (...) in binary
	// then it is not a valid expression
	right := parser.ParseExpression(PrecLowest)
	if rightExp, ok := right.(*ast.IdentExpression); !ok {
		panic(fmt.Sprintf("Expected an identifier, got %s", rightExp.String()))
	} else {
		return &ast.InExpression{
			TokenValue:  token,
			ActualValue: left,
			Unit:        rightExp,
		}
	}
}

type CallParselet struct {
}

func NewCallParselet() *CallParselet {
	return &CallParselet{}
}

func (cp *CallParselet) Parse(parser *Parser, left ast.Expression, token *ast.Token) ast.Expression {
	var args []ast.Expression

	// if next token is right paren, consume it and forward
	if parser.Peek(0).Type == ast.RPAREN {
		parser.Consume()
		return &ast.CallExpression{FunctionNameExpression: left, Args: args}
	}

	// otherwise loop and collect expressions delimited by a comma until right paren is encountered.
	for {
		args = append(args, parser.ParseExpression(PrecLowest))
		peeked := parser.Peek(0)
		if peeked.Type == ast.COMMA {
			parser.Consume()
			continue
		}
		consumed := parser.Consume()
		if consumed.Type == ast.RPAREN {
			break
		}
		panic("Missing right paren")
	}

	return &ast.CallExpression{FunctionNameExpression: left, Args: args}
}
