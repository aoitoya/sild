package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF   TokenType = "EOF"
	IDENT TokenType = "IDENT"

	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	COLON        TokenType = ":"
	SEMICOLON    TokenType = ";"
	PLUS         TokenType = "+"
	ASSIGN       TokenType = "="
	DOUBLE_QUOTE TokenType = `"`
	TRUE         TokenType = "true"
	FALSE        TokenType = "false"

	LET TokenType = "LET"

	TYPE_NUMBER  TokenType = "TYPE_NUMBER"
	TYPE_STRING  TokenType = "TYPE_STRING"
	TYPE_BOOLEAN TokenType = "TYPE_BOOLEAN"
)

var keywords = map[string]TokenType{
	"let":   LET,
	"true":  TRUE,
	"false": FALSE,
}

var types = map[string]TokenType{
	"number":  TYPE_NUMBER,
	"string":  TYPE_STRING,
	"boolean": TYPE_BOOLEAN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

func LookupType(t string) TokenType {
	if tok, ok := types[t]; ok {
		return tok
	}

	return IDENT
}
