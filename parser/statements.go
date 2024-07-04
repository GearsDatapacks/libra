package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseTopLevelStatement() (ast.Statement, *diagnostics.Diagnostic) {
	attributes := []ast.Attribute{}
	for p.next().Kind == token.ATTRIBUTE_NAME {
		found := false
		name := p.next().Value[1:]
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
			p.Diagnostics.Report(diagnostics.InvalidAttribute(p.next().Location, p.next().Value[1:]))
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
					p.Diagnostics.Report(diagnostics.CannotAttribute(stmt.Tokens()[0].Location, attribute.GetName()))
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

	equals := p.expect(token.EQUALS)
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.VariableDeclaration{
		Keyword:    keyword,
		Identifier: identifier,
		Type:       typeAnnotation,
		Equals:     equals,
		Value:      value,
	}, nil
}

func (p *parser) parserOptionalDefaultValue() (*ast.DefaultValue, *diagnostics.Diagnostic) {
	if p.next().Kind != token.EQUALS {
		return nil, nil
	}

	equals := p.consume()
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.DefaultValue{
		Equals: equals,
		Value:  value,
	}, nil
}

func (p *parser) parseTypeOrIdent(declare bool) (*ast.TypeOrIdent, *diagnostics.Diagnostic) {
	var name, colon *token.Token
	var ty ast.Expression

	initial, err := p.parseTypeExpression()
	if err != nil {
		return nil, err
	}
	if ident, ok := initial.(*ast.Identifier); ok {
		name = &ident.Token
		if declare {
			p.identifiers[ident.Token.Value] = ident.Token.Location
		}

		if p.next().Kind == token.COLON {
			tok := p.consume()
			colon = &tok
			ty, err = p.parseTypeExpression()
			if err != nil {
				return nil, err
			}
		}
	} else {
		ty = initial
	}

	return &ast.TypeOrIdent{
		Name:  name,
		Colon: colon,
		Type:  ty,
	}, nil
}

func (p *parser) parseParameter() (*ast.Parameter, *diagnostics.Diagnostic) {
	var mutable *token.Token
	var value *ast.DefaultValue
	if p.isKeyword("mut") {
		tok := p.consume()
		mutable = &tok
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

	} else if mutable != nil {
		p.Diagnostics.Report(diagnostics.MutWithoutParamName(mutable.Location))
	}

	return &ast.Parameter{
		Mutable:     mutable,
		TypeOrIdent: *name,
		Default:     value,
	}, nil
}

func (p *parser) parseFunctionDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	var methodOf *ast.MethodOf
	var memberOf *ast.MemberOf

	if p.next().Kind == token.LEFT_PAREN {
		leftParen := p.consume()
		var mutable *token.Token
		if p.isKeyword("mut") {
			tok := p.consume()
			mutable = &tok
		}
		ty, err := p.parseTypeExpression()
		if err != nil {
			return nil, err
		}

		rightParen := p.expect(token.RIGHT_PAREN)

		methodOf = &ast.MethodOf{
			LeftParen:  leftParen,
			Mutable:    mutable,
			Type:       ty,
			RightParen: rightParen,
		}
	}

	name := p.expect(token.IDENTIFIER)

	if p.next().Kind == token.DOT {
		if methodOf != nil {
			p.Diagnostics.Report(diagnostics.MemberAndMethodNotAllowed(name.Location))
		}

		dot := p.consume()
		memberOf = &ast.MemberOf{
			Name: name,
			Dot:  dot,
		}
		name = p.expect(token.IDENTIFIER)
	} else if methodOf == nil {
		p.identifiers[name.Value] = name.Location
	}

	leftParen := p.expect(token.LEFT_PAREN)
	defer p.exitScope(p.enterScope())
	params, rightParen := parseDerefExprList(p, token.RIGHT_PAREN, p.parseParameter)

	if len(params) > 0 {
		lastParam := params[len(params)-1]
		if lastParam.TypeOrIdent.Type == nil && lastParam.Default == nil {
			p.Diagnostics.ReportMany(diagnostics.LastParameterMustHaveType(lastParam.TypeOrIdent.Name.Location, name.Location))
		}
	}

	returnType, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}
	body, err := p.parseBlock(true)
	if err != nil {
		return nil, err
	}

	return &ast.FunctionDeclaration{
		Keyword:    keyword,
		MethodOf:   methodOf,
		MemberOf:   memberOf,
		Name:       name,
		LeftParen:  leftParen,
		Parameters: params,
		RightParen: rightParen,
		ReturnType: returnType,
		Body:       body,
	}, nil
}

func (p *parser) parseReturnStatement() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	var value ast.Expression
	if !p.eof() && p.canContinue() {
		var err *diagnostics.Diagnostic
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}
	return &ast.ReturnStatement{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *parser) parseYieldStatement() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.YieldStatement{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *parser) parseBreakStatement() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	var value ast.Expression
	if !p.eof() && p.canContinue() {
		var err *diagnostics.Diagnostic
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}
	return &ast.BreakStatement{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *parser) parseTypeDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	name := p.delcareIdentifier()
	equals := p.expect(token.EQUALS)
	ty, err := p.parseTypeExpression()
	if err != nil {
		return nil, err
	}

	return &ast.TypeDeclaration{
		Keyword: keyword,
		Name:    name,
		Equals:  equals,
		Type:    ty,
	}, nil
}

