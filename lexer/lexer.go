package lexer

import (
	"errors"
	"log"
	"strings"

	"github.com/gearsdatapacks/libra/token"
)

type lexer struct{
	code []byte
	pos int
}

func NewLexer(code []byte) *lexer {
	return &lexer{ code: code }
}

func (l *lexer) Tokenise() []token.Token {
	tokens := []token.Token{}

	for !l.eof() {
		if isWhitespace(l.peek()) {
			l.next()
			continue
		}

		tokens = append(tokens, l.parseToken())
	}

	tokens = append(tokens, l.createToken(token.EOF, '\u0000'))
	tokens[len(tokens)-1].Value = "EndOfFile"

	return tokens
}

func (l *lexer) parseToken() token.Token {
	nextChar := l.peek()

	if isNumeric(nextChar) {
		number := []rune{}
		for isNumeric(l.peek()) {
			number = append(number, l.next())
		}
		return l.createToken(token.NUMBER, number...)
	} else if sym, ok := l.parseSymbol(); ok {
		return sym
	} else {
		log.Fatalf("lexer: Unexpected token: %q", nextChar)
	}

	return token.Token{}
}

func (l *lexer) parseSymbol() (token.Token, bool) {
	for symbol, tokenType := range token.Symbols {
		if l.startsWith(symbol) {
			l.pos += len(symbol)
			tok := l.createToken(tokenType, []rune(symbol)...)
			return tok, true
		}
	}

	return token.Token{}, false
}

func (l *lexer) createToken(tokenType token.Type, value ...rune) token.Token {
	return token.New(
		l.pos - len(value),
		l.pos,
		tokenType,
		value,
	)
}

func (l *lexer) next() rune {
	if l.pos > len(l.code) {
		log.Fatal(errors.New("lexer: expected more charactes, got EOF"))
	}

	nextByte := l.code[l.pos]
	l.pos++
	return rune(nextByte)
}

func (l *lexer) peek() rune {
	if l.pos > len(l.code) {
		log.Fatal(errors.New("lexer: expected more charactes, got EOF"))
	}
	return rune(l.code[l.pos])
}

func (l *lexer) startsWith(prefix string) bool {
	workingString := string(l.code[l.pos:])
	return strings.HasPrefix(workingString, prefix)
} 

func (l *lexer) eof() bool {
	return l.pos >= len(l.code)
}
