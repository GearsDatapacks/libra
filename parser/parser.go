package parser

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"

	"github.com/gearsdatapacks/libra/text"
)

type parser struct {
	tokens       []token.Token
	pos          int
	nudFns       map[token.Kind]nudFn
	ledOps       []lookupFn
	keywords     []keyword
	identifiers  map[string]text.Location
	noBraces     bool
	typeExpr     bool
	bracketLevel uint
	Diagnostics  diagnostics.Manager
}

func New(tokens []token.Token, diagnostics diagnostics.Manager) *parser {
	p := &parser{
		tokens:       tokens,
		pos:          0,
		nudFns:       map[token.Kind]nudFn{},
		ledOps:       []lookupFn{},
		keywords:     []keyword{},
		identifiers:  map[string]text.Location{},
		noBraces:     false,
		bracketLevel: 0,
		Diagnostics:  diagnostics,
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

		if !p.eof() {
			p.expectNewline()
		}

	}

	return program
}

type nudFn func() ast.Expression
type ledFn func(ast.Expression) ast.Expression
type lookupFn func(ast.Expression) (opInfo, bool)

type kwdKind int

const (
	expr kwdKind = iota
	stmt
	decl
)

type keyword struct {
	Name     string
	StmtName string
	Fn       func() ast.Statement
	Kind     kwdKind
}

type opInfo struct {
	leftPrecedence  int
	rightPrecedence int
	parseFn         ledFn
	typeSyntax      bool
}

func (p *parser) registerNudFn(kind token.Kind, fn nudFn) {
	p.nudFns[kind] = fn
}

func (p *parser) registerKeyword(kwd string, fn func() ast.Statement, kind kwdKind, name ...string) {
	stmtName := append(name, "")[0]

	p.keywords = append(p.keywords, keyword{
		Name:     kwd,
		StmtName: stmtName,
		Fn:       fn,
		Kind:     kind,
	})
}

func (p *parser) registerLedOp(kind token.Kind, precedence int, fn ledFn, extra ...bool) {
	rightassociative := false
	typeSyntax := false
	if len(extra) > 0 {
		rightassociative = extra[0]
	}
	if len(extra) > 1 {
		typeSyntax = extra[1]
	}

	leftPrecedence := precedence
	rightPrecedence := precedence
	if rightassociative {
		rightPrecedence -= 1
	}

	p.registerLedLookup(func(foo ast.Expression) (opInfo, bool) {
		if p.next().Kind == kind {
			return opInfo{
				leftPrecedence:  leftPrecedence,
				rightPrecedence: rightPrecedence,
				parseFn:         fn,
				typeSyntax:      typeSyntax,
			}, true
		}

		return opInfo{}, false
	})
}

func (p *parser) registerLedLookup(fn lookupFn) {
	p.ledOps = append(p.ledOps, fn)
}

func (p *parser) lookupNudFn() nudFn {
	for _, kwd := range p.keywords {
		if kwd.Kind == expr && p.isKeyword(kwd.Name) {
			return func() ast.Expression { return kwd.Fn().(ast.Expression) }
		}
	}
	fn, ok := p.nudFns[p.next().Kind]
	if ok {
		return fn
	}
	return nil
}

func (p *parser) lookupLedOp(left ast.Expression) (opInfo, bool) {
	for _, lookup := range p.ledOps {
		info, ok := lookup(left)
		if ok {
			return info, true
		}
	}

	return opInfo{}, false
}

func (p *parser) lookupLedFn(left ast.Expression) ledFn {
	info, ok := p.lookupLedOp(left)
	if !ok {
		return nil
	}
	if p.typeExpr && !info.typeSyntax {
		return nil
	}
	return info.parseFn
}

func (p *parser) leftPrecedence(left ast.Expression) int {
	info, ok := p.lookupLedOp(left)
	if !ok {
		return Lowest
	}
	return info.leftPrecedence
}

func (p *parser) rightPrecedence(left ast.Expression) int {
	info, ok := p.lookupLedOp(left)
	if !ok {
		return Lowest
	}
	return info.rightPrecedence
}

