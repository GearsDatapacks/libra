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
}

func (p *parser) next() token.Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}

  for p.tokens[p.pos].Kind == token.NEWLINE {
    p.pos++
  }

	return p.tokens[p.pos]
}

func (p *parser) nextWithNewlines() token.Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}

	return p.tokens[p.pos]
}

func (p *parser) consume() token.Token {
	next := p.next()
	p.pos++
	return next
}

func (p *parser) eof() bool {
  return p.next().Kind == token.EOF
}

