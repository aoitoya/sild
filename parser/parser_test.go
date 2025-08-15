package parser

import (
	"fmt"
	"testing"

	"github.com/toyaAoi/sild/scanner"
)

func TestParseProgram(t *testing.T) {
	tests := []struct {
		input         string
		expectedName  string
		expectedType  string
		expectedValue string
	}{
		{"let x: number = 42;", "x", "number", "42"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			p := New(scanner.New(tt.input))

			statements := p.ParseProgram().Statements

			if len(statements) == 0 {
				t.Fatalf("ParseProgram() returned no statements")
			}

			if statements[0] == nil {
				t.Fatalf("statements[0] is nil")
			}

			expected := fmt.Sprintf("name: %q, type: %q, value: %q",
				tt.expectedName, tt.expectedType, tt.expectedValue)
			got := statements[0].String()

			if expected != got {
				t.Fatalf("input: %s\nexpected: %s\ngot: %s",
					tt.input, expected, got)
			}

			fmt.Printf("✅ Test %d passed: %s\n", i, tt.input)
		})
	}
}
