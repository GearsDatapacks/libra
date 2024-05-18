package token

import "github.com/gearsdatapacks/libra/text"

type Kind int

type Token struct {
	Kind     Kind
	Value    string
	Location text.Location
}

const (
	EOF Kind = iota
	INVALID
	NEWLINE

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
	DOUBLE_DOT

	BANG

	PIPE
	AMPERSAND
	ARROW
	TILDE
	SEMICOLON

	COMMENT

	PLUS_EQUALS =  PLUS + EQUALS
	MINUS_EQUALS = MINUS + EQUALS
	STAR_EQUALS = STAR + EQUALS
	SLASH_EQUALS = SLASH + EQUALS
	PERCENT_EQUALS = PERCENT + EQUALS
)

func New(kind Kind, value string, span text.Location) Token {
	return Token{
		Kind:     kind,
		Value:    value,
		Location: span,
	}
}

func (kind Kind) String() string {
	switch kind {
	case LEFT_PAREN:
		return "`(`"
	case RIGHT_PAREN:
		return "`)`"
	case LEFT_BRACE:
		return "`{`"
	case RIGHT_BRACE:
		return "`}`"
	case LEFT_SQUARE:
		return "`[`"
	case RIGHT_SQUARE:
		return "`]`"
	case COMMA:
		return "`,`"
	case DOT:
		return "`.`"
	case COLON:
		return "`:`"
	case QUESTION:
		return "`?`"
	case EQUALS:
		return "`=`"
	case PLUS_EQUALS:
		return "`+=`"
	case MINUS_EQUALS:
		return "`-=`"
	case STAR_EQUALS:
		return "`*=`"
	case SLASH_EQUALS:
		return "`/=`"
	case PERCENT_EQUALS:
		return "`%=`"
	case DOUBLE_AMPERSAND:
		return "`&&`"
	case DOUBLE_PIPE:
		return "`||`"
	case LEFT_ANGLE:
		return "`<`"
	case RIGHT_ANGLE:
		return "`>`"
	case LEFT_ANGLE_EQUALS:
		return "`<=`"
	case RIGHT_ANGLE_EQUALS:
		return "`>=`"
	case DOUBLE_EQUALS:
		return "`==`"
	case BANG_EQUALS:
		return "`!=`"
	case DOUBLE_LEFT_ANGLE:
		return "`<<`"
	case DOUBLE_RIGHT_ANGLE:
		return "`>>`"
	case PLUS:
		return "`+`"
	case MINUS:
		return "`-`"
	case STAR:
		return "`*`"
	case SLASH:
		return "`/`"
	case PERCENT:
		return "`%`"
	case DOUBLE_STAR:
		return "`**`"
	case DOUBLE_PLUS:
		return "`++`"
	case DOUBLE_MINUS:
		return "`--`"
	case DOUBLE_DOT:
		return "`..`"
	case BANG:
		return "`!`"
	case PIPE:
		return "`|`"
	case ARROW:
		return "`->`"
	case AMPERSAND:
		return "`&`"
	case TILDE:
		return "`~`"
	case SEMICOLON:
		return "`;`"
	case EOF:
		return "<Eof>"
	case INVALID:
		return "?"
	case NEWLINE:
		return "<newline>"
	case INTEGER:
		return "integer"
	case FLOAT:
		return "float"
	case STRING:
		return "string"
	case IDENTIFIER:
		return "identifier"
	case COMMENT:
		return "comment"
	default:
		return "?"
	}
}
