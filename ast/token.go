package ast

type TokenType string

const (
	PIPE = "|"

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	NUMBER = "NUMBER"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	LPAREN = "("
	RPAREN = ")"

	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

// Puter scanner scans expression line by line so no need to store line information here.
type Token struct {
	Type     TokenType
	Literal  string
	StartPos int
}

// Inclusive start, exclusive end.
// "cat" is start = 0 and end = 3
func (t *Token) EndPos() int {
	return t.StartPos + len(t.Literal)
}
