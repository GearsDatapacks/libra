package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
)

func (p *parser) parseTopLevelStatement() (ast.Statement, *diagnostics.Diagnostic) {
	attributes := []ast.Attribute{}
	for p.next().Kind == token.ATTRIBUTE_NAME {
		found := false
		name := p.next().ExtraValue
		for _, attr := range p.attributes {
			if name == attr.Name {
				attribute, err := attr.Fn()
				if err != nil {
					p.Diagnostics.Report(err)
					p.consumeUntil(token.NEWLINE, token.SEMICOLON)
				} else {
					attributes = append(attributes, attribute)
				}
				found = true
				break
			}
		}

		if !found {
			p.Diagnostics.Report(diagnostics.InvalidAttribute(p.next().Location, p.next().ExtraValue))
			p.consumeUntil(token.NEWLINE, token.SEMICOLON)
		}
	}

	for _, kwd := range p.keywords {
		if p.isKeyword(kwd.Name) {
			stmt, err := kwd.Fn()
			if err != nil {
				return nil, err
			}

			for _, attribute := range attributes {
				if !ast.TryAddAttribute(stmt, attribute) {
					p.Diagnostics.Report(diagnostics.CannotAttribute(stmt.GetLocation(), attribute.GetName()))
				}
			}

			return stmt, nil
		}
	}

	return p.parseStatement()
}

func (p *parser) parseStatement() (ast.Statement, *diagnostics.Diagnostic) {
	for _, kwd := range p.keywords {
		if p.isKeyword(kwd.Name) {
			if kwd.Kind == decl {
				p.Diagnostics.Report(diagnostics.OnlyTopLevelStatement(p.next().Location, kwd.StmtName))
			}
			return kwd.Fn()
		}
	}

	return p.parseExpression()
}

func (p *parser) parseVariableDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	identifier := p.delcareIdentifier()

	typeAnnotation, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}

	p.expect(token.EQUALS)
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.VariableDeclaration{
		Keyword:      keyword,
		NameLocation: identifier.Location,
		Name:         identifier.Value,
		Type:         typeAnnotation,
		Value:        value,
	}, nil
}

func (p *parser) parserOptionalDefaultValue() (ast.Expression, *diagnostics.Diagnostic) {
	if p.next().Kind != token.EQUALS {
		return nil, nil
	}

	p.consume()
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (p *parser) parseTypeOrIdent(declare bool) (*ast.TypeOrIdent, *diagnostics.Diagnostic) {
	var name *string
	var ty ast.Expression
	var location text.Location

	initial, err := p.parseTypeExpression()
	if err != nil {
		return nil, err
	}
	if ident, ok := initial.(*ast.Identifier); ok {
		name = &ident.Name
		location = ident.Location
		if declare {
			p.identifiers[ident.Name] = ident.Location
		}

		if p.next().Kind == token.COLON {
			p.consume()
			ty, err = p.parseTypeExpression()
			if err != nil {
				return nil, err
			}
		}
	} else {
		location = initial.GetLocation()
		ty = initial
	}

	return &ast.TypeOrIdent{
		Location: location,
		Name:     name,
		Type:     ty,
	}, nil
}

func (p *parser) parseParameter() (*ast.Parameter, *diagnostics.Diagnostic) {
	var value ast.Expression
	var location text.Location
	mutable := p.isKeyword("mut")

	if mutable {
		location = p.consume().Location
	}

	name, err := p.parseTypeOrIdent(true)
	if err != nil {
		return nil, err
	}

	if name.Name != nil {
		value, err = p.parserOptionalDefaultValue()
		if err != nil {
			return nil, err
		}

	} else if mutable {
		p.Diagnostics.Report(diagnostics.MutWithoutParamName(location))
	}

	if !mutable {
		location = name.Location
	}

	return &ast.Parameter{
		Location:    location,
		Mutable:     mutable,
		TypeOrIdent: *name,
		Default:     value,
	}, nil
}

func (p *parser) parseFunctionDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	var methodOf *ast.MethodOf
	var memberOf *ast.MemberOf

	if p.next().Kind == token.LEFT_PAREN {
		p.consume()
		mutable := p.isKeyword("mut")
		if mutable {
			p.consume()
		}
		ty, err := p.parseTypeExpression()
		if err != nil {
			return nil, err
		}

		p.expect(token.RIGHT_PAREN)

		methodOf = &ast.MethodOf{
			Mutable: mutable,
			Type:    ty,
		}
	}

	name := p.expect(token.IDENTIFIER)

	if p.next().Kind == token.DOT {
		if methodOf != nil {
			p.Diagnostics.Report(diagnostics.MemberAndMethodNotAllowed(name.Location))
		}

		p.consume()
		memberOf = &ast.MemberOf{
			Location: name.Location,
			Name:     name.Value,
		}
		name = p.expect(token.IDENTIFIER)
	} else if methodOf == nil {
		p.identifiers[name.Value] = name.Location
	}

	p.expect(token.LEFT_PAREN)
	defer p.exitScope(p.enterScope())
	params := parseDerefExprList(p, token.RIGHT_PAREN, p.parseParameter)

	if len(params) > 0 {
		lastParam := params[len(params)-1]
		if lastParam.TypeOrIdent.Type == nil && lastParam.Default == nil {
			p.Diagnostics.ReportMany(diagnostics.LastParameterMustHaveType(lastParam.TypeOrIdent.Location, name.Location))
		}
	}

	returnType, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}

	var body *ast.Block

	if p.canContinue() && p.next().Kind == token.LEFT_BRACE {
		block, err := p.parseBlock(true)
		if err != nil {
			return nil, err
		}
		body = block
	}

	return &ast.FunctionDeclaration{
		Location:     location,
		NameLocation: name.Location,
		MethodOf:     methodOf,
		MemberOf:     memberOf,
		Name:         name.Value,
		Parameters:   params,
		ReturnType:   returnType,
		Body:         body,
	}, nil
}

