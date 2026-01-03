package parser

import (
	"fmt"
	ast "puter/ast"
	s "puter/scanner"
)

type Parser struct {
	prefixParseFns map[ast.TokenType]PrefixParselet
	infixParseFns  map[ast.TokenType]InfixParselet
	scanner        *s.Scanner
}

func NewParser() *Parser {
	parser := &Parser{
		prefixParseFns: make(map[ast.TokenType]PrefixParselet),
		infixParseFns:  make(map[ast.TokenType]InfixParselet),
		scanner:        s.NewScanner(""),
	}

	// Special parselets
	parser.prefixParseFns[ast.IDENT] = NewIdentParselet()
	parser.prefixParseFns[ast.LPAREN] = NewGroupParselet()
	parser.prefixParseFns[ast.NUMBER] = NewNumberParselet()
	parser.prefixParseFns[ast.TRUE] = NewBooleanParselet()
	parser.prefixParseFns[ast.FALSE] = NewBooleanParselet()
	parser.infixParseFns[ast.ASSIGN] = NewAsssignParselet()
	parser.infixParseFns[ast.LPAREN] = NewCallParselet()

	// Simple parselets
	parser.prefixParseFns[ast.MINUS] = NewPrefixOperatorParselet(PrecPrefix)
	parser.prefixParseFns[ast.BANG] = NewPrefixOperatorParselet(PrecPrefix)
	parser.infixParseFns[ast.EQ] = NewbinaryOperatorParselet(PrecEquals, false)
	parser.infixParseFns[ast.NOT_EQ] = NewbinaryOperatorParselet(PrecEquals, false)
	parser.infixParseFns[ast.PLUS] = NewbinaryOperatorParselet(PrecSum, false)
	parser.infixParseFns[ast.IN] = NewbinaryOperatorParselet(PrecIn, false)
	parser.infixParseFns[ast.MINUS] = NewbinaryOperatorParselet(PrecSum, false)
	parser.infixParseFns[ast.ASTERISK] = NewbinaryOperatorParselet(PrecProduct, false)
	parser.infixParseFns[ast.SLASH] = NewbinaryOperatorParselet(PrecProduct, false)
	parser.infixParseFns[ast.CARET] = NewbinaryOperatorParselet(PrecExponent, true)
	parser.infixParseFns[ast.GT] = NewbinaryOperatorParselet(PrecLessGreater, false)
	parser.infixParseFns[ast.LT] = NewbinaryOperatorParselet(PrecLessGreater, false)
	parser.infixParseFns[ast.LTE] = NewbinaryOperatorParselet(PrecLessGreater, false)
	parser.infixParseFns[ast.GTE] = NewbinaryOperatorParselet(PrecLessGreater, false)
	parser.infixParseFns[ast.LOGICAL_AND] = NewbinaryOperatorParselet(PrecLogical, false)
	parser.infixParseFns[ast.LOGICAL_OR] = NewbinaryOperatorParselet(PrecLogical, false)
	// Any identifier after an expression is included here.
	// x usd
	// x percent
	// etc
	parser.infixParseFns[ast.IDENT] = NewPostfixOperatorParselet(PrecIn)

	return parser
}

func (p *Parser) Parse(text string) ast.Expression {
	p.scanner.SetState(0, text)
	return p.parseExpression(0)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	token := p.Consume()

	nud, ok := p.prefixParseFns[token.Type]
	if !ok {
		panic(fmt.Sprintf("Could not parse %s", token.Literal))
	}

	left := nud.Parse(p, token)

	for precedence < p.getNextPrecedence() {
		token = p.Consume()

		led := p.infixParseFns[token.Type]
		if led == nil {
			panic(fmt.Sprintf("Missing infix parselet for %s", token.Type))
		}
		left = led.Parse(p, left, token)
	}

	return left
}

func (p *Parser) Consume() *ast.Token {
	next := p.scanner.Next()
	return next
}

func (p *Parser) Peek(offset int) *ast.Token {
	peeked := p.scanner.Peek(offset)
	return peeked
}

func (p *Parser) getNextPrecedence() int {
	peeked := p.Peek(0)
	if res, ok := p.infixParseFns[peeked.Type]; ok {
		return res.Precedence()
	}
	return 0
}
