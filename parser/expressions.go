package parser

import (
	"fmt"
	"strconv"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseExpression() (ast.Expression, error) {
	return p.parseAssignmentExpression()
}

// Orders of precedence

// Assignment
// Logical operators
// Comparison
// Addition/Subtraction
// Multiplication/Division
// Member access
// Function call
// Unary operation
// Literal

func (p *parser) parseAssignmentExpression() (ast.Expression, error) {
	assignee, err := p.parseLogicalExpression()
	if err != nil {
		return nil, err
	}

	if !p.canContinue() || !p.next().Is(token.AssignmentOperator) {
		return assignee, nil
	}

	operation := p.consume()

	value, err := p.parseAssignmentExpression()
	if err != nil {
		return nil, err
	}

	return &ast.AssignmentExpression{
		Assignee:  assignee,
		Value:     value,
		Operation: operation.Value,
		BaseNode:  &ast.BaseNode{Token: assignee.GetToken()},
	}, nil
}

func (p *parser) parseLogicalExpression() (ast.Expression, error) {
	left, err := p.parseComparisonExpression()
	if err != nil {
		return nil, err
	}

	for p.canContinue() && p.next().Is(token.LogicalOperator) {
		operator := p.consume().Value
		right, err := p.parseComparisonExpression()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryOperation{
			Left:     left,
			Operator: operator,
			Right:    right,
			BaseNode: &ast.BaseNode{Token: left.GetToken()},
		}
	}

	return left, nil
}

func (p *parser) parseComparisonExpression() (ast.Expression, error) {
	left, err := p.parseAdditiveExpression()
	if err != nil {
		return nil, err
	}

	for p.canContinue() && p.next().Is(token.ComparisonOperator) {
		operator := p.consume().Value
		right, err := p.parseAdditiveExpression()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryOperation{
			Left:     left,
			Operator: operator,
			Right:    right,
			BaseNode: &ast.BaseNode{Token: left.GetToken()},
		}
	}

	return left, nil
}

func (p *parser) parseAdditiveExpression() (ast.Expression, error) {
	left, err := p.parseMultiplicativeExpression()
	if err != nil {
		return nil, err
	}

	for p.canContinue() && p.next().Is(token.AdditiveOperator) {
		operator := p.consume().Value
		right, err := p.parseMultiplicativeExpression()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryOperation{
			Left:     left,
			Operator: operator,
			Right:    right,
			BaseNode: &ast.BaseNode{Token: left.GetToken()},
		}
	}

	return left, nil
}

func (p *parser) parseMultiplicativeExpression() (ast.Expression, error) {
	left, err := p.parseExponentialExpression()
	if err != nil {
		return nil, err
	}

	for p.canContinue() && p.next().Is(token.MultiplicativeOperator) {
		operator := p.consume().Value
		right, err := p.parseExponentialExpression()
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryOperation{
			Left:     left,
			Operator: operator,
			Right:    right,
			BaseNode: &ast.BaseNode{Token: left.GetToken()},
		}
	}

	return left, nil
}

func (p *parser) parseExponentialExpression() (ast.Expression, error) {
	left, err := p.parsePrefixOperation()
	if err != nil {
		return nil, err
	}

	if !p.canContinue() || p.next().Type != token.POWER {
		return left, nil
	}

	operator := p.consume().Value

	right, err := p.parseExponentialExpression()
	if err != nil {
		return nil, err
	}

	return &ast.BinaryOperation{
		Left:     left,
		Operator: operator,
		Right:    right,
		BaseNode: &ast.BaseNode{Token: left.GetToken()},
	}, nil
}

func (p *parser) parsePrefixOperation() (ast.Expression, error) {
	if !p.next().Is(token.PrefixOperator) {
		return p.parsePostfixOperation()
	}

	operator := p.consume()
	value, err := p.parsePrefixOperation()
	if err != nil {
		return nil, err
	}

	return &ast.UnaryOperation{
		Operator: operator.Value,
		Value:    value,
		BaseNode: &ast.BaseNode{Token: operator},
		Postfix:  false,
	}, nil
}

func (p *parser) parsePostfixOperation() (ast.Expression, error) {
	value, err := p.parseIndexExpression()
	if err != nil {
		return nil, err
	}

	for p.next().Is(token.PostfixOperator) {
		value = &ast.UnaryOperation{
			Value:    value,
			Operator: p.consume().Value,
			BaseNode: &ast.BaseNode{Token: value.GetToken()},
			Postfix:  true,
		}
	}

	return value, nil
}

func (p *parser) parseIndexExpression() (ast.Expression, error) {
	left, err := p.parseLiteral()
	if err != nil {
		return nil, err
	}

	if p.next().Type != token.LEFT_SQUARE {
		return left, nil
	}

	p.consume()
	index, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	p.expect(token.RIGHT_SQUARE, "Expected bracket after index expression")

	return &ast.IndexExpression{
		Left: left,
		Index: index,
		BaseNode: &ast.BaseNode{Token: left.GetToken()},
	}, nil
}

func (p *parser) parseFunctionCall() (ast.Expression, error) {
	token := p.consume()

	args, err := p.parseArgumentList()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionCall{
		Name:     token.Value,
		Args:     args,
		BaseNode: &ast.BaseNode{Token: token},
	}, nil
}

func (p *parser) parseList() (ast.Expression, error) {
	tok := p.consume()
	values := []ast.Expression{}

	for p.next().Type != token.RIGHT_SQUARE && !p.eof() {
		nextExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		values = append(values, nextExpr)

		if p.next().Type != token.RIGHT_SQUARE {
			_, err := p.expect(token.COMMA, "Expected comma after list entry, got %q")
			if err != nil {
				return nil, err
			}
		}
	}

	_, err := p.expect(token.RIGHT_SQUARE, "Expected closing bracket after list, got %q")
	if err != nil {
		return nil, err
	}

	return &ast.ListLiteral{
		Elements: values,
		BaseNode: &ast.BaseNode{Token: tok},
	}, nil
}

func (p *parser) parseIdentifier() (ast.Expression, error) {
	if p.isKeyword("true") {
		tok := p.consume()
		return &ast.BooleanLiteral{
			Value:    true,
			BaseNode: &ast.BaseNode{Token: tok},
		}, nil
	}

	if p.isKeyword("false") {
		tok := p.consume()
		return &ast.BooleanLiteral{
			Value:    false,
			BaseNode: &ast.BaseNode{Token: tok},
		}, nil
	}

	if p.isKeyword("null") {
		tok := p.consume()
		return &ast.NullLiteral{
			BaseNode: &ast.BaseNode{Token: tok},
		}, nil
	}

	tok := p.consume()
	return &ast.Identifier{
		Symbol:   tok.Value,
		BaseNode: &ast.BaseNode{Token: tok},
	}, nil
}

func (p *parser) parseLiteral() (ast.Expression, error) {
	switch p.next().Type {
	case token.INTEGER:
		tok := p.consume()
		value, _ := strconv.ParseInt(tok.Value, 10, 32)
		return &ast.IntegerLiteral{
			Value:    int(value),
			BaseNode: &ast.BaseNode{Token: tok},
		}, nil

	case token.FLOAT:
		tok := p.consume()
		value, _ := strconv.ParseFloat(tok.Value, 64)
		return &ast.FloatLiteral{
			Value:    value,
			BaseNode: &ast.BaseNode{Token: tok},
		}, nil

	case token.STRING:
		tok := p.consume()
		return &ast.StringLiteral{
			Value:    tok.Value,
			BaseNode: &ast.BaseNode{Token: tok},
		}, nil

	case token.IDENTIFIER:
		switch p.tokens[1].Type {
		case token.LEFT_PAREN:
			return p.parseFunctionCall()

		default:
			return p.parseIdentifier()
		}

	case token.LEFT_PAREN:
		p.consume()
		p.bracketLevel++
		expression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(token.RIGHT_PAREN, "Expected closing parentheses after bracketed expression, got %q")
		if err != nil {
			return nil, err
		}

		p.bracketLevel--
		return expression, nil

	case token.LEFT_SQUARE:
		return p.parseList()

	default:
		return nil, p.error(fmt.Sprintf("Expected expression, got %q", p.next().Value), p.next())
	}
}