func (p *parser) parseReturnStatement() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	var value ast.Expression
	if !p.eof() && p.canContinue() {
		var err *diagnostics.Diagnostic
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}
	return &ast.ReturnStatement{
		Location: location,
		Value:    value,
	}, nil
}

func (p *parser) parseYieldStatement() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.YieldStatement{
		Location: location,
		Value:    value,
	}, nil
}

func (p *parser) parseBreakStatement() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	var value ast.Expression
	if !p.eof() && p.canContinue() {
		var err *diagnostics.Diagnostic
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}
	return &ast.BreakStatement{
		Location: location,
		Value:    value,
	}, nil
}

func (p *parser) parseTypeDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	name := p.delcareIdentifier().Value
	p.expect(token.EQUALS)
	ty, err := p.parseTypeExpression()
	if err != nil {
		return nil, err
	}

	return &ast.TypeDeclaration{
		Location:   location,
		Name:       name,
		Type:       ty,
		Tag:        nil,
		Attributes: ast.DeclarationAttributes{},
	}, nil
}

func (p *parser) parseStructField() (*ast.StructField, *diagnostics.Diagnostic) {
	var location text.Location
	pub := p.isKeyword("pub")
	if pub {
		location = p.consume().Location
	}
	typeOrIdent, err := p.parseTypeOrIdent(false)
	if err != nil {
		return nil, err
	}

	if !pub {
		location = typeOrIdent.Location
	}

	return &ast.StructField{
		Location:    location,
		Pub:         pub,
		TypeOrIdent: *typeOrIdent,
	}, nil
}

func (p *parser) parseStructDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	name := p.delcareIdentifier()
	var body []ast.StructField

	if p.canContinue() && p.next().Kind == token.LEFT_BRACE {
		p.consume()
		body = parseDerefExprList(p, token.RIGHT_BRACE, p.parseStructField)
	}

	return &ast.StructDeclaration{
		Location:     location,
		NameLocation: name.Location,
		Name:         name.Value,
		Body:         body,
	}, nil
}

func (p *parser) parseInterfaceMember() (*ast.InterfaceMember, *diagnostics.Diagnostic) {
	name := p.expect(token.IDENTIFIER).Value
	p.expect(token.LEFT_PAREN)
	params := parseDelimExprList(p, token.RIGHT_PAREN, p.parseExpression)
	returnType, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}

	return &ast.InterfaceMember{
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
	}, nil
}

func (p *parser) parseInterfaceDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	name := p.delcareIdentifier().Value
	p.expect(token.LEFT_BRACE)
	members := parseDerefExprList(p, token.RIGHT_BRACE, p.parseInterfaceMember)

	return &ast.InterfaceDeclaration{
		Location:   location,
		Name:       name,
		Members:    members,
		Attributes: ast.DeclarationAttributes{},
	}, nil
}

