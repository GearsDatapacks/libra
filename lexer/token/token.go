package token

type Type string

type Token struct {
	Type   Type
	Value  string
	Line   int
	Offset int
}

const (
	EOF = "EOF"

	INTEGER = "INTEGER"

	LEFT_PAREN  = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	OPERATOR    = "OPERATOR"
)

var Symbols = map[string]Type{
	"(": LEFT_PAREN,
	")": RIGHT_PAREN,

	"+": OPERATOR,
	"-": OPERATOR,
	"*": OPERATOR,
	"/": OPERATOR,
	"%": OPERATOR,
}

func New(line int, offset int, tokenType Type, value []rune) Token {
	return Token{
		Type:   tokenType,
		Value:  string(value),
		Line:   line,
		Offset: offset,
	}
}
