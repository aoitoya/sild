package codegen

import (
	"fmt"
	"strings"

	"github.com/toyaAoi/sild/ast"
	"github.com/toyaAoi/sild/token"
)

type Program = ast.Program
type Statement = ast.Statement

type Generator struct {
	output strings.Builder
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(p *Program) string {
	g.output.Reset()

	g.output.WriteString("package main\n\n")

	for _, stmt := range p.Statements {
		if _, ok := stmt.(*ast.FunctionDeclaration); ok {
			g.output.WriteString(g.generateFunctionDeclaration(stmt.(*ast.FunctionDeclaration)) + "\n")
		}
	}

	g.output.WriteString("func main() {\n")

	for _, stmt := range p.Statements {
		if _, ok := stmt.(*ast.FunctionDeclaration); ok {
			continue
		}

		g.output.WriteString("    " + g.generateStatement(stmt) + "\n")
	}

	g.output.WriteString("}\n")

	return g.output.String()
}

func (g *Generator) generateStatement(stmt Statement) string {
	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		return g.generateVariableDeclaration(s)
	case *ast.FunctionDeclaration:
		return g.generateFunctionDeclaration(s)
	case *ast.ReturnStatement:
		return g.generateReturnStatement(s)
	default:
		return ""
	}
}

func (g *Generator) generateVariableDeclaration(varDec *ast.VariableDeclaration) string {
	if varDec.Type == "string" {
		return fmt.Sprintf("%s := %q", varDec.Name, varDec.Expr.String())
	}
	return fmt.Sprintf("%s := %s", varDec.Name, varDec.Expr.String())
}

func returnType(t token.TokenType) string {
	switch t {
	case token.TYPE_NUMBER:
		return "int"
	case token.TYPE_STRING:
		return "string"
	case token.TYPE_BOOLEAN:
		return "bool"
	case token.TYPE_VOID:
		return ""
	default:
		return ""
	}
}

func (g *Generator) generateFunctionDeclaration(fn *ast.FunctionDeclaration) string {
	builder := strings.Builder{}

	params := ""
	for i, p := range fn.Params {
		if i > 0 {
			params += ", "
		}

		params += p.String()
	}
	builder.WriteString("func ")
	builder.WriteString(fn.Name.Literal)
	builder.WriteString("(")
	builder.WriteString(params)
	builder.WriteString(")")
	if fn.ReturnType.Type != token.TYPE_VOID {
		builder.WriteString(" ")
		builder.WriteString(returnType(fn.ReturnType.Type))
	}
	builder.WriteString(" {\n")

	for _, stmt := range fn.Body {
		builder.WriteString("    " + g.generateStatement(stmt) + "\n")
	}

	builder.WriteString("}\n")

	return builder.String()
}

func (g *Generator) generateReturnStatement(stmt *ast.ReturnStatement) string {
	return fmt.Sprintf("return %s", stmt.Value.String())
}
