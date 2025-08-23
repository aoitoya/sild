package scanner

import (
	"github.com/toyaAoi/sild/token"
)

type Scanner struct {
	input   string
	pos     int
	readPos int
	ch      byte
	pastTok token.Token
}

func (s *Scanner) readChar() {
	if s.readPos >= len(s.input) {
		s.ch = 0
	} else {
		s.ch = s.input[s.readPos]
	}

	s.pos = s.readPos
	s.readPos++
}

func (s *Scanner) NextToken() token.Token {
	s.skipWhiteSpaces()

	var tok token.Token

	switch s.ch {
	case ':':
		tok = s.newToken(token.COLON)
	case ';':
		tok = s.newToken(token.SEMICOLON)
	case '=':
		tok = s.newToken(token.ASSIGN)
	case '+':
		tok = s.newToken(token.PLUS)
	case '-':
		tok = s.newToken(token.MINUS)
	case '*':
		tok = s.newToken(token.MUL)
	case '/':
		tok = s.newToken(token.DIV)
	case '!':
		tok = s.newToken(token.BANG)
	case '"':
		s.readChar()

		tok.Type = token.STRING
		tok.Literal = s.readStr()

		s.readChar()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(s.ch) {
			tok.Literal = s.readIdent()

			// type token
			if s.pastTok.Type == token.COLON && s.peakNextChar() == '=' {
				tok.Type = token.LookupType(tok.Literal)
			} else {
				tok.Type = token.LookupIdent(tok.Literal)
			}
		} else if isDigit(s.ch) {
			tok.Type = token.NUMBER
			tok.Literal = s.readInt()
		}
	}

	s.pastTok = tok
	return tok
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}

func (s *Scanner) skipWhiteSpaces() {
	for isWhiteSpace(s.ch) {
		s.readChar()
	}
}

func (s *Scanner) peakNextChar() byte {
	idx := s.readPos
	for isWhiteSpace(s.input[idx]) {
		idx++
	}

	return s.input[idx]
}

func (s *Scanner) newToken(tok token.TokenType) token.Token {
	token := token.Token{Type: tok, Literal: string(s.ch)}
	s.readChar()
	return token
}

func (s *Scanner) readIdent() string {
	pos := s.pos
	for isLetter(s.ch) {
		s.readChar()
	}

	return s.input[pos:s.pos]
}

func (s *Scanner) readInt() string {
	pos := s.pos
	for isDigit(s.ch) {
		s.readChar()
	}

	return s.input[pos:s.pos]
}

func (s *Scanner) readStr() string {
	pos := s.pos
	for s.ch != '"' && s.ch != 0 {
		s.readChar()
	}
	return s.input[pos:s.pos]
}

func isLetter(ch byte) bool {
	return ('a' <= ch && 'z' >= 'a') || ('A' <= ch && 'Z' >= ch)
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func New(input string) *Scanner {
	s := &Scanner{input: input}
	s.readChar()

	return s
}
