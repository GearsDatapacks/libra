package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseTopLevelStatement() ast.Statement {
	if p.isKeyword("fn") {
		return p.parseFunctionDeclaration()
	}

	if p.isKeyword("type") {
		return p.parseTypeDeclaration()
	}

	if p.isKeyword("struct") {
		return p.parseStructDeclaration()
	}

	if p.isKeyword("interface") {
		return p.parseInterfaceDeclaration()
	}

	if p.isKeyword("import") {
		return p.parseImportStatement()
	}

	if p.isKeyword("enum") || p.isKeyword("union") {
		return p.parseEnumDeclaration()
	}

	return p.parseStatement()
}

func (p *parser) parseStatement() ast.Statement {
	if p.next().Kind == token.LEFT_BRACE {
		return p.parseBlockStatement()
	}

	if p.isKeyword("const") || p.isKeyword("let") || p.isKeyword("mut") {
		return p.parseVariableDeclaration()
	}

	if p.isKeyword("if") {
		return p.parseIfStatement()
	}

	if p.isKeyword("else") {
		p.Diagnostics.Report(diagnostics.ElseStatementWithoutIf(p.next().Location))
	}

	if p.isKeyword("while") {
		return p.parseWhileLoop()
	}

	if p.isKeyword("for") {
		return p.parseForLoop()
	}

	if p.isKeyword("break") {
		return &ast.BreakStatement{
			Keyword: p.consume(),
		}
	}

	if p.isKeyword("continue") {
		return &ast.ContinueStatement{
			Keyword: p.consume(),
		}
	}

	if p.isKeyword("fn") {
		p.Diagnostics.Report(diagnostics.OnlyTopLevelStatement(p.next().Location, "Function declaration"))
		return p.parseFunctionDeclaration()
	}

	if p.isKeyword("return") {
		return p.parseReturnStatement()
	}

	if p.isKeyword("type") {
		p.Diagnostics.Report(diagnostics.OnlyTopLevelStatement(p.next().Location, "Type declaration"))
		return p.parseTypeDeclaration()
	}

	if p.isKeyword("struct") {
		p.Diagnostics.Report(diagnostics.OnlyTopLevelStatement(p.next().Location, "Struct declaration"))
		return p.parseStructDeclaration()
	}

	if p.isKeyword("interface") {
		p.Diagnostics.Report(diagnostics.OnlyTopLevelStatement(p.next().Location, "Interface declaration"))
		return p.parseInterfaceDeclaration()
	}

	if p.isKeyword("import") {
		p.Diagnostics.Report(diagnostics.OnlyTopLevelStatement(p.next().Location, "Import statement"))
		return p.parseImportStatement()
	}

	if p.isKeyword("enum") || p.isKeyword("union") {
		p.Diagnostics.Report(diagnostics.OnlyTopLevelStatement(p.next().Location, p.next().Value+" declaration"))
		return p.parseEnumDeclaration()
	}

	return &ast.ExpressionStatement{
		Expression: p.parseExpression(),
	}
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

func (p *parser) parseBlockStatement(noScope ...bool) *ast.BlockStatement {
	leftBrace := p.expect(token.LEFT_BRACE)
	if len(noScope) == 0 || !noScope[0] {
		defer p.exitScope(p.enterScope())
	}
	statements, rightBrace := parseDelimStmtList(p, token.RIGHT_BRACE, p.parseStatement)

	return &ast.BlockStatement{
		LeftBrace:  leftBrace,
		Statements: statements,
		RightBrace: rightBrace,
	}
}

func (p *parser) parseIfStatement() ast.Statement {
	keyword := p.consume()

	p.noBraces = true
	p.bracketLevel++
	condition := p.parseSubExpression(Lowest)
	p.noBraces = false
	p.bracketLevel--

	body := p.parseBlockStatement()
	var elseBranch *ast.ElseBranch

	if p.isKeyword("else") {
		elseBranch = &ast.ElseBranch{}
		elseBranch.ElseKeyword = p.consume()
		if p.isKeyword("if") {
			elseBranch.Statement = p.parseIfStatement()
		} else {
			elseBranch.Statement = p.parseBlockStatement()
		}
	}

	return &ast.IfStatement{
		Keyword:    keyword,
		Condition:  condition,
		Body:       body,
		ElseBranch: elseBranch,
	}
}

func (p *parser) parseWhileLoop() ast.Statement {
	keyword := p.consume()

	p.noBraces = true
	p.bracketLevel++
	condition := p.parseSubExpression(Lowest)
	p.bracketLevel--
	p.noBraces = false

	body := p.parseBlockStatement()

	return &ast.WhileLoop{
		Keyword:   keyword,
		Condition: condition,
		Body:      body,
	}
}

func (p *parser) parseForLoop() ast.Statement {
	forKeyword := p.consume()
	defer p.exitScope(p.enterScope())

	variable := p.delcareIdentifier()
	inKeyword := p.expectKeyword("in")

	p.noBraces = true
	p.bracketLevel++
	iterator := p.parseSubExpression(Lowest)
	p.bracketLevel--
	p.noBraces = false

	body := p.parseBlockStatement(true)

	return &ast.ForLoop{
		ForKeyword: forKeyword,
		Variable:   variable,
		InKeyword:  inKeyword,
		Iterator:   iterator,
		Body:       body,
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

func (p *parser) parseParameter() ast.Parameter {
	var mutable *token.Token
	if p.isKeyword("mut") {
		tok := p.consume()
		mutable = &tok
	}

	name := p.delcareIdentifier()
	ty := p.parseOptionalTypeAnnotation()
	value := p.parserOptionalDefaultValue()

	return ast.Parameter{
		Mutable: mutable,
		Name:    name,
		Type:    ty,
		Default: value,
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
		ty := p.parseType()
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
		if lastParam.Type == nil && lastParam.Default == nil {
			p.Diagnostics.Report(diagnostics.LastParameterMustHaveType(lastParam.Name.Location, name.Location)...)
		}
	}

	returnType := p.parseOptionalTypeAnnotation()
	body := p.parseBlockStatement(true)

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

func (p *parser) parseTypeDeclaration() ast.Statement {
	keyword := p.consume()
	name := p.delcareIdentifier()
	equals := p.expect(token.EQUALS)
	ty := p.parseType()

	return &ast.TypeDeclaration{
		Keyword: keyword,
		Name:    name,
		Equals:  equals,
		Type:    ty,
	}
}

func (p *parser) parseStructField() ast.StructField {
	name := p.expect(token.IDENTIFIER)
	ty := p.parseOptionalTypeAnnotation()

	return ast.StructField{
		Name: name,
		Type: ty,
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

		if len(fields) > 0 {
			last := fields[len(fields)-1]
			if last.Type == nil {
				p.Diagnostics.Report(diagnostics.LastStructFieldMustHaveType(last.Name.Location, name.Location)...)
			}
		}

		structDecl.StructType = &ast.Struct{
			LeftBrace:  leftBrace,
			Fields:     fields,
			RightBrace: rightBrace,
		}
	} else if p.canContinue() && p.next().Kind == token.LEFT_PAREN {
		leftParen := p.consume()
		types, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseType)

		structDecl.TupleType = &ast.TupleStruct{
			LeftParen:  leftParen,
			Types:      types,
			RightParen: rightParen,
		}
	}

	return structDecl
}

func (p *parser) parseInterfaceMember() ast.InterfaceMember {
	name := p.expect(token.IDENTIFIER)
	leftParen := p.expect(token.LEFT_PAREN)
	params, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseType)
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

func (p *parser) parseEnumMember(isUnion bool) ast.EnumMember {
	var name token.Token
	var types *ast.TypeList
	var structType *ast.StructBody
	var value *ast.ValueAssignment

	if isUnion {
		name = p.delcareIdentifier()
	} else {
		name = p.expect(token.IDENTIFIER)
	}

	switch p.next().Kind {
	case token.LEFT_PAREN:
		types = &ast.TypeList{}
		types.LeftParen = p.consume()
		types.Types, types.RightParen = parseDelimExprList(p, token.RIGHT_PAREN, p.parseType)

	case token.LEFT_BRACE:
		structType = &ast.StructBody{}
		structType.LeftBrace = p.consume()
		structType.Fields, structType.RightBrace = parseDelimExprList(p, token.RIGHT_BRACE, p.parseStructField)
	case token.EQUALS:
		value = &ast.ValueAssignment{}
		value.Equals = p.consume()
		value.Value = p.parseExpression()
	}

	return ast.EnumMember{
		Name:   name,
		Types:  types,
		Struct: structType,
		Value:  value,
	}
}

func (p *parser) parseEnumDeclaration() ast.Statement {
	keyword := p.consume()
	isUnion := keyword.Value == "union"
	name := p.delcareIdentifier()
	valueType := p.parseOptionalTypeAnnotation()
	leftBrace := p.expect(token.LEFT_BRACE)
	members, rightBrace := parseDelimExprList(p,
		token.RIGHT_BRACE,
		func() ast.EnumMember { return p.parseEnumMember(isUnion) },
	)

	return &ast.EnumDeclaration{
		Keyword:    keyword,
		Name:       name,
		ValueType:  valueType,
		LeftBrace:  leftBrace,
		Members:    members,
		RightBrace: rightBrace,
	}
}
