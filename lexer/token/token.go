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
	FLOAT      = "FLOAT"
	IDENTIFIER = "IDENTIFIER"

	LEFT_PAREN  = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	COMMA       = "COMMA"

	ASSIGNMENT_OPERATOR     = "ASSIGNMENT_OPERATOR"
	LOGICAL_OPERATOR        = "LOGICAL_OPERATOR"
	COMPARISON_OPERATOR     = "COMPARISON_OPERATOR"
	ADDITIVE_OPERATOR       = "ADDITIVE_OPERATOR"
	MULTIPLICATIVE_OPERATOR = "MULTIPLICATIVE_OPERATOR"
)

var Symbols = map[string]Type{
	"(": LEFT_PAREN,
	")": RIGHT_PAREN,
	",": COMMA,

	"+": ADDITIVE_OPERATOR,
	"-": ADDITIVE_OPERATOR,
	"*": MULTIPLICATIVE_OPERATOR,
	"/": MULTIPLICATIVE_OPERATOR,
	"%": MULTIPLICATIVE_OPERATOR,

	">=": COMPARISON_OPERATOR,
	"<=": COMPARISON_OPERATOR,
	">":  COMPARISON_OPERATOR,
	"<":  COMPARISON_OPERATOR,

	"||": LOGICAL_OPERATOR,
	"&&": LOGICAL_OPERATOR,

	"=": ASSIGNMENT_OPERATOR,
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
