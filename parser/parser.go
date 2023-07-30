package parser

import (
	"log"

	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/lexer/token"
)

type parser struct {
	tokens []token.Token
}

func New() *parser {
	return &parser{}
}

func (p *parser) eof() bool {
	return p.tokens[0].Type == token.EOF
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
		log.Fatalf(fString, nextToken.Value)
	}

	return nextToken
}

func (p *parser) Parse(tokens []token.Token) ast.Program {
	p.tokens = tokens

	program := ast.Program{}

	for !p.eof() {
		program.Body = append(program.Body, p.parseStatement())
	}

	return program
}
