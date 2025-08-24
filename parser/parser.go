package parser

import (
	"fmt"
	"slices"

	"github.com/toyaAoi/sild/ast"
	"github.com/toyaAoi/sild/scanner"
	"github.com/toyaAoi/sild/token"
)

type Parser struct {
	s       *scanner.Scanner
	currTok token.Token
	peekTok token.Token
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currTok.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt == nil {
			return program
		}
		program.Statements = append(program.Statements, stmt)

		if p.peekTok.Type != token.EOF {
			p.nextTok()
		} else {
			break
		}
	}

	return program
}

func (p *Parser) nextTok() token.Token {
	tok := p.currTok
	p.currTok = p.peekTok
	p.peekTok = p.s.NextToken()
	return tok
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	return p.peekTok.Type == t
}

func (p *Parser) expectPeekValueType() bool {
	switch p.peekTok.Type {
	case token.TYPE_NUMBER, token.TYPE_STRING, token.TYPE_BOOLEAN:
		return true
	default:
		return false
	}
}

func New(sc *scanner.Scanner) *Parser {
	p := &Parser{s: sc}

	p.nextTok()
	p.nextTok()

	return p
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currTok.Type {
	case token.LET:
		stmt := p.parseVariableDeclaration()
		if stmt == nil {
			return nil
		}
		return stmt
	default:
		return nil
	}
}

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	stmt := &ast.VariableDeclaration{}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	p.nextTok()

	stmt.Name = p.currTok.Literal

	if !p.expectPeek(token.COLON) {
		return nil
	}
	p.nextTok()
	if !p.expectPeekValueType() {
		return nil
	}
	p.nextTok()
	stmt.Type = p.currTok.Literal

	p.nextTok()
	if p.currTok.Type != token.ASSIGN {
		return nil
	}
	p.nextTok()
	expr := p.parseExpression()
	if expr == nil {
		fmt.Println(107, p)
		return nil
	}
	stmt.Expr = expr

	if p.currTok.Type != token.SEMICOLON {
		fmt.Println(111, p)
		return nil
	}

	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	// fmt.Println(120, p)
	return p.parseTerm()
}

func (p *Parser) parseTerm() ast.Expression {
	expr := p.parseFactor()
	if expr == nil {
		fmt.Println("parseFactor returned nil in parseTerm")
		return nil
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.currTok
		p.nextTok()
		right := p.parseFactor()

		if right == nil {
			fmt.Println("parseFactor returned nil in parseTerm loop")
			return nil
		}
		expr = &ast.BinaryExpression{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) parseFactor() ast.Expression {
	expr := p.parseUnary()
	if expr == nil {
		fmt.Println("parseUnary returned nil in parseFactor")
		return nil
	}

	for p.match(token.MUL, token.DIV) {
		operator := p.currTok
		p.nextTok()
		right := p.parseUnary()

		if right == nil {
			fmt.Println("parseUnary returned nil in parseFactor loop")
			return nil
		}
		expr = &ast.BinaryExpression{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) parseUnary() ast.Expression {

	if p.match(token.BANG, token.MINUS) {
		operator := p.currTok
		p.nextTok()
		right := p.parseUnary()
		if right == nil {
			return nil
		}
		return &ast.UnaryExpression{Operator: operator, Right: right}
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ast.Expression {
	switch p.currTok.Type {
	case token.NUMBER:

		return &ast.NumberLiteral{Token: p.nextTok()}
	case token.STRING:

		return &ast.StringLiteral{Token: p.nextTok()}
	case token.BOOLEAN:
		return &ast.BooleanLiteral{Token: p.nextTok()}
	case token.LEFT_PAREN:
		p.nextTok()
		expr := p.parseExpression()
		if expr == nil {
			return nil
		}
		if p.nextTok().Type != token.RIGHT_PAREN {
			return nil
		}
		return expr
	default:
		return nil
	}
}

func (p *Parser) match(types ...token.TokenType) bool {
	return slices.Contains(types, p.currTok.Type)
}
