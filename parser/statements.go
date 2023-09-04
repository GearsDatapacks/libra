package parser

import (
	"fmt"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseStatement(inline ...bool) (ast.Statement, error) {
	var statement ast.Statement
	var err error = nil

	if p.isKeyword("var") || p.isKeyword("const") {
		statement, err = p.parseVariableDeclaration()
	} else if p.isKeyword("fn") {
		statement, err = p.parseFunctionDeclaration()
	} else if p.isKeyword("return") {
		statement, err = p.parseReturnStatement()
	} else if p.isKeyword("if") {
		statement, err = p.parseIfStatement()
	} else if p.isKeyword("else") {
		return nil, p.error("Cannot use else statement without preceding if", p.next())
	} else if p.isKeyword("while") {
		statement, err = p.parseWhileLoop()
	} else if p.isKeyword("for") {
		statement, err = p.parseForLoop()
	} else {
		statement, err = p.parseExpressionStatement()
	}

	if len(inline) != 0 && inline[0] {
		return statement, err
	}

	if !p.eof() && !p.next().LeadingNewline {
		return nil, p.error(fmt.Sprintf("Expected new line after statement, got %q", p.next().Value), p.next())
	}

	return statement, err
}

func (p *parser) parseExpressionStatement() (ast.Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{
		Expression: expr,
		BaseNode:   &ast.BaseNode{Token: expr.GetToken()},
	}, nil
}

func (p *parser) parseVariableDeclaration() (ast.Statement, error) {
	tok := p.consume()
	isConstant := tok.Value == "const"
	name, err := p.expect(
		token.IDENTIFIER,
		"ParseError: Expected identifier for variable declaration, got %q",
	)
	if err != nil {
		return nil, err
	}

	var dataType ast.TypeExpression = &ast.InferType{}

	if p.canContinue() && p.next().Type != token.ASSIGNMENT_OPERATOR {
		dataType, err = p.parseType()
		if err != nil {
			return nil, err
		}
	}

	// TODO: add possibility for `var x string`

	if !p.canContinue() || p.next().Type != token.ASSIGNMENT_OPERATOR {
		if isConstant {
			return nil, p.error(fmt.Sprintf("Cannot leave constant %q uninitialised", name.Value), p.next())
		}

		if dataType.Type() == "Infer" {
			return nil, p.error(fmt.Sprintf("Cannot declare uninitialised variable %q without type annotation", name.Value), p.next())
		}

		return &ast.VariableDeclaration{
			Constant: isConstant,
			Name:     name.Value,
			BaseNode: &ast.BaseNode{Token: tok},
			Value:    nil,
			DataType: dataType,
		}, nil
	}

	_, err = p.expect(
		token.ASSIGNMENT_OPERATOR,
		"ParseError: Missing initialiser in variable declaration",
	)
	if err != nil {
		return nil, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	p.usedSymbols = append(p.usedSymbols, name.Value)

	return &ast.VariableDeclaration{
		Constant: isConstant,
		Name:     name.Value,
		BaseNode: &ast.BaseNode{Token: tok},
		Value:    value,
		DataType: dataType,
	}, nil
}

func (p *parser) parseFunctionDeclaration() (ast.Statement, error) {
	tok := p.consume()

	name, err := p.expect(token.IDENTIFIER, "Expected function name, got %q")
	if err != nil {
		return nil, err
	}

	p.usedSymbols = append(p.usedSymbols, name.Value)

	parameters, err := p.parseParameterList()
	if err != nil {
		return nil, err
	}

	outerSymbols := make([]string, len(p.usedSymbols))
	copy(outerSymbols, p.usedSymbols)

	for _, param := range parameters {
		p.usedSymbols = append(p.usedSymbols, param.Name)
	}

	var returnType ast.TypeExpression = &ast.VoidType{}

	if p.next().Type != token.LEFT_BRACE {
		returnType, err = p.parseType()
		if err != nil {
			return nil, err
		}
	}

	code, err := p.parseCodeBlock()
	if err != nil {
		return nil, err
	}

	p.usedSymbols = outerSymbols

	return &ast.FunctionDeclaration{
		Name:       name.Value,
		Parameters: parameters,
		Body:       code,
		ReturnType: returnType,
		BaseNode:   &ast.BaseNode{Token: tok},
	}, nil
}

func (p *parser) parseReturnStatement() (ast.Statement, error) {
	token := p.consume()

	var value ast.Expression = &ast.VoidValue{}

	if p.canContinue() {
		var err error = nil
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	return &ast.ReturnStatement{
		Value:    value,
		BaseNode: &ast.BaseNode{Token: token},
	}, nil
}

func (p *parser) parseIfStatement() (*ast.IfStatement, error) {
	tok := p.consume()

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	body, err := p.parseCodeBlock()
	if err != nil {
		return nil, err
	}
	var elseStatement ast.IfElseStatement = nil

	if p.isKeyword("else") {
		elseToken := p.consume()
		if p.next().Type == token.LEFT_BRACE {
			code, err := p.parseCodeBlock()
			if err != nil {
				return nil, err
			}
			elseStatement = &ast.ElseStatement{
				Body: code,
				BaseNode: &ast.BaseNode{Token: elseToken},
			}
		} else {
			elseStatement, err = p.parseIfStatement()	
		}

		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStatement{
		Condition: condition,
		Body: body,
		BaseNode: &ast.BaseNode{Token: tok},
		Else: elseStatement,
	}, nil
}

func (p *parser) parseWhileLoop() (ast.Statement, error) {
	tok := p.consume()

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	body, err := p.parseCodeBlock()
	if err != nil {
		return nil, err
	}

	return &ast.WhileLoop{
		Condition: condition,
		Body: body,
		BaseNode: &ast.BaseNode{Token: tok},
	}, nil
}

func (p *parser) parseForLoop() (ast.Statement, error) {
	tok := p.consume()

	outerSymbols := make([]string, len(p.usedSymbols))
	copy(outerSymbols, p.usedSymbols)

	initial, err := p.parseStatement(true)
	if err != nil {
		return nil, err
	}

	_, err = p.expect(token.SEMICOLON, "Expected for loop condition")
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	_, err = p.expect(token.SEMICOLON, "Expected for loop update statement")
	if err != nil {
		return nil, err
	}

	update, err := p.parseStatement(true)
	if err != nil {
		return nil, err
	}

	body, err := p.parseCodeBlock()
	if err != nil {
		return nil, err
	}

	p.usedSymbols = outerSymbols

	return &ast.ForLoop{
		Initial: initial,
		Condition: condition,
		Update: update,
		Body: body,
		BaseNode: &ast.BaseNode{Token: tok},
	}, nil
}
