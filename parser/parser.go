package parser

import (
	"fmt"
	ast "puter/ast"
	s "puter/scanner"
)

var precedences = map[ast.TokenType]int{
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

	parser.prefixParseFns[ast.IDENT] = NewIdentParselet()
	parser.prefixParseFns[ast.LPAREN] = NewGroupParselet()

	parser.infixParseFns[ast.ASSIGN] = NewAsssignParselet()

	for _, token := range []ast.Token{
		ast.PLUS,
		ast.MINUS,
	} {
	}

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

	for precedence < getPrecedence(token.Type) {
		token = p.Consume()

		infix := p.infixParseFns[token.Type]
		left = infix.Parse(p, left, token)
	}

	return left
}

func (p *Parser) Consume() *ast.Token {
	next := p.scanner.Next(p.line)
	return next
}

func getPrecedence(tokenType ast.TokenType) int {
	res, ok := precedences[tokenType]
	if !ok {
		panic(fmt.Sprintf("Precedence does not support token type %s", tokenType))
	}
	return res
}
