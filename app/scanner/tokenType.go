package scanner

type tokenType int

const (
	WORD tokenType = iota
)

var tokenTypeNames = map[tokenType]string{
	WORD: "WORD",
}

func (t tokenType) String() string {
	return tokenTypeNames[t]
}
