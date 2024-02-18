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

func (p *parser) parseBlockStatement() *ast.BlockStatement {
	leftBrace := p.expect(token.LEFT_BRACE)
	defer p.exitScope(p.enterScope())
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
	condition := p.parseSubExpression(Lowest)
	p.noBraces = false

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
	condition := p.parseSubExpression(Lowest)
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
	iterator := p.parseSubExpression(Lowest)
	p.noBraces = false

	body := p.parseBlockStatement()

	return &ast.ForLoop{
		ForKeyword: forKeyword,
		Variable:   variable,
		InKeyword:  inKeyword,
		Iterator:   iterator,
		Body:       body,
	}
}
