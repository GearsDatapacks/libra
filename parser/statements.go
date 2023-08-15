package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseStatement() ast.Statement {
	var statement ast.Statement

	if p.isKeyword("var") || p.isKeyword("const") {
		statement = p.parseVariableDeclaration()
	} else {
		statement = p.parseExpressionStatement()
	}

	if !p.eof() && !p.next().LeadingNewline {
		p.error("Expected new line after statement", p.next())
	}

	return statement
}

func (p *parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression()
	return &ast.ExpressionStatement{
		Expression: expr,
		BaseNode:   &ast.BaseNode{Token: expr.GetToken()},
	}
}

func (p *parser) parseVariableDeclaration() ast.Statement {
	tok := p.consume()
	isConstant := tok.Value == "const"
	name := p.expect(
		token.IDENTIFIER,
		"ParseError: Expected identifier for variable declaration, got %q",
	).Value

	// TODO: add possibility for `var x string`

	p.expect(
		token.ASSIGNMENT_OPERATOR,
		"ParseError: Missing initialiser in variable declaration",
	)

	value := p.parseExpression()

	return &ast.VariableDeclaration{
		Constant: isConstant,
		Name:     name,
		BaseNode: &ast.BaseNode{Token: tok},
		Value:    value,
	}
}
