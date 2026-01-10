package ast

type TokenType string

const (
	PIPE = "|"

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	PERCENT = "%"

	IDENT  = "IDENT"
	NUMBER = "NUMBER"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	CARET    = "^"

	LOGICAL_AND = "&&"
	LOGICAL_OR  = "||"

	LT  = "<"
	LTE = "<="
	GT  = ">"
	GTE = ">="

	EQ     = "=="
	NOT_EQ = "!="

	LPAREN = "("
	RPAREN = ")"

	COMMA = ","

	// Keywords
	TRUE  = "TRUE"
	FALSE = "FALSE"
	IN    = "IN"
)

// Puter scanner scans expression line by line so no need to store line information here.
type Token struct {
	Type     TokenType
	Literal  string
	StartPos int
}

func NewToken(tokenType TokenType, literal string, pos int) *Token {
	return &Token{
		Type:     tokenType,
		Literal:  literal,
		StartPos: pos,
	}
}

// Inclusive start, exclusive end.
// "cat" is start = 0 and end = 3
func (t *Token) EndPos() int {
	return t.StartPos + len(t.Literal)
}
