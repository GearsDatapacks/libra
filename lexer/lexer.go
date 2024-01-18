package lexer

import "github.com/gearsdatapacks/libra/lexer/token"

type lexer struct {
	src  string
	pos  int
	line int
	col  int
}

func New(src string) *lexer {
	return &lexer{
		src:  src,
		pos:  0,
		line: 0,
		col:  0,
	}
}

func (l *lexer) Tokenise() []token.Token {
	tokens := []token.Token{}

	for {
		nextToken := l.nextToken()
		tokens = append(tokens, nextToken)
		if nextToken.Kind == token.EOF {
			break
		}
	}

	return tokens
}

func (l *lexer) nextToken() token.Token {
	nextToken := token.New(token.INVALID, "", token.NewSpan(l.pos, l.pos, l.line, l.col))

	next := l.next()
	pos := l.pos

	if next == 0 {
		nextToken.Kind = token.EOF
		nextToken.Value = "\x00"
    return nextToken
	} else if kind, ok := l.parsePunctuation(); ok {
		nextToken.Kind = kind
	} else {
		l.consume()
	}

	nextToken.Value = l.src[pos:l.pos]

	return nextToken
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
  case '\n':
    kind = token.NEWLINE
    l.consume()

	default:
		return token.INVALID, false
	}

	return kind, true
}

func (l *lexer) next() byte {
	if l.pos >= len(l.src) {
		return 0
	}
	return l.src[l.pos]
}

func (l *lexer) consume() {
	l.pos++
}
