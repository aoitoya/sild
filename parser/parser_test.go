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
