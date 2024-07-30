package parser

import (
	"strconv"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
)

// Precedence
const (
	Lowest = iota
	Assignment
	Logical
	Comparison
	Range
	Bitwise
	Additive
	Multiplicative
	Exponential
	Typecheck
	Prefix
	Postfix
	Cast
)

func (p *parser) parseExpression() (ast.Expression, *diagnostics.Diagnostic) {
	old := p.noBraces
	p.noBraces = false
	expr, diag := p.parseSubExpression(Lowest)
	p.noBraces = old
	return expr, diag
}

func (p *parser) parseSubExpression(precedence int) (ast.Expression, *diagnostics.Diagnostic) {
	nudFn := p.lookupNudFn()

	if nudFn == nil {
		if p.typeExpr {
			return nil, diagnostics.ExpectedType(p.next().Location, p.next().Kind)
		} else {
			return nil, diagnostics.ExpectedExpression(p.next().Location, p.next().Kind)
		}
	}

	left, err := nudFn()
	if err != nil {
		return nil, err
	}

	ledFn := p.lookupLedFn(left)
	for p.canContinue() &&
		ledFn != nil &&
		precedence < p.leftPrecedence(left) {

		left, err = ledFn(left)
		if err != nil {
			return nil, err
		}

		ledFn = p.lookupLedFn(left)
	}

	return left, nil
}

func (p *parser) parseBinaryExpression(left ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	precedence := p.rightPrecedence(left)
	operator := p.consume()
	right, err := p.parseSubExpression(precedence)
	if err != nil {
		return nil, err
	}

	return &ast.BinaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}, nil
}

func (p *parser) parseAssignmentExpression(assignee ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	precedence := p.rightPrecedence(assignee)
	operator := p.consume()
	value, err := p.parseSubExpression(precedence)
	if err != nil {
		return nil, err
	}

	return &ast.AssignmentExpression{
		Assignee: assignee,
		Operator: operator,
		Value:    value,
	}, nil
}

func (p *parser) parsePrefixExpression() (ast.Expression, *diagnostics.Diagnostic) {
	operator := p.consume()
	operand, err := p.parseSubExpression(Prefix)
	if err != nil {
		return nil, err
	}

	return &ast.PrefixExpression{
		Location: operator.Location,
		Operator: operator.Kind,
		Operand:  operand,
	}, nil
}

func (p *parser) parsePostfixExpression(operand ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	operator := p.consume()

	return &ast.PostfixExpression{
		OperatorLocation: operator.Location,
		Operand:          operand,
		Operator:         operator.Kind,
	}, nil
}

func (p *parser) parseDerefExpression(operand ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	p.consume()

	return &ast.DerefExpression{
		Operand: operand,
	}, nil
}

func (p *parser) parsePtrOrRef() (ast.Expression, *diagnostics.Diagnostic) {
	operator := p.consume()
	mut := p.isKeyword("mut")
	if mut {
		p.consume()
	}
	operand, err := p.parseSubExpression(Prefix)
	if err != nil {
		return nil, err
	}

	if operator.Kind == token.STAR {
		return &ast.PointerType{
			Location: operator.Location,
			Mutable:  mut,
			Operand:  operand,
		}, nil
	}

	return &ast.RefExpression{
		Location: operator.Location,
		Mutable:  mut,
		Operand:  operand,
	}, nil
}

func (p *parser) parseOptionType() (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	operand, nil := p.parseSubExpression(Prefix)

	return &ast.OptionType{
		Location: location,
		Operand:  operand,
	}, nil
}

func (p *parser) parseFunctionCall(callee ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	p.consume()
	arguments := parseDelimExprList(p, token.RIGHT_PAREN, p.parseExpression)

	return &ast.FunctionCall{
		Callee:    callee,
		Arguments: arguments,
	}, nil
}

func (p *parser) parseIndexExpression(left ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	var index ast.Expression

	if p.next().Kind != token.RIGHT_SQUARE {
		var err *diagnostics.Diagnostic
		index, err = p.parseExpression()
		if err != nil {
			p.Diagnostics.Report(err)
			p.consumeUntil(token.RIGHT_SQUARE)
		}
	}
	p.expect(token.RIGHT_SQUARE)

	return &ast.IndexExpression{
		Left:     left,
		Location: location,
		Index:    index,
	}, nil
}

func (p *parser) parseMember(left ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	member := p.expect(token.IDENTIFIER)

	return &ast.MemberExpression{
		Location:       location,
		MemberLocation: member.Location,
		Left:           left,
		Member:         member.Value,
	}, nil
}

