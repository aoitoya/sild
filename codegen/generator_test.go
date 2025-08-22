// codegen/generator_test.go
package codegen

import (
	"testing"

	"github.com/toyaAoi/sild/ast"
)

type VariableDeclaration = ast.VariableDeclaration

func createProgram(statements ...Statement) *Program {
	return &Program{Statements: statements}
}

func TestBasicVariableDeclaration(t *testing.T) {
	tests := []struct {
		name     string
		input    *Program
		expected string
	}{
		{
			name: "number variable",
			input: createProgram(&VariableDeclaration{
				Name:  "x",
				Type:  "number",
				Value: "42",
			}),
			expected: `package main

func main() {
    x := 42
}
`,
		},
		{
			name: "string variable",
			input: createProgram(&VariableDeclaration{
				Name:  "name",
				Type:  "string",
				Value: "hello",
			}),
			expected: `package main

func main() {
    name := "hello"
}
`,
		},
		{
			name: "boolean variable",
			input: createProgram(&VariableDeclaration{
				Name:  "isActive",
				Type:  "boolean",
				Value: "true",
			}),
			expected: `package main

func main() {
    isActive := true
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			result := g.Generate(tt.input)

			if result != tt.expected {
				t.Errorf("Test %s failed\nExpected:\n%s\nGot:\n%s",
					tt.name, tt.expected, result)
			}
		})
	}
}
