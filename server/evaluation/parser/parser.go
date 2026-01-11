package parser

import (
	"fmt"
	ast "puter/evaluation/ast"
	s "puter/evaluation/scanner"
)

type Parser struct {
	prefixParseFns map[ast.TokenType]PrefixParselet
	infixParseFns  map[ast.TokenType]InfixParselet
	scanner        *s.Scanner
	diagnostics    []*ast.Diagnostic
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
	// x %
	// x%
	// etc
	parser.infixParseFns[ast.IDENT] = NewPostfixOperatorParselet(PrecIn)
	parser.infixParseFns[ast.PERCENT] = NewPostfixOperatorParselet(PrecPercent)

	return parser
}

func (p *Parser) Parse(text string) (ast.Expression, *ast.Diagnostic) {
	p.scanner.SetState(0, text)
	return p.parseExpression(0)
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, *ast.Diagnostic) {
	token := p.Consume()

	nud, ok := p.prefixParseFns[token.Type]
	if !ok {
		diag := ast.NewDiagnostic(
			fmt.Sprintf("Unrecognized prefix token %s", token.Type),
			token.StartPos(),
			token.EndPos(),
		)
		return nil, diag
	}

	left, err := nud.Parse(p, token)
	if err != nil {
		return nil, err
	}

	for precedence < p.getNextPrecedence() {
		token = p.Consume()

		led := p.infixParseFns[token.Type]
		if led == nil {
			diag := ast.NewDiagnostic(
				fmt.Sprintf("Unrecognized infix token %s", token.Type),
				token.StartPos(),
				token.EndPos(),
			)
			return nil, diag
		}
		left, err = led.Parse(p, left, token)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
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

func (p *Parser) GetDiagnostics() {

}
