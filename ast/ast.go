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
	Value string
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
	return fmt.Sprintf("%s %s %s", b.Left.String(), b.Operator.Literal, b.Right.String())
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
