package ir

import (
	"bytes"
	"os"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/printer"
)

type LoweredPackage struct {
	Modules map[string]*LoweredModule
}

func (p *LoweredPackage) String() string {
	var text bytes.Buffer

	irPrinter := printer.New(&text, false)
	for _, kv := range printer.SortMap(p.Modules) {
		irPrinter.Node(kv.Value)
	}
	irPrinter.Print()

	return text.String()
}

func (p *LoweredPackage) Print() {
	irPrinter := printer.New(os.Stdout, true)
	for _, kv := range printer.SortMap(p.Modules) {
		irPrinter.Node(kv.Value)
	}
	irPrinter.Print()
}

type LoweredModule struct {
	Name         string
	Types        []*TypeDeclaration
	MainFunction *FunctionDeclaration
	Functions    []*FunctionDeclaration
	Globals      []*VariableDeclaration
}

func (m *LoweredModule) Print(node *printer.Node) {
	node.Text(
		"%sMODULE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		m.Name,
	)

	printer.Nodes(node, m.Types)
	printer.Nodes(node, m.Functions)
	printer.Nodes(node, m.Globals)
}

type Label struct {
	Name string
}

func (l *Label) Print(node *printer.Node) {
	node.Text(
		"%sLABEL %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		l.Name,
	)
}

type Goto struct {
	Label string
}

func (g *Goto) Print(node *printer.Node) {
	node.Text(
		"%sGOTO %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		g.Label,
	)
}

type GotoIf struct {
	Condition Expression
	Label     string
}

func (g *GotoIf) Print(node *printer.Node) {
	node.
		Text(
			"%sGOTO_IF %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			g.Label,
		).
		Node(g.Condition)
}