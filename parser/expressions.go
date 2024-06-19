package parser

import (
	"strconv"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
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
	p.noBraces = false
	return p.parseSubExpression(Lowest)
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
		Operator: operator,
		Operand:  operand,
	}, nil
}

func (p *parser) parsePostfixExpression(operand ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	operator := p.consume()

	return &ast.PostfixExpression{
		Operand:  operand,
		Operator: operator,
	}, nil
}

func (p *parser) parseDerefExpression(operand ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	operator := p.consume()

	return &ast.DerefExpression{
		Operand:  operand,
		Operator: operator,
	}, nil
}

func (p *parser) parsePtrOrRef() (ast.Expression, *diagnostics.Diagnostic) {
	operator := p.consume()
	var mut *token.Token
	if p.isKeyword("mut") {
		tok := p.consume()
		mut = &tok
	}
	operand, err := p.parseSubExpression(Prefix)
	if err != nil {
		return nil, err
	}

	if operator.Kind == token.STAR {
		return &ast.PointerType{
			Operator: operator,
			Mutable:  mut,
			Operand:  operand,
		}, nil
	}

	return &ast.RefExpression{
		Operator: operator,
		Mutable:  mut,
		Operand:  operand,
	}, nil
}

func (p *parser) parseOptionType() (ast.Expression, *diagnostics.Diagnostic) {
	operator := p.consume()
	operand, nil := p.parseSubExpression(Prefix)

	return &ast.OptionType{
		Operator: operator,
		Operand:  operand,
	}, nil
}

func (p *parser) parseFunctionCall(callee ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	leftParen := p.consume()
	arguments, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseExpression)

	return &ast.FunctionCall{
		Callee:     callee,
		LeftParen:  leftParen,
		Arguments:  arguments,
		RightParen: rightParen,
	}, nil
}

func (p *parser) parseIndexExpression(left ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	leftSquare := p.consume()
	var index ast.Expression

	if p.next().Kind != token.RIGHT_SQUARE {
		var err *diagnostics.Diagnostic
		index, err = p.parseExpression()
		if err != nil {
			p.Diagnostics.Report(err)
			p.consumeUntil(token.RIGHT_SQUARE)
		}
	}
	rightSquare := p.expect(token.RIGHT_SQUARE)

	return &ast.IndexExpression{
		Left:        left,
		LeftSquare:  leftSquare,
		Index:       index,
		RightSquare: rightSquare,
	}, nil
}

func (p *parser) parseMember(left ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	dot := p.consume()
	member := p.expect(token.IDENTIFIER)

	return &ast.MemberExpression{
		Left:   left,
		Dot:    dot,
		Member: member,
	}, nil
}

func (p *parser) parseInferredTypeExpression() (ast.Expression, *diagnostics.Diagnostic) {
	dot := p.consume()
	if p.next().Kind == token.IDENTIFIER {
		member := p.consume()
		return &ast.MemberExpression{
			Left:   nil,
			Dot:    dot,
			Member: member,
		}, nil
	}

	if p.next().Kind == token.LEFT_BRACE {
		return p.parseStructExpression(&ast.InferredExpression{Token: dot})
	}

	return nil, diagnostics.ExpectedMemberOrStructBody(p.next().Location, p.next())
}

func (p *parser) parseStructMember() (*ast.StructMember, *diagnostics.Diagnostic) {
	var name, colon *token.Token
	var value ast.Expression

	initial, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if ident, ok := initial.(*ast.Identifier); ok {
		name = &ident.Token

		if p.next().Kind == token.COLON {
			tok := p.consume()
			colon = &tok
			value, err = p.parseExpression()
			if err != nil {
				return nil, err
			}
		}
	} else {
		value = initial
	}

	return &ast.StructMember{
		Name:  name,
		Colon: colon,
		Value: value,
	}, nil
}

