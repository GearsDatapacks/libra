package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseType() ast.TypeExpression {
	ty := p.parsePrimaryType()

	if p.next().Kind == token.PIPE {
		types := []ast.TypeExpression{ty}
		
		for p.canContinue() && p.next().Kind == token.PIPE {
			p.consume()
			types = append(types, p.parsePrimaryType())
		}

		ty = &ast.Union{
			Types: types,
		}
	}

	return ty
}

func (p *parser) parsePrimaryType() ast.TypeExpression {
	return p.parseTypeName()
}

func (p *parser) parseTypeName() ast.TypeExpression {
	name := p.expect(token.IDENTIFIER)

	return &ast.TypeName{
		Name: name,
	}
}
