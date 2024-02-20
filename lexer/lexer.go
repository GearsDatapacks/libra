package lexer

import (
	"bytes"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
)

type lexer struct {
	src         string
	pos         int
	line        int
	col         int
	Diagnostics diagnostics.Manager
}

func New(src string, file string) *lexer {
	return &lexer{
		src:         src,
		pos:         0,
		line:        0,
		col:         0,
		Diagnostics: diagnostics.New(file, src),
	}
}

func (l *lexer) Tokenise() []token.Token {
	tokens := []token.Token{}

	for {
		nextToken := l.nextToken()
		if nextToken.Kind == token.INVALID {
			continue
		}

		tokens = append(tokens, nextToken)

		if nextToken.Kind == token.EOF {
			break
		}
	}

	return tokens
}

func (l *lexer) nextToken() token.Token {
	l.skipWhitespace()
	nextToken := token.New(token.INVALID, "", token.NewSpan(l.line, l.col, l.col))

	next := l.next()
	pos := l.pos

	if l.eof() {
		nextToken.Kind = token.EOF
		nextToken.Value = "\x00"
		return nextToken
	} else if kind, ok := l.parsePunctuation(); ok {
		nextToken.Kind = kind
	} else if isNumber(next) {
		nextToken.Kind, nextToken.Value = l.parseNumber()
	} else if isIdentifierStart(next) {
		nextToken.Kind = token.IDENTIFIER
		for isIdentifierMiddle(l.next()) {
			l.consume()
		}
	} else if next == '"' {
		nextToken.Kind = token.STRING
		nextToken.Value = l.parseString()
	} else {
		l.Diagnostics.ReportInvalidCharacter(token.NewSpan(l.line, l.col, l.col+1), next)
		l.consume()
	}

	nextToken.Span.End = l.col
	if nextToken.Value == "" {
		nextToken.Value = l.src[pos:l.pos]
	}

	return nextToken
}

func (l *lexer) parseNumber() (token.Kind, string) {
	kind := token.INTEGER
	text := bytes.NewBuffer([]byte{})
	// TODO: 0x, 0b, etc.
	for isNumber(l.next()) || l.next() == '_' {
		if l.next() != '_' {
			text.WriteByte(l.next())
		}

		l.consume()
	}

	if l.next() == '.' && isNumber(l.peek(1)) {
		if l.peek(-1) == '_' {
			l.Diagnostics.ReportNumbersCannotEndWithSeparator(
				token.NewSpan(l.line, l.col-1, l.col))
		}

		kind = token.FLOAT
		text.WriteByte(l.next())
		l.consume()

		for isNumber(l.next()) || l.next() == '_' {
			if l.next() != '_' {
				text.WriteByte(l.next())
			}

			l.consume()
		}
	}

	if l.peek(-1) == '_' {
		l.Diagnostics.ReportNumbersCannotEndWithSeparator(
			token.NewSpan(l.line, l.col-1, l.col))
	}

	return kind, text.String()
}

func (l *lexer) parseString() string {
	pos := l.col
	l.consume()
	text := bytes.NewBuffer([]byte{})

	for !l.eof() && l.next() != '"' {
		if l.next() == '\\' {
			l.consume()
			char, ok := escape(l.next())
			if !ok {
				l.Diagnostics.ReportInvalidEscapeSequence(
					token.NewSpan(l.line, l.col-1, l.col+1), l.next())
			}
			text.WriteByte(char)
		} else {
			text.WriteByte(l.next())
		}

		l.consume()
	}

	if l.eof() {
		l.Diagnostics.ReportUnterminatedString(token.NewSpan(l.line, pos, l.col))
	}

	l.consume()
	return text.String()
}

