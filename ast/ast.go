package ast

import (
	"fmt"

	"github.com/toyaAoi/sild/token"
)

type Node interface {
	TokenLiteral()
	String() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type Program struct {
	Statements []Statement
}

type Identifier struct {
	Token token.Token
	Value string
}

type ValueType token.Token

type VariableDeclaration struct {
	Name  string
	Type  string
	Value string
}

func (v *VariableDeclaration) TokenLiteral() {}
func (v *VariableDeclaration) String() string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("name: %q, type: %q, value: %q", v.Name, v.Type, v.Value)
}
