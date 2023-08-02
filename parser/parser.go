package parser

import (
	"fmt"
	"log"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

type parser struct {
	tokens       []token.Token
	bracketLevel int
}

func New() *parser {
	return &parser{}
}

func (p *parser) eof() bool {
	return len(p.tokens) == 0 || p.tokens[0].Type == token.EOF
}

func (p *parser) next() token.Token {
	return p.tokens[0]
}

func (p *parser) consume() token.Token {
	nextToken := p.tokens[0]
	p.tokens = p.tokens[1:]
	return nextToken
}

func (p *parser) expect(tokenType token.Type, fString string) token.Token {
	nextToken := p.consume()

	if nextToken.Type != tokenType {
		p.error(fmt.Sprintf(fString, nextToken.Value), nextToken)
	}

	return nextToken
}

func (p *parser) isKeyword(keyword string) bool {
	return p.next().Type == token.IDENTIFIER && p.next().Value == keyword
}

/*
func (p *parser) expectKeyword(keyword string, fString string) token.Token {
	if !p.isKeyword(keyword) {
		log.Fatalf(fString, p.next().Value)
	}

	return p.consume()
} */

// Only continue parsing if there is not a newline, or we're in an expression
func (p *parser) canContinue() bool {
	return !p.next().LeadingNewline || p.bracketLevel != 0
}

func (p *parser) error(message string, errorToken token.Token) {
	log.Fatalf("SyntaxError at line %d, column %d: %s", errorToken.Line, errorToken.Column, message)
}

func (p *parser) Parse(tokens []token.Token) ast.Program {
	p.tokens = tokens

	program := ast.Program{}

	for !p.eof() {
		program.Body = append(program.Body, p.parseStatement())
	}

	return program
}