func (p *parser) parseStructExpression(instanceOf ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	leftBrace := p.consume()

	members, rightBrace := parseDerefExprList(p, token.RIGHT_BRACE, p.parseStructMember)

	return &ast.StructExpression{
		Struct:     instanceOf,
		LeftBrace:  leftBrace,
		Members:    members,
		RightBrace: rightBrace,
	}, nil
}

func (p *parser) parseCastExpression(left ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	arrow := p.consume()
	toType, err := p.parseTypeExpression()
	if err != nil {
		return nil, err
	}

	return &ast.CastExpression{
		Left:  left,
		Arrow: arrow,
		Type:  toType,
	}, nil
}

func (p *parser) parseTypeCheckExpression(left ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	operator := p.consume()
	ty, err := p.parseTypeExpression()
	if err != nil {
		return nil, err
	}

	return &ast.TypeCheckExpression{
		Left:     left,
		Operator: operator,
		Type:     ty,
	}, nil
}

func (p *parser) parseRangeExpression(start ast.Expression) (ast.Expression, *diagnostics.Diagnostic) {
	operator := p.consume()
	end, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.RangeExpression{
		Start:    start,
		Operator: operator,
		End:      end,
	}, nil
}

func (p *parser) parseTuple() (ast.Expression, *diagnostics.Diagnostic) {
	leftParen := p.consume()

	values, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseExpression)

	if len(values) == 1 {
		return &ast.ParenthesisedExpression{
			LeftParen:  leftParen,
			Expression: values[0],
			RightParen: rightParen,
		}, nil
	}

	return &ast.TupleExpression{
		LeftParen:  leftParen,
		Values:     values,
		RightParen: rightParen,
	}, nil
}

func (p *parser) parseIdentifier() (ast.Expression, *diagnostics.Diagnostic) {
	tok := p.consume()

	switch tok.Value {
	case "true":
		return &ast.BooleanLiteral{
			Token: tok,
			Value: true,
		}, nil

	case "false":
		return &ast.BooleanLiteral{
			Token: tok,
			Value: false,
		}, nil

	default:
		return &ast.Identifier{
			Token: tok,
			Name:  tok.Value,
		}, nil
	}
}

func (p *parser) parseInteger() (ast.Expression, *diagnostics.Diagnostic) {
	tok := p.consume()
	value, _ := strconv.ParseInt(tok.Value, 10, 64)
	return &ast.IntegerLiteral{
		Token: tok,
		Value: value,
	}, nil
}

func (p *parser) parseFloat() (ast.Expression, *diagnostics.Diagnostic) {
	tok := p.consume()
	value, _ := strconv.ParseFloat(tok.Value, 64)
	return &ast.FloatLiteral{
		Token: tok,
		Value: value,
	}, nil
}

func (p *parser) parseString() (ast.Expression, *diagnostics.Diagnostic) {
	tok := p.consume()
	return &ast.StringLiteral{
		Token: tok,
		Value: tok.Value,
	}, nil
}

func (p *parser) parseList() (ast.Expression, *diagnostics.Diagnostic) {
	leftSquare := p.consume()
	values, rightSquare := parseDelimExprList(p, token.RIGHT_SQUARE, p.parseExpression)

	return &ast.ListLiteral{
		LeftSquare:  leftSquare,
		Values:      values,
		RightSquare: rightSquare,
	}, nil
}

func (p *parser) parseKeyValue() (*ast.KeyValue, *diagnostics.Diagnostic) {
	key, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	colon := p.expect(token.COLON)
	value, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.KeyValue{
		Key:   key,
		Colon: colon,
		Value: value,
	}, nil
}