func (p *parser) register() {
	p.registerKeyword("fn", p.parseFunctionDeclaration, decl, "Function declaration")
	p.registerKeyword("fn", func() ast.Statement { return p.parseFunctionExpression() }, expr)
	p.registerKeyword("type", p.parseTypeDeclaration, decl, "Type declaration")
	p.registerKeyword("struct", p.parseStructDeclaration, decl, "Struct declaration")
	p.registerKeyword("interface", p.parseInterfaceDeclaration, decl, "Interface declaration")
	p.registerKeyword("import", p.parseImportStatement, decl, "Import")
	p.registerKeyword("enum", p.parseEnumDeclaration, decl, "Enum declaration")
	p.registerKeyword("union", p.parseUnionDeclaration, decl, "Union declaration")
	p.registerKeyword("tag", p.parseTagDeclaration, decl, "Tag declaration")

	p.registerKeyword("const", p.parseVariableDeclaration, stmt)
	p.registerKeyword("let", p.parseVariableDeclaration, stmt)
	p.registerKeyword("mut", p.parseVariableDeclaration, stmt)
	p.registerKeyword("do", func() ast.Statement { return p.parseBlockExpression() }, expr)
	p.registerKeyword("if", func() ast.Statement { return p.parseIfExpression() }, expr)
	p.registerKeyword("else", func() ast.Statement {
		p.Diagnostics.Report(diagnostics.ElseStatementWithoutIf(p.next().Location))
		return &ast.ErrorNode{}
	}, expr)
	p.registerKeyword("while", func() ast.Statement { return p.parseWhileLoop() }, expr)
	p.registerKeyword("for", func() ast.Statement { return p.parseForLoop() }, expr)
	p.registerKeyword("return", p.parseReturnStatement, stmt)
	p.registerKeyword("yield", p.parseYieldStatement, stmt)
	p.registerKeyword("break", p.parseBreakStatement, stmt)
	p.registerKeyword("continue", func() ast.Statement {
		return &ast.ContinueStatement{
			Keyword: p.consume(),
		}
	}, stmt)

	// Literals
	p.registerNudFn(token.INTEGER, p.parseInteger)
	p.registerNudFn(token.FLOAT, p.parseFloat)
	p.registerNudFn(token.STRING, p.parseString)
	p.registerNudFn(token.IDENTIFIER, p.parseIdentifier)
	p.registerNudFn(token.LEFT_SQUARE, p.parseList)
	p.registerNudFn(token.LEFT_BRACE, p.parseMap)

	p.registerNudFn(token.LEFT_PAREN, p.parseTuple)
	p.registerNudFn(token.DOT, p.parseInferredTypeExpression)

	// Postfix expressions
	p.registerLedOp(token.LEFT_PAREN, Postfix, p.parseFunctionCall)
	p.registerLedOp(token.LEFT_SQUARE, Postfix, p.parseIndexExpression, false, true)
	p.registerLedOp(token.DOT, Postfix, p.parseMember, false, true)
	p.registerLedOp(token.ARROW, Postfix, p.parseCastExpression)
	p.registerLedOp(token.DOUBLE_DOT, Range, p.parseRangeExpression)

	p.registerLedLookup(func(left ast.Expression) (opInfo, bool) {
		if p.next().Kind != token.LEFT_BRACE {
			return opInfo{}, false
		}

		if p.noBraces {
			return opInfo{}, false
		}

		_, isIdent := left.(*ast.Identifier)
		_, isMember := left.(*ast.MemberExpression)

		if !isIdent && !isMember {
			return opInfo{}, false
		}

		return opInfo{
			leftPrecedence:  Postfix,
			rightPrecedence: Postfix,
			parseFn:         p.parseStructExpression,
		}, true
	})

	// Postfix operators
	p.registerLedOp(token.DOUBLE_PLUS, Postfix, p.parsePostfixExpression)
	p.registerLedOp(token.DOUBLE_MINUS, Postfix, p.parsePostfixExpression)
	p.registerLedOp(token.QUESTION, Postfix, p.parsePostfixExpression, false, true)
	p.registerLedOp(token.BANG, Postfix, p.parsePostfixExpression, false, true)

	// Prefix operators
	p.registerNudFn(token.MINUS, p.parsePrefixExpression)
	p.registerNudFn(token.PLUS, p.parsePrefixExpression)
	p.registerNudFn(token.BANG, p.parsePrefixExpression)
	p.registerNudFn(token.STAR, p.parseRefDeref)
	p.registerNudFn(token.AMPERSAND, p.parseRefDeref)
	p.registerNudFn(token.TILDE, p.parsePrefixExpression)

	// Assignment
	p.registerLedOp(token.EQUALS, Assignment, p.parseAssignmentExpression, true)
	p.registerLedOp(token.PLUS_EQUALS, Assignment, p.parseAssignmentExpression, true)
	p.registerLedOp(token.MINUS_EQUALS, Assignment, p.parseAssignmentExpression, true)
	p.registerLedOp(token.STAR_EQUALS, Assignment, p.parseAssignmentExpression, true)
	p.registerLedOp(token.SLASH_EQUALS, Assignment, p.parseAssignmentExpression, true)
	p.registerLedOp(token.PERCENT_EQUALS, Assignment, p.parseAssignmentExpression, true)

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
	p.registerLedOp(token.PIPE, Bitwise, p.parseBinaryExpression, false, true)
	p.registerLedOp(token.AMPERSAND, Bitwise, p.parseBinaryExpression)

	p.registerLedOp(token.PLUS, Additive, p.parseBinaryExpression)
	p.registerLedOp(token.MINUS, Additive, p.parseBinaryExpression)

	p.registerLedOp(token.STAR, Multiplicative, p.parseBinaryExpression)
	p.registerLedOp(token.SLASH, Multiplicative, p.parseBinaryExpression)
	p.registerLedOp(token.PERCENT, Multiplicative, p.parseBinaryExpression)

	p.registerLedOp(token.DOUBLE_STAR, Exponential, p.parseBinaryExpression, true)

	p.registerLedLookup(func(_ ast.Expression) (opInfo, bool) {
		if p.isKeyword("is") {
			return opInfo{
				leftPrecedence:  Comparison,
				rightPrecedence: Comparison,
				parseFn:         p.parseTypeCheckExpression,
			}, true
		}

		return opInfo{}, false
	})
}

