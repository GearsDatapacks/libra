package parser

import (
	"strconv"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

// Precedence
const (
	Lowest = iota
	Assignment
	Logical
	Comparison
	Bitwise
	Additive
	Multiplicative
	Exponential
	Typecheck
	Prefix
	Postfix
	Cast
)

func (p *parser) parseExpression() ast.Expression {
	return p.parseSubExpression(Lowest)
}

func (p *parser) parseSubExpression(precedence int) ast.Expression {
	nudFn := p.lookupNudFn(p.next().Kind)

	if nudFn == nil {
		p.Diagnostics.ReportExpectedExpression(p.next().Span, p.next().Kind)
		return &ast.ErrorExpression{}
	}

	left := nudFn()

	ledFn := p.lookupLedFn(p.next().Kind)
	for ledFn != nil && precedence < p.leftPrecedence(p.next().Kind) {
		left = ledFn(left)

		ledFn = p.lookupLedFn(p.next().Kind)
	}

	return left
}

func (p *parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	operator := p.consume()
	right := p.parseSubExpression(p.rightPrecedence(operator.Kind))

	return &ast.BinaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (p *parser) parsePrefixExpression() ast.Expression {
	operator := p.consume()
	operand := p.parseSubExpression(Prefix)

	return &ast.PrefixExpression{
		Operator: operator,
		Operand:  operand,
	}
}

func (p *parser) parsePostfixExpression(operand ast.Expression) ast.Expression {
	operator := p.consume()

	return &ast.PostfixExpression{
		Operand:  operand,
		Operator: operator,
	}
}

func (p *parser) parseParenthesisedExpression() ast.Expression {
	leftParen := p.consume()
	expr := p.parseExpression()
	rightParen := p.expect(token.RIGHT_PAREN)
	return &ast.ParenthesisedExpression{
		LeftParen:  leftParen,
		Expression: expr,
		RightParen: rightParen,
	}
}

func (p *parser) parseIdentifier() ast.Expression {
	tok := p.consume()

	switch tok.Value {
	case "true":
		return &ast.BooleanLiteral{
			Token: tok,
			Value: true,
		}

	case "false":
		return &ast.BooleanLiteral{
			Token: tok,
			Value: false,
		}

	default:
		return &ast.Identifier{
			Token: tok,
			Name:  tok.Value,
		}
	}
}

func (p *parser) parseInteger() ast.Expression {
	tok := p.consume()
	value, _ := strconv.ParseInt(tok.Value, 10, 64)
	return &ast.IntegerLiteral{
		Token: tok,
		Value: value,
	}
}

func (p *parser) parseFloat() ast.Expression {
	tok := p.consume()
	value, _ := strconv.ParseFloat(tok.Value, 64)
	return &ast.FloatLiteral{
		Token: tok,
		Value: value,
	}
}

func (p *parser) parseString() ast.Expression {
	tok := p.consume()
	return &ast.StringLiteral{
		Token: tok,
		Value: tok.Value,
	}
}

func (p *parser) parseList() ast.Expression {
	leftSquare := p.consume()
	values, rightSquare := parseDelemitedList(p, token.RIGHT_SQUARE, p.parseExpression)

	return &ast.ListLiteral{
		LeftSquare:  leftSquare,
		Values:      values,
		RightSquare: rightSquare,
	}
}

func (p *parser) parseKeyValue() ast.KeyValue {
	key := p.parseExpression()
	colon := p.expect(token.COLON)
	value := p.parseExpression()

	return ast.KeyValue{
		Key:   key,
		Colon: colon,
		Value: value,
	}
}

func (p *parser) parseMap() ast.Expression {
	leftBrace := p.consume()
	keyValues, rightBrace := parseDelemitedList(p, token.RIGHT_BRACE, p.parseKeyValue)

	return &ast.MapLiteral{
		LeftBrace:  leftBrace,
		KeyValues:  keyValues,
		RightBrace: rightBrace,
	}
}
