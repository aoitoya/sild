package parser

import (
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
	case token.TYPE_NUMBER, token.TYPE_STRING, token.TYPE_BOOLEAN, token.TYPE_VOID:
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
	case token.FUNCTION:
		stmt := p.parseFunctionDeclaration()
		if stmt == nil {
			return nil
		}
		return stmt
	case token.RETURN:
		stmt := p.parseReturnStatement()
		if stmt == nil {
			return nil
		}
		return stmt
	default:
		return nil
	}
}

func (p *Parser) parseFunctionDeclaration() *ast.FunctionDeclaration {
	fn := &ast.FunctionDeclaration{}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	p.nextTok()

	fn.Name = p.currTok

	if !p.expectPeek(token.LEFT_PAREN) {
		return nil
	}
	p.nextTok()

	fn.Params = []ast.FunctionParam{}
	for p.currTok.Type != token.RIGHT_PAREN {
		if p.expectPeek(token.RIGHT_PAREN) {
			p.nextTok()
			break
		}

		if !p.expectPeek(token.IDENT) {
			return nil
		}
		p.nextTok()

		paramName := p.currTok
		if !p.expectPeek(token.COLON) {
			return nil

		}
		p.nextTok()

		if !p.expectPeekValueType() {
			return nil
		}
		p.nextTok()

		fn.Params = append(fn.Params, ast.FunctionParam{Type: p.currTok, Name: paramName})
		if !p.expectPeek(token.COMMA) {
			p.nextTok()
			break
		}

		p.nextTok()
	}

	if !p.expectPeek(token.COLON) {
		return nil
	}
	p.nextTok()
	if !p.expectPeekValueType() {
		return nil
	}
	p.nextTok()
	fn.ReturnType = p.currTok

	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}
	p.nextTok()

	fn.Body = []ast.Statement{}
	for p.currTok.Type != token.RIGHT_BRACE {
		if p.currTok.Type == token.EOF {
			return nil
		}

		if p.peekTok.Type == token.RIGHT_BRACE {
			p.nextTok()
			break
		}
		p.nextTok()
		fn.Body = append(fn.Body, p.parseStatement())
	}
	// p.nextTok()

	return fn
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
		return nil
	}
	stmt.Expr = expr

	if p.currTok.Type != token.SEMICOLON {
		return nil
	}

	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseTerm()
}

func (p *Parser) parseTerm() ast.Expression {
	expr := p.parseFactor()
	if expr == nil {
		return nil
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.currTok
		p.nextTok()
		right := p.parseFactor()

		if right == nil {
			return nil
		}
		expr = &ast.BinaryExpression{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) parseFactor() ast.Expression {
	expr := p.parseUnary()
	if expr == nil {
		return nil
	}

	for p.match(token.MUL, token.DIV) {
		operator := p.currTok
		p.nextTok()
		right := p.parseUnary()

		if right == nil {
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

func (p *Parser) parseFunctionCall() ast.Expression {
	expr := p.currTok
	
	if !p.expectPeek(token.LEFT_PAREN) {
		return p.parsePrimary()
	}

	p.nextTok()
	
	args := []ast.Expression{}
	for p.currTok.Type != token.RIGHT_PAREN {
		// if p.expectPeek(token.RIGHT_PAREN) {
			p.nextTok()
		// 	break
		// }
		args = append(args, p.parseExpression())
		if p.currTok.Type != token.COMMA {
			p.nextTok()
			break
		}
	}
	
	 ex := &ast.FunctionCallExpression{Token: expr, Args: args}
	 return ex
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currTok}

	p.nextTok()

	stmt.Value = p.parseExpression()
	if p.peekTok.Type == token.SEMICOLON {
		p.nextTok()
	}

	return stmt
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
	case token.IDENT: {
		if p.expectPeek(token.LEFT_PAREN) {
			return p.parseFunctionCall()
		}
		
		return &ast.VariableExpression{Token: p.nextTok()}
	}
	default:
		return nil
	}
}

func (p *Parser) match(types ...token.TokenType) bool {
	return slices.Contains(types, p.currTok.Type)
}
