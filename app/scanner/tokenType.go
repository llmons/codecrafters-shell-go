package scanner

type tokenType int

const (
	WORD tokenType = iota
	STRING
)

var tokenTypeNames = map[tokenType]string{
	WORD:   "WORD",
	STRING: "STRING",
}

func (t tokenType) String() string {
	return tokenTypeNames[t]
}
