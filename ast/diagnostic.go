package ast

type Diagnostic struct {
	Message  string
	StartPos int
	EndPos   int
}

func NewDiagnostic(message string, startPos int, endPos int) *Diagnostic {
	return &Diagnostic{message, startPos, endPos}
}

func NewDiagnosticAtToken(message string, token *Token) *Diagnostic {
	return &Diagnostic{message, token.StartPos(), token.EndPos()}
}
