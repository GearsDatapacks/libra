package parser

import (
	"fmt"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseStatement(inline ...bool) ast.Statement {
	var statement ast.Statement

	if p.isKeyword("var") || p.isKeyword("const") {
		statement = p.parseVariableDeclaration()
	} else if p.isKeyword("function") {
		statement = p.parseFunctionDeclaration()
	} else if p.isKeyword("return") {
		statement = p.parseReturnStatement()
	} else if p.isKeyword("if") {
		statement = p.parseIfStatement()
	} else if p.isKeyword("else") {
		p.error("Cannot use else statement without preceding if", p.next())
	} else if p.isKeyword("while") {
		statement = p.parseWhileLoop()
	} else if p.isKeyword("for") {
		statement = p.parseForLoop()
	} else {
		statement = p.parseExpressionStatement()
	}

	if len(inline) != 0 && inline[0] {
		return statement
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
		Name:       name,
		Parameters: parameters,
		Body:       code,
		ReturnType: returnType,
		BaseNode:   &ast.BaseNode{Token: tok},
	}
}

func (p *parser) parseReturnStatement() ast.Statement {
	token := p.consume()

	var value ast.Expression = &ast.VoidValue{}

	if p.canContinue() {
		value = p.parseExpression()
	}

	return &ast.ReturnStatement{
		Value:    value,
		BaseNode: &ast.BaseNode{Token: token},
	}
}

func (p *parser) parseIfStatement() *ast.IfStatement {
	tok := p.consume()

	condition := p.parseExpression()
	body := p.parseCodeBlock()
	var elseStatement ast.IfElseStatement = nil

	if p.isKeyword("else") {
		elseToken := p.consume()
		if p.next().Type == token.LEFT_BRACE {
			elseStatement = &ast.ElseStatement{
				Body: p.parseCodeBlock(),
				BaseNode: &ast.BaseNode{Token: elseToken},
			}
		} else {
			elseStatement = p.parseIfStatement()
		}
	}

	return &ast.IfStatement{
		Condition: condition,
		Body: body,
		BaseNode: &ast.BaseNode{Token: tok},
		Else: elseStatement,
	}
}

func (p *parser) parseWhileLoop() ast.Statement {
	tok := p.consume()

	condition := p.parseExpression()
	body := p.parseCodeBlock()

	return &ast.WhileLoop{
		Condition: condition,
		Body: body,
		BaseNode: &ast.BaseNode{Token: tok},
	}
}

func (p *parser) parseForLoop() ast.Statement {
	tok := p.consume()

	initial := p.parseStatement(true)
	p.expect(token.SEMICOLON, "Expected for loop condition")
	condition := p.parseExpression()
	p.expect(token.SEMICOLON, "Expected for loop update statement")
	update := p.parseStatement(true)
	body := p.parseCodeBlock()

	return &ast.ForLoop{
		Initial: initial,
		Condition: condition,
		Update: update,
		Body: body,
		BaseNode: &ast.BaseNode{Token: tok},
	}
}
