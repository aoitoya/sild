package codegen

import (
	"strings"
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

func TestMultipleVariableDeclarations(t *testing.T) {
	program := createProgram(
		&VariableDeclaration{Name: "x", Type: "number", Value: "42"},
		&VariableDeclaration{Name: "name", Type: "string", Value: "hello"},
		&VariableDeclaration{Name: "isActive", Type: "boolean", Value: "true"},
	)

	expected := `package main

func main() {
    x := 42
    name := "hello"
    isActive := true
}
`

	g := New()
	result := g.Generate(program)

	if result != expected {
		t.Errorf("Multiple variables test failed\nExpected:\n%s\nGot:\n%s",
			expected, result)
	}
}

func TestEmptyProgram(t *testing.T) {
	program := &Program{Statements: []Statement{}}

	expected := `package main

func main() {
}
`

	g := New()
	result := g.Generate(program)

	if result != expected {
		t.Errorf("Empty program test failed\nExpected:\n%s\nGot:\n%s",
			expected, result)
	}
}

func TestVariableNaming(t *testing.T) {
	tests := []struct {
		name         string
		variableName string
		shouldPass   bool
	}{
		{"simple name", "x", true},
		{"camelCase", "firstName", true},
		{"with numbers", "var1", true},
		{"underscore", "my_var", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := createProgram(&VariableDeclaration{
				Name:  tt.variableName,
				Type:  "number",
				Value: "42",
			})

			g := New()
			result := g.Generate(program)

			if !strings.Contains(result, tt.variableName+" := 42") {
				if tt.shouldPass {
					t.Errorf("Expected variable name %s to be in output", tt.variableName)
				}
			}
		})
	}
}

func TestErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		input   *Program
		wantErr bool
	}{
		{
			name:    "nil program",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()

			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("Generate() panicked unexpectedly: %v", r)
				}
			}()

			result := g.Generate(tt.input)

			if tt.input == nil && result != "" {
				t.Errorf("Expected empty result for nil input, got: %s", result)
			}
		})
	}
}

func TestIndentation(t *testing.T) {
	program := createProgram(
		&VariableDeclaration{Name: "x", Type: "number", Value: "42"},
	)

	g := New()
	result := g.Generate(program)

	lines := strings.SplitSeq(result, "\n")

	for line := range lines {
		if strings.Contains(line, "x := 42") {
			if !strings.HasPrefix(line, "    ") {
				t.Errorf("Expected 4-space indentation, got: '%s'", line)
			}
		}
	}
}
