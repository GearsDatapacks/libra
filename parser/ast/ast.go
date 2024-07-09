package ast

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/text"
)

const INDENT_STEP = "  "

type Statement interface {
	printable
	Tokens() []token.Token
}

type Expression interface {
	Statement
	expressionNode()
	Location() text.Location
}

type Program struct {
	Statements []Statement
}

func (p *Program) Tokens() []token.Token {
	tokens := []token.Token{}

	for _, statement := range p.Statements {
		tokens = append(tokens, statement.Tokens()...)
	}

	return tokens
}

func (p *Program) String() string {
	var text bytes.Buffer

	context := printContext{
		writer:    &text,
		indent:    0,
		useColour: false,
	}
	writeNodeList(context, p.Statements)

	return text.String()
}

func (p *Program) Print() {
	context := printContext{
		writer:    os.Stdout,
		indent:    0,
		useColour: true,
	}
	writeNodeList(context, p.Statements)
}

type printable interface {
	String(printContext)
}

type printContext struct {
	writer    io.Writer
	indent    uint32
	useColour bool
}

func (p *printContext) write(format string, values ...any) {
	fmt.Fprintf(p.writer, format, values...)
}

func (p *printContext) writeNest() {
	fmt.Fprintf(p.writer, "\n%s", strings.Repeat(INDENT_STEP, int(p.indent)))
}

func (p *printContext) withNest() printContext {
	context := p.nested()
	context.writeNest()
	return context
}

func (p *printContext) nested() printContext {
	context := *p
	context.indent++
	return context
}

func (p *printContext) writeNode(node printable) {
	node.String(p.withNest())
}

func writeNodeList[T printable](p printContext, nodes []T) {
	for i, node := range nodes {
		if i != 0 {
			p.writeNest()
		}
		node.String(p)
	}
}

func (p *printContext) colour(colour colour.Colour) string {
	if p.useColour {
		return string(colour)
	}
	return ""
}