func (l *lexer) parsePunctuation() (token.Kind, bool) {
	kind := token.INVALID
	switch l.next() {
	case '+':
		kind = token.PLUS
		l.consume()
		if l.next() == '+' {
			kind = token.DOUBLE_PLUS
			l.consume()
		} else if l.next() == '=' {
			kind = token.PLUS_EQUALS
			l.consume()
		}
	case '(':
		kind = token.LEFT_PAREN
		l.consume()
	case ')':
		kind = token.RIGHT_PAREN
		l.consume()
	case '{':
		kind = token.LEFT_BRACE
		l.consume()
	case '}':
		kind = token.RIGHT_BRACE
		l.consume()
	case '[':
		kind = token.LEFT_SQUARE
		l.consume()
	case ']':
		kind = token.RIGHT_SQUARE
		l.consume()
	case ',':
		kind = token.COMMA
		l.consume()
	case '.':
		kind = token.DOT
		l.consume()
		if l.next() == '.' {
			kind = token.DOUBLE_DOT
			l.consume()
		}
	case ':':
		kind = token.COLON
		l.consume()
	case '?':
		kind = token.QUESTION
		l.consume()
	case '=':
		kind = token.EQUALS
		l.consume()
		if l.next() == '=' {
			kind = token.DOUBLE_EQUALS
			l.consume()
		}
	case '<':
		kind = token.LEFT_ANGLE
		l.consume()
		if l.next() == '=' {
			kind = token.LEFT_ANGLE_EQUALS
			l.consume()
		} else if l.next() == '<' {
			kind = token.DOUBLE_LEFT_ANGLE
			l.consume()
		}
	case '>':
		kind = token.RIGHT_ANGLE
		l.consume()
		if l.next() == '=' {
			kind = token.RIGHT_ANGLE_EQUALS
			l.consume()
		} else if l.next() == '>' {
			kind = token.DOUBLE_RIGHT_ANGLE
			l.consume()
		}
	case '-':
		kind = token.MINUS
		l.consume()
		if l.next() == '=' {
			kind = token.MINUS_EQUALS
			l.consume()
		} else if l.next() == '-' {
			kind = token.DOUBLE_MINUS
			l.consume()
		} else if l.next() == '>' {
			kind = token.ARROW
			l.consume()
		}
	case '*':
		kind = token.STAR
		l.consume()
		if l.next() == '=' {
			kind = token.STAR_EQUALS
			l.consume()
		} else if l.next() == '*' {
			kind = token.DOUBLE_STAR
			l.consume()
		}
	case '/':
		kind = token.SLASH
		l.consume()
		if l.next() == '=' {
			kind = token.SLASH_EQUALS
			l.consume()
		}
	case '%':
		kind = token.PERCENT
		l.consume()
		if l.next() == '=' {
			kind = token.PERCENT_EQUALS
			l.consume()
		}
	case '!':
		kind = token.BANG
		l.consume()
		if l.next() == '=' {
			kind = token.BANG_EQUALS
			l.consume()
		}
	case '|':
		kind = token.PIPE
		l.consume()
		if l.next() == '|' {
			kind = token.DOUBLE_PIPE
			l.consume()
		}
	case '&':
		kind = token.AMPERSAND
		l.consume()
		if l.next() == '&' {
			kind = token.DOUBLE_AMPERSAND
			l.consume()
		}
	case '~':
		kind = token.TILDE
		l.consume()
	case ';':
		kind = token.SEMICOLON
		l.consume()
	case '\n':
		kind = token.NEWLINE
		l.consume()

	default:
		return token.INVALID, false
	}

	return kind, true
}

func (l *lexer) skipWhitespace() {
	for isWhitespace(l.next()) {
		l.consume()
	}
}

func (l *lexer) next() byte {
	return l.peek(0)
}

func (l *lexer) peek(n int) byte {
	if l.pos+n >= len(l.src) {
		return 0
	}
	return l.src[l.pos+n]
}

func (l *lexer) consume() {
	next := l.next()
	l.col++
	l.pos++
	if next == '\n' {
		l.line++
		l.col = 0
	}
}

func (l *lexer) eof() bool {
	return l.pos >= len(l.src)
}
