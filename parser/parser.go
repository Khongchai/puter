package parser

import (
	"fmt"
	ast "puter/ast"
	s "puter/scanner"
)

var precedences = map[ast.TokenType]int{
	ast.ASSIGN:   PrecAssignment,
	ast.EQ:       PrecEquals,
	ast.NOT_EQ:   PrecEquals,
	ast.LT:       PrecLessGreater,
	ast.GT:       PrecLessGreater,
	ast.PLUS:     PrecSum,
	ast.MINUS:    PrecSum,
	ast.SLASH:    PrecProduct,
	ast.ASTERISK: PrecProduct,
	ast.LPAREN:   PrecCall,
}

type Parser struct {
	prefixParseFns map[ast.TokenType]PrefixParselet
	infixParseFns  map[ast.TokenType]InfixParselet
	scanner        *s.Scanner
	line           int
}

func NewParser(text string) *Parser {
	parser := &Parser{
		prefixParseFns: make(map[ast.TokenType]PrefixParselet),
		infixParseFns:  make(map[ast.TokenType]InfixParselet),
		scanner:        s.NewScanner(text),
		line:           0,
	}

	// Special parselets
	parser.prefixParseFns[ast.IDENT] = NewIdentParselet()
	parser.prefixParseFns[ast.LPAREN] = NewGroupParselet()
	parser.prefixParseFns[ast.NUMBER] = NewNumberParselet()
	parser.infixParseFns[ast.ASSIGN] = NewAsssignParselet()
	parser.infixParseFns[ast.LPAREN] = NewCallParselet()

	// Simple parselets
	parser.prefixParseFns[ast.MINUS] = NewPrefixOperatorParselet(PrecPrefix)
	parser.prefixParseFns[ast.BANG] = NewPrefixOperatorParselet(PrecPrefix)
	parser.infixParseFns[ast.PLUS] = NewbinaryOperatorParselet(PrecSum, false)
	parser.infixParseFns[ast.MINUS] = NewbinaryOperatorParselet(PrecSum, false)
	parser.infixParseFns[ast.ASTERISK] = NewbinaryOperatorParselet(PrecProduct, false)
	parser.infixParseFns[ast.SLASH] = NewbinaryOperatorParselet(PrecProduct, false)
	parser.infixParseFns[ast.CARET] = NewbinaryOperatorParselet(PrecExponent, true)

	return parser
}

// Only parses math expression for now.
func (p *Parser) Parse() ast.Expression {
	return p.ParseExpression(0)
}

// new line not handled yet.
func (p *Parser) ParseExpression(precedence int) ast.Expression {
	token := p.Consume()

	prefixParselet, ok := p.prefixParseFns[token.Type]
	if !ok {
		panic(fmt.Sprintf("Could not parse %s", token.Literal))
	}

	left := prefixParselet.Parse(p, token)

	for precedence < p.getNextPrecedence() {
		token = p.Consume()

		infix := p.infixParseFns[token.Type]
		left = infix.Parse(p, left, token)
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
	if res, ok := precedences[peeked.Type]; ok {
		return res
	}
	return 0
}
