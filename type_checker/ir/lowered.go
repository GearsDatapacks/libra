package ir

import (
	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/printer"
)

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
