package ir

import (
	"bytes"
	"os"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
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
	Imports      []*ImportStatement
	Types        []*TypeDeclaration
	MainFunction *FunctionDeclaration
	Functions    []*FunctionDeclaration
	Globals      []*VariableDeclaration
	// For ABI passes
	FunctionCalls []*FunctionCall
}

func (m *LoweredModule) Print(node *printer.Node) {
	node.Text(
		"%sMODULE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		m.Name,
	)

	printer.Nodes(node, m.Imports)
	printer.Nodes(node, m.Types)
	printer.Nodes(node, m.Functions)
	printer.Nodes(node, m.Globals)
}

type Label struct {
	Location text.Location
	Name string
}

func (l *Label) GetLocation() text.Location {
	return l.Location
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
	Location text.Location
	Label string
}

func (g *Goto) GetLocation() text.Location {
	return g.Location
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
	Location text.Location
	Condition Expression
	Label     string
}

func (g *GotoIf) GetLocation() text.Location {
	return g.Location
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

type GotoUnless struct {
	Location text.Location
	Condition Expression
	Label     string
}

func (g *GotoUnless) GetLocation() text.Location {
	return g.Location
}

func (g *GotoUnless) Print(node *printer.Node) {
	node.
		Text(
			"%sGOTO_UNLESS %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			g.Label,
		).
		Node(g.Condition)
}

type Branch struct {
	Location text.Location
	Condition Expression
	IfLabel   string
	ElseLabel string
}

func (b *Branch) GetLocation() text.Location {
	return b.Location
}

func (b *Branch) Print(node *printer.Node) {
	node.
		Text(
			"%sBRANCH %s%s %selse %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			b.IfLabel,
			node.Colour(colour.Attribute),
			node.Colour(colour.Name),
			b.ElseLabel,
		).
		Node(b.Condition)
}

type BitCast struct {
	expression
	Location text.Location
	Value Expression
	To    types.Type
}

func (b *BitCast) GetLocation() text.Location {
	return b.Location
}

func (b *BitCast) Print(node *printer.Node) {
	node.
		Text("%sBIT_CAST", node.Colour(colour.NodeName)).
		Node(b.Value).
		Node(b.To)
}

func (b *BitCast) Type() types.Type {
	return b.To
}

func (b *BitCast) IsConst() bool {
	return false
}

func (b *BitCast) ConstValue() values.ConstValue {
	return nil
}
