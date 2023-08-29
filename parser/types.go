package parser

import (
	"fmt"

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

	if p.next().Type != token.LEFT_SQUARE {
		return elemType
	}

	p.consume()
	p.expect(token.RIGHT_SQUARE, "List types must have empty brackets")
	return &ast.ListType{
		ElementType: elemType,
		BaseNode: &ast.BaseNode{Token: elemType.GetToken()},
	}
}

func (p *parser) parsePrimaryType() ast.TypeExpression {
	switch p.next().Type {
	case token.IDENTIFIER:
		tok := p.consume()
		return &ast.TypeName{ Name: tok.Value, BaseNode: &ast.BaseNode{ Token: tok } }
		
	default:
		p.error(fmt.Sprintf("Expected type, got %q", p.next().Value), p.next())
		return &ast.TypeName{}
	}
}
