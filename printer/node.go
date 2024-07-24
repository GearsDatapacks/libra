package printer

import (
	"cmp"
	"fmt"
	"reflect"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/text"
)

type Node struct {
	text     string
	children []Node
	printer  *printer
	rejected bool
}

func (n *Node) FakeNode(text string, callback func(*Node), values ...any) *Node {
	n.node(fmt.Sprintf(text, values...), callback)
	return n
}

func (n *Node) Node(p Printable) *Node {
	n.node("", p.Print)
	return n
}

func isNil(p Printable) bool {
	if p == nil {
		return true
	}
	value := reflect.ValueOf(p)
	return value.Kind() == reflect.Pointer && value.IsNil()
}

func (n *Node) OptionalNode(p Printable) *Node {
	if !isNil(p) {
		n.Node(p)
	}
	return n
}

func (n *Node) node(text string, callback func(*Node)) {
	newNode := Node{text: text, printer: n.printer}
	if callback != nil {
		callback(&newNode)
	}
	if !newNode.rejected {
		n.children = append(n.children, newNode)
	}
}

func Nodes[T Printable](n *Node, nodes []T) {
	for _, node := range nodes {
		n.Node(node)
	}
}

func Map[K cmp.Ordered, V Printable](n *Node, m map[K]V) {
	nodes := SortMap(m)
	for _, node := range nodes {
		n.Node(node.Value)
	}
}

func (n *Node) Text(text string, values ...any) *Node {
	n.text += fmt.Sprintf(text, values...)
	return n
}

func (n *Node) TextIf(condition bool, text string, values ...any) *Node {
	if condition {
		n.Text(text, values...)
	}
	return n
}

func (n *Node) Location(node interface{ GetLocation() text.Location }) *Node {
	location := node.GetLocation()
	n.Text(
		" %s(%d:%d)",
		n.Colour(colour.Location),
		location.Span.Start,
		location.Span.End,
	)
	return n
}

func (n *Node) Colour(colour colour.Colour) string {
	return n.printer.colour(colour)
}

func (n *Node) Reject() {
	n.rejected = true
}
