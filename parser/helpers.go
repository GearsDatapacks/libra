package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func parseDelimExprList[Elem any](p *parser, delim token.Kind, elemFn func() (Elem, *diagnostics.Diagnostic)) (result []Elem, delimToken token.Token) {
	result = []Elem{}
	p.bracketLevel++

	for !p.eof() && p.next().Kind != delim {
		nextElem, err := elemFn()
		if err != nil {
			p.Diagnostics.Report(err)
			p.consumeUntil(token.COMMA, delim)
			continue
		}

		result = append(result, nextElem)

		if p.next().Kind == token.COMMA {
			p.consume()
		} else {
			break
		}
	}

	delimToken = p.expect(delim)
	p.bracketLevel--

	return result, delimToken
}

func parseDerefExprList[Elem any](p *parser, delim token.Kind, elemFn func() (*Elem, *diagnostics.Diagnostic)) (result []Elem, delimToken token.Token) {
	result = []Elem{}
	p.bracketLevel++

	for !p.eof() && p.next().Kind != delim {
		nextElem, err := elemFn()
		if err != nil {
			p.Diagnostics.Report(err)
			p.consumeUntil(token.COMMA, delim)
			continue
		}

		result = append(result, *nextElem)

		if p.next().Kind == token.COMMA {
			p.consume()
		} else {
			break
		}
	}

	delimToken = p.expect(delim)
	p.bracketLevel--

	return result, delimToken
}

func extendDelimExprList[Elem any](p *parser, first Elem, delim token.Kind, elemFn func() (*Elem, *diagnostics.Diagnostic)) (result []Elem, delimToken token.Token) {
	result = []Elem{first}
	p.bracketLevel++

	for !p.eof() && p.next().Kind != delim {
		if p.next().Kind == token.COMMA {
			p.consume()
		} else {
			break
		}

		if !p.eof() && p.next().Kind != delim {
			nextElem, err := elemFn()
			if err != nil {
				p.Diagnostics.Report(err)
				p.consumeUntil(token.COMMA, delim)
				continue
			}
	
			result = append(result, *nextElem)
		}
	}

	delimToken = p.expect(delim)
	p.bracketLevel--

	return result, delimToken
}

func parseDelimStmtList[Elem any](p *parser, delim token.Kind, elemFn func() (Elem, *diagnostics.Diagnostic)) (result []Elem, delimToken token.Token) {
	result = []Elem{}

	for !p.eof() && p.next().Kind != delim {
		nextElem, err := elemFn()
		if err != nil {
			p.Diagnostics.Report(err)
			p.consumeUntil(token.NEWLINE, token.SEMICOLON, delim)
			continue
		}

		result = append(result, nextElem)

		if p.next().Kind == token.COMMA {
			p.consume()
		} else {
			break
		}

		next := p.nextWithNewlines().Kind
		if next == token.NEWLINE || next == token.SEMICOLON {
			p.consumeNewlines()
		} else {
			break
		}
	}

	delimToken = p.expect(delim)

	return result, delimToken
}

func extendDelimStmtList[Elem any](p *parser, first Elem, delim token.Kind, elemFn func() (Elem, *diagnostics.Diagnostic)) (result []Elem, delimToken token.Token) {
	result = []Elem{first}

	for !p.eof() && p.next().Kind != delim {
		next := p.nextWithNewlines().Kind
		if next == token.NEWLINE || next == token.SEMICOLON {
			p.consumeNewlines()
		} else {
			break
		}

		if !p.eof() && p.next().Kind != delim {
			nextElem, err := elemFn()
			if err != nil {
				p.Diagnostics.Report(err)
				p.consumeUntil(token.NEWLINE, token.SEMICOLON, delim)
				continue
			}
	
			result = append(result, nextElem)
		}
	}

	delimToken = p.expect(delim)

	return result, delimToken
}

func (p *parser) parseOptionalTypeAnnotation() (*ast.TypeAnnotation, *diagnostics.Diagnostic) {
	if p.next().Kind != token.COLON {
		return nil, nil
	}

	colon := p.consume()
	ty, err := p.parseTypeExpression()
	if err != nil {
		return nil, err
	}

	return &ast.TypeAnnotation{
		Colon: colon,
		Type:  ty,
	}, nil
}

func (p *parser) parseTypeExpression() (ast.Expression, *diagnostics.Diagnostic) {
	old := p.typeExpr
	p.typeExpr = true
	ty, err := p.parseExpression()
	p.typeExpr = old
	return ty, err
}
