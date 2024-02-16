package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseType() ast.TypeExpression {
	return p.parseTypeName()
}

func (p *parser) parseTypeName() ast.TypeExpression {
	name := p.expect(token.IDENTIFIER)

	return &ast.TypeName{
		Name: name,
	}
}
