package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseStatement() ast.Statement {
	if p.isKeyword("let") || p.isKeyword("const") {
		return p.parseVariableDeclaration()
	} else {
		return p.parseExpressionStatement()
	}
}

func (p *parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression()
	return &ast.ExpressionStatement{
		Expression: expr,
		BaseNode: &ast.BaseNode{Token: expr.GetToken()},
	}
}

func (p *parser) parseVariableDeclaration() ast.Statement {
	tok := p.consume()
	isConstant := tok.Value == "const"
	name := p.expect(
		token.IDENTIFIER,
		"ParseError: Expected identifier for variable declaration, got %q",
	).Value

	// TODO: add possibility for `let x string`

	p.expect(
		token.EQUALS,
		"ParseError: Missing initialiser in variable declaration",
	)

	value := p.parseExpression()

	return &ast.VariableDeclaration{
		Constant: isConstant,
		Name: name,
		BaseNode: &ast.BaseNode{Token: tok},
		Value: &value,
	}
}
