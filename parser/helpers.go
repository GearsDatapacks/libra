package parser

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func (p *parser) parseArgumentList() ([]ast.Expression, error) {
	_, err := p.expect(token.LEFT_PAREN, "Expected '(' to open argument list, got %q")
	if err != nil {
		return nil, err
	}

	args := []ast.Expression{}

	if p.next().Type != token.RIGHT_PAREN {
		args, err = p.parseArgs()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.expect(token.RIGHT_PAREN, "Expected comma or end of argument list")
	if err != nil {
		return nil, err
	}

	return args, nil
}

func (p *parser) parseArgs() ([]ast.Expression, error) {
	initialExpr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	args := []ast.Expression{ initialExpr }

	for p.next().Type == token.COMMA {
		p.consume()

		nextExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, nextExpr)
	}

	return args, nil
}

func (p *parser) parseParameterList() ([]ast.Parameter, error) {
	_, err := p.expect(token.LEFT_PAREN, "Expected '(' to open parameter list, got %q")
	if err != nil {
		return nil, err
	}

	params := []ast.Parameter{}

	if p.next().Type != token.RIGHT_PAREN {
		params, err = p.parseParameters()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.expect(token.RIGHT_PAREN, "Expected comma or end of parameter list")
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (p *parser) parseParameters() ([]ast.Parameter, error) {
	initialParam, err := p.parseParameter()
	if err != nil {
		return nil, err
	}
	params := []ast.Parameter{ initialParam }

	for p.next().Type == token.COMMA {
		p.consume()

		nextParam, err := p.parseParameter()
		if err != nil {
			return nil, err
		}
		params = append(params, nextParam)
	}

	return params, nil
}

func (p *parser) parseParameter() (ast.Parameter, error) {
	name, err := p.expect(token.IDENTIFIER, "Invalid parameter name %q")
	if err != nil {
		return ast.Parameter{}, err
	}
	dataType, err := p.parseType()
	if err != nil {
		return ast.Parameter{}, err
	}

	return ast.Parameter{Name: name.Value, Type: dataType}, nil
}

func (p *parser) parseCodeBlock() ([]ast.Statement, error) {
	outerSymbols := make([]string, len(p.usedSymbols))
	copy(outerSymbols, p.usedSymbols)

	_, err := p.expect(token.LEFT_BRACE, "Unexpected %q, expected code block")
	if err != nil {
		return nil, err
	}

	code := []ast.Statement{}

	for p.next().Type != token.RIGHT_BRACE {
		nextStmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		code = append(code, nextStmt)
	}

	p.usedSymbols = outerSymbols

	_, err = p.expect(token.RIGHT_BRACE, "Unexpected %q, expected '}")
	if err != nil {
		return nil, err
	}

	return code, nil
}
