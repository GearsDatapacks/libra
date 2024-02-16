package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseTopLevelStatement() ast.Statement {
	switch p.next().Kind {
	default:
		return p.parseStatement()
	}
}

func (p *parser) parseStatement() ast.Statement {
	if p.isKeyword("const") || p.isKeyword("let") || p.isKeyword("mut") {
		return p.parseVariableDeclaration()
	}

	return &ast.ExpressionStatement{
		Expression: p.parseExpression(),
	}
}

func (p *parser) parseVariableDeclaration() ast.Statement {
	keyword := p.consume()
	identifier := p.delcareIdentifier()

	typeAnnotation := p.parseOptionalTypeAnnotation()

	equals := p.expect(token.EQUALS)
	value := p.parseExpression()

	return &ast.VariableDeclaration{
		Keyword:    keyword,
		Identifier: identifier,
		Type:       typeAnnotation,
		Equals:     equals,
		Value:      value,
	}
}
