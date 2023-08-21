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
	} else if p.isKeyword("function") {
		statement = p.parseFunctionDeclaration()
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

func (p *parser) parseFunctionDeclaration() ast.Statement {
	tok := p.consume()

	name := p.expect(token.IDENTIFIER, "Expected function name, got %q").Value

	parameters := p.parseParameterList()

	returnType := ""

	if p.next().Type == token.IDENTIFIER {
		returnType = p.consume().Value
	}

	code := p.parseCodeBlock()

	return &ast.FunctionDeclaration{
		Name: name,
		Parameters: parameters,
		Body: code,
		ReturnType: returnType,
		BaseNode: &ast.BaseNode{Token: tok},
	}
}
