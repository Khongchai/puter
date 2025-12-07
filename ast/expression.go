package ast

import "fmt"

type Expression interface {
	String() string
	Token() Token
}

type NumberExpression struct {
	ActualValue float64
	TokenValue  Token
}

func (ne *NumberExpression) String() string {
	return fmt.Sprintf("%f", ne.ActualValue)
}

func (ne *NumberExpression) Token() Token {
	return ne.TokenValue
}

type PrefixExpression struct {
	Right      Expression
	TokenValue Token // the operator token
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("%s%s", pe.TokenValue.Literal, pe.Right.String())
}

func (pe *PrefixExpression) Token() Token {
	return pe.TokenValue
}
