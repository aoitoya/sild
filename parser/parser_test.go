package parser

import (
	"fmt"
	"testing"

	"github.com/toyaAoi/sild/scanner"
)

func TestParseProgram(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedName  string
		expectedType  string
		expectedValue string
		expectError   bool
	}{
		{
			name:          "number assignment",
			input:         "let x: number = 42;",
			expectedName:  "x",
			expectedType:  "number",
			expectedValue: "42",
		},
		{
			name:          "string assignment",
			input:         `let name: string = "hello";`,
			expectedName:  "name",
			expectedType:  "string",
			expectedValue: "hello",
		},
		{
			name:          "boolean true assignment",
			input:         "let isActive: boolean = true;",
			expectedName:  "isActive",
			expectedType:  "boolean",
			expectedValue: "true",
		},
		{
			name:          "boolean false assignment",
			input:         "let isDone: boolean = false;",
			expectedName:  "isDone",
			expectedType:  "boolean",
			expectedValue: "false",
		},
		{
			name:          "zero value number",
			input:         "let count: number = 0;",
			expectedName:  "count",
			expectedType:  "number",
			expectedValue: "0",
		},
		{
			name:          "empty string",
			input:         `let message: string = "";`,
			expectedName:  "message",
			expectedType:  "string",
			expectedValue: "",
		},
		{
			name:          "negative number",
			input:         "let temperature: number = -5;",
			expectedName:  "temperature",
			expectedType:  "number",
			expectedValue: "-5",
		},
		{
			name:        "missing semicolon",
			input:       "let x: number = 42",
			expectError: true,
		},
		{
			name:        "missing type annotation",
			input:       "let x = 42;",
			expectError: true,
		},
		{
			name:        "missing value",
			input:       "let x: number = ;",
			expectError: true,
		},
		{
			name:        "invalid type",
			input:       "let x: invalid = 42;",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(scanner.New(tt.input))

			program := p.ParseProgram()

			if tt.expectError {
				if len(program.Statements) > 0 && program.Statements[0] != nil {
					t.Fatalf("expected error but got valid statement: %s", program.Statements[0].String())
				}
				return
			}

			if len(program.Statements) == 0 {
				t.Fatalf("ParseProgram() returned no statements")
			}

			if program.Statements[0] == nil {
				t.Fatalf("statements[0] is nil")
			}

			expected := fmt.Sprintf("name: %q, type: %q, value: %q",
				tt.expectedName, tt.expectedType, tt.expectedValue)
			got := program.Statements[0].String()

			if expected != got {
				t.Fatalf("input: %s\nexpected: %s\ngot: %s",
					tt.input, expected, got)
			}

			t.Logf("✅ Test passed: %s", tt.name)
		})
	}
}

func TestMultipleStatements(t *testing.T) {
	input := `
			let x: number = 42;
			let name: string = "test";
			let flag: boolean = true;
		`

	p := New(scanner.New(input))
	program := p.ParseProgram()

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedName  string
		expectedType  string
		expectedValue string
	}{
		{"x", "number", "42"},
		{"name", "string", "test"},
		{"flag", "boolean", "true"},
	}

	for i, tt := range tests {
		expected := fmt.Sprintf("name: %q, type: %q, value: %q",
			tt.expectedName, tt.expectedType, tt.expectedValue)
		got := program.Statements[i].String()

		if expected != got {
			t.Errorf("statement %d: expected %s, got %s", i, expected, got)
		}
	}
}

func TestWhitespaceHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"single space", "let x: number = 42;"},
		{"multiple spaces", "let  x  :  number  =  42  ;"},
		{"tabs", "let\tx\t:\tnumber\t=\t42\t;"},
		{"newlines", "let\nx\n:\nnumber\n=\n42\n;"},
		{"mixed whitespace", "let\t x\n: \t number \t =\n 42 \t ;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(scanner.New(tt.input))
			program := p.ParseProgram()

			if len(program.Statements) == 0 {
				t.Fatal("no statements parsed")
			}

			got := program.Statements[0].String()
			expected := "name: \"x\", type: \"number\", value: \"42\""

			if got != expected {
				t.Errorf("expected %q, got %q", expected, got)
			}
		})
	}
}

func TestExpressionParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"simple addition",
			"let result: number = 1 + 2;",
			"name: \"result\", type: \"number\", value: \"(1 + 2)\"",
		},
		{
			"operator precedence",
			"let result: number = 1 + 2 * 3;",
			"name: \"result\", type: \"number\", value: \"(1 + (2 * 3))\"",
		},
		{
			"parentheses",
			"let result: number = (1 + 2) * 3;",
			"name: \"result\", type: \"number\", value: \"((1 + 2) * 3)\"",
		},
		{
			"division",
			"let result: number = 10 / 2;",
			"name: \"result\", type: \"number\", value: \"(10 / 2)\"",
		},
		{
			"mixed operations",
			"let result: number = 1 + 2 * 3 / 2 - 1;",
			"name: \"result\", type: \"number\", value: \"((1 + ((2 * 3) / 2)) - 1)\"",
		},
		{
			"unary minus",
			"let result: number = -5 + 10;",
			"name: \"result\", type: \"number\", value: \"(-5 + 10)\"",
		},
		{
			"nested parentheses",
			"let result: number = ((1 + 2) * (3 - 4)) / 5;",
			"name: \"result\", type: \"number\", value: \"(((1 + 2) * (3 - 4)) / 5)\"",
		},
		{
			"complex expression",
			"let result: number = (10 + (2 * 3) - 8) / 2 * 5;",
			"name: \"result\", type: \"number\", value: \"((((10 + (2 * 3)) - 8) / 2) * 5)\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(scanner.New(tt.input))
			program := p.ParseProgram()

			if len(program.Statements) == 0 {
				t.Fatal("no statements parsed")
			}

			got := program.Statements[0].String()

			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
