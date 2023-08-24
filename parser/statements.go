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
	} else if p.isKeyword("return") {
		statement = p.parseReturnStatement()
	} else {
		statement = p.parseExpressionStatement()
	}

	if !p.eof() && !p.next().LeadingNewline {
		p.error(fmt.Sprintf("Expected new line after statement, got %q", p.next().Value), p.next())
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
	var dataType ast.TypeExpression = &ast.InferType{}
	
	if p.canContinue() && p.next().Type != token.ASSIGNMENT_OPERATOR {
		dataType = p.parseType()
	}

	// TODO: add possibility for `var x string`

	if !p.canContinue() || p.next().Type != token.ASSIGNMENT_OPERATOR {
		if isConstant {
			p.error(fmt.Sprintf("Cannot leave constant %q uninitialised", name), p.next())
		}

		if dataType.Type() == "Infer" {
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

	p.usedSymbols = append(p.usedSymbols, name)

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
	p.usedSymbols = append(p.usedSymbols, name)

	parameters := p.parseParameterList()

	var returnType ast.TypeExpression = &ast.VoidType{}

	if p.next().Type != token.LEFT_BRACE {
		returnType = p.parseType()
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

func (p *parser) parseReturnStatement() ast.Statement {
	token := p.consume()

	var value ast.Expression = &ast.VoidValue{}

	if p.canContinue() {
		value = p.parseExpression()
	}

	return &ast.ReturnStatement{
		Value: value,
		BaseNode: &ast.BaseNode{Token: token},
	}
}
