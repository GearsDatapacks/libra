package token

import "github.com/gearsdatapacks/libra/utils"

type Type int

type Token struct {
	Type           Type
	Value          string
	Line           int
	Column         int
	LeadingNewline bool
}

const (
	EOF = iota

	INTEGER
	FLOAT
	STRING
	IDENTIFIER

	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_SQUARE
	RIGHT_SQUARE
	COMMA
	DOT
	// SEMICOLON
	COLON
	QUESTION

	EQUALS
	PLUS_EQUALS
	MINUS_EQUALS
	STAR_EQUALS
	SLAHS_EQUALS
	PERCENT_EQUALS

	DOUBLE_AMPERSAND
	DOUBLE_PIPE

	LEFT_ANGLE
	RIGHT_ANGLE
	LEFT_ANGLE_EQUALS
	RIGHT_ANGLE_EQUALS
	DOUBLE_EQUALS
	BANG_EQUALS

	DOUBLE_LEFT_ANGLE
	DOUBLE_RIGHT_ANGLE

	PLUS
	MINUS

	STAR
	SLASH
	PERCENT

	DOUBLE_STAR

	DOUBLE_PLUS
	DOUBLE_MINUS

	BANG

	PIPE
	AMPERSAND
	ARROW
)

var Symbols = map[string]Type{
	"(": LEFT_PAREN,
	")": RIGHT_PAREN,
	"{": LEFT_BRACE,
	"}": RIGHT_BRACE,
	"[": LEFT_SQUARE,
	"]": RIGHT_SQUARE,
	",": COMMA,
	".": DOT,
	// ";": SEMICOLON,
	":": COLON,
	"?": QUESTION,

	"+":  PLUS,
	"-":  MINUS,
	"*":  STAR,
	"/":  SLASH,
	"%":  PERCENT,
	"**": DOUBLE_STAR,

	"=":  EQUALS,
	"+=": PLUS_EQUALS,
	"-=": MINUS_EQUALS,
	"*=": STAR_EQUALS,
	"/=": SLAHS_EQUALS,
	"%=": PERCENT_EQUALS,

	"<<": DOUBLE_LEFT_ANGLE,
	">>": DOUBLE_RIGHT_ANGLE,

	"<":  LEFT_ANGLE,
	"<=": LEFT_ANGLE_EQUALS,
	">":  RIGHT_ANGLE,
	">=": RIGHT_ANGLE_EQUALS,
	"==": DOUBLE_EQUALS,
	"!=": BANG_EQUALS,

	"|": PIPE,
	"&": AMPERSAND,

	"||": DOUBLE_PIPE,
	"&&": DOUBLE_AMPERSAND,

	"++": DOUBLE_PLUS,
	"--": DOUBLE_MINUS,
	"!":  BANG,
	"->": ARROW,
}

var AssignmentOperator = []Type{
	EQUALS,
	PLUS_EQUALS,
	MINUS_EQUALS,
	STAR_EQUALS,
	SLAHS_EQUALS,
	PERCENT_EQUALS,
}

var PrefixOperator = []Type{
	MINUS,
	BANG,
	AMPERSAND,
	STAR,
}

var PostfixOperator = []Type{
	DOUBLE_PLUS,
	DOUBLE_MINUS,
	QUESTION,
	BANG,
}

var BinOpInfo = map[Type]struct {
	Precedence       int
	RightAssociative bool
}{
	DOUBLE_AMPERSAND: {Precedence: 1},
	DOUBLE_PIPE:  {Precedence: 1},

	LEFT_ANGLE:       {Precedence: 2},
	LEFT_ANGLE_EQUALS:    {Precedence: 2},
	RIGHT_ANGLE:    {Precedence: 2},
	RIGHT_ANGLE_EQUALS: {Precedence: 2},
	DOUBLE_EQUALS:           {Precedence: 2},
	BANG_EQUALS:       {Precedence: 2},

	DOUBLE_LEFT_ANGLE:  {Precedence: 3},
	DOUBLE_RIGHT_ANGLE: {Precedence: 3, RightAssociative: true},

	PLUS:  {Precedence: 4},
	MINUS: {Precedence: 4},

	STAR:    {Precedence: 5},
	SLASH:   {Precedence: 5},
	PERCENT: {Precedence: 5},

	DOUBLE_STAR: {Precedence: 6, RightAssociative: true},
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