func (p *parser) parseInferredTypeExpression() (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	left := &ast.InferredExpression{Location: location}

	if p.next().Kind == token.IDENTIFIER {
		member := p.consume()
		return &ast.MemberExpression{
			Location:       location,
			MemberLocation: member.Location,
			Left:           left,
			Member:         member.Value,
		}, nil
	}

	if p.next().Kind == token.LEFT_BRACE {
		return p.parseStructExpression(left)
	}

	return nil, diagnostics.ExpectedMemberOrStructBody(p.next().Location, p.next())
}

func (p *parser) parseStructMember() (*ast.StructMember, *diagnostics.Diagnostic) {
	var name *string
	var value ast.Expression
	var location text.Location

	initial, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if ident, ok := initial.(*ast.Identifier); ok {
		name = &ident.Name
		location = ident.Location

		if p.next().Kind == token.COLON {
			p.consume()
			value, err = p.parseExpression()
			if err != nil {
				return nil, err
			}
		}
	} else {
		value = initial
		location = initial.GetLocation()
	}

	return &ast.StructMember{
		Location: location,
		Name:     name,
		Value:    value,
	}, nil
}

func (p *parser) parseStructExpression(instanceOf ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	p.consume()

	members := parseDerefExprList(p, token.RIGHT_BRACE, p.parseStructMember)

	return &ast.StructExpression{
		Struct:  instanceOf,
		Members: members,
	}, nil
}

func (p *parser) parseCastExpression(left ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	toType, err := p.parseTypeExpression()
	if err != nil {
		return nil, err
	}

	return &ast.CastExpression{
		Location: location,
		Left:     left,
		Type:     toType,
	}, nil
}

func (p *parser) parseTypeCheckExpression(left ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	ty, err := p.parseTypeExpression()
	if err != nil {
		return nil, err
	}

	return &ast.TypeCheckExpression{
		Location: location,
		Left:     left,
		Type:     ty,
	}, nil
}

func (p *parser) parseRangeExpression(start ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	end, err := p.parseSubExpression(Range)
	if err != nil {
		return nil, err
	}

	return &ast.RangeExpression{
		Location: location,
		Start:    start,
		End:      end,
	}, nil
}

func (p *parser) parseTuple() (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location

	values := parseDelimExprList(p, token.RIGHT_PAREN, p.parseExpression)

	if len(values) == 1 {
		return &ast.ParenthesisedExpression{
			Location:   location,
			Expression: values[0],
		}, nil
	}

	return &ast.TupleExpression{
		Location: location,
		Values:   values,
	}, nil
}

func (p *parser) parseIdentifier() (ast.Expression, *diagnostics.Diagnostic) {
	tok := p.consume()

	switch tok.Value {
	case "true":
		return &ast.BooleanLiteral{
			Location: tok.Location,
			Value:    true,
		}, nil

	case "false":
		return &ast.BooleanLiteral{
			Location: tok.Location,
			Value:    false,
		}, nil

	default:
		return &ast.Identifier{
			Location: tok.Location,
			Name:     tok.Value,
		}, nil
	}
}

func (p *parser) parseInteger() (ast.Expression, *diagnostics.Diagnostic) {
	tok := p.consume()
	radix := 10

	if len(tok.Value) >= 2 {
		switch tok.Value[:2] {
		case "0b": radix = 2
		case "0o": radix = 8
		case "0x": radix = 16
		}
	}

	value, _ := strconv.ParseInt(tok.ExtraValue, radix, 64)
	return &ast.IntegerLiteral{
		Token: tok,
		Value: value,
	}, nil
}

func (p *parser) parseFloat() (ast.Expression, *diagnostics.Diagnostic) {
	tok := p.consume()
	value, _ := strconv.ParseFloat(tok.ExtraValue, 64)
	return &ast.FloatLiteral{
		Token: tok,
		Value: value,
	}, nil
}

func (p *parser) parseString() (ast.Expression, *diagnostics.Diagnostic) {
	tok := p.consume()
	return &ast.StringLiteral{
		Token: tok,
		Value: tok.ExtraValue,
	}, nil
}

func (p *parser) parseList() (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	values := parseDelimExprList(p, token.RIGHT_SQUARE, p.parseExpression)

	return &ast.ListLiteral{
		Location: location,
		Values:   values,
	}, nil
}

