package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
)

type parser struct {
	tokens      []token.Token
	pos         int
	nudFns      map[token.Kind]nudFn
	ledOps      map[token.Kind]opInfo
	Diagnostics diagnostics.Manager
}

func New(tokens []token.Token, diagnostics diagnostics.Manager) *parser {
	p := &parser{
		tokens:      tokens,
		pos:         0,
		nudFns:      map[token.Kind]nudFn{},
		ledOps:      map[token.Kind]opInfo{},
		Diagnostics: diagnostics,
	}

	p.register()

	return p
}

func (p *parser) Parse() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for !p.eof() {
		pos := p.pos

		program.Statements = append(program.Statements, p.parseTopLevelStatement())

		if p.pos == pos {
			p.consume()
		}

		if !p.eof() && p.nextWithNewlines().Kind != token.NEWLINE {
			p.Diagnostics.ReportExpectedNewline(p.next().Span, p.next().Kind)
		}

		p.consumeNewlines()
	}

	return program
}

type nudFn func() ast.Expression
type ledFn func(ast.Expression) ast.Expression

type opInfo struct {
	leftPrecedence  int
	rightPrecedence int
	parseFn         ledFn
}

func (p *parser) registerNudFn(kind token.Kind, fn nudFn) {
	p.nudFns[kind] = fn
}

func (p *parser) registerLedOp(kind token.Kind, precedence int, fn ledFn, rightAssociative ...bool) {
	isRightassociative := false
	if len(rightAssociative) != 0 {
		isRightassociative = rightAssociative[0]
	}

	leftPrecedence := precedence
	rightPrecedence := precedence
	if isRightassociative {
		rightPrecedence -= 1
	}

	p.ledOps[kind] = opInfo{
		leftPrecedence:  leftPrecedence,
		rightPrecedence: rightPrecedence,
		parseFn:         fn,
	}
}

func (p *parser) lookupNudFn(kind token.Kind) nudFn {
	fn, ok := p.nudFns[kind]
	if !ok {
		return nil
	}
	return fn
}

func (p *parser) lookupLedFn(kind token.Kind) ledFn {
	info, ok := p.ledOps[kind]
	if !ok {
		return nil
	}
	return info.parseFn
}

func (p *parser) leftPrecedence(kind token.Kind) int {
	info, ok := p.ledOps[kind]
	if !ok {
		return Lowest
	}
	return info.leftPrecedence
}

func (p *parser) rightPrecedence(kind token.Kind) int {
	info, ok := p.ledOps[kind]
	if !ok {
		return Lowest
	}
	return info.rightPrecedence
}

func (p *parser) register() {
	// Literals
	p.registerNudFn(token.INTEGER, p.parseInteger)
	p.registerNudFn(token.FLOAT, p.parseFloat)
	p.registerNudFn(token.STRING, p.parseString)
	p.registerNudFn(token.IDENTIFIER, p.parseIdentifier)
	p.registerNudFn(token.LEFT_SQUARE, p.parseList)
	p.registerNudFn(token.LEFT_BRACE, p.parseMap)

	p.registerNudFn(token.LEFT_PAREN, p.parseParenthesisedExpression)

	// Postfix expressions
	p.registerLedOp(token.LEFT_PAREN, Postfix, p.parseFunctionCall)
	p.registerLedOp(token.LEFT_SQUARE, Postfix, p.parseIndexExpression)

	// Postfix operators
	p.registerLedOp(token.DOUBLE_PLUS, Postfix, p.parsePostfixExpression)
	p.registerLedOp(token.DOUBLE_MINUS, Postfix, p.parsePostfixExpression)
	p.registerLedOp(token.QUESTION, Postfix, p.parsePostfixExpression)
	p.registerLedOp(token.BANG, Postfix, p.parsePostfixExpression)

	// Prefix operators
	p.registerNudFn(token.MINUS, p.parsePrefixExpression)
	p.registerNudFn(token.PLUS, p.parsePrefixExpression)
	p.registerNudFn(token.BANG, p.parsePrefixExpression)
	p.registerNudFn(token.STAR, p.parsePrefixExpression)
	p.registerNudFn(token.AMPERSAND, p.parsePrefixExpression)
	p.registerNudFn(token.TILDE, p.parsePrefixExpression)

	// Binary operators
	p.registerLedOp(token.DOUBLE_AMPERSAND, Logical, p.parseBinaryExpression)
	p.registerLedOp(token.DOUBLE_PIPE, Logical, p.parseBinaryExpression)

	p.registerLedOp(token.LEFT_ANGLE, Comparison, p.parseBinaryExpression)
	p.registerLedOp(token.LEFT_ANGLE_EQUALS, Comparison, p.parseBinaryExpression)
	p.registerLedOp(token.RIGHT_ANGLE, Comparison, p.parseBinaryExpression)
	p.registerLedOp(token.RIGHT_ANGLE_EQUALS, Comparison, p.parseBinaryExpression)
	p.registerLedOp(token.DOUBLE_EQUALS, Comparison, p.parseBinaryExpression)
	p.registerLedOp(token.BANG_EQUALS, Comparison, p.parseBinaryExpression)

	p.registerLedOp(token.DOUBLE_LEFT_ANGLE, Bitwise, p.parseBinaryExpression)
	p.registerLedOp(token.DOUBLE_RIGHT_ANGLE, Bitwise, p.parseBinaryExpression)
	p.registerLedOp(token.PIPE, Bitwise, p.parseBinaryExpression)
	p.registerLedOp(token.AMPERSAND, Bitwise, p.parseBinaryExpression)

	p.registerLedOp(token.PLUS, Additive, p.parseBinaryExpression)
	p.registerLedOp(token.MINUS, Additive, p.parseBinaryExpression)

	p.registerLedOp(token.STAR, Multiplicative, p.parseBinaryExpression)
	p.registerLedOp(token.SLASH, Multiplicative, p.parseBinaryExpression)
	p.registerLedOp(token.PERCENT, Multiplicative, p.parseBinaryExpression)

	p.registerLedOp(token.DOUBLE_STAR, Exponential, p.parseBinaryExpression, true)
}

func (p *parser) next() token.Token {
	return p.peek(0)
}

func (p *parser) peek(offset int) token.Token {
	if p.pos+offset >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}

	for p.tokens[p.pos+offset].Kind == token.NEWLINE {
		offset++
	}

	return p.tokens[p.pos+offset]
}

func (p *parser) nextWithNewlines() token.Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}

	return p.tokens[p.pos]
}

func (p *parser) consume() token.Token {
	next := p.next()
	for p.tokens[p.pos].Kind == token.NEWLINE {
		p.pos++
	}

	p.pos++
	return next
}

func (p *parser) consumeNewlines() {
	for p.tokens[p.pos].Kind == token.NEWLINE {
		p.pos++
	}
}

func (p *parser) expect(kind token.Kind) token.Token {
	if p.next().Kind == kind {
		return p.consume()
	}
	p.Diagnostics.ReportExpectedToken(p.next().Span, kind, p.next().Kind)
	tok := token.New(kind, "", token.NewSpan(0, 0, 0))
	return tok
}

func (p *parser) eof() bool {
	return p.next().Kind == token.EOF
}