func (p *parser) parseStructField() (*ast.StructField, *diagnostics.Diagnostic) {
	var pub *token.Token
	if p.isKeyword("pub") {
		tok := p.consume()
		pub = &tok
	}
	typeOrIdent, err := p.parseTypeOrIdent(false)
	if err != nil {
		return nil, err
	}

	return &ast.StructField{
		Pub:         pub,
		TypeOrIdent: *typeOrIdent,
	}, nil
}

func (p *parser) parseStructDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	name := p.delcareIdentifier()
	structDecl := &ast.StructDeclaration{
		Keyword: keyword,
		Name:    name,
	}

	if p.canContinue() && p.next().Kind == token.LEFT_BRACE {
		leftBrace := p.consume()
		fields, rightBrace := parseDerefExprList(p, token.RIGHT_BRACE, p.parseStructField)

		structDecl.Body = &ast.StructBody{
			LeftBrace:  leftBrace,
			Fields:     fields,
			RightBrace: rightBrace,
		}
	}

	return structDecl, nil
}

func (p *parser) parseInterfaceMember() (*ast.InterfaceMember, *diagnostics.Diagnostic) {
	name := p.expect(token.IDENTIFIER)
	leftParen := p.expect(token.LEFT_PAREN)
	params, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseExpression)
	returnType, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}

	return &ast.InterfaceMember{
		Name:       name,
		LeftParen:  leftParen,
		Parameters: params,
		RightParen: rightParen,
		ReturnType: returnType,
	}, nil
}

func (p *parser) parseInterfaceDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	name := p.delcareIdentifier()
	leftBrace := p.expect(token.LEFT_BRACE)
	members, rightBrace := parseDerefExprList(p, token.RIGHT_BRACE, p.parseInterfaceMember)

	return &ast.InterfaceDeclaration{
		Keyword:    keyword,
		Name:       name,
		LeftBrace:  leftBrace,
		Members:    members,
		RightBrace: rightBrace,
	}, nil
}

func (p *parser) parseImportStatement() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	var symbols *ast.ImportedSymbols

	if p.next().Kind == token.LEFT_BRACE {
		symbols = &ast.ImportedSymbols{}

		symbols.LeftBrace = p.consume()
		symbols.Symbols, symbols.RightBrace = parseDelimExprList(
			p, token.RIGHT_BRACE,
			func() (token.Token, *diagnostics.Diagnostic) { return p.delcareIdentifier(), nil },
		)
		symbols.From = p.expectKeyword("from")
	}

	var all *ast.ImportAll

	if p.next().Kind == token.STAR {
		if symbols != nil {
			p.Diagnostics.Report(diagnostics.OneImportModifierAllowed(p.next().Location))
		}
		all = &ast.ImportAll{}

		all.Star = p.consume()
		all.From = p.expectKeyword("from")
	}

	module := p.expect(token.STRING)

	var alias *ast.ImportAlias

	if p.canContinue() && p.isKeyword("as") {
		if symbols != nil || all != nil {
			p.Diagnostics.Report(diagnostics.OneImportModifierAllowed(p.next().Location))
		}
		alias = &ast.ImportAlias{}
		alias.As = p.consume()
		alias.Alias = p.delcareIdentifier()
	}

	return &ast.ImportStatement{
		Keyword: keyword,
		Symbols: symbols,
		All:     all,
		Module:  module,
		Alias:   alias,
	}, nil
}

func (p *parser) parseEnumMember() (*ast.EnumMember, *diagnostics.Diagnostic) {
	name := p.expect(token.IDENTIFIER)
	value, err := p.parserOptionalDefaultValue()
	if err != nil {
		return nil, err
	}

	return &ast.EnumMember{
		Name:  name,
		Value: value,
	}, nil
}

func (p *parser) parseEnumDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	name := p.delcareIdentifier()
	valueType, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}

	leftBrace := p.expect(token.LEFT_BRACE)
	members, rightBrace := parseDerefExprList(p, token.RIGHT_BRACE, p.parseEnumMember)

	return &ast.EnumDeclaration{
		Keyword:    keyword,
		Name:       name,
		ValueType:  valueType,
		LeftBrace:  leftBrace,
		Members:    members,
		RightBrace: rightBrace,
	}, nil
}

func (p *parser) parseUnionMember() (*ast.UnionMember, *diagnostics.Diagnostic) {
	name := p.expect(token.IDENTIFIER)
	ty, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}

	var compound *ast.StructBody
	if ty == nil && p.next().Kind == token.LEFT_BRACE {
		leftBrace := p.consume()
		fields, rightBrace := parseDerefExprList(p, token.RIGHT_BRACE, p.parseStructField)
		compound = &ast.StructBody{
			LeftBrace:  leftBrace,
			Fields:     fields,
			RightBrace: rightBrace,
		}
	}

	return &ast.UnionMember{
		Name:     name,
		Type:     ty,
		Compound: compound,
	}, nil
}

func (p *parser) parseUnionDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	name := p.delcareIdentifier()
	leftBrace := p.expect(token.LEFT_BRACE)
	members, rightBrace := parseDerefExprList(p, token.RIGHT_BRACE, p.parseUnionMember)

	return &ast.UnionDeclaration{
		Keyword:    keyword,
		Name:       name,
		LeftBrace:  leftBrace,
		Members:    members,
		RightBrace: rightBrace,
	}, nil
}

func (p *parser) parseTagDeclaration() (ast.Statement, *diagnostics.Diagnostic) {
	keyword := p.consume()
	name := p.delcareIdentifier()

	return &ast.TagDeclaration{
		Keyword: keyword,
		Name:    name,
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
