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

func (p *parser) parseAssignmentExpression() (ast.Expression, error) {
	assignee, err := p.parseBinaryOperation(0)
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

func (p *parser) parseBinaryOperation(minPrecedence int) (ast.Expression, error) {
	left, err := p.parsePrefixOperation()
	if err != nil {
		return nil, err
	}

	for {
		opInfo, isOp := token.BinOpInfo[p.next().Type]
		if !isOp || opInfo.Precedence < minPrecedence {
			break
		}

		op := p.consume().Value

		newMinPrec := opInfo.Precedence + 1
		if opInfo.RightAssociative {
			newMinPrec = opInfo.Precedence
		}

		right, err := p.parseBinaryOperation(newMinPrec)
		if err != nil {
			return nil, err
		}

		left = &ast.BinaryOperation{
			Left:     left,
			Operator: op,
			Right:    right,
			BaseNode: &ast.BaseNode{Token: left.GetToken()},
		}
	}

	return left, nil
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
	value, err := p.parseCallMemberExpression()
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

func (p *parser) parseCallMemberExpression() (ast.Expression, error) {
	left, err := p.parseLiteral()
	if err != nil {
		return nil, err
	}

	for p.next().Type == token.LEFT_PAREN || p.next().Type == token.LEFT_SQUARE || p.next().Type == token.DOT {
		if p.next().Type == token.LEFT_PAREN {
			left, err = p.parseFunctionCall(left)
		} else if p.next().Type == token.LEFT_SQUARE {
			left, err = p.parseIndexExpression(left)
		} else if p.next().Type == token.DOT {
			left, err = p.parseMemberExpression(left)
		}

		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *parser) parseMemberExpression(left ast.Expression) (ast.Expression, error) {
	p.consume()
	memberName, err := p.expect(token.IDENTIFIER, "Invalid member name %q")
	if err != nil {
		return nil, err
	}

	return &ast.MemberExpression{
		Left:     left,
		Member:    memberName.Value,
		BaseNode: &ast.BaseNode{Token: left.GetToken()},
	}, nil
}

func (p *parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	p.consume()
	index, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	_, err = p.expect(token.RIGHT_SQUARE, "Unexpected token %q, expecting ']'")
	if err != nil {
		return nil, err
	}

	return &ast.IndexExpression{
		Left:     left,
		Index:    index,
		BaseNode: &ast.BaseNode{Token: left.GetToken()},
	}, nil
}

func (p *parser) parseFunctionCall(left ast.Expression) (ast.Expression, error) {
	args, err := p.parseArgumentList()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionCall{
		Left:     left,
		Args:     args,
		BaseNode: &ast.BaseNode{Token: left.GetToken()},
	}, nil
}

func (p *parser) parseStructExpression() (ast.Expression, error) {
	name := p.consume()
	p.consume()

	members := map[string]ast.Expression{}

	for !p.eof() && p.next().Type != token.RIGHT_BRACE {
		memberName, err := p.expect(token.IDENTIFIER, "Invalid struct member name %q")
		if err != nil {
			return nil, err
		}

		_, err = p.expect(token.COLON, "Unexpected %q, expected ':'")
		if err != nil {
			return nil, err
		}

		memberValue, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		members[memberName.Value] = memberValue

		if p.next().Type != token.RIGHT_BRACE {
			_, err := p.expect(token.COMMA, "Expected comma or end of struct body")
			if err != nil {
				return nil, err
			}
		}
	}

	p.expect(token.RIGHT_BRACE, "Unexpected EOF, expected '}'")

	return &ast.StructExpression{
		BaseNode: &ast.BaseNode{Token: name},
		Name:     name.Value,
		Members:  members,
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
			_, err := p.expect(token.COMMA, "Expected comma or end of list")
			if err != nil {
				return nil, err
			}
		}
	}

	_, err := p.expect(token.RIGHT_SQUARE, "Unexpected EOF, expecting ']'")
	if err != nil {
		return nil, err
	}

	return &ast.ListLiteral{
		Elements: values,
		BaseNode: &ast.BaseNode{Token: tok},
	}, nil
}

func (p *parser) parseMap() (ast.Expression, error) {
	tok := p.consume()

	values := map[ast.Expression]ast.Expression{}

	for p.next().Type != token.RIGHT_BRACE && !p.eof() {
		keyExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(token.COLON, "Unexpected %q, expecting ':'")
		if err != nil {
			return nil, err
		}

		valueExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		values[keyExpr] = valueExpr

		if p.next().Type != token.RIGHT_BRACE {
			_, err = p.expect(token.COMMA, "Expected comma or end of map")
			if err != nil {
				return nil, err
			}
		}
	}

	_, err := p.expect(token.RIGHT_BRACE, "Unexpected EOF, expecting '}'")
	if err != nil {
		return nil, err
	}

	return &ast.MapLiteral{
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
		if p.tokens[1].Type == token.LEFT_BRACE && !p.noBraces {
			return p.parseStructExpression()
		} else {
			return p.parseIdentifier()
		}

	case token.LEFT_PAREN:
		p.consume()
		p.bracketLevel++
		noBraces := p.noBraces
		p.noBraces = false
		expression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(token.RIGHT_PAREN, "Unexpected %q, expecting ')'")
		if err != nil {
			return nil, err
		}

		p.bracketLevel--
		p.noBraces = noBraces
		return expression, nil

	case token.LEFT_SQUARE:
		return p.parseList()

	case token.LEFT_BRACE:
		return p.parseMap()

	default:
		return nil, p.error(fmt.Sprintf("Expected expression, got %q", p.next().Value), p.next())
	}
}
