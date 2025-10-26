package scanner

import (
	"io"

	"github.com/toyaAoi/sild/token"
)

type Scanner struct {
	r        io.Reader
	buf      []byte
	pos      int
	ch       byte
	pastTok  token.Token
	eof      bool
	readErr  error
}

func (s *Scanner) readChar() {
	if s.pos >= len(s.buf) {
		s.buf = make([]byte, 1024)
		n, err := s.r.Read(s.buf)
		if err != nil && err != io.EOF {
			s.readErr = err
			s.ch = 0
			return
		}
		s.buf = s.buf[:n]
		s.pos = 0

		if n == 0 {
			s.ch = 0
			s.eof = true
			return
		}
	}

	s.ch = s.buf[s.pos]
	s.pos++
}

func (s *Scanner) NextToken() token.Token {
	s.skipWhiteSpaces()

	var tok token.Token

	switch s.ch {
	case ',':
		tok = s.newToken(token.COMMA)
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
	case '(':
		tok = s.newToken(token.LEFT_PAREN)
	case ')':
		tok = s.newToken(token.RIGHT_PAREN)
	case '{':
		tok = s.newToken(token.LEFT_BRACE)
	case '}':
		tok = s.newToken(token.RIGHT_BRACE)
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

			// type token - either after colon in type annotation or in variable declaration
			if s.pastTok.Type == token.COLON {
				// Check if it's a type name
				if typ := token.LookupType(tok.Literal); typ != "" {
					tok.Type = typ
				} else {
					tok.Type = token.LookupIdent(tok.Literal)
				}
			} else if s.peakNextChar() == '=' {
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
	// Save current position and character
	oldPos := s.pos
	oldCh := s.ch
	oldBuf := make([]byte, len(s.buf))
	copy(oldBuf, s.buf)

	// Read next non-whitespace character
	for isWhiteSpace(s.ch) && !s.eof && s.readErr == nil {
		s.readChar()
	}

	ch := s.ch

	// Restore position and character
	s.pos = oldPos
	s.ch = oldCh
	s.buf = oldBuf

	return ch
}

func (s *Scanner) newToken(tok token.TokenType) token.Token {
	token := token.Token{Type: tok, Literal: string(s.ch)}
	s.readChar()
	return token
}

func (s *Scanner) readIdent() string {
	var buf []byte
	for isLetter(s.ch) {
		buf = append(buf, s.ch)
		s.readChar()
	}
	return string(buf)
}

func (s *Scanner) readInt() string {
	var buf []byte
	for isDigit(s.ch) {
		buf = append(buf, s.ch)
		s.readChar()
	}
	return string(buf)
}

func (s *Scanner) readStr() string {
	var buf []byte
	for s.ch != '"' && s.ch != 0 {
		buf = append(buf, s.ch)
		s.readChar()
	}
	return string(buf)
}

func isLetter(ch byte) bool {
	return ('a' <= ch && 'z' >= 'a') || ('A' <= ch && 'Z' >= ch)
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func New(r io.Reader) *Scanner {
	s := &Scanner{r: r, buf: make([]byte, 0, 1024)}
	s.readChar()
	return s
}
