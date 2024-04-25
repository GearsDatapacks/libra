package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseType() ast.TypeExpression {
	ty := p.parsePostfixType()

	if p.next().Kind == token.PIPE {
		types := []ast.TypeExpression{ty}

		for p.canContinue() && p.next().Kind == token.PIPE {
			p.consume()
			types = append(types, p.parsePostfixType())
		}

		ty = &ast.Union{
			Types: types,
		}
	}

	return ty
}

func (p *parser) parsePostfixType() ast.TypeExpression {
	left := p.parsePrefixType()

	done := false
	for !done {
		switch p.next().Kind {
		case token.LEFT_SQUARE:
			left = p.parseArrayType(left)
		case token.QUESTION:
			left = &ast.OptionType{
				Type:     left,
				Question: p.consume(),
			}
		case token.BANG:
			left = &ast.ErrorType{
				Type: left,
				Bang: p.consume(),
			}
		case token.DOT:
			return p.parseMemberType(left)
		default:
			done = true
		}
	}

	return left
}

func (p *parser) parsePrefixType() ast.TypeExpression {
	switch p.next().Kind {
	case token.STAR:
		return p.parsePointerType()
	default:
		return p.parsePrimaryType()
	}
}

func (p *parser) parsePrimaryType() ast.TypeExpression {
	switch p.next().Kind {
	case token.IDENTIFIER:
		return p.parseTypeName()
	case token.BANG:
		return &ast.ErrorType{
			Type: nil,
			Bang: p.consume(),
		}
	case token.LEFT_PAREN:
		return p.parseTupleType()
	case token.LEFT_BRACE:
		return p.parseMapType()
	default:
		p.Diagnostics.ReportExpectedType(p.next().Location, p.next().Kind)
		return &ast.ErrorNode{}
	}
}

func (p *parser) parseArrayType(ty ast.TypeExpression) ast.TypeExpression {
	leftSquare := p.consume()
	var count ast.Expression
	if p.next().Kind != token.RIGHT_SQUARE {
		if p.next().Kind == token.IDENTIFIER && p.next().Value == "_" {
			count = &ast.InferredExpression{
				Token: p.consume(),
			}
		} else {
			count = p.parseExpression()
		}
	}
	rightSquare := p.expect(token.RIGHT_SQUARE)

	return &ast.ArrayType{
		Type:        ty,
		LeftSquare:  leftSquare,
		Count:       count,
		RightSquare: rightSquare,
	}
}

func (p *parser) parseMemberType(left ast.TypeExpression) ast.TypeExpression {
	dot := p.consume()
	member := p.expect(token.IDENTIFIER)

	return &ast.MemberType{
		Left:   left,
		Dot:    dot,
		Member: member,
	}
}

func (p *parser) parsePointerType() ast.TypeExpression {
	star := p.consume()
	var mut *token.Token
	if p.isKeyword("mut") {
		tok := p.consume()
		mut = &tok
	}
	ty := p.parsePrefixType()

	return &ast.PointerType{
		Star: star,
		Mut:  mut,
		Type: ty,
	}
}

func (p *parser) parseTupleType() ast.TypeExpression {
	leftParen := p.consume()
	types, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseType)

	if len(types) == 1 {
		return &ast.ParenthesisedType{
			LeftParen:  leftParen,
			Type:       types[0],
			RightParen: rightParen,
		}
	}

	return &ast.TupleType{
		LeftParen:  leftParen,
		Types:      types,
		RightParen: rightParen,
	}
}

func (p *parser) parseMapType() ast.TypeExpression {
	leftBrace := p.consume()
	keyType := p.parseType()
	colon := p.expect(token.COLON)
	valueType := p.parseType()
	rightBrace := p.expect(token.RIGHT_BRACE)

	return &ast.MapType{
		LeftBrace:  leftBrace,
		KeyType:    keyType,
		Colon:      colon,
		ValueType:  valueType,
		RightBrace: rightBrace,
	}
}

func (p *parser) parseTypeName() ast.TypeExpression {
	name := p.expect(token.IDENTIFIER)

	return &ast.TypeName{
		Name: name,
	}
}
