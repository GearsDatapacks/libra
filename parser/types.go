package parser

import (
	"fmt"
	"strconv"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseType() ast.TypeExpression {
	return p.parseUnion()
}

func (p *parser) parseUnion() ast.TypeExpression {
	left := p.parseListType()
	
	if p.next().Type != token.BITWISE_OR {
		return left
	}
	
	types := []ast.TypeExpression{left}
	
	for p.next().Type == token.BITWISE_OR {
		p.consume()

		types = append(types, p.parseListType())
	}

	return &ast.Union{
		ValidTypes: types,
		BaseNode: &ast.BaseNode{ Token: types[0].GetToken() },
	}
}

func (p *parser) parseListType() ast.TypeExpression {
	elemType := p.parsePrimaryType()

	for p.next().Type == token.LEFT_SQUARE || p.next().Type == token.LEFT_BRACE {
		if nextTok := p.consume().Type; nextTok == token.LEFT_SQUARE {
			p.expect(token.RIGHT_SQUARE, "List types must have empty brackets")
			elemType = &ast.ListType{
				ElementType: elemType,
				BaseNode: &ast.BaseNode{Token: elemType.GetToken()},
			}
			continue
		}
			
		if p.next().Type == token.RIGHT_BRACE {
			p.consume()
			elemType = &ast.ArrayType{
				ElementType: elemType,
				Length: -1,
				BaseNode: &ast.BaseNode{Token: elemType.GetToken()},
			}
			continue
		}
		lengthTok := p.expect(token.INTEGER, "Array types must have length of an integer value")
		length, _ := strconv.ParseInt(lengthTok.Value, 10, 32)
		intLength := int(length)
		p.expect(token.RIGHT_BRACE, "Array types must contain one entry")
		elemType = &ast.ArrayType{
			ElementType: elemType,
			Length: intLength,
			BaseNode: &ast.BaseNode{Token: elemType.GetToken()},
		}
	}
	return elemType
}

func (p *parser) parsePrimaryType() ast.TypeExpression {
	switch p.next().Type {
	case token.IDENTIFIER:
		tok := p.consume()
		return &ast.TypeName{ Name: tok.Value, BaseNode: &ast.BaseNode{ Token: tok } }

	case token.LEFT_PAREN:
		p.consume()
		expr := p.parseType()
		p.expect(token.RIGHT_PAREN, "Expected closing bracket after expression, got %q")
		return expr
		
	default:
		p.error(fmt.Sprintf("Expected type, got %q", p.next().Value), p.next())
		return &ast.TypeName{}
	}
}
