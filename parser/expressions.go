package parser

import (
	"log"
	"strconv"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseExpression() ast.Expression {
	return p.parseAdditiveExpression()
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

func (p *parser) parseAdditiveExpression() ast.Expression {
	left := p.parseMultiplicativeExpression()

	for p.next().Value == "+" || p.next().Value == "-" {
		operator := p.consume().Value
		right := p.parseMultiplicativeExpression()
		left = &ast.BinaryOperation{
			Left:     left,
			Operator: operator,
			Right:    right,
			BaseNode: &ast.BaseNode{Token: left.GetToken()},
		}
	}

	return left
}

func (p *parser) parseMultiplicativeExpression() ast.Expression {
	left := p.parseLiteral()

	for p.next().Value == "*" || p.next().Value == "/" {
		operator := p.consume().Value
		right := p.parseLiteral()
		left = &ast.BinaryOperation{
			Left:     left,
			Operator: operator,
			Right:    right,
			BaseNode: &ast.BaseNode{Token: left.GetToken()},
		}
	}

	return left
}

func (p *parser) parseLiteral() ast.Expression {
	switch p.next().Type {
	case token.INTEGER:
		tok := p.consume()
		value, _ := strconv.ParseInt(tok.Value, 10, 32)
		return &ast.IntegerLiteral{
			Value:    int(value),
			BaseNode: &ast.BaseNode{Token: tok},
		}
	
	case token.IDENTIFIER:
		tok := p.consume()
		return &ast.Identifier{
			Symbol: tok.Value,
			BaseNode: &ast.BaseNode{Token: tok},
		}

	case token.LEFT_PAREN:
		p.consume()
		expression := p.parseExpression()
		p.expect(token.RIGHT_PAREN, "Expected closing parentheses after bracketed expression, got %q")
		return expression
	default:
		log.Fatalf("ParseError: Unexpected token %q", p.next().Value)
		return &ast.IntegerLiteral{}
	}
}
