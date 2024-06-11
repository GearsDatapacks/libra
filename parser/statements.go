package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseTopLevelStatement() ast.Statement {
	for _, kwd := range p.keywords {
		if p.isKeyword(kwd.Name) {
			if kwd.Kind == decl {
				return kwd.Fn()
			} else {
				return kwd.Fn()
			}
		}
	}

	return p.parseStatement()
}

func (p *parser) parseStatement() ast.Statement {
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

func (p *parser) parseVariableDeclaration() ast.Statement {
	keyword := p.consume()
	identifier := p.delcareIdentifier()

	typeAnnotation := p.parseOptionalTypeAnnotation()

	equals := p.expect(token.EQUALS)
	value := p.parseExpression()

	return &ast.VariableDeclaration{
		Keyword:    keyword,
		Identifier: identifier,
		Type:       typeAnnotation,
		Equals:     equals,
		Value:      value,
	}
}

func (p *parser) parserOptionalDefaultValue() *ast.DefaultValue {
	if p.next().Kind != token.EQUALS {
		return nil
	}

	equals := p.consume()
	value := p.parseExpression()
	return &ast.DefaultValue{
		Equals: equals,
		Value:  value,
	}
}

func (p *parser) parseTypeOrIdent(declare bool) ast.TypeOrIdent {
	var name, colon *token.Token
	var ty ast.Expression

	initial := p.parseTypeExpression()
	if ident, ok := initial.(*ast.Identifier); ok {
		name = &ident.Token
		if declare {
			p.identifiers[ident.Token.Value] = ident.Token.Location
		}

		if p.next().Kind == token.COLON {
			tok := p.consume()
			colon = &tok
			ty = p.parseTypeExpression()
		}
	} else {
		ty = initial
	}

	return ast.TypeOrIdent{
		Name:  name,
		Colon: colon,
		Type:  ty,
	}
}

func (p *parser) parseParameter() ast.Parameter {
	var mutable *token.Token
	var value *ast.DefaultValue
	if p.isKeyword("mut") {
		tok := p.consume()
		mutable = &tok
	}

	name := p.parseTypeOrIdent(true)
	if name.Name != nil {
		value = p.parserOptionalDefaultValue()
	} else if mutable != nil {
		p.Diagnostics.Report(diagnostics.MutWithoutParamName(mutable.Location))
	}

	return ast.Parameter{
		Mutable:     mutable,
		TypeOrIdent: name,
		Default:     value,
	}
}

func (p *parser) parseFunctionDeclaration() ast.Statement {
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
		ty := p.parseTypeExpression()
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
	params, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseParameter)

	if len(params) > 0 {
		lastParam := params[len(params)-1]
		if lastParam.TypeOrIdent.Type == nil && lastParam.Default == nil {
			p.Diagnostics.Report(diagnostics.LastParameterMustHaveType(lastParam.TypeOrIdent.Name.Location, name.Location)...)
		}
	}

	returnType := p.parseOptionalTypeAnnotation()
	body := p.parseBlock(true)

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
	}
}