func (p *parser) parseMapOrBlock() (ast.Expression, *diagnostics.Diagnostic) {
	leftBrace := p.consume()
	if p.next().Kind == token.RIGHT_BRACE {
		rightBrace := p.consume()
		return &ast.MapLiteral{
			LeftBrace:  leftBrace,
			KeyValues:  []ast.KeyValue{},
			RightBrace: rightBrace,
		}, nil
	}

	first, err := p.parseStatement()
	if err != nil {
		p.Diagnostics.Report(err)
		p.consumeUntil(token.LEFT_BRACE)
		rightBrace := p.expect(token.LEFT_BRACE)

		return &ast.MapLiteral{
			LeftBrace:  leftBrace,
			KeyValues:  []ast.KeyValue{},
			RightBrace: rightBrace,
		}, nil
	}

	if key, ok := first.(ast.Expression); ok && p.next().Kind == token.COLON {
		colon := p.consume()
		value, err := p.parseExpression()
		if err != nil {
			p.consumeUntil(token.LEFT_BRACE)
			rightBrace := p.expect(token.LEFT_BRACE)

			return &ast.MapLiteral{
				LeftBrace:  leftBrace,
				KeyValues:  []ast.KeyValue{},
				RightBrace: rightBrace,
			}, nil
		}

		keyValues, rightBrace := extendDelimExprList(
			p,
			ast.KeyValue{Key: key, Colon: colon, Value: value},
			token.RIGHT_BRACE, p.parseKeyValue,
		)

		return &ast.MapLiteral{
			LeftBrace:  leftBrace,
			KeyValues:  keyValues,
			RightBrace: rightBrace,
		}, nil
	}

	stmts, rightBrace := extendDelimStmtList(p, first, token.RIGHT_BRACE, p.parseStatement)
	return &ast.Block{
		LeftBrace:  leftBrace,
		Statements: stmts,
		RightBrace: rightBrace,
	}, nil
}

func (p *parser) parseFunctionExpression() (ast.Expression, *diagnostics.Diagnostic) {
	keyword := p.consume()

	leftParen := p.expect(token.LEFT_PAREN)
	defer p.exitScope(p.enterScope())
	params, rightParen := parseDerefExprList(p, token.RIGHT_PAREN, p.parseParameter)

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
		Keyword:    keyword,
		LeftParen:  leftParen,
		Parameters: params,
		RightParen: rightParen,
		ReturnType: returnType,
		Body:       body,
	}, nil
}

func (p *parser) parseBlock(noScope ...bool) (*ast.Block, *diagnostics.Diagnostic) {
	leftBrace := p.expect(token.LEFT_BRACE)
	if len(noScope) == 0 || !noScope[0] {
		defer p.exitScope(p.enterScope())
	}
	statements, rightBrace := parseDelimStmtList(p, token.RIGHT_BRACE, p.parseStatement)

	return &ast.Block{
		LeftBrace:  leftBrace,
		Statements: statements,
		RightBrace: rightBrace,
	}, nil
}

func (p *parser) parseIfExpression() (ast.Expression, *diagnostics.Diagnostic) {
	keyword := p.consume()

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

	var elseBranch *ast.ElseBranch

	if p.isKeyword("else") {
		elseBranch = &ast.ElseBranch{}
		elseBranch.ElseKeyword = p.consume()
		if p.isKeyword("if") {
			elseBranch.Statement, err = p.parseIfExpression()
			if err != nil {
				return nil, err
			}
		} else {
			elseBranch.Statement, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
		}
	}

	return &ast.IfExpression{
		Keyword:    keyword,
		Condition:  condition,
		Body:       body,
		ElseBranch: elseBranch,
	}, nil
}

func (p *parser) parseWhileLoop() (ast.Expression, *diagnostics.Diagnostic) {
	keyword := p.consume()

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
		Keyword:   keyword,
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *parser) parseForLoop() (ast.Expression, *diagnostics.Diagnostic) {
	forKeyword := p.consume()
	defer p.exitScope(p.enterScope())

	variable := p.delcareIdentifier()
	inKeyword := p.expectKeyword("in")

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
		ForKeyword: forKeyword,
		Variable:   variable,
		InKeyword:  inKeyword,
		Iterator:   iterator,
		Body:       body,
	}, nil
}
