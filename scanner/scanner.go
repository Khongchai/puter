package scanner

import (
	"iter"
	ast "puter/ast"
)

func Scan(text string, line int) iter.Seq[ast.Token] {
	pos := 0

	newToken := func(tokenType ast.TokenType, literal string) *ast.Token {
		return &ast.Token{
			Type:     tokenType,
			Literal:  literal,
			StartPos: pos,
		}
	}

	skipWhitespace := func() {
		ch := text[pos]
		for ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			pos++
		}
	}

	ch := func(offset int) byte {
		if pos+offset >= len(text) {
			return 0
		} else {
			return text[pos+offset]
		}
	}

	yielder := func(yield func(ast.Token) bool) {
		skipWhitespace()

		var token *ast.Token

		for {
			switch ch(0) {
			case '|':
				token = newToken(ast.PIPE, string(ch(0)))
				pos++
			case '=':
				if ch(1) == '=' {
					token = newToken(ast.EQ, "==")
					pos += 2
				} else {
					token = newToken(ast.ASSIGN, string(ch(0)))
					pos++
				}
			case '+':
				token = newToken(ast.PLUS, string(ch(0)))
				pos++
			case '-':
				token = newToken(ast.MINUS, string(ch(0)))
				pos++
			case '!':
				if ch(0) == '=' {
					token = newToken(ast.NOT_EQ, "!=")
					pos += 2
				} else {
					token = newToken(ast.BANG, string(ch(0)))
					pos++
				}
			case '/':
				token = newToken(ast.SLASH, string(ch(0)))
				pos++
			case '*':
				token = newToken(ast.ASTERISK, string(ch(0)))
				pos++
			case '<':
				token = newToken(ast.LT, string(ch(0)))
				pos++
			case '>':
				token = newToken(ast.GT, string(ch(0)))
				pos++
			case '(':
				token = newToken(ast.LPAREN, string(ch(0)))
				pos++
			case ')':
				token = newToken(ast.RPAREN, string(ch(0)))
				pos++
			case 0:
				token = newToken(ast.EOF, "")
				pos++
			default:
				if isLetter(ch(0)) {
					i := 1
					for isLetter(ch(i)) {
						i++
					}
					token = newToken(ast.IDENT, text[pos:pos+i])
					pos += i
				} else if isDigit(ch(0)) {
					i := 1
					for isDigit(ch(i)) {
						i++
					}
					token = newToken(ast.NUMBER, text[pos:pos+i])
					pos += i
				} else {
					token = newToken(ast.ILLEGAL, string(ch(0)))
					pos++
				}
			}

			if !yield(*token) {
				return
			}
		}
	}

	return yielder
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
