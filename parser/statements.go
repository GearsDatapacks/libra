package parser

import (
	"github.com/gearsdatapacks/libra/ast"
)

func (p *parser) parseStatement() ast.Statement {
	return p.parseExpressionStatement()
}

func (p *parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression()
	return &ast.ExpressionStatement{
		Expression: expr,
		BaseNode: &ast.BaseNode{Token: expr.GetToken()},
	}
}
