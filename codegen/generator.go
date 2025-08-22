package codegen

import (
	"fmt"
	"strings"

	"github.com/toyaAoi/sild/ast"
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
	g.output.WriteString("func main() {\n")

	for _, stmt := range p.Statements {
		g.output.WriteString("    " + g.generateStatement(stmt) + "\n")
	}

	g.output.WriteString("}\n")

	return g.output.String()
}

func (g *Generator) generateStatement(stmt Statement) string {
	switch s := stmt.(type) {
	case *ast.VariableDeclaration:
		return g.generateVariableDeclaration(s)
	default:
		return ""
	}
}

func (g *Generator) generateVariableDeclaration(varDec *ast.VariableDeclaration) string {
	if varDec.Type == "string" {
		return fmt.Sprintf("%s := %q", varDec.Name, varDec.Value)
	}
	return fmt.Sprintf("%s := %s", varDec.Name, varDec.Value)
}
