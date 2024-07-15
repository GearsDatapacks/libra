package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/gearsdatapacks/libra/colour"
)

const INDENT_STEP = "  "

type Printable interface {
	String(Context)
}

type Context struct {
	writer    io.Writer
	indent    uint32
	useColour bool
}

func New(writer io.Writer, useColour bool) Context {
	return Context{
		writer:    writer,
		indent:    0,
		useColour: useColour,
	}
}

func (p *Context) Write(format string, values ...any) {
	fmt.Fprintf(p.writer, format, values...)
}

func (p *Context) WriteNest() {
	fmt.Fprintf(p.writer, "\n%s", strings.Repeat(INDENT_STEP, int(p.indent)))
}

func (p *Context) WithNest() Context {
	context := p.Nested()
	context.WriteNest()
	return context
}

func (p *Context) Nested() Context {
	context := *p
	context.indent++
	return context
}

func (p *Context) WriteNode(node Printable) {
	node.String(p.WithNest())
}

func WriteNodeList[T Printable](p Context, nodes []T) {
	for i, node := range nodes {
		if i != 0 {
			p.WriteNest()
		}
		node.String(p)
	}
}

func (p *Context) Colour(colour colour.Colour) string {
	if p.useColour {
		return string(colour)
	}
	return ""
}

func (p Context) WithWriter(writer io.Writer) Context {
	p.writer = writer
	return p
}
