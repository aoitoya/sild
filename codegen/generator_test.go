package codegen

import (
	"testing"

	"github.com/toyaAoi/sild/ast"
	"github.com/toyaAoi/sild/token"
)

type VariableDeclaration = ast.VariableDeclaration

type testCase struct {
	name     string
	program  *Program
	expected string
}

func createProgram(statements ...Statement) *Program {
	return &Program{Statements: statements}
}

func createVariableDeclaration(name, typ, value string) *VariableDeclaration {
	return &VariableDeclaration{
		Name: name,
		Type: typ,
		Expr: &ast.NumberLiteral{Token: token.Token{Type: token.NUMBER, Literal: value}},
	}
}

func runTests(t *testing.T, tests []testCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			result := g.Generate(tt.program)

			if result != tt.expected {
				t.Errorf("Test %s failed\nExpected:\n%q\nGot:\n%q",
					tt.name, tt.expected, result)
			}
		})
	}
}

func TestVariableDeclarations(t *testing.T) {
	tests := []testCase{
		{
			name: "number variable",
			program: createProgram(
				createVariableDeclaration("x", "number", "42"),
			),
			expected: "package main\n\nfunc main() {\n    x := 42\n}\n",
		},
		{
			name: "multiple variables",
			program: createProgram(
				createVariableDeclaration("x", "number", "42"),
				createVariableDeclaration("name", "string", "hello"),
				createVariableDeclaration("active", "boolean", "true"),
			),
			expected: "package main\n\nfunc main() {\n    x := 42\n    name := \"hello\"\n    active := true\n}\n",
		},
	}

	runTests(t, tests)
}

func TestExpressions(t *testing.T) {
	tests := []testCase{
		{
			name: "simple addition",
			program: createProgram(
				&VariableDeclaration{
					Name: "result",
					Type: "number",
					Expr: &ast.BinaryExpression{
						Left:     &ast.NumberLiteral{Token: token.Token{Type: token.NUMBER, Literal: "10"}},
						Operator: token.Token{Type: token.PLUS, Literal: "+"},
						Right:    &ast.NumberLiteral{Token: token.Token{Type: token.NUMBER, Literal: "20"}},
					},
				},
			),
			expected: "package main\n\nfunc main() {\n    result := (10 + 20)\n}\n",
		},
		{
			name: "nested expressions",
			program: createProgram(
				&VariableDeclaration{
					Name: "result",
					Type: "number",
					Expr: &ast.BinaryExpression{
						Left: &ast.NumberLiteral{Token: token.Token{Type: token.NUMBER, Literal: "10"}},
						Operator: token.Token{Type: token.PLUS, Literal: "+"},
						Right: &ast.BinaryExpression{
							Left:     &ast.NumberLiteral{Token: token.Token{Type: token.NUMBER, Literal: "5"}},
							Operator: token.Token{Type: token.MUL, Literal: "*"},
							Right:    &ast.NumberLiteral{Token: token.Token{Type: token.NUMBER, Literal: "3"}},
						},
					},
				},
			),
			expected: "package main\n\nfunc main() {\n    result := (10 + (5 * 3))\n}\n",
		},
	}

	runTests(t, tests)
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		program     *Program
		expectPanic bool
	}{
		{
			name:        "empty program",
			program:     &Program{Statements: []Statement{}},
			expectPanic: false,
		},
		{
			name:        "nil program",
			program:     nil,
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()

			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic, but none occurred")
					}
				}()
			}

			_ = g.Generate(tt.program)
		})
	}
}
