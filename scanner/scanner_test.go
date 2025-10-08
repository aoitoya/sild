package scanner

import (
	"testing"

	"github.com/toyaAoi/sild/token"
)

func TestNextToken(t *testing.T) {
	input := `let x: number = 42;
            let name: string = "hello";
            let isActive: boolean = true;
            let count: number = 0;
            let message: string = "world";`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// let x: number = 42;
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.TYPE_NUMBER, "number"},
		{token.ASSIGN, "="},
		{token.NUMBER, "42"},
		{token.SEMICOLON, ";"},

		// let name: string = "hello";
		{token.LET, "let"},
		{token.IDENT, "name"},
		{token.COLON, ":"},
		{token.TYPE_STRING, "string"},
		{token.ASSIGN, "="},
		{token.STRING, `hello`},
		{token.SEMICOLON, ";"},

		// let isActive: boolean = true;
		{token.LET, "let"},
		{token.IDENT, "isActive"},
		{token.COLON, ":"},
		{token.TYPE_BOOLEAN, "boolean"},
		{token.ASSIGN, "="},
		{token.BOOLEAN, "true"},
		{token.SEMICOLON, ";"},

		// let count: number = 0;
		{token.LET, "let"},
		{token.IDENT, "count"},
		{token.COLON, ":"},
		{token.TYPE_NUMBER, "number"},
		{token.ASSIGN, "="},
		{token.NUMBER, "0"},
		{token.SEMICOLON, ";"},

		// let message: string = "world";
		{token.LET, "let"},
		{token.IDENT, "message"},
		{token.COLON, ":"},
		{token.TYPE_STRING, "string"},
		{token.ASSIGN, "="},
		{token.STRING, `world`},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}

	sc := New(input)

	for i, tt := range tests {
		tok := sc.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tokenLiteral wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenSingleLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			tokenType token.TokenType
			literal   string
		}
	}{
		{
			name:  "simple variable declaration",
			input: "let x: number = 42;",
			expected: []struct {
				tokenType token.TokenType
				literal   string
			}{
				{token.LET, "let"},
				{token.IDENT, "x"},
				{token.COLON, ":"},
				{token.TYPE_NUMBER, "number"},
				{token.ASSIGN, "="},
				{token.NUMBER, "42"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			name:  "string declaration",
			input: `let greeting: string = "hello world";`,
			expected: []struct {
				tokenType token.TokenType
				literal   string
			}{
				{token.LET, "let"},
				{token.IDENT, "greeting"},
				{token.COLON, ":"},
				{token.TYPE_STRING, "string"},
				{token.ASSIGN, "="},
				{token.STRING, "hello world"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			name:  "boolean declaration",
			input: "let flag: boolean = false;",
			expected: []struct {
				tokenType token.TokenType
				literal   string
			}{
				{token.LET, "let"},
				{token.IDENT, "flag"},
				{token.COLON, ":"},
				{token.TYPE_BOOLEAN, "boolean"},
				{token.ASSIGN, "="},
				{token.BOOLEAN, "false"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := New(tt.input)

			for i, expected := range tt.expected {
				tok := sc.NextToken()

				if tok.Type != expected.tokenType {
					t.Fatalf("test %s[%d] - tokenType wrong. expected=%q, got=%q",
						tt.name, i, expected.tokenType, tok.Type)
				}

				if tok.Literal != expected.literal {
					t.Fatalf("test %s[%d] - tokenLiteral wrong. expected=%q, got=%q",
						tt.name, i, expected.literal, tok.Literal)
				}
			}
		})
	}
}
