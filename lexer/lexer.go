package lexer

import (
	"bytes"
	"strconv"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/text"
)

type lexer struct {
	file        *text.SourceFile
	pos         int
	line        int
	col         int
	Diagnostics diagnostics.Manager
}

func New(file *text.SourceFile) *lexer {
	return &lexer{
		file:        file,
		pos:         0,
		line:        0,
		col:         0,
		Diagnostics: diagnostics.Manager{},
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
	nextToken := token.New(token.INVALID, "", l.getLocation(l.line, l.line, l.col, l.col))

	next := l.next()
	pos := l.pos

	if l.eof() {
		nextToken.Kind = token.EOF
		nextToken.Value = "\x00"
		return nextToken
	} else if next == '/' && l.peek(1) == '/' {
		nextToken.Kind = token.COMMENT
		nextToken.Value = l.parseLineComment()
	} else if next == '/' && l.peek(1) == '*' {
		nextToken.Kind = token.COMMENT
		nextToken.Value = l.parseBlockComment()
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
		l.Diagnostics.Report(diagnostics.InvalidCharacter(l.getLocation(l.line, l.line, l.col, l.col+1), next))
		l.consume()
	}

	nextToken.Location.Span.End = l.col
	nextToken.Location.Span.EndLine = l.line
	if nextToken.Value == "" {
		nextToken.Value = l.file.Text[pos:l.pos]
	}

	return nextToken
}

func (l *lexer) parseNumber() (token.Kind, string) {
	kind := token.INTEGER
	str := bytes.NewBuffer([]byte{})
	// TODO: 0x, 0b, etc.
	for isNumber(l.next()) || l.next() == '_' {
		if l.next() != '_' {
			str.WriteByte(l.next())
		}

		l.consume()
	}

	if l.next() == '.' && isNumber(l.peek(1)) {
		if l.peek(-1) == '_' {
			l.Diagnostics.Report(diagnostics.NumbersCannotEndWithSeparator(
				l.getLocation(l.line, l.line, l.col-1, l.col)))
		}

		kind = token.FLOAT
		str.WriteByte(l.next())
		l.consume()

		for isNumber(l.next()) || l.next() == '_' {
			if l.next() != '_' {
				str.WriteByte(l.next())
			}

			l.consume()
		}
	}

	if l.peek(-1) == '_' {
		l.Diagnostics.Report(diagnostics.NumbersCannotEndWithSeparator(
			l.getLocation(l.line, l.line, l.col-1, l.col)))
	}

	return kind, str.String()
}

func (l *lexer) parseString() string {
	startLine := l.line
	pos := l.col
	l.consume()
	str := bytes.NewBuffer([]byte{})

	for !l.eof() && l.next() != '"' {
		if l.next() == '\\' {
			l.consume()
			if l.next() != '\n' {
				char, ok := l.escape(l.next())
				if !ok {
					l.Diagnostics.Report(diagnostics.InvalidEscapeSequence(
						l.getLocation(l.line, l.line, l.col-1, l.col+1), l.next()))
				}
				str.WriteByte(char)
			}
		} else {
			str.WriteByte(l.next())
		}

		l.consume()
	}

	if l.eof() {
		l.Diagnostics.Report(diagnostics.UnterminatedString(l.getLocation(startLine, l.line, pos, l.col)))
	}

	l.consume()
	return str.String()
}

func (l *lexer) escape(c byte) (char byte, ok bool) {
	switch c {
	case '\\':
		char = '\\'
	case '"':
		char = '"'
	case 'a':
		char = '\a'
	case 'b':
		char = '\b'
	case 'f':
		char = '\f'
	case 'n':
		char = '\n'
	case 'r':
		char = '\r'
	case 't':
		char = '\t'
	case 'v':
		char = '\v'
	case '0':
		char = 0
	case 'x':
		if l.pos+3 >= len(l.file.Text) {
			l.Diagnostics.Report(diagnostics.ExpectedEscapeSequence(l.getLocation(l.line, l.line, l.col-1, l.col+1)))
			return 0, true
		}

		l.consume()
		nextTwoChars := string(l.next()) + string(l.peek(1))
		c, e := strconv.ParseUint(nextTwoChars, 16, 8)
		if e != nil {
			l.Diagnostics.Report(diagnostics.InvalidAsciiSequence(l.getLocation(l.line, l.line, l.col-2, l.col+2), nextTwoChars))
		} else {
			l.consume()
		}
		char = byte(c)
	case 'u':
		if l.pos+5 >= len(l.file.Text) {
			l.Diagnostics.Report(diagnostics.ExpectedEscapeSequence(l.getLocation(l.line, l.line, l.col-1, l.col+1)))
			return 0, true
		}

		l.consume()
		nextFourChars := string(l.next()) + string(l.peek(1)) + string(l.peek(2)) + string(l.peek(3))
		// TODO: support more than one byte utf8 sequences
		c, e := strconv.ParseUint(nextFourChars, 16, 8)
		if e != nil {
			l.Diagnostics.Report(diagnostics.InvalidUnicodeSequence(l.getLocation(l.line, l.line, l.col-2, l.col+4), nextFourChars))
		} else {
			l.consumeMany(3)
		}
		char = byte(c)
	default:
		return c, false
	}

	ok = true
	return
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
		} else if l.next() == '*' {
			kind = token.DOT_STAR
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

func (l *lexer) parseLineComment() string {
	l.consume()
	l.consume()
	l.skipWhitespace()
	var result bytes.Buffer
	for !l.eof() && l.next() != '\n' {
		result.WriteByte(l.consume())
	}

	return result.String()
}

func (l *lexer) parseBlockComment() string {
	startLine := l.line
	startCol := l.col
	l.consume()
	l.consume()
	l.skipWhitespace()
	nestLevel := 1
	var result bytes.Buffer
	for !l.eof() && nestLevel > 0 {
		if l.next() == '/' && l.peek(1) == '*' {
			nestLevel++
			result.WriteByte(l.consume())
			result.WriteByte(l.consume())
		} else if l.next() == '*' && l.peek(1) == '/' {
			nestLevel--
			result.WriteByte(l.consume())
			result.WriteByte(l.consume())
		} else {
			result.WriteByte(l.consume())
		}
	}

	if nestLevel > 0 {
		l.Diagnostics.Report(diagnostics.UnterminatedComment(l.getLocation(startLine, l.line, startCol, l.col)))
	}

	return result.String()
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
	if l.pos+n >= len(l.file.Text) {
		return 0
	}
	return l.file.Text[l.pos+n]
}

func (l *lexer) consume() byte {
	if l.pos >= len(l.file.Text) {
		return l.next()
	}
	next := l.next()
	l.col++
	l.pos++
	if next == '\n' {
		l.line++
		l.col = 0
	}
	return next
}

func (l *lexer) consumeMany(n int) {
	for range n {
		l.consume()
	}
}

func (l *lexer) eof() bool {
	return l.pos >= len(l.file.Text)
}

func (l *lexer) getLocation(startLine, endLine, col, end int) text.Location {
	span := text.NewSpan(startLine, endLine, col, end)
	return text.Location{
		Span: span,
		File: l.file,
	}
}