func (p *parser) parseImportStatement() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	var symbols []ast.ImportedSymbol

	if p.next().Kind == token.LEFT_BRACE {
		p.consume()
		symbols = parseDelimExprList(
			p, token.RIGHT_BRACE,
			func() (ast.ImportedSymbol, *diagnostics.Diagnostic) {
				ident := p.delcareIdentifier()
				return ast.ImportedSymbol{
					Location: ident.Location,
					Name:     ident.Value,
				}, nil
			},
		)
		p.expectKeyword("from")
	}

	all := p.next().Kind == token.STAR

	if all {
		if symbols != nil {
			p.Diagnostics.Report(diagnostics.OneImportModifierAllowed(p.next().Location))
		}

		p.consume()
		p.expectKeyword("from")
	}

	module := p.expect(token.STRING)

	var alias *string

	if p.canContinue() && p.isKeyword("as") {
		if symbols != nil || all {
			p.Diagnostics.Report(diagnostics.OneImportModifierAllowed(p.next().Location))
		}
		p.consume()
		name := p.delcareIdentifier().Value
		alias = &name
	}

	return &ast.ImportStatement{
		Location: location,
		Symbols:  symbols,
		All:      all,
		Module:   module,
		Alias:    alias,
	}, nil
}

func (p *parser) parseEnumMember() (*ast.EnumMember, *diagnostics.Diagnostic) {
	name := p.expect(token.IDENTIFIER)
	value, err := p.parserOptionalDefaultValue()
	if err != nil {
		return nil, err
	}

	return &ast.EnumMember{
		Name:     name.Value,
		Location: name.Location,
		Value:    value,
	}, nil
}

func (p *parser) parseEnumDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	name := p.delcareIdentifier().Value
	valueType, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}

	p.expect(token.LEFT_BRACE)
	members := parseDerefExprList(p, token.RIGHT_BRACE, p.parseEnumMember)

	return &ast.EnumDeclaration{
		Location:  location,
		Name:      name,
		ValueType: valueType,
		Members:   members,
	}, nil
}

func (p *parser) parseUnionMember() (*ast.UnionMember, *diagnostics.Diagnostic) {
	name := p.expect(token.IDENTIFIER)
	ty, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}

	var compound []ast.StructField
	if ty == nil && p.next().Kind == token.LEFT_BRACE {
		p.consume()
		compound = parseDerefExprList(p, token.RIGHT_BRACE, p.parseStructField)
	}

	return &ast.UnionMember{
		NameLocation: name.Location,
		Name:         name.Value,
		Type:         ty,
		Compound:     compound,
	}, nil
}

func (p *parser) parseUnionDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	name := p.delcareIdentifier().Value
	p.expect(token.LEFT_BRACE)
	members := parseDerefExprList(p, token.RIGHT_BRACE, p.parseUnionMember)

	return &ast.UnionDeclaration{
		Location: location,
		Name:     name,
		Members:  members,
		Untagged: false,
	}, nil
}

func (p *parser) parseTagDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	location := p.consume().Location
	name := p.delcareIdentifier().Value
	var body []ast.Expression

	if p.canContinue() && p.next().Kind == token.LEFT_BRACE {
		p.consume()
		body = parseDelimExprList(p, token.RIGHT_BRACE, p.parseTypeExpression)
	}

	return &ast.TagDeclaration{
		Location: location,
		Name:     name,
		Body:     body,
	}, nil
}

func (p *parser) parsePubStatement() (ast.Statement, *diagnostics.Diagnostic) {
	pub := p.consume()
	stmt, err := p.parseTopLevelStatement()
	if err != nil {
		return nil, err
	}

	if decl, ok := stmt.(ast.Declaration); ok {
		decl.MarkExport()
		return decl, nil
	}
	p.Diagnostics.Report(diagnostics.CannotExport(pub.Location))
	return stmt, nil
}

func (p *parser) parseExplStatement() (ast.Statement, *diagnostics.Diagnostic) {
	kwd := p.consume()
	stmt, err := p.parseTopLevelStatement()
	if err != nil {
		return nil, err
	}
	if expl, ok := stmt.(ast.Explicit); ok {
		expl.MarkExplicit()
		return expl, nil
	}
	p.Diagnostics.Report(diagnostics.CannotExplicit(kwd.Location))
	return stmt, nil
}
