package token

type Type string

type Token struct {
	Type Type
	Value string
	Start int
	End int
}

const (
	EOF = "EOF"

	NUMBER = "NUMBER"

	LEFT_PAREN = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	OPERATOR = "OPERATOR"
)

var Symbols = map[string]Type {
	"(": LEFT_PAREN,
	")": RIGHT_PAREN,

	"+": OPERATOR,
	"-": OPERATOR,
	"*": OPERATOR,
	"/": OPERATOR,
	"%": OPERATOR,
}

func New(start int, end int, tokenType Type, value []rune) Token {
	return Token{
		Type: tokenType,
		Value: string(value),
		Start: start,
		End: end,
	}
}
