package lexer

import (
	"fmt"
	"log"
	"strings"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type lexer struct {
	code      []byte
	pos       int
	oldLine   int
	oldColumn int
	line      int
	column    int
}

func New(code []byte) *lexer {
	return &lexer{
		code: code,
		line: 1, oldLine: 1,
		column: 1, oldColumn: 1,
	}
}

func (l *lexer) Tokenise() ([]token.Token, error) {
	tokens := []token.Token{}

	for {
		nextToken, err := l.parseToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, nextToken)

		if nextToken.Type == token.EOF {
			break
		}
	}

	return tokens, nil
}

func (l *lexer) parseToken() (token.Token, error) {
	leadingNewline := l.skip()

	l.oldLine = l.line
	l.oldColumn = l.column

	if l.eof() {
		tok := l.createToken(token.EOF, []rune{'\u0000'}, leadingNewline)
		tok.Value = "EndOfFile"
		return tok, nil
	}

	nextChar := l.next()

	if isNumeric(nextChar, 10) {
		return l.parseNumber(leadingNewline), nil
	} else if sym, ok := l.parseSymbol(); ok {
		sym.LeadingNewline = leadingNewline
		return sym, nil
	} else if isAlphabetic(nextChar) {
		ident := []rune{}
		for !l.eof() && isAlphanumeric(l.next()) {
			ident = append(ident, l.consume())
		}
		return l.createToken(token.IDENTIFIER, ident, leadingNewline), nil
	} else if nextChar == '"' {
		stringValue := []rune{}
		l.consume()
		for l.next() != '"' {
			stringValue = append(stringValue, l.consume())
		}
		l.consume()
		return l.createToken(token.STRING, stringValue, leadingNewline), nil
	} else {
		return token.Token{}, l.error(fmt.Sprintf("Unexpected token: %q", nextChar))
	}
}

func (l *lexer) parseSymbol() (token.Token, bool) {
	// Prevent prefering "+" over "+="
	longestSymbol := ""
	var symbolType token.Type
	for symbol, tokenType := range token.Symbols {
		if l.startsWith(symbol) && len(symbol) > len(longestSymbol) {
			longestSymbol = symbol
			symbolType = tokenType
		}
	}

	if longestSymbol != "" {
		l.pos += len(longestSymbol)
		tok := l.createToken(symbolType, []rune(longestSymbol), false)
		return tok, true
	}

	return token.Token{}, false
}

func (l *lexer) parseNumber(leadingNewline bool) token.Token {
	number := []rune{}
	var radix int32 = 10
	if l.next() == '0' {
		r := GetRadix(l.peek(1))
		if r != -1 {
			radix = r
			number = append(number, l.consume(), l.consume())
		}
	}

	if !isNumeric(l.next(), radix) {
		log.Fatal("Invalid token: empty integer literal")
	}

	for !l.eof() && isNumeric(l.next(), radix) || l.next() == '_' {
		if l.next() != '_' {
			number = append(number, l.next())
		}
		l.consume()
	}

	if l.next() == '.' && isNumeric(l.peek(1), 10) && radix == 10 {
		number = append(number, l.consume())
		for !l.eof() && isNumeric(l.next(), radix) || l.next() == '_' {
			if l.next() != '_' {
				number = append(number, l.next())
			}
			l.consume()
		}

		return l.createToken(token.FLOAT, number, leadingNewline)
	}

	return l.createToken(token.INTEGER, number, leadingNewline)
}

func (l *lexer) skip() bool {
	leadingNewline := false

	for l.isSkippable() {
		leadingNewline = l.skipWhitespace()

		if l.next() == '#' {
			for !l.eof() && l.next() != '\n' {
				l.consume()
			}
		}

		if l.startsWith("//") {
			for !l.eof() && !l.startsWith("\\\\") {
				l.consume()
			}
			l.consume()
			l.consume()
		}
	}

	return leadingNewline
}

func (l *lexer) isSkippable() bool {
	if l.eof() {
		return false
	}
	return isWhitespace(l.next()) || isNewline(l.next()) || l.next() == '#' || l.startsWith("//")
}

func (l *lexer) skipWhitespace() (leadingNewline bool) {
	for isNewline(l.next()) {
		l.consume()
		leadingNewline = true
	}

	for !l.eof() && isWhitespace(l.next()) {
		l.consume()
	}

	for isNewline(l.next()) {
		l.consume()
		leadingNewline = true
	}

	return leadingNewline
}

func (l *lexer) createToken(tokenType token.Type, value []rune, leadingNewline bool) token.Token {
	tok := token.New(
		l.oldLine,
		l.oldColumn,
		tokenType,
		value,
		leadingNewline,
	)

	l.oldLine = l.line
	l.oldColumn = l.column

	return tok
}

func (l *lexer) consume() rune {
	if l.eof() {
		return '\u0000'
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

func (l *lexer) peek(offset int) rune {
	if l.pos+offset >= len(l.code) {
		return '\u0000'
	}
	return rune(l.code[l.pos+offset])
}

func (l *lexer) next() rune {
	return l.peek(0)
}

func (l *lexer) startsWith(prefix string) bool {
	workingString := string(l.code[l.pos:])
	return strings.HasPrefix(workingString, prefix)
}

func (l *lexer) eof() bool {
	return l.pos >= len(l.code)
}

func (l *lexer) error(message string) error {
	return fmt.Errorf("SyntaxError at line %d, column %d: %s", l.line, l.column, message)
}
