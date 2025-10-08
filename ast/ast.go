package ast

import (
	"fmt"

	"github.com/toyaAoi/sild/token"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type Identifier struct {
	Token token.Token
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Token.Literal
}

type VariableDeclaration struct {
	Name string
	Type string
	Expr Expression
}

func (v *VariableDeclaration) statementNode() {}
func (v *VariableDeclaration) String() string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("name: %q, type: %q, value: %q", v.Name, v.Type, v.Expr.String())
}

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (b *BinaryExpression) expressionNode() {}
func (b *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Operator.Literal, b.Right.String())
}

type NumberLiteral struct {
	Token token.Token
}

func (s *NumberLiteral) expressionNode() {}
func (s *NumberLiteral) String() string {
	return s.Token.Literal
}

type StringLiteral struct {
	Token token.Token
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) String() string {
	return s.Token.Literal
}

type BooleanLiteral struct {
	Token token.Token
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) String() string {
	return b.Token.Literal
}

type UnaryExpression struct {
	Operator token.Token
	Right    Expression
}

func (u *UnaryExpression) expressionNode() {}
func (u *UnaryExpression) String() string {
	return fmt.Sprintf("%s%s", u.Operator.Literal, u.Right.String())
}

type ParenthesizedExpression struct {
	Expression Expression
}

func (p *ParenthesizedExpression) expressionNode() {}
func (p *ParenthesizedExpression) String() string {
	return fmt.Sprintf("(%s)", p.Expression.String())
}

type FunctionParam struct {
	Type token.Token
	Name token.Token
}

func (fp *FunctionParam) String() string {
	return fmt.Sprintf("%s %s", fp.Name.Literal, mapType(fp.Type.Literal))
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) String() string {
	if r == nil {
		return "<nil>"
	}
	return "return " + r.Value.String()
}

type FunctionDeclaration struct {
    Name       token.Token
    Params     []FunctionParam
    Body       []Statement
    ReturnType token.Token
}

func (f *FunctionDeclaration) statementNode() {}
func (f *FunctionDeclaration) String() string {
	if f == nil {
		return "<nil>"
	}

	var params []string
	for _, param := range f.Params {
		params = append(params, param.String())
	}

	var body []string
	for _, stmt := range f.Body {
		body = append(body, stmt.String())
	}

	return fmt.Sprintf("name: %q, params: %q, body: %q, return type: %q", f.Name.Literal, params, body, f.ReturnType.Literal)
}

func mapType(tsType string) string {

	typeMap := map[string]string{
		"number":  "int",
		"string":  "string",
		"boolean": "bool",
		"void":    "",
	}

	if goType, exists := typeMap[tsType]; exists {
		return goType
	}

	return "any"
}

type VariableExpression struct {
	Token token.Token
}

func (v *VariableExpression) expressionNode() {}
func (v *VariableExpression) String() string {
	return v.Token.Literal
}
