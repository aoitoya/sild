package codegen

import (
	"strings"
	"testing"

	"github.com/toyaAoi/sild/ast"
	"github.com/toyaAoi/sild/parser"
	"github.com/toyaAoi/sild/scanner"
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


func TestFunctionCodeGeneration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple_void_function",
			input: `function greet(): void {}`,
			expected: `package main

func greet() {
}

func main() {
}
`,
		},
		{
			name: "function_with_return_type",
			input: `function getNumber(): number {
    return 42;
}`,
			expected: `package main

func getNumber() int {
    return 42
}

func main() {
}
`,
		},
		{
			name: "function_with_single_parameter",
			input: `function double(x: number): number {
    return x * 2;
}`,
			expected: `package main

func double(x int) int {
    return (x * 2)
}

func main() {
}
`,
		},
		{
			name: "function_with_multiple_parameters",
			input: `function add(a: number, b: number): number {
    return a + b;
}`,
			expected: `package main

func add(a int, b int) int {
    return (a + b)
}

func main() {
}
`,
		},
		{
			name: "function_with_string_parameter",
			input: `function greet(name: string): string {
    return name;
}`,
			expected: `package main

func greet(name string) string {
    return name
}

func main() {
}
`,
		},
		{
			name: "function_with_boolean_parameter",
			input: `function isValid(active: boolean): boolean {
    return active;
}`,
			expected: `package main

func isValid(active bool) bool {
    return active
}

func main() {
}
`,
		},
		{
			name: "function_with_mixed_parameters",
			input: `function process(name: string, age: number, active: boolean): string {
    return name;
}`,
			expected: `package main

func process(name string, age int, active bool) string {
    return name
}

func main() {
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := scanner.New(tt.input)
			parser := parser.New(scanner)
			program := parser.ParseProgram()
			
			if program == nil {
				t.Fatalf("Failed to parse input: %s", tt.input)
			}
			
			generator := New()
			output := generator.Generate(program)
			
			if !compareOutput(output, tt.expected) {
				t.Errorf("Output mismatch\nExpected:\n%s\nGot:\n%s", tt.expected, output)
			}
		})
	}
}

func TestFunctionAndVariableMix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "variable_before_function",
			input: `let x: number = 42;
function getValue(): number {
    return x;
}`,
			expected: `package main

func getValue() int {
    return x
}

func main() {
    x := 42
}
`,
		},
		{
			name: "function_before_variable",
			input: `function getValue(): number {
    return 42;
}
let x: number = 10;`,
			expected: `package main

func getValue() int {
    return 42
}

func main() {
    x := 10
}
`,
		},
		{
			name: "multiple_functions_and_variables",
			input: `function add(a: number, b: number): number {
    return a + b;
}
let x: number = 5;
function multiply(a: number, b: number): number {
    return a * b;
}
let y: number = 10;`,
			expected: `package main

func add(a int, b int) int {
    return (a + b)
}

func multiply(a int, b int) int {
    return (a * b)
}

func main() {
    x := 5
    y := 10
}
`,
		},
		{
			name: "function_call_in_main",
			input: `function double(x: number): number {
    return x * 2;
}
let result: number = double(21);`,
			expected: `package main

func double(x int) int {
    return (x * 2)
}

func main() {
    result := double(21)
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := scanner.New(tt.input)
			parser := parser.New(scanner)
			program := parser.ParseProgram()
			
			if program == nil {
				t.Fatalf("Failed to parse input: %s", tt.input)
			}
			
			generator := New()
			output := generator.Generate(program)
			
			if !compareOutput(output, tt.expected) {
				t.Errorf("Output mismatch\nExpected:\n%s\nGot:\n%s", tt.expected, output)
			}
		})
	}
}

func TestFunctionWithComplexBody(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "function_with_multiple_statements",
			input: `function calculate(x: number): number {
    let doubled: number = x * 2;
    let result: number = doubled + 10;
    return result;
}`,
			expected: `package main

func calculate(x int) int {
    doubled := (x * 2)
    result := (doubled + 10)
    return result
}

func main() {
}
`,
		},
		{
			name: "function_with_arithmetic",
			input: `function compute(a: number, b: number): number {
    let sum: number = a + b;
    let product: number = sum * 2;
    return product;
}`,
			expected: `package main

func compute(a int, b int) int {
    sum := (a + b)
    product := (sum * 2)
    return product
}

func main() {
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := scanner.New(tt.input)
			parser := parser.New(scanner)
			program := parser.ParseProgram()
			
			if program == nil {
				t.Fatalf("Failed to parse input: %s", tt.input)
			}
			
			generator := New()
			output := generator.Generate(program)
			
			if !compareOutput(output, tt.expected) {
				t.Errorf("Output mismatch\nExpected:\n%s\nGot:\n%s", tt.expected, output)
			}
		})
	}
}

func TestMultipleFunctions(t *testing.T) {
	input := `function add(a: number, b: number): number {
    return a + b;
}

function subtract(a: number, b: number): number {
    return a - b;
}

function multiply(a: number, b: number): number {
    return a * b;
}

let result: number = add(10, 5);`

	expected := `package main

func add(a int, b int) int {
    return (a + b)
}

func subtract(a int, b int) int {
    return (a - b)
}

func multiply(a int, b int) int {
    return (a * b)
}

func main() {
    result := add(10, 5)
}
`

	scanner := scanner.New(input)
	parser := parser.New(scanner)
	program := parser.ParseProgram()
	
	if program == nil {
		t.Fatalf("Failed to parse input")
	}
	
	generator := New()
	output := generator.Generate(program)
	
	if !compareOutput(output, expected) {
		t.Errorf("Output mismatch\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestEmptyFunctionBody(t *testing.T) {
	input := `function doNothing(): void {
}`

	expected := `package main

func doNothing() {
}

func main() {
}
`

	scanner := scanner.New(input)
	parser := parser.New(scanner)
	program := parser.ParseProgram()
	
	if program == nil {
		t.Fatalf("Failed to parse input")
	}
	
	generator := New()
	output := generator.Generate(program)
	
	if !compareOutput(output, expected) {
		t.Errorf("Output mismatch\nExpected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestFunctionReturningExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "return_arithmetic",
			input: `function calculate(): number {
    return 10 + 20;
}`,
			expected: `package main

func calculate() int {
    return (10 + 20)
}

func main() {
}
`,
		},
		{
			name: "return_complex_expression",
			input: `function compute(x: number): number {
    return x * 2 + 10;
}`,
			expected: `package main

func compute(x int) int {
    return ((x * 2) + 10)
}

func main() {
}
`,
		},
		{
			name: "return_parenthesized",
			input: `function calc(a: number, b: number): number {
    return (a + b) * 2;
}`,
			expected: `package main

func calc(a int, b int) int {
    return ((a + b) * 2)
}

func main() {
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := scanner.New(tt.input)
			parser := parser.New(scanner)
			program := parser.ParseProgram()
			
			if program == nil {
				t.Fatalf("Failed to parse input: %s", tt.input)
			}
			
			generator := New()
			output := generator.Generate(program)
			
			if !compareOutput(output, tt.expected) {
				t.Errorf("Output mismatch\nExpected:\n%s\nGot:\n%s", tt.expected, output)
			}
		})
	}
}

// Helper function to compare outputs ignoring minor whitespace differences
func compareOutput(got, expected string) bool {
	got = strings.TrimSpace(got)
	expected = strings.TrimSpace(expected)
	
	// Normalize line endings
	got = strings.ReplaceAll(got, "\r\n", "\n")
	expected = strings.ReplaceAll(expected, "\r\n", "\n")
	
	return got == expected
}