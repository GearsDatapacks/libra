package parser

import (
	"fmt"
	"strconv"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseType() (ast.TypeExpression, error) {
	return p.parseUnion()
}

func (p *parser) parseUnion() (ast.TypeExpression, error) {
	left, err := p.parseArrayType()
	if err != nil {
		return nil, err
	}
	
	if p.next().Type != token.BITWISE_OR {
		return left, nil
	}
	
	types := []ast.TypeExpression{left}
	
	for p.next().Type == token.BITWISE_OR {
		p.consume()

		nextType, err := p.parseArrayType()
		if err != nil {
			return nil, err
		}
		types = append(types, nextType)
	}

	return &ast.Union{
		ValidTypes: types,
		BaseNode: &ast.BaseNode{ Token: types[0].GetToken() },
	}, nil
}

func (p *parser) parseArrayType() (ast.TypeExpression, error) {
	var elemType ast.TypeExpression
	var tok token.Token
	if p.next().Type == token.LEFT_SQUARE {
		elemType = &ast.InferType{BaseNode: &ast.BaseNode{Token: p.next()}}
		tok = p.next()
	} else {
		var err error = nil
		elemType, err = p.parsePrimaryType()
		if err != nil {
			return nil, err
		}
		tok = elemType.GetToken()
	}

	for p.next().Type == token.LEFT_SQUARE {
		p.consume()

		if p.next().Type == token.RIGHT_SQUARE {
			p.consume()
			elemType = &ast.ListType{
				ElementType: elemType,
				BaseNode: &ast.BaseNode{Token: tok},
			}
			continue
		}
		
		if p.isKeyword("_") {
			p.consume()
			_, err := p.expect(token.RIGHT_SQUARE, "Array types must contain one entry")
			if err != nil {
				return nil, err
			}
			elemType = &ast.ArrayType{
				ElementType: elemType,
				Length: -1,
				BaseNode: &ast.BaseNode{Token: tok},
			}
			continue
		}
		
		lengthTok, err := p.expect(token.INTEGER, "Invalid array length: %q")
		if err != nil {
			return nil, err
		}

		length, _ := strconv.ParseInt(lengthTok.Value, 10, 32)
		intLength := int(length)
		_, err = p.expect(token.RIGHT_SQUARE, "Array types must contain one entry")
		if err != nil {
			return nil, err
		}
		elemType = &ast.ArrayType{
			ElementType: elemType,
			Length: intLength,
			BaseNode: &ast.BaseNode{Token: tok},
		}
	}

	return elemType, nil
}

func (p *parser) parsePrimaryType() (ast.TypeExpression, error) {
	switch p.next().Type {
	case token.IDENTIFIER:
		tok := p.consume()
		return &ast.TypeName{ Name: tok.Value, BaseNode: &ast.BaseNode{ Token: tok } }, nil

	case token.LEFT_PAREN:
		p.consume()
		expr, err := p.parseType()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(token.RIGHT_PAREN, "Expected closing bracket after expression, got %q")
		return expr, err
		
	default:
		return nil, p.error(fmt.Sprintf("Expected type, got %q", p.next().Value), p.next())
	}
}
