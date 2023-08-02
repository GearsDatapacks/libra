package lexer

import (
	"fmt"
	"log"
	"strings"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type lexer struct {
	code   []byte
	pos    int
	line   int
	column int
}

func New(code []byte) *lexer {
	return &lexer{code: code, line: 1, column: 1}
}

func (l *lexer) Tokenise() []token.Token {
	tokens := []token.Token{}

	for {
		nextToken := l.parseToken()
		tokens = append(tokens, nextToken)

		if nextToken.Type == token.EOF {
			break
		}
	}

	return tokens
}

func (l *lexer) parseToken() token.Token {
	l.skipWhitespace()

	if l.startsWith("//") {
		for !l.eof() && l.next() != '\n' {
			l.consume()
		}
	}

	l.skipWhitespace()

	leadingNewline := l.skipNewlines()

	if l.eof() {
		tok := l.createToken(token.EOF, []rune{'\u0000'}, leadingNewline)
		tok.Value = "EndOfFile"
		return tok
	}

	nextChar := l.next()

	if isNumeric(nextChar) {
		number := []rune{}
		for !l.eof() && isNumeric(l.next()) {
			number = append(number, l.consume())
		}
		return l.createToken(token.INTEGER, number, leadingNewline)
	} else if sym, ok := l.parseSymbol(); ok {
		sym.LeadingNewline = leadingNewline
		return sym
	} else if isAlphabetic(nextChar) {
		ident := []rune{}
		for !l.eof() && isAlphanumeric(l.next()) {
			ident = append(ident, l.consume())
		}
		return l.createToken(token.IDENTIFIER, ident, leadingNewline)
	} else {
		l.error(fmt.Sprintf("Unexpected token: %q", nextChar))
	}

	return token.Token{}
}

func (l *lexer) parseSymbol() (token.Token, bool) {
	for symbol, tokenType := range token.Symbols {
		if l.startsWith(symbol) {
			l.pos += len(symbol)
			tok := l.createToken(tokenType, []rune(symbol), false)
			return tok, true
		}
	}

	return token.Token{}, false
}

func (l *lexer) skipWhitespace() {
	for !l.eof() && isWhitespace(l.next()) {
		l.consume()
	}
}

func (l *lexer) skipNewlines() bool {
	if l.eof() || !isNewline(l.next()) {
		return false
	}

	for !l.eof() && isNewline(l.next()) {
		l.consume()
	}
	return true
}

func (l *lexer) createToken(tokenType token.Type, value []rune, leadingNewline bool) token.Token {
	return token.New(
		l.line,
		l.column,
		tokenType,
		value,
		leadingNewline,
	)
}

func (l *lexer) consume() rune {
	if l.eof() {
		l.error("Expected more charactes, got EOF")
	}

	nextByte := l.code[l.pos]
	l.pos++
	l.column++

	if nextByte == '\n' {
		l.line++
		l.column = 1
	}

	return rune(nextByte)
}

func (l *lexer) next() rune {
	if l.eof() {
		l.error("Expected more charactes, got EOF")
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

func (l *lexer) error(message string) {
	log.Fatalf("SyntaxError at line %d, column %d: %s", l.line, l.column, message)
}
