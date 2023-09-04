package token

import "github.com/gearsdatapacks/libra/utils"

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

	LEFT_PAREN   = "LEFT_PAREN"
	RIGHT_PAREN  = "RIGHT_PAREN"
	LEFT_BRACE   = "LEFT_BRACE"
	RIGHT_BRACE  = "RIGHT_BRACE"
	LEFT_SQUARE  = "LEFT_SQUARE"
	RIGHT_SQUARE = "RIGHT_SQUARE"
	COMMA        = "COMMA"
	SEMICOLON    = "SEMICOLON"

	ASSIGN     = "ASSIGN"
	ADD_ASSIGN = "ADD_ASSIGN"
	SUB_ASSIGN = "SUB_ASSIGN"
	MUL_ASSIGN = "MUL_ASSIGN"
	DIV_ASSIGN = "DIV_ASSIGN"
	MOD_ASSIGN = "MOD_ASSIGN"

	LOGICAL_AND = "LOGICAL_AND"
	LOGICAL_OR  = "LOGICAL_OR"

	LESS_THAN       = "LESS_THAN"
	GREATER_THAN    = "GREATER_THAN"
	LESS_THAN_EQ    = "LESS_THAN_EQ"
	GREATER_THAN_EQ = "GREATER_THAN_EQ"
	EQUAL           = "EQUAL"
	NOT_EQUAL       = "NOT_EQUAL"

	ADD      = "ADD"
	SUBTRACT = "SUBTRACT"

	MULTIPLY = "MULTIPLY"
	DIVIDE   = "DIVIDE"
	MODULO   = "MODULO"

	POWER = "POWER"

	INCREMENT = "INCREMENT"
	DECREMENT = "DECREMENT"

	LOGICAL_NOT = "LOGICAL_NOT"

	BITWISE_OR = "BITWISE_OR"
)

var Symbols = map[string]Type{
	"(": LEFT_PAREN,
	")": RIGHT_PAREN,
	"{": LEFT_BRACE,
	"}": RIGHT_BRACE,
	"[": LEFT_SQUARE,
	"]": RIGHT_SQUARE,
	",": COMMA,
	";": SEMICOLON,

	"+":  ADD,
	"-":  SUBTRACT,
	"*":  MULTIPLY,
	"/":  DIVIDE,
	"%":  MODULO,
	"**": POWER,

	"=":  ASSIGN,
	"+=": ADD_ASSIGN,
	"-=": SUB_ASSIGN,
	"*=": MUL_ASSIGN,
	"/=": DIV_ASSIGN,
	"%=": MOD_ASSIGN,

	"<":  LESS_THAN,
	"<=": LESS_THAN_EQ,
	">":  GREATER_THAN,
	">=": GREATER_THAN_EQ,
	"==": EQUAL,
	"!=": NOT_EQUAL,

	"|": BITWISE_OR,

	"||": LOGICAL_OR,
	"&&": LOGICAL_AND,

	"++": INCREMENT,
	"--": DECREMENT,
	"!":  LOGICAL_NOT,
}

var ComparisonOperator = []Type{
	LESS_THAN,
	LESS_THAN_EQ,
	GREATER_THAN,
	GREATER_THAN_EQ,
	EQUAL,
	NOT_EQUAL,
}

var AssignmentOperator = []Type{
	ASSIGN,
	ADD_ASSIGN,
	SUB_ASSIGN,
	MUL_ASSIGN,
	DIV_ASSIGN,
	MOD_ASSIGN,
}

var LogicalOperator = []Type{
	LOGICAL_AND,
	LOGICAL_OR,
}

var AdditiveOperator = []Type{ADD, SUBTRACT}

var MultiplicativeOperator = []Type{
	MULTIPLY,
	DIVIDE,
	MODULO,
}

var PrefixOperator = []Type{
	SUBTRACT,
	LOGICAL_NOT,
}

var PostfixOperator = []Type{
	INCREMENT,
	DECREMENT,
}

func (tokenType Type) Is(opGroup []Type) bool {
	return utils.Contains(opGroup, tokenType)
}

func (token Token) Is(opGroup []Type) bool {
	return utils.Contains(opGroup, token.Type)
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
