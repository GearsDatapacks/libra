package printer

import (
	"fmt"
	"io"

	"github.com/gearsdatapacks/libra/colour"
)

type Printable interface {
	Print(*Node)
}

type printer struct {
	writer    io.Writer
	useColour bool
	isFirst   bool
	root      Node
}

func New(writer io.Writer, useColour bool) *printer {
	p := &printer{
		writer:    writer,
		useColour: useColour,
		isFirst:   true,
	}

	p.root = Node{printer: p}

	return p
}

func (printer *printer) Node(printable Printable) {
	printer.root.Node(printable)
}

func (p *printer) Print() {
	for _, node := range p.root.children {
		p.doPrintNode(node, "", "", "")
	}
}

func (p *printer) write(format string, values ...any) {
	fmt.Fprintf(p.writer, format, values...)
}

func (p *printer) colour(colour colour.Colour) string {
	if p.useColour {
		return string(colour)
	}
	return ""
}

func (p *printer) doPrintNode(
	node Node,
	prefix, tree, prefixAddition string,
) {
	if p.isFirst {
		p.isFirst = false
	} else {
		p.write("\n")
	}

	p.write(
		"%s%s%s%s",
		p.colour(colour.Symbol),
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
