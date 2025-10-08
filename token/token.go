package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF   TokenType = "EOF"
	IDENT TokenType = "IDENT"

	NUMBER  TokenType = "NUMBER"
	STRING  TokenType = "STRING"
	BOOLEAN TokenType = "BOOLEAN"

	COMMA        TokenType = ","
	COLON        TokenType = ":"
	SEMICOLON    TokenType = ";"
	PLUS         TokenType = "+"
	MINUS        TokenType = "-"
	MUL          TokenType = "*"
	DIV          TokenType = "/"
	BANG         TokenType = "!"
	LEFT_PAREN   TokenType = "("
	RIGHT_PAREN  TokenType = ")"
	LEFT_BRACE   TokenType = "{"
	RIGHT_BRACE  TokenType = "}"
	ASSIGN       TokenType = "="
	DOUBLE_QUOTE TokenType = `"`

	LET      TokenType = "LET"
	FUNCTION TokenType = "FUNCTION"
	RETURN   TokenType = "RETURN"

	TYPE_NUMBER  TokenType = "TYPE_NUMBER"
	TYPE_STRING  TokenType = "TYPE_STRING"
	TYPE_BOOLEAN TokenType = "TYPE_BOOLEAN"
	TYPE_VOID    TokenType = "TYPE_VOID"
)

var keywords = map[string]TokenType{
	"let":      LET,
	"true":     BOOLEAN,
	"false":    BOOLEAN,
	"function": FUNCTION,
	"return":   RETURN,
}

var types = map[string]TokenType{
	"number":  TYPE_NUMBER,
	"string":  TYPE_STRING,
	"boolean": TYPE_BOOLEAN,
	"void":    TYPE_VOID,
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
