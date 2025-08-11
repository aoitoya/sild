package scanner

type Scanner struct {
	input   string
	pos     int
	readPos int
	ch      byte
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

func (s *Scanner) NextToken() Token {
	s.skipWhitespaces()

	var tok Token

	switch s.ch {
	case ':':
		tok = s.newToken(COLON)
	case ';':
		tok = s.newToken(SEMICOLON)
	case '=':
		tok = s.newToken(EQUAL)
	case '+':
		tok = s.newToken(PLUS)
	default:
		if isLetter(s.ch) {
			tok.Literal = s.readIdent()
			tok.Type = lookupIdent(tok.Literal)
		} else if isDigit(s.ch) {
			tok.Type = INT
			tok.Literal = s.readInt()
		}
	}

	return tok
}

func (s *Scanner) skipWhitespaces() {
	for s.ch == ' ' || s.ch == '\n' || s.ch == '\t' || s.ch == '\r' {
		s.readChar()
	}
}

func (s *Scanner) newToken(tok string) Token {
	token := Token{tok, string(s.ch)}
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

func isLetter(ch byte) bool {
	return ('a' <= ch && 'z' >= 'a') || ('A' <= ch && 'Z' >= ch) || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func New(input string) *Scanner {
	s := &Scanner{input: input}
	s.readChar()

	return s
}
