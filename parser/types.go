package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseType() ast.TypeExpression {
	ty := p.parsePrimaryType()

	if p.next().Kind == token.PIPE {
		types := []ast.TypeExpression{ty}

		for p.canContinue() && p.next().Kind == token.PIPE {
			p.consume()
			types = append(types, p.parsePrimaryType())
		}

		ty = &ast.Union{
			Types: types,
		}
	}

	return ty
}

func (p *parser) parsePrimaryType() ast.TypeExpression {
	switch p.next().Kind {
	case token.LEFT_SQUARE:
		return p.parseArrayType()
	case token.STAR:
		return p.parsePointerType()
	case token.IDENTIFIER:
		return p.parseTypeName()
	default:
		p.Diagnostics.ReportExpectedType(p.next().Span, p.next().Kind)
		return &ast.ErrorNode{}
	}
}

func (p *parser) parseArrayType() ast.TypeExpression {
	leftSquare := p.consume()
	var count ast.Expression
	if p.next().Kind != token.RIGHT_SQUARE {
		count = p.parseExpression()
	}
	rightSquare := p.expect(token.RIGHT_SQUARE)
	ty := p.parsePrimaryType()

	return &ast.ArrayType{
		LeftSquare:  leftSquare,
		Count:       count,
		RightSquare: rightSquare,
		Type:        ty,
	}
}

func (p *parser) parsePointerType() ast.TypeExpression {
	star := p.consume()
	var mut *token.Token
	if p.isKeyword("mut") {
		tok := p.consume()
		mut = &tok
	}
	ty := p.parsePrimaryType()
	
	return &ast.PointerType{
		Star: star,
		Mut:  mut,
		Type: ty,
	}
}

func (p *parser) parseTypeName() ast.TypeExpression {
	name := p.expect(token.IDENTIFIER)

	return &ast.TypeName{
		Name: name,
	}
}