func (p *parser) parseKeyValue() (*ast.KeyValue, *diagnostics.Diagnostic) {
	key, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	p.expect(token.COLON)
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.KeyValue{
		Key:   key,
		Value: value,
	}, nil
}

func (p *parser) parseMapOrBlock() (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	if p.next().Kind == token.RIGHT_BRACE {
		p.consume()
		return &ast.MapLiteral{
			Location:  location,
			KeyValues: []ast.KeyValue{},
		}, nil
	}

	first, err := p.parseStatement()
	if err != nil {
		p.Diagnostics.Report(err)
		p.consumeUntil(token.RIGHT_BRACE)
		p.expect(token.RIGHT_BRACE)

		return &ast.MapLiteral{
			Location:  location,
			KeyValues: []ast.KeyValue{},
		}, nil
	}

	if key, ok := first.(ast.Expression); ok && p.next().Kind == token.COLON {
		p.consume()
		value, err := p.parseExpression()
		if err != nil {
			p.consumeUntil(token.RIGHT_BRACE)
			p.expect(token.RIGHT_BRACE)

			return &ast.MapLiteral{
				Location:  location,
				KeyValues: []ast.KeyValue{},
			}, nil
		}

		keyValues := extendDelimExprList(
			p,
			ast.KeyValue{Key: key, Value: value},
			token.RIGHT_BRACE, p.parseKeyValue,
		)

		return &ast.MapLiteral{
			Location:  location,
			KeyValues: keyValues,
		}, nil
	}

	stmts := extendDelimStmtList(p, first, token.RIGHT_BRACE, p.parseStatement)
	return &ast.Block{
		Location:   location,
		Statements: stmts,
	}, nil
}

func (p *parser) parseFunctionExpression() (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location

	p.expect(token.LEFT_PAREN)
	defer p.exitScope(p.enterScope())
	params := parseDerefExprList(p, token.RIGHT_PAREN, p.parseParameter)

	returnType, err := p.parseOptionalTypeAnnotation()
	if err != nil {
		return nil, err
	}

	var body *ast.Block
	if p.canContinue() && p.next().Kind == token.LEFT_BRACE {
		var err *diagnostics.Diagnostic
		body, err = p.parseBlock(true)
		if err != nil {
			return nil, err
		}
	}

	return &ast.FunctionExpression{
		Location:   location,
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
	}, nil
}

func (p *parser) parseBlock(noScope ...bool) (*ast.Block, *diagnostics.Diagnostic) {
	location := p.expect(token.LEFT_BRACE).Location
	if len(noScope) == 0 || !noScope[0] {
		defer p.exitScope(p.enterScope())
	}
	statements := parseDelimStmtList(p, token.RIGHT_BRACE, p.parseStatement)

	return &ast.Block{
		Location:   location,
		Statements: statements,
	}, nil
}

func (p *parser) parseIfExpression() (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location

	p.noBraces = true
	p.bracketLevel++
	condition, err := p.parseSubExpression(Lowest)
	if err != nil {
		return nil, err
	}

	p.noBraces = false
	p.bracketLevel--

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Statement

	if p.isKeyword("else") {
		p.consume()
		if p.isKeyword("if") {
			elseBranch, err = p.parseIfExpression()
			if err != nil {
				return nil, err
			}
		} else {
			elseBranch, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
		}
	}

	return &ast.IfExpression{
		Location:   location,
		Condition:  condition,
		Body:       body,
		ElseBranch: elseBranch,
	}, nil
}

func (p *parser) parseWhileLoop() (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location

	p.noBraces = true
	p.bracketLevel++
	condition, err := p.parseSubExpression(Lowest)
	if err != nil {
		return nil, err
	}

	p.bracketLevel--
	p.noBraces = false

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.WhileLoop{
		Location:  location,
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *parser) parseForLoop() (ast.Expression, *diagnostics.Diagnostic) {
	location := p.consume().Location
	defer p.exitScope(p.enterScope())

	variable := p.delcareIdentifier().Value
	p.expectKeyword("in")

	p.noBraces = true
	p.bracketLevel++
	iterator, err := p.parseSubExpression(Lowest)
	if err != nil {
		return nil, err
	}

	p.bracketLevel--
	p.noBraces = false

	body, err := p.parseBlock(true)
	if err != nil {
		return nil, err
	}

	return &ast.ForLoop{
		LLocation: location,
		Variable:  variable,
		Iterator:  iterator,
		Body:      body,
	}, nil
}
