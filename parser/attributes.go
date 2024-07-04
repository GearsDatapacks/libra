package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseTagAttribute() (ast.Attribute, *diagnostics.Diagnostic) {
	tok := p.consume()
	name := p.expect(token.IDENTIFIER)
	return &ast.TagAttribute{Token: tok, Name: name.Value}, nil
}

func (p *parser) parseImplAttribute() (ast.Attribute, *diagnostics.Diagnostic) {
	tok := p.consume()
	name := p.expect(token.IDENTIFIER)
	return &ast.ImplAttribute{Token: tok, Name: name.Value}, nil
}

func (p *parser) parseUntaggedAttribue() (ast.Attribute, *diagnostics.Diagnostic) {
	tok := p.consume()
	return &ast.UntaggedAttribute{Token: tok}, nil
}

func (p *parser) parseTodoAttribue() (ast.Attribute, *diagnostics.Diagnostic) {
	tok := p.consume()
	text := ""
	if p.next().Kind == token.ATTRIBUTE_BODY {
		text = p.consume().Value
	}
	return &ast.TodoAttribute{Token:   tok, Message: text}, nil
}

func (p *parser) parseDocAttribue() (ast.Attribute, *diagnostics.Diagnostic) {
	tok := p.consume()
	text := ""
	if p.next().Kind == token.ATTRIBUTE_BODY {
		text = p.consume().Value
	}
	return &ast.DocAttribute{Token:   tok, Message: text}, nil
}

func (p *parser) parseDeprecatedAttribue() (ast.Attribute, *diagnostics.Diagnostic) {
	tok := p.consume()
	text := ""
	if p.next().Kind == token.ATTRIBUTE_BODY {
		text = p.consume().Value
	}
	return &ast.DeprecatedAttribute{Token:   tok, Message: text}, nil
}
