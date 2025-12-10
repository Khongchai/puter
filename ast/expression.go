package ast

import (
	"fmt"
	"strings"
)

type Expression interface {
	String() string
	Token() *Token
}

// any number in puter
type NumberExpression struct {
	ActualValue float64
	TokenValue  *Token
}

func (ne *NumberExpression) String() string {
	return fmt.Sprintf("%f", ne.ActualValue)
}

func (ne *NumberExpression) Token() *Token {
	return ne.TokenValue
}

type OperatorExpression struct {
	Left     Expression
	Operator *Token
	Right    Expression
}

func (oe *OperatorExpression) String() string {
	printed := fmt.Sprintf("(%s %s %s)", (oe.Left).String(), oe.Operator.Literal, (oe.Right).String())
	return printed
}

func (oe *OperatorExpression) Token() *Token {
	return oe.Operator
}

// Variable name
type IdentExpression struct {
	ActualValue string
	TokenValue  *Token
}

func (ne *IdentExpression) String() string {
	return ne.ActualValue
}

func (ne *IdentExpression) Token() *Token {
	return ne.TokenValue
}

type PrefixExpression struct {
	Right      Expression
	TokenValue *Token // the operator token
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("%s%s", pe.TokenValue.Literal, (pe.Right).String())
}

func (pe *PrefixExpression) Token() *Token {
	return pe.TokenValue
}

type AssignExpression struct {
	Name  Expression
	Right Expression
}

func (ae *AssignExpression) String() string {
	return fmt.Sprintf("%s = %s", ae.Name.String(), ae.Right.String())
}

func (pe *AssignExpression) Token() *Token {
	return pe.Name.Token()
}

type CallExpression struct {
	FunctionNameExpression Expression
	Args                   []Expression
}

func (ce *CallExpression) String() string {
	var names []string
	for _, a := range ce.Args {
		names = append(names, a.String())
	}
	joined := strings.Join(names, ", ")
	s := fmt.Sprintf("%s(%s)", ce.FunctionNameExpression.String(), joined)
	return s
}
func (ce *CallExpression) Token() *Token {
	return ce.FunctionNameExpression.Token()
}
