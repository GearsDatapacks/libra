package token

type Kind int

type Token struct {
	Kind  Kind
	Value string
	Span  Span
}

type Span struct {
	Start int
	End   int
	Line  int
	Col   int
}

const (
	EOF Kind = iota
	INVALID

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

func New(kind Kind, value string, span Span) Token {
	return Token{
		Kind:  kind,
		Value: value,
		Span:  span,
	}
}

func NewSpan(start, end, line, col int) Span {
	return Span{
		Start: start,
		End:   end,
		Line:  line,
		Col:   col,
	}
}
