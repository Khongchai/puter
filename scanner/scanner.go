package scanner

import (
	ast "puter/ast"
)

type Scanner struct {
	pos  int
	text string
}

func NewScanner(text string) *Scanner {
	return &Scanner{
		pos:  0,
		text: text,
	}
}

func (s *Scanner) SetState(pos int, text string) {
	s.pos = 0
	s.text = text
}

func (s *Scanner) Next() *ast.Token {
	s.skipWhitespace()

	var token *ast.Token

	switch s.ch(0) {
	case '|':
		if s.ch(1) == '|' {
			token = ast.NewToken(ast.LOGICAL_OR, "||", s.pos)
			s.pos += 2
		} else {
			token = ast.NewToken(ast.PIPE, string(s.ch(0)), s.pos)
			s.pos++
		}
	case '&':
		if s.ch(1) == '&' {
			token = ast.NewToken(ast.LOGICAL_AND, "&&", s.pos)
			s.pos += 2
		} else {
			token = ast.NewToken(ast.ILLEGAL, string(s.ch(0)), s.pos)
			s.pos++
		}
	case '=':
		if s.ch(1) == '=' {
			token = ast.NewToken(ast.EQ, "==", s.pos)
			s.pos += 2
		} else {
			token = ast.NewToken(ast.ASSIGN, string(s.ch(0)), s.pos)
			s.pos++
		}
	case '+':
		token = ast.NewToken(ast.PLUS, string(s.ch(0)), s.pos)
		s.pos++
	case '-':
		token = ast.NewToken(ast.MINUS, string(s.ch(0)), s.pos)
		s.pos++
	case '!':
		if s.ch(1) == '=' {
			token = ast.NewToken(ast.NOT_EQ, "!=", s.pos)
			s.pos += 2
		} else {
			token = ast.NewToken(ast.BANG, string(s.ch(0)), s.pos)
			s.pos++
		}
	case ',':
		token = ast.NewToken(ast.COMMA, string(s.ch(0)), s.pos)
		s.pos++
	case '/':
		token = ast.NewToken(ast.SLASH, string(s.ch(0)), s.pos)
		s.pos++
	case '*':
		token = ast.NewToken(ast.ASTERISK, string(s.ch(0)), s.pos)
		s.pos++
	case '<':
		token = ast.NewToken(ast.LT, string(s.ch(0)), s.pos)
		s.pos++
	case '>':
		token = ast.NewToken(ast.GT, string(s.ch(0)), s.pos)
		s.pos++
	case '(':
		token = ast.NewToken(ast.LPAREN, string(s.ch(0)), s.pos)
		s.pos++
	case ')':
		token = ast.NewToken(ast.RPAREN, string(s.ch(0)), s.pos)
		s.pos++
	case 0:
		token = ast.NewToken(ast.EOF, "", s.pos)
		s.pos++
	default:
		if isLetter(s.ch(0)) {
			i := 1
			for isLetter(s.ch(i)) {
				i++
			}
			text := s.text[s.pos : s.pos+i]
			tokenType := func() ast.TokenType {
				switch text {
				case "true":
					return ast.TRUE
				case "false":
					return ast.FALSE
				case "in":
					return ast.IN
				default:
					return ast.IDENT
				}
			}()
			token = ast.NewToken(tokenType, text, s.pos)
			s.pos += i
		} else if isDigit(s.ch(0)) {
			i := 1
			for {
				if isDigit(s.ch(i)) || (s.ch(i) == '.' && isDigit(s.ch(i+1))) {
					i++
					continue
				}
				break
			}
			token = ast.NewToken(ast.NUMBER, s.text[s.pos:s.pos+i], s.pos)
			s.pos += i
		} else {
			token = ast.NewToken(ast.ILLEGAL, string(s.ch(0)), s.pos)
			s.pos++
		}
	}

	return token
}

func (s *Scanner) Peek(offset int) *ast.Token {
	prevPos := s.pos

	var tok *ast.Token
	for range max(0, offset-1) { // avoid unnecessary assignment
		s.Next()
	}
	tok = s.Next()

	s.pos = prevPos

	return tok
}

func (s *Scanner) skipWhitespace() {
	for {
		if s.pos >= len(s.text) {
			return
		}
		ch := s.text[s.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			s.pos++
			continue
		}
		return
	}
}

func (s *Scanner) ch(offset int) byte {
	if s.pos+offset >= len(s.text) {
		return 0
	}
	return s.text[s.pos+offset]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
