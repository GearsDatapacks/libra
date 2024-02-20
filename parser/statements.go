package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseTopLevelStatement() ast.Statement {
	switch p.next().Kind {
	default:
		return p.parseStatement()
	}
}

func (p *parser) parseStatement() ast.Statement {
	if p.isKeyword("const") || p.isKeyword("let") || p.isKeyword("mut") {
		return p.parseVariableDeclaration()
	}

	if p.isKeyword("if") {
		return p.parseIfStatement()
	}

	if p.isKeyword("else") {
		p.Diagnostics.ReportElseStatementWithoutIf(p.next().Span)
	}

	if p.isKeyword("while") {
		return p.parseWhileLoop()
	}

	if p.isKeyword("for") {
		return p.parseForLoop()
	}

	if p.isKeyword("fn") {
		return p.parseFunctionDeclaration()
	}

	if p.isKeyword("return") {
		return p.parseReturnStatement()
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

func (p *parser) parseParameter() ast.Parameter {
	name := p.delcareIdentifier()
	ty := p.parseOptionalTypeAnnotation()

	return ast.Parameter{
		Name: name,
		Type: ty,
	}
}

func (p *parser) parseFunctionDeclaration() ast.Statement {
	keyword := p.consume()
	var methodOf *ast.MethodOf
	var memberOf *ast.MemberOf

	if p.next().Kind == token.LEFT_PAREN {
		leftParen := p.consume()
		ty := p.parseType()
		rightParen := p.expect(token.RIGHT_PAREN)

		methodOf = &ast.MethodOf{
			LeftParen:  leftParen,
			Type:       ty,
			RightParen: rightParen,
		}
	}

	name := p.expect(token.IDENTIFIER)

	if p.next().Kind == token.DOT {
		if methodOf != nil {
			p.Diagnostics.ReportMemberAndMethodNotAllowed(name.Span)
		}

		dot := p.consume()
		memberOf = &ast.MemberOf{
			Name: name,
			Dot:  dot,
		}
		name = p.expect(token.IDENTIFIER)
	} else if methodOf == nil {
		p.identifiers[name.Value] = name.Span
	}

	leftParen := p.expect(token.LEFT_PAREN)
	defer p.exitScope(p.enterScope())
	params, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseParameter)

	if len(params) > 0 {
		lastParam := params[len(params)-1]
		if lastParam.Type == nil {
			p.Diagnostics.ReportLastParameterMustHaveType(lastParam.Name.Span)
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

func (p *parser) parseReturnStatement() ast.Statement	{
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
