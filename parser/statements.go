package parser

import (
	"fmt"

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
	dataType := ""
	
	if p.canContinue() && p.next().Type == token.IDENTIFIER {
		dataType = p.consume().Value
	}

	// TODO: add possibility for `var x string`

	if !p.canContinue() || p.next().Type != token.ASSIGNMENT_OPERATOR {
		if isConstant {
			p.error(fmt.Sprintf("Cannot leave constant %q uninitialised", name), p.next())
		}

		if dataType == "" {
			p.error(fmt.Sprintf("Cannot declare uninitialised variable %q without type annotation", name), p.next())
		}

		return &ast.VariableDeclaration{
			Constant: isConstant,
			Name:     name,
			BaseNode: &ast.BaseNode{Token: tok},
			Value:    nil,
			DataType: dataType,
		}
	}

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
		DataType: dataType,
	}
}
