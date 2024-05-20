package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseTopLevelStatement() ast.Statement {
	for _, kwd := range p.keywords {
		if p.isKeyword(kwd.Name) {
			return kwd.Fn()
		}
	}

	return p.parseStatement()
}

func (p *parser) parseStatement() ast.Statement {
	if p.next().Kind == token.LEFT_BRACE {
		return p.parseBlock()
	}

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

func (p *parser) parseBlock(noScope ...bool) *ast.Block {
	leftBrace := p.expect(token.LEFT_BRACE)
	if len(noScope) == 0 || !noScope[0] {
		defer p.exitScope(p.enterScope())
	}
	statements, rightBrace := parseDelimStmtList(p, token.RIGHT_BRACE, p.parseStatement)

	return &ast.Block{
		LeftBrace:  leftBrace,
		Statements: statements,
		RightBrace: rightBrace,
	}
}

func (p *parser) parseBlockExpression() ast.Expression {
	if p.isKeyword("do") {
		p.consume()
	}
	return p.parseBlock()
}

func (p *parser) parseIfExpression() ast.Expression {
	keyword := p.consume()

	p.noBraces = true
	p.bracketLevel++
	condition := p.parseSubExpression(Lowest)
	p.noBraces = false
	p.bracketLevel--

	body := p.parseBlock()
	var elseBranch *ast.ElseBranch

	if p.isKeyword("else") {
		elseBranch = &ast.ElseBranch{}
		elseBranch.ElseKeyword = p.consume()
		if p.isKeyword("if") {
			elseBranch.Statement = p.parseIfExpression()
		} else {
			elseBranch.Statement = p.parseBlock()
		}
	}

	return &ast.IfExpression{
		Keyword:    keyword,
		Condition:  condition,
		Body:       body,
		ElseBranch: elseBranch,
	}
}

func (p *parser) parseWhileLoop() ast.Expression {
	keyword := p.consume()

	p.noBraces = true
	p.bracketLevel++
	condition := p.parseSubExpression(Lowest)
	p.bracketLevel--
	p.noBraces = false

	body := p.parseBlock()

	return &ast.WhileLoop{
		Keyword:   keyword,
		Condition: condition,
		Body:      body,
	}
}

func (p *parser) parseForLoop() ast.Expression {
	forKeyword := p.consume()
	defer p.exitScope(p.enterScope())

	variable := p.delcareIdentifier()
	inKeyword := p.expectKeyword("in")

	p.noBraces = true
	p.bracketLevel++
	iterator := p.parseSubExpression(Lowest)
	p.bracketLevel--
	p.noBraces = false

	body := p.parseBlock(true)

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

	return ast.UnionMember{
		Name: name,
		Type: ty,
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
