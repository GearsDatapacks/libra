package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/text"
)

const INDENT_STEP = "  "

type Printable interface {
	Print(*Printer)
}

type node struct {
	indent   uint32
	text     string
	children []node
	rejected bool
}

type Printer struct {
	writer    io.Writer
	indent    uint32
	useColour bool
	stack     []node
}

func New(writer io.Writer, useColour bool) *Printer {
	return &Printer{
		writer:    writer,
		indent:    0,
		useColour: useColour,
		stack:     []node{},
	}
}

func (p *Printer) write(format string, values ...any) {
	fmt.Fprintf(p.writer, format, values...)
}

func (printer *Printer) QueueNode(printable Printable, noIndent ...bool) {
	indent := !append(noIndent, false)[0]
	if indent {
		printer.Nest()
	}

	printer.stack = append(printer.stack, node{
		text:   "",
		indent: printer.indent,
	})

	printable.Print(printer)
	printer.completeNode()
	if indent {
		printer.UnNest()
	}
}

func (p *Printer) QueueInfo(info string, values ...any) {
	lastNode := &p.stack[len(p.stack)-1]
	lastNode.children = append(lastNode.children, node{
		indent: p.indent,
		text:   fmt.Sprintf(info, values...),
	})
}

func (p *Printer) AddInfo(info string, values ...any) {
	lastNode := &p.stack[len(p.stack)-1]
	lastNode.children[len(lastNode.children)-1].text += fmt.Sprintf(info, values...)
}

func (p *Printer) AddLocation(node interface{GetLocation()text.Location}) {
	location := node.GetLocation()
	p.AddInfo(
		" %s(%d:%d)",
		p.Colour(colour.Location),
		location.Span.Start,
		location.Span.End,
	)
}

func QueueNodeList[T Printable](p *Printer, nodes []T, noIndent ...bool) {
	for _, node := range nodes {
		p.QueueNode(node, noIndent...)
	}
}

func (p *Printer) Colour(colour colour.Colour) string {
	if p.useColour {
		return string(colour)
	}
	return ""
}

func (p *Printer) Nest() {
	p.indent++
}

func (p *Printer) UnNest() {
	p.indent--
}

func (p *Printer) RejectNode() {
	p.stack[len(p.stack)-1].rejected = true
}

func pop[T any](slice *[]T) T {
	s := *slice
	var value T
	value, *slice = s[len(s)-1], s[:len(s)-1]
	return value
}

func (p *Printer) completeNode() {
	node := pop(&p.stack)
	if len(p.stack) == 0 {
		p.printNode(node)
	} else if !node.rejected {
		lastNode := &p.stack[len(p.stack)-1]
		lastNode.children = append(lastNode.children, node)
	}
}

func (p *Printer) printNode(node node) {
	if len(node.text) != 0 {
		p.write("\n%s%s", strings.Repeat(INDENT_STEP, int(node.indent)), node.text)
	}
	
	for _, node := range node.children {
		p.printNode(node)
	}
}
