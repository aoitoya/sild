package parser

import (
	"fmt"
	"testing"

	"github.com/toyaAoi/sild/ast"
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
					fmt.Printf("%T", program.Statements[0])
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

func TestFunctionDeclarationParsing(t *testing.T) {
	tests := []struct {
		name               string
		tsInput            string
		expectedName       string
		expectedParams     []string
		expectedBody       []string
		expectedReturnType string
	}{
		{
			name:         "simple function with number params and return",
			tsInput:      "function add(x: number, y: number): number { return x + y }",
			expectedName: "add",
			expectedParams: []string{
				"x int",
				"y int",
			},
			expectedReturnType: "number",
			expectedBody:       []string{"return (x + y)"},
		},
		// {
		// 	name:         "function with string param and return",
		// 	tsInput:      "function getName(name: string): string { return name }",
		// 	expectedName: "getName",
		// 	expectedParams: []string{
		// 		"name string",
		// 	},
		// 	expectedReturnType: "string",
		// 	expectedBody:       []string{"return name"},
		// },
		// {
		// 	name:         "function with boolean param and return",
		// 	tsInput:      "function isActive(active: boolean): boolean { return active }",
		// 	expectedName: "isActive",
		// 	expectedParams: []string{
		// 		"active bool",
		// 	},
		// 	expectedReturnType: "boolean",
		// 	expectedBody:       []string{"return active"},
		// },
		// {
		// 	name:               "function with no params and void return",
		// 	tsInput:            "function greet(): void {}",
		// 	expectedName:       "greet",
		// 	expectedParams:     []string{},
		// 	expectedReturnType: "void",
		// 	expectedBody:       []string{},
		// },
		// {
		// 	name:               "function with variable declaration",
		// 	tsInput:            "function test(): void { let x: number = 42; }",
		// 	expectedName:       "test",
		// 	expectedParams:     []string{},
		// 	expectedReturnType: "void",
		// 	expectedBody:       []string{"x := 42"},
		// },
		// {
		// 	name:         "function with multiple statements",
		// 	tsInput:      "function complex(a: number, b: number): number { let x: number = 10 * a; let y: number = 20 * b; let z: number = x + y; return x + y + z}",
		// 	expectedName: "complex",
		// 	expectedParams: []string{
		// 		"a int",
		// 		"b int",
		// 	},
		// 	expectedReturnType: "number",
		// 	expectedBody: []string{
		// 		"x := 10 * a",
		// 		"y := 20 * b",
		// 		"z := x + y",
		// 		"return x + y + z",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(scanner.New(tt.tsInput))
			program := p.ParseProgram()

			if len(program.Statements) == 0 {
				t.Fatal("no statements parsed")
			}

			fnDecl, ok := program.Statements[0].(*ast.FunctionDeclaration)
			if !ok {
				t.Fatalf("expected FunctionDeclaration, got %T", program.Statements[0])
			}

			if fnDecl.Name.Literal != tt.expectedName {
				t.Errorf("expected function name %q, got %q", tt.expectedName, fnDecl.Name.Literal)
			}

			if len(fnDecl.Params) != len(tt.expectedParams) {
				t.Fatalf("expected %d params, got %d", len(tt.expectedParams), len(fnDecl.Params))
			}

			for i, param := range fnDecl.Params {
				expectedParam := tt.expectedParams[i]
				paramStr := param.String()
				if paramStr != expectedParam {
					t.Errorf("param %d: expected %q, got %q", i, expectedParam, paramStr)
				}
			}

			if fnDecl.ReturnType.Literal != tt.expectedReturnType {
				t.Errorf("expected return type %q, got %q", tt.expectedReturnType, fnDecl.ReturnType.Literal)
			}

			for i, stmt := range fnDecl.Body {
				expectedStmt := tt.expectedBody[i]
				stmtStr := stmt.String()
				if stmtStr != expectedStmt {
					t.Errorf("stmt %d: expected %q, got %q", i, expectedStmt, stmtStr)
				}
			}
		})
	}
}

func TestFunctionErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		desc  string
	}{
		{
			name:  "missing_function_keyword",
			input: "greet(): void {}",
			desc:  "should fail without function keyword",
		},
		{
			name:  "missing_function_name",
			input: "function (): void {}",
			desc:  "should fail without function name",
		},
		{
			name:  "missing_opening_paren",
			input: "function greet): void {}",
			desc:  "should fail without opening parenthesis",
		},
		{
			name:  "missing_closing_paren",
			input: "function greet(: void {}",
			desc:  "should fail without closing parenthesis",
		},
		{
			name:  "missing_return_type_colon",
			input: "function greet() void {}",
			desc:  "should fail without colon before return type",
		},
		{
			name:  "missing_return_type",
			input: "function greet(): {}",
			desc:  "should fail without return type",
		},
		{
			name:  "missing_opening_brace",
			input: "function greet(): void }",
			desc:  "should fail without opening brace",
		},
		{
			name:  "missing_closing_brace",
			input: "function greet(): void {",
			desc:  "should fail without closing brace",
		},
		{
			name:  "invalid_parameter_syntax",
			input: "function add(a number): number {}",
			desc:  "should fail without colon in parameter",
		},
		{
			name:  "missing_parameter_type",
			input: "function add(a:): number {}",
			desc:  "should fail without parameter type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New(scanner.New(tt.input))
			program := parser.ParseProgram()

			if program != nil && len(program.Statements) > 0 {
				if funcDecl, ok := program.Statements[0].(*ast.FunctionDeclaration); ok {
					if isValidFunctionDeclaration(funcDecl) {
						t.Errorf("❌ Expected parsing to fail for %q (%s), but got valid function", tt.input, tt.desc)
					}
				}
			}

			t.Logf("✅ Error case handled correctly: %s", tt.desc)
		})
	}
}

func isValidFunctionDeclaration(funcDecl *ast.FunctionDeclaration) bool {
	return funcDecl.Name.Literal != "" && funcDecl.ReturnType.Literal != "" && funcDecl.Body != nil
}
