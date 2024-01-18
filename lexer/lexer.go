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

  switch l.next() {
  case 0:
    nextToken.Kind = token.EOF
    nextToken.Value = l.consume()

  default:
    nextToken.Value = l.consume()
  }

  return nextToken
}

func (l *lexer) next() byte {
  if l.pos >= len(l.src) {
    return 0
  }
  return l.src[l.pos]
}

func (l *lexer) consume() string {
  next := l.next()
  l.pos++
  return string(next)
}

