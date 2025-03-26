package scanner

type Token struct {
	tokenType tokenType
	Lexeme    string
}

func (t Token) String() string {
	return t.Lexeme
}
