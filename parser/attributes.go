package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseIdentifierAttribute() (ast.Attribute, *diagnostics.Diagnostic) {
	tok := p.consume()
	name := p.expect(token.IDENTIFIER)
	return &ast.TextAttribute{Token: tok, Text: name.Value}, nil
}

func (p *parser) parseFlagAttribute() (ast.Attribute, *diagnostics.Diagnostic) {
	tok := p.consume()
	return &ast.FlagAttribute{Token: tok}, nil
}

func (p *parser) parseAttributeWithOptionalBody() (ast.Attribute, *diagnostics.Diagnostic) {
	tok := p.consume()
	text := ""
	if p.next().Kind == token.ATTRIBUTE_BODY {
		text = p.consume().Value
	}
	return &ast.TextAttribute{Token: tok, Text: text}, nil
}
