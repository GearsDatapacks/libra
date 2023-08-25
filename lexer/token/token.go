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
	STRING     = "STRING"
	IDENTIFIER = "IDENTIFIER"

	LEFT_PAREN  = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	LEFT_BRACE  = "LEFT_BRACE"
	RIGHT_BRACE = "RIGHT_BRACE"
	COMMA       = "COMMA"
	SEMICOLON   = "SEMICOLON"

	ASSIGNMENT_OPERATOR     = "ASSIGNMENT_OPERATOR"
	LOGICAL_OPERATOR        = "LOGICAL_OPERATOR"
	COMPARISON_OPERATOR     = "COMPARISON_OPERATOR"
	ADDITIVE_OPERATOR       = "ADDITIVE_OPERATOR"
	MULTIPLICATIVE_OPERATOR = "MULTIPLICATIVE_OPERATOR"
	EXPONENTIAL_OPERATOR    = "EXPONENTIAL_OPERATOR"
	BITWISE_OR              = "BITWISE_OR"
)

var Symbols = map[string]Type{
	"(": LEFT_PAREN,
	")": RIGHT_PAREN,
	"{": LEFT_BRACE,
	"}": RIGHT_BRACE,
	",": COMMA,
	";": SEMICOLON,

	"+":  ADDITIVE_OPERATOR,
	"-":  ADDITIVE_OPERATOR,
	"*":  MULTIPLICATIVE_OPERATOR,
	"/":  MULTIPLICATIVE_OPERATOR,
	"%":  MULTIPLICATIVE_OPERATOR,
	"**": EXPONENTIAL_OPERATOR,

	"=":  ASSIGNMENT_OPERATOR,
	"+=": ASSIGNMENT_OPERATOR,
	"*=": ASSIGNMENT_OPERATOR,
	"/=": ASSIGNMENT_OPERATOR,
	"%=": ASSIGNMENT_OPERATOR,

	"<":  COMPARISON_OPERATOR,
	"<=": COMPARISON_OPERATOR,
	">":  COMPARISON_OPERATOR,
	">=": COMPARISON_OPERATOR,
	"==": COMPARISON_OPERATOR,
	"!=": COMPARISON_OPERATOR,

	"|": BITWISE_OR,

	"||": LOGICAL_OPERATOR,
	"&&": LOGICAL_OPERATOR,
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
