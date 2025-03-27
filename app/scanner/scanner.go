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
	case '\n':
	case '\'':
		s.stringSingleQuote()
	case '"':
		s.stringDoubleQuote()
	default:
		s.word()
	}
}

func (s *Scanner) stringSingleQuote() {
	for s.peek() != '\'' && !s.isAtEnd() {
		s.advance()
	}

	if s.isAtEnd() {
		util.Error("Unterminated string.")
		return
	}
	s.advance()

	// handle link quotes
	value := s.source[s.start+1 : s.current-1]
	for s.peek() == '\'' && !s.isAtEnd() {
		s.advance()

		start := s.current
		for s.peek() != '\'' && !s.isAtEnd() {
			s.advance()
		}
		if s.isAtEnd() {
			util.Error("Unterminated string.")
			return
		}
		s.advance()
		value += s.source[start : s.current-1]
	}
	s.addToken(STRING, value)
}

func (s *Scanner) stringDoubleQuote() {
	for s.peek() != '"' && !s.isAtEnd() {
		s.advance()
	}
	if s.isAtEnd() {
		util.Error("Unterminated string.")
		return
	}
	s.advance()

	// handle link quotes
	value := s.source[s.start+1 : s.current-1]
	for s.peek() == '"' && !s.isAtEnd() {
		s.advance()

		start := s.current
		for s.peek() != '"' && !s.isAtEnd() {
			s.advance()
		}
		if s.isAtEnd() {
			util.Error("Unterminated string.")
			return
		}
		s.advance()
		value += s.source[start : s.current-1]
	}
	s.addToken(STRING, value)
}

func (s *Scanner) word() {
	for s.peek() != ' ' && s.peek() != '\t' && s.peek() != '\r' && s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
	s.addToken(WORD, s.source[s.start:s.current])
}

func (s *Scanner) addToken(tokenType tokenType, lexeme string) {
	s.tokens = append(s.tokens, Token{
		tokenType: tokenType,
		Lexeme:    lexeme,
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
