package parser

import "github.com/gearsdatapacks/libra/parser/ast"

func (p *parser) parseTopLevelStatement() ast.Statement {
	switch p.next().Kind {
	default:
		return p.parseStatement()
	}
}

func (p *parser) parseStatement() ast.Statement {
	switch p.next().Kind {
	default:
		return &ast.ExpressionStatement{
			Expression: p.parseExpression(),
		}
	}
}
