package scanner

type Token struct {
	Type    string
	Literal string
}

const (
	IDENT = "IDENT"
	INT   = "INT"

	COLON     = ":"
	SEMICOLON = ";"
	PLUS      = "+"
	EQUAL     = "="

	LET = "LET"
)

var keywords = map[string]string{
	"let": LET,
}

func lookupIdent(ident string) string {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
