package scanner

import "github.com/codecrafters-io/shell-starter-go/app/util"

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case ' ':
	case '\t':
	case '\r':
	case '\'', '"':
		s.string()
	default:
		s.word()
	}
}

func (s *Scanner) string() {
	for s.peek() != '\'' && s.peek() != '"' && !s.isAtEnd() {
		s.advance()
	}
	if s.isAtEnd() {
		util.Error("Unterminated string.")
		return
	}
	s.advance()
}

func (s *Scanner) word() {
	for s.peek() != ' ' && s.peek() != '\t' && s.peek() != '\r' && !s.isAtEnd() {
		s.advance()
	}
	s.addToken(WORD)
}

func (s *Scanner) addToken(tokenType tokenType) {
	s.tokens = append(s.tokens, Token{
		tokenType: tokenType,
		Lexeme:    s.source[s.start:s.current],
	})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}
