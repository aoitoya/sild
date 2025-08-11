package scanner

import (
	"fmt"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let x: number = 42;`

	tests := []struct {
		expectedType    string
		expectedLiteral string
	}{
		{LET, "let"},
		{IDENT, "x"},
		{COLON, ":"},
		{IDENT, "number"},
		{EQUAL, "="},
		{INT, "42"},
		{SEMICOLON, ";"},
	}

	sc := New(input)

	for i, tt := range tests {
		tok := sc.NextToken()

		fmt.Println(tok)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tokeLiteral wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
