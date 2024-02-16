package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func parseDelemitedList[Elem any](p *parser, delim token.Kind, elemFn func() Elem) ( result []Elem, delimToken token.Token) {
	result = []Elem{}

	for !p.eof() && p.next().Kind != delim {
		result = append(result, elemFn())

		if p.next().Kind == token.COMMA{
			p.consume()
		} else {
			break
		}
	}

	delimToken = p.expect(delim)

	return result, delimToken
}

func (p *parser) parseOptionalTypeAnnotation() *ast.TypeAnnotation {
	if p.next().Kind != token.COLON {
		return nil
	}

	colon := p.consume()
	ty := p.parseType()

	return &ast.TypeAnnotation{
		Colon: colon,
		Type:  ty,
	}
}
