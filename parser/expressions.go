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

func (p *parser) parseExpression() ast.Expression {
	p.noBraces = false
	return p.parseSubExpression(Lowest)
}

func (p *parser) parseSubExpression(precedence int) ast.Expression {
	nudFn := p.lookupNudFn()

	if nudFn == nil {
		if p.typeExpr {
			p.Diagnostics.Report(diagnostics.ExpectedType(p.next().Location, p.next().Kind))
		} else {
			p.Diagnostics.Report(diagnostics.ExpectedExpression(p.next().Location, p.next().Kind))
		}
		return &ast.ErrorNode{}
	}

	left := nudFn()

	ledFn := p.lookupLedFn(left)
	for p.canContinue() &&
		ledFn != nil &&
		precedence < p.leftPrecedence(left) {

		left = ledFn(left)

		ledFn = p.lookupLedFn(left)
	}

	return left
}

func (p *parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	precedence := p.rightPrecedence(left)
	operator := p.consume()
	right := p.parseSubExpression(precedence)

	return &ast.BinaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (p *parser) parseAssignmentExpression(assignee ast.Expression) ast.Expression {
	precedence := p.rightPrecedence(assignee)
	operator := p.consume()
	value := p.parseSubExpression(precedence)

	return &ast.AssignmentExpression{
		Assignee: assignee,
		Operator: operator,
		Value:    value,
	}
}

func (p *parser) parsePrefixExpression() ast.Expression {
	operator := p.consume()
	operand := p.parseSubExpression(Prefix)

	return &ast.PrefixExpression{
		Operator: operator,
		Operand:  operand,
	}
}

func (p *parser) parsePostfixExpression(operand ast.Expression) ast.Expression {
	operator := p.consume()

	return &ast.PostfixExpression{
		Operand:  operand,
		Operator: operator,
	}
}

func (p *parser) parseDerefExpression(operand ast.Expression) ast.Expression {
	operator := p.consume()

	return &ast.DerefExpression{
		Operand:  operand,
		Operator: operator,
	}
}

func (p *parser) parsePtrOrRef() ast.Expression {
	operator := p.consume()
	var mut *token.Token
	if p.isKeyword("mut") {
		tok := p.consume()
		mut = &tok
	}
	operand := p.parseSubExpression(Prefix)

	if operator.Kind == token.STAR {
		return &ast.PointerType{
			Operator: operator,
			Mutable:  mut,
			Operand:  operand,
		}
	}
	return &ast.RefExpression{
		Operator: operator,
		Mutable:  mut,
		Operand:  operand,
	}
}

func (p *parser) parseOptionType() ast.Expression {
	operator := p.consume()
	operand := p.parseExpression()

	return &ast.OptionType{
		Operator: operator,
		Operand:  operand,
	}
}

func (p *parser) parseFunctionCall(callee ast.Expression) ast.Expression {
	leftParen := p.consume()
	arguments, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseExpression)

	return &ast.FunctionCall{
		Callee:     callee,
		LeftParen:  leftParen,
		Arguments:  arguments,
		RightParen: rightParen,
	}
}

func (p *parser) parseIndexExpression(left ast.Expression) ast.Expression {
	leftSquare := p.consume()
	var index ast.Expression

	if p.next().Kind != token.RIGHT_SQUARE {
		index = p.parseExpression()
	}
	rightSquare := p.expect(token.RIGHT_SQUARE)

	return &ast.IndexExpression{
		Left:        left,
		LeftSquare:  leftSquare,
		Index:       index,
		RightSquare: rightSquare,
	}
}

func (p *parser) parseMember(left ast.Expression) ast.Expression {
	dot := p.consume()
	member := p.expect(token.IDENTIFIER)

	return &ast.MemberExpression{
		Left:   left,
		Dot:    dot,
		Member: member,
	}
}

func (p *parser) parseInferredTypeExpression() ast.Expression {
	dot := p.consume()
	if p.next().Kind == token.IDENTIFIER {
		member := p.consume()
		return &ast.MemberExpression{
			Left:   nil,
			Dot:    dot,
			Member: member,
		}
	}

	if p.next().Kind == token.LEFT_BRACE {
		return p.parseStructExpression(&ast.InferredExpression{Token: dot})
	}

	p.Diagnostics.Report(diagnostics.ExpectedMemberOrStructBody(p.next().Location, p.next()))
	return &ast.InferredExpression{Token: dot}
}

func (p *parser) parseStructMember() ast.StructMember {
	var name, colon *token.Token
	var value ast.Expression

	initial := p.parseExpression()
	if ident, ok := initial.(*ast.Identifier); ok {
		name = &ident.Token

		if p.next().Kind == token.COLON {
			tok := p.consume()
			colon = &tok
			value = p.parseExpression()
		}
	} else {
		value = initial
	}

	return ast.StructMember{
		Name:  name,
		Colon: colon,
		Value: value,
	}
}