func (p *parser) isKeyword(value string) bool {
	if p.next().Kind != token.IDENTIFIER || p.next().Value != value {
		return false
	}

	_, isDeclared := p.identifiers[p.next().Value]

	return !isDeclared
}

func (p *parser) delcareIdentifier() token.Token {
	ident := p.expect(token.IDENTIFIER)
	p.identifiers[ident.Value] = ident.Location
	return ident
}

func (p *parser) enterScope() map[string]text.Location {
	oldScope := p.identifiers
	p.identifiers = make(map[string]text.Location, len(oldScope))

	for ident, span := range oldScope {
		p.identifiers[ident] = span
	}

	return oldScope
}

func (p *parser) exitScope(scope map[string]text.Location) {
	p.identifiers = scope
}

func (p *parser) next() token.Token {
	return p.peek(0)
}

func (p *parser) peek(offset int) token.Token {
	if p.pos+offset >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}

	for p.tokens[p.pos+offset].Kind == token.NEWLINE || p.tokens[p.pos+offset].Kind == token.COMMENT {
		offset++
	}

	return p.tokens[p.pos+offset]
}

func (p *parser) nextWithNewlines() token.Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}

	offset := 0
	for p.tokens[p.pos+offset].Kind == token.COMMENT {
		offset++
	}

	return p.tokens[p.pos+offset]
}

func (p *parser) canContinue() bool {
	return p.nextWithNewlines().Kind != token.NEWLINE ||
		p.next().Kind == token.DOT ||
		p.bracketLevel > 0
}

func (p *parser) consume() token.Token {
	p.consumeNewlines()
	next := p.next()

	p.pos++
	return next
}

func (p *parser) consumeNewlines() {
	for p.tokens[p.pos].Kind == token.NEWLINE ||
		p.tokens[p.pos].Kind == token.SEMICOLON ||
		p.tokens[p.pos].Kind == token.COMMENT {
		p.pos++
	}
}

func (p *parser) expect(kind token.Kind) token.Token {
	if p.next().Kind == kind {
		return p.consume()
	}
	p.Diagnostics.Report(diagnostics.ExpectedToken(p.next().Location, kind, p.next().Kind))
	span := text.NewSpan(0, 0, 0, 0)
	location := text.Location{
		Span: span,
		File: p.next().Location.File,
	}
	tok := token.New(kind, "", location)
	return tok
}

func (p *parser) expectKeyword(keyword string) token.Token {
	if p.isKeyword(keyword) {
		return p.consume()
	}

	if p.next().Kind == token.IDENTIFIER && p.next().Value == keyword {
		p.Diagnostics.Report(diagnostics.KeywordOverwritten(p.next().Location, keyword, p.identifiers[keyword])...)
		return p.consume()
	}

	p.Diagnostics.Report(diagnostics.ExpectedKeyword(p.next().Location, keyword, p.next()))
	span := text.NewSpan(0, 0, 0, 0)
	location := text.Location{
		Span: span,
		File: p.next().Location.File,
	}
	tok := token.New(token.IDENTIFIER, "", location)
	return tok
}

func (p *parser) expectNewline() {
	if p.nextWithNewlines().Kind != token.NEWLINE &&
		p.next().Kind != token.SEMICOLON {
		p.Diagnostics.Report(diagnostics.ExpectedNewline(p.next().Location, p.next().Kind))
	}

	p.consumeNewlines()
}

func (p *parser) eof() bool {
	return p.next().Kind == token.EOF
}
