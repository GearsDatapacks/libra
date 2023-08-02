package token

type Type string

type Token struct {
	Type           Type
	Value          string
	Line           int
	Column         int
	LeadingNewline bool
}

const (
	EOF = "EOF"

	INTEGER    = "INTEGER"
	IDENTIFIER = "IDENTIFIER"

	LEFT_PAREN  = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	OPERATOR    = "OPERATOR"
	EQUALS      = "EQUALS"
)

var Symbols = map[string]Type{
	"(": LEFT_PAREN,
	")": RIGHT_PAREN,

	"+": OPERATOR,
	"-": OPERATOR,
	"*": OPERATOR,
	"/": OPERATOR,
	"%": OPERATOR,

	"=": EQUALS,
}

func New(line int, offset int, tokenType Type, value []rune, leadingNewline bool) Token {
	return Token{
		Type:           tokenType,
		Value:          string(value),
		Line:           line,
		Column:         offset,
		LeadingNewline: leadingNewline,
	}
}