func (p *parser) parseStructExpression(instanceOf ast.Expression) ast.Expression {
	leftBrace := p.consume()

	members, rightBrace := parseDelimExprList(p, token.RIGHT_BRACE, p.parseStructMember)

	return &ast.StructExpression{
		Struct:     instanceOf,
		LeftBrace:  leftBrace,
		Members:    members,
		RightBrace: rightBrace,
	}
}

func (p *parser) parseCastExpression(left ast.Expression) ast.Expression {
	arrow := p.consume()
	toType := p.parseTypeExpression()

	return &ast.CastExpression{
		Left:  left,
		Arrow: arrow,
		Type:  toType,
	}
}

func (p *parser) parseTypeCheckExpression(left ast.Expression) ast.Expression {
	operator := p.consume()
	ty := p.parseTypeExpression()

	return &ast.TypeCheckExpression{
		Left:     left,
		Operator: operator,
		Type:     ty,
	}
}

func (p *parser) parseRangeExpression(start ast.Expression) ast.Expression {
	operator := p.consume()
	end := p.parseExpression()

	return &ast.RangeExpression{
		Start:    start,
		Operator: operator,
		End:      end,
	}
}

func (p *parser) parseTuple() ast.Expression {
	leftParen := p.consume()

	values, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseExpression)

	if len(values) == 1 {
		return &ast.ParenthesisedExpression{
			LeftParen:  leftParen,
			Expression: values[0],
			RightParen: rightParen,
		}
	}

	return &ast.TupleExpression{
		LeftParen:  leftParen,
		Values:     values,
		RightParen: rightParen,
	}
}

func (p *parser) parseIdentifier() ast.Expression {
	tok := p.consume()

	switch tok.Value {
	case "true":
		return &ast.BooleanLiteral{
			Token: tok,
			Value: true,
		}

	case "false":
		return &ast.BooleanLiteral{
			Token: tok,
			Value: false,
		}

	default:
		return &ast.Identifier{
			Token: tok,
			Name:  tok.Value,
		}
	}
}

func (p *parser) parseInteger() ast.Expression {
	tok := p.consume()
	value, _ := strconv.ParseInt(tok.Value, 10, 64)
	return &ast.IntegerLiteral{
		Token: tok,
		Value: value,
	}
}

func (p *parser) parseFloat() ast.Expression {
	tok := p.consume()
	value, _ := strconv.ParseFloat(tok.Value, 64)
	return &ast.FloatLiteral{
		Token: tok,
		Value: value,
	}
}

func (p *parser) parseString() ast.Expression {
	tok := p.consume()
	return &ast.StringLiteral{
		Token: tok,
		Value: tok.Value,
	}
}

func (p *parser) parseList() ast.Expression {
	leftSquare := p.consume()
	values, rightSquare := parseDelimExprList(p, token.RIGHT_SQUARE, p.parseExpression)

	return &ast.ListLiteral{
		LeftSquare:  leftSquare,
		Values:      values,
		RightSquare: rightSquare,
	}
}

func (p *parser) parseKeyValue() ast.KeyValue {
	key := p.parseExpression()
	colon := p.expect(token.COLON)
	value := p.parseExpression()

	return ast.KeyValue{
		Key:   key,
		Colon: colon,
		Value: value,
	}
}

func (p *parser) parseMapOrBlock() ast.Expression {
	leftBrace := p.consume()
	if p.next().Kind == token.RIGHT_BRACE {
		rightBrace := p.consume()
		return &ast.MapLiteral{
			LeftBrace:  leftBrace,
			KeyValues:  []ast.KeyValue{},
			RightBrace: rightBrace,
		}
	}

	first := p.parseStatement()
	if key, ok := first.(ast.Expression); ok && p.next().Kind == token.COLON {
		colon := p.consume()
		value := p.parseExpression()
		keyValues, rightBrace := extendDelimExprList(
			p,
			ast.KeyValue{Key: key, Colon: colon, Value: value},
			token.RIGHT_BRACE, p.parseKeyValue,
		)

		return &ast.MapLiteral{
			LeftBrace:  leftBrace,
			KeyValues:  keyValues,
			RightBrace: rightBrace,
		}
	}

	stmts, rightBrace := extendDelimStmtList(p, first, token.RIGHT_BRACE, p.parseStatement)
	return &ast.Block{
		LeftBrace:  leftBrace,
		Statements: stmts,
		RightBrace: rightBrace,
	}
}

func (p *parser) parseFunctionExpression() ast.Expression {
	keyword := p.consume()

	leftParen := p.expect(token.LEFT_PAREN)
	defer p.exitScope(p.enterScope())
	params, rightParen := parseDelimExprList(p, token.RIGHT_PAREN, p.parseParameter)

	returnType := p.parseOptionalTypeAnnotation()
	var body *ast.Block
	if p.canContinue() && p.next().Kind == token.LEFT_BRACE {
		body = p.parseBlock(true)
	}

	return &ast.FunctionExpression{
		Keyword:    keyword,
		LeftParen:  leftParen,
		Parameters: params,
		RightParen: rightParen,
		ReturnType: returnType,
		Body:       body,
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