func (p *parser) parseReturnStatement() ast.Statement {
	keyword := p.consume()
	var value ast.Expression
	if !p.eof() && p.canContinue() {
		value = p.parseExpression()
	}
	return &ast.ReturnStatement{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *parser) parseYieldStatement() ast.Statement {
	keyword := p.consume()
	value := p.parseExpression()

	return &ast.YieldStatement{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *parser) parseBreakStatement() ast.Statement {
	keyword := p.consume()
	var value ast.Expression
	if !p.eof() && p.canContinue() {
		value = p.parseExpression()
	}
	return &ast.BreakStatement{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *parser) parseTypeDeclaration() ast.Statement {
	keyword := p.consume()
	name := p.delcareIdentifier()
	equals := p.expect(token.EQUALS)
	ty := p.parseTypeExpression()

	return &ast.TypeDeclaration{
		Keyword: keyword,
		Name:    name,
		Equals:  equals,
		Type:    ty,
	}
}

func (p *parser) parseStructField() ast.StructField {
	var pub *token.Token
	if p.isKeyword("pub") {
		tok := p.consume()
		pub = &tok
	}
	typeOrIdent := p.parseTypeOrIdent(false)

	return ast.StructField{
		Pub:         pub,
		TypeOrIdent: typeOrIdent,
	}
}

func (p *parser) parseStructDeclaration() ast.Statement {
	keyword := p.consume()
	name := p.delcareIdentifier()
	structDecl := &ast.StructDeclaration{
		Keyword: keyword,
		Name:    name,
	}

	if p.canContinue() && p.next().Kind == token.LEFT_BRACE {
		leftBrace := p.consume()
		fields, rightBrace := parseDelimExprList(p, token.RIGHT_BRACE, p.parseStructField)

		structDecl.Body = &ast.StructBody{
			LeftBrace:  leftBrace,
			Fields:     fields,
			RightBrace: rightBrace,
		}
	}

	return structDecl
}

func (p *parser) parseInterfaceMember() ast.InterfaceMember {
	name := p.expect(token.IDENTIFIER)
	leftParen := p.expect(token.LEFT_PAREN)
	params, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseExpression)
	returnType := p.parseOptionalTypeAnnotation()

	return ast.InterfaceMember{
		Name:       name,
		LeftParen:  leftParen,
		Parameters: params,
		RightParen: rightParen,
		ReturnType: returnType,
	}
}

func (p *parser) parseInterfaceDeclaration() ast.Statement {
	keyword := p.consume()
	name := p.delcareIdentifier()
	leftBrace := p.expect(token.LEFT_BRACE)
	members, rightBrace := parseDelimExprList(p, token.RIGHT_BRACE, p.parseInterfaceMember)

	return &ast.InterfaceDeclaration{
		Keyword:    keyword,
		Name:       name,
		LeftBrace:  leftBrace,
		Members:    members,
		RightBrace: rightBrace,
	}
}

func (p *parser) parseImportStatement() ast.Statement {
	keyword := p.consume()
	var symbols *ast.ImportedSymbols

	if p.next().Kind == token.LEFT_BRACE {
		symbols = &ast.ImportedSymbols{}

		symbols.LeftBrace = p.consume()
		symbols.Symbols, symbols.RightBrace = parseDelimExprList(p, token.RIGHT_BRACE, p.delcareIdentifier)
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
	}
}

func (p *parser) parseEnumMember() ast.EnumMember {
	name := p.expect(token.IDENTIFIER)
	value := p.parserOptionalDefaultValue()

	return ast.EnumMember{
		Name:  name,
		Value: value,
	}
}

func (p *parser) parseEnumDeclaration() ast.Statement {
	keyword := p.consume()
	name := p.delcareIdentifier()
	valueType := p.parseOptionalTypeAnnotation()
	leftBrace := p.expect(token.LEFT_BRACE)
	members, rightBrace := parseDelimExprList(p, token.RIGHT_BRACE, p.parseEnumMember)

	return &ast.EnumDeclaration{
		Keyword:    keyword,
		Name:       name,
		ValueType:  valueType,
		LeftBrace:  leftBrace,
		Members:    members,
		RightBrace: rightBrace,
	}
}

func (p *parser) parseUnionMember() ast.UnionMember {
	name := p.expect(token.IDENTIFIER)
	ty := p.parseOptionalTypeAnnotation()
	var compound *ast.StructBody
	if ty == nil && p.next().Kind == token.LEFT_BRACE {
		leftBrace := p.consume()
		fields, rightBrace := parseDelimExprList(p, token.RIGHT_BRACE, p.parseStructField)
		compound = &ast.StructBody{
			LeftBrace:  leftBrace,
			Fields:     fields,
			RightBrace: rightBrace,
		}
	}

	return ast.UnionMember{
		Name:     name,
		Type:     ty,
		Compound: compound,
	}
}

func (p *parser) parseUnionDeclaration() ast.Statement {
	keyword := p.consume()
	name := p.delcareIdentifier()
	leftBrace := p.expect(token.LEFT_BRACE)
	members, rightBrace := parseDelimExprList(p, token.RIGHT_BRACE, p.parseUnionMember)

	return &ast.UnionDeclaration{
		Keyword:    keyword,
		Name:       name,
		LeftBrace:  leftBrace,
		Members:    members,
		RightBrace: rightBrace,
	}
}

func (p *parser) parseTagDeclaration() ast.Statement {
	keyword := p.consume()
	name := p.delcareIdentifier()

	return &ast.TagDeclaration{
		Keyword: keyword,
		Name:    name,
	}
}

func (p *parser) parsePubStatement() ast.Statement {
	pub := p.consume()
	stmt := p.parseTopLevelStatement()
	if decl, ok := stmt.(ast.Declaration); ok {
		decl.MarkExport()
		return decl
	}
	p.Diagnostics.Report(diagnostics.CannotExport(pub.Location))
	return stmt
}

func (p *parser) parseExplStatement() ast.Statement {
	kwd := p.consume()
	stmt := p.parseTopLevelStatement()
	if expl, ok := stmt.(ast.Explicit); ok {
		expl.MarkExplicit()
		return expl
	}
	p.Diagnostics.Report(diagnostics.CannotExplicit(kwd.Location))
	return stmt
}
