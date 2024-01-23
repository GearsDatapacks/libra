package parser

import (
	"strconv"

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


