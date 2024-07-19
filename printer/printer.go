package printer

import (
	"fmt"
	"io"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/text"
)

type Printable interface {
	Print(*Printer)
}

type node struct {
	text     string
	children []node
	rejected bool
}

type Printer struct {
	writer    io.Writer
	useColour bool
	isFirst bool
	stack     []node
}

func New(writer io.Writer, useColour bool) *Printer {
	return &Printer{
		writer:    writer,
		useColour: useColour,
		isFirst: true,
		stack:     []node{},
	}
}

func (p *Printer) write(format string, values ...any) {
	fmt.Fprintf(p.writer, format, values...)
}

func (printer *Printer) QueueNode(printable Printable) {
	printer.queueNode("", printable.Print)
}

func (printer *Printer) queueNode(text string, callback func(*Printer)) {
	printer.stack = append(printer.stack, node{
		text: text,
	})

	if callback != nil {
		callback(printer)
	}
	printer.completeNode()
}

func (p *Printer) QueueInfo(info string, callback func(*Printer), values ...any) {
	p.queueNode(fmt.Sprintf(info, values...), callback)
}

func (p *Printer) AddInfo(info string, values ...any) {
	p.stack[len(p.stack)-1].text += fmt.Sprintf(info, values...)
}

func (p *Printer) AddLocation(node interface{ GetLocation() text.Location }) {
	location := node.GetLocation()
	p.AddInfo(
		" %s(%d:%d)",
		p.Colour(colour.Location),
		location.Span.Start,
		location.Span.End,
	)
}

func QueueNodeList[T Printable](p *Printer, nodes []T) {
	for _, node := range nodes {
		p.QueueNode(node)
	}
}

func (p *Printer) Colour(colour colour.Colour) string {
	if p.useColour {
		return string(colour)
	}
	return ""
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
	p.doPrintNode(node, "", "", "  ")
}

func (p *Printer) doPrintNode(
	node node,
	prefix, tree, prefixAddition string,
) {
	if p.isFirst {
		p.isFirst = false
	} else {
		p.write("\n")
	}

	p.write(
		"%s%s%s%s",
		p.Colour(colour.Symbol),
		prefix,
		tree,
		node.text,
	)

	prefix += prefixAddition
	numNodes := len(node.children)

	for i, node := range node.children {
		var nextTree, nextAddition string
		if i == numNodes-1 {
			nextTree = "└─"
			nextAddition = "  "
		} else {
			nextTree = "├─"
			nextAddition = "│ "
		}

		p.doPrintNode(node, prefix, nextTree, nextAddition)
	}
}
