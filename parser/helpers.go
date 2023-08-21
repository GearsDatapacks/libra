package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseArgumentList() []ast.Expression {
	p.expect(token.LEFT_PAREN, "Expected '(' to open argument list, got %q")

	args := []ast.Expression{}

	if p.next().Type != token.RIGHT_PAREN {
		args = p.parseArgs()
	}

	p.expect(token.RIGHT_PAREN, "Expected ')' after argument list, got %q")

	return args
}

func (p *parser) parseArgs() []ast.Expression {
	args := []ast.Expression{ p.parseExpression() }

	for p.next().Type == token.COMMA {
		p.consume()

		args = append(args, p.parseExpression())
	}

	return args
}

func (p *parser) parseParameterList() [][2]string {
	p.expect(token.LEFT_PAREN, "Expected '(' to open parameter list, got %q")

	params := [][2]string{}

	if p.next().Type != token.RIGHT_PAREN {
		params = p.parseParameters()
	}

	p.expect(token.RIGHT_PAREN, "Expected ')' after parameter list, got %q")

	return params
}

func (p *parser) parseParameters() [][2]string {
	params := [][2]string{ p.parseParameter() }

	for p.next().Type == token.COMMA {
		p.consume()

		params = append(params, p.parseParameter())
	}

	return params
}

func (p *parser) parseParameter() [2]string {
	name := p.expect(token.IDENTIFIER, "Expected identifier for parameter name")
	dataType := p.expect(token.IDENTIFIER, "Expected identifier for parameter type")

	return [2]string{name.Value, dataType.Value}
}

func (p *parser) parseCodeBlock() []ast.Statement {
	p.expect(token.LEFT_BRACE, "Expected '{' to begin code block")

	code := []ast.Statement{}

	for p.next().Type != token.RIGHT_BRACE {
		code = append(code, p.parseStatement())
	}

	p.expect(token.RIGHT_BRACE, "Expected '}' after code block")

	return code
}
