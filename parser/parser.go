package parser

import (
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
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextTok()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currTok.Type {
	case token.LET:
		return p.parseVariableDeclaration()
	default:
		return nil
	}
}

func (p *Parser) nextTok() {
	p.currTok = p.peekTok
	p.peekTok = p.s.NextToken()
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

func (p *Parser) expectPeekValue() bool {
	switch p.peekTok.Type {
	case token.NUMBER, token.STRING, token.BOOLEAN:
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

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	stmt := &ast.VariableDeclaration{}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	p.nextTok()

	stmt.Name = p.currTok.Literal

	if p.expectPeek(token.COLON) {
		p.nextTok()
		if p.expectPeekValueType() {
			p.nextTok()
			stmt.Type = p.currTok.Literal
		} else {
			return nil
		}
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextTok()

	if !p.expectPeekValue() {
		return nil
	}
	p.nextTok()

	stmt.Value = p.currTok.Literal

	if !p.expectPeek(token.SEMICOLON) {
		p.nextTok()
	}

	return stmt
}
