package parser

import (
	"fmt"
	"strings"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/utils"
)

type parser struct {
	tokens         []token.Token
	bracketLevel   int
	noBraces       bool
	requireNewline bool
	usedSymbols    []string
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

func (p *parser) expect(tokenType token.Type, fString string) (token.Token, error) {
	nextToken := p.consume()

	if nextToken.Type != tokenType {
		err := fString
		if strings.Contains(err, "%") {
			err = fmt.Sprintf(fString, nextToken.Value)
		}
		return token.Token{}, p.error(err, nextToken)
	}

	return nextToken, nil
}

func (p *parser) isKeyword(keyword string) bool {
	isKeyword := p.next().Type == token.IDENTIFIER && p.next().Value == keyword
	notUsed := !utils.Contains(p.usedSymbols, keyword)
	return isKeyword && notUsed
}

func (p *parser) expectKeyword(keyword string, fString string) (token.Token, error) {
	if p.isKeyword(keyword) {
		return p.consume(), nil
	}

	if strings.Contains(fString, "%") {
		fString = fmt.Sprintf(fString, p.next().Value)
	}
	return p.next(), p.error(fString, p.next())
}

// Only continue parsing if there is not a newline, or we're in an expression
func (p *parser) canContinue() bool {
	return !p.next().LeadingNewline || p.bracketLevel != 0
}

func (p *parser) error(message string, errorToken token.Token) error {
	return fmt.Errorf("SyntaxError at line %d, column %d: %s", errorToken.Line, errorToken.Column, message)
}

func (p *parser) needsNewline() bool {
	needsNewline := p.requireNewline
	p.requireNewline = false
	return needsNewline
}

func (p *parser) Parse(tokens []token.Token) (ast.Program, error) {
	p.tokens = tokens

	program := ast.Program{}
	p.usedSymbols = []string{}
	p.requireNewline = false

	for !p.eof() {
		nextStatement, err := p.parseStatement()
		if err != nil {
			return ast.Program{}, err
		}
		program.Body = append(program.Body, nextStatement)
	}

	return program, nil
}
