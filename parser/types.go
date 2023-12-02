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
	left, err := p.parseSuffixType()
	if err != nil {
		return nil, err
	}
	
	if p.next().Type != token.BITWISE_OR {
		return left, nil
	}
	
	types := []ast.TypeExpression{left}
	
	for p.next().Type == token.BITWISE_OR {
		p.consume()

		nextType, err := p.parseSuffixType()
		if err != nil {
			return nil, err
		}
		types = append(types, nextType)
	}

	return &ast.Union{
		ValidTypes: types,
		BaseNode: ast.BaseNode{ Token: types[0].GetToken() },
	}, nil
}

func (p *parser) parseSuffixType() (ast.TypeExpression, error) {
	leftType, err := p.parsePrimaryType()
	if err != nil {
		return nil, err
	}

	for {
		switch p.next().Type {
		case token.LEFT_SQUARE:
			leftType, err = p.parseArrayType(leftType, leftType.GetToken())
			if err != nil {
				return nil, err
			}

		case token.LOGICAL_NOT:
			p.consume()
			leftType = &ast.ErrorType{
				ResultType: leftType,
				BaseNode:   ast.BaseNode{Token: leftType.GetToken()},
			}

		default:
			return leftType, nil
		}
	}
}

func (p *parser) parseArrayType(elemType ast.TypeExpression, tok token.Token) (ast.TypeExpression, error) {
	if p.next().Type == token.LEFT_SQUARE {
		
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
				BaseNode: ast.BaseNode{Token: tok},
			}
			continue
		}
		
		if p.isKeyword("_") {
			p.consume()
			_, err := p.expect(token.RIGHT_SQUARE, "Unexpected %q, expecting ']'")
			if err != nil {
				return nil, err
			}
			elemType = &ast.ArrayType{
				ElementType: elemType,
				Length: -1,
				BaseNode: ast.BaseNode{Token: tok},
			}
			continue
		}
		
		lengthTok, err := p.expect(token.INTEGER, "Invalid array length: %q")
		if err != nil {
			return nil, err
		}

		length, _ := strconv.ParseInt(lengthTok.Value, 10, 32)
		intLength := int(length)
		_, err = p.expect(token.RIGHT_SQUARE, "Unexpected %q, expecting ']'")
		if err != nil {
			return nil, err
		}
		elemType = &ast.ArrayType{
			ElementType: elemType,
			Length: intLength,
			BaseNode: ast.BaseNode{Token: tok},
		}
	}

	return elemType, nil
}

func (p *parser) parseMapType() (ast.TypeExpression, error) {
	tok := p.consume()

	keyType, err := p.parseType()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(token.COLON, "Unexpected %q, expecting ':'")
	if err != nil {
		return nil, err
	}

	valueType, err := p.parseType()
	if err != nil {
		return nil, err
	}

	_, err = p.expect(token.RIGHT_BRACE, "Unexpected %q, expecting '}'")
	if err != nil {
		return nil, err
	}

	return &ast.MapType{
		BaseNode:  ast.BaseNode{Token: tok},
		KeyType:   keyType,
		ValueType: valueType,
	}, nil
}

func (p *parser) parsePrimaryType() (ast.TypeExpression, error) {
	switch p.next().Type {
	case token.IDENTIFIER:
		tok := p.consume()
		return &ast.TypeName{ Name: tok.Value, BaseNode: ast.BaseNode{ Token: tok } }, nil

	case token.LEFT_PAREN:
		tok := p.consume()
		expr, err := p.parseType()
		if err != nil {
			return nil, err
		}

		if p.next().Type == token.COMMA {
			members := []ast.TypeExpression{expr}
			for p.next().Type == token.COMMA {
				p.consume()
				nextExpr, err := p.parseType()
				if err != nil {
					return nil, err
				}
				members = append(members, nextExpr)
			}
			expr = &ast.TupleType{Members: members, BaseNode: ast.BaseNode{Token: tok}}
		}

		_, err = p.expect(token.RIGHT_PAREN, "Expected comma or end of tuple type")
		return expr, err
	
	case token.LEFT_BRACE:
		return p.parseMapType()
	
	case token.LEFT_SQUARE:
		elemType := &ast.InferType{BaseNode: ast.BaseNode{Token: p.next()}}
		tok := p.next()
		return p.parseArrayType(elemType, tok)
	
	case token.LOGICAL_NOT:
		return &ast.ErrorType{
			BaseNode:   ast.BaseNode{ Token: p.consume() },
			ResultType: &ast.VoidType{},
		}, nil
		
	default:
		return nil, p.error(fmt.Sprintf("Expected type, got %q", p.next().Value), p.next())
	}
}
