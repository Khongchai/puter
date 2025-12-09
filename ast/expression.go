package ast

import "fmt"

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
	printed := fmt.Sprintf("(%s%s%s)", (oe.Left).String(), oe.Operator.Literal, (oe.Right).String())
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
	TokenValue *Token // the name token
	Right      Expression
}

func (ae *AssignExpression) String() string {
	return fmt.Sprintf("%s = %s", ae.TokenValue.Literal, ae.Right.String())
}

func (pe *AssignExpression) Token() *Token {
	return pe.TokenValue
}
