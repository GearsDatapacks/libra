package ir

import (
	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/printer"
)

type VariableDeclaration struct {
	Name  string
	Value Expression
}

func (v *VariableDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sVAR_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			v.Name,
		).
		Node(v.Value)
}

type FunctionDeclaration struct {
	Name       string
	Parameters []string
	Body       *Block
}

func (f *FunctionDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sFUNC_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			f.Name,
		)

	for _, param := range f.Parameters {
		node.Text(" %s%s", node.Colour(colour.Name), param)
	}

	node.Node(f.Body)
}

type ReturnStatement struct {
	Value Expression
}

func (r *ReturnStatement) Print(node *printer.Node) {
	node.
		Text(
			"%sRETURN",
			node.Colour(colour.NodeName),
		).
		OptionalNode(r.Value)
}

type BreakStatement struct {
	Value Expression
}

func (b *BreakStatement) Print(node *printer.Node) {
	node.
		Text(
			"%sBREAK",
			node.Colour(colour.NodeName),
		).
		OptionalNode(b.Value)
}

type YieldStatement struct {
	Value Expression
}

func (y *YieldStatement) Print(node *printer.Node) {
	node.
		Text(
			"%sYIELD",
			node.Colour(colour.NodeName),
		).
		Node(y.Value)
}

type ContinueStatement struct{}

func (*ContinueStatement) Print(node *printer.Node) {
	node.
		Text(
			"%sCONTINUE",
			node.Colour(colour.NodeName),
		)
}

type ImportStatement struct {
	Module    string
	Name      string
	Symbols   []string
	ImportAll bool
}

func (i *ImportStatement) Print(node *printer.Node) {
	node.
		Text(
			"%sIMPORT",
			node.Colour(colour.NodeName),
		).
		TextIf(
			i.ImportAll,
			" %s*",
			node.Colour(colour.Symbol),
		).
		Text(
			" %s%q",
			node.Colour(colour.Literal),
			i.Module,
		).
		Text(
			" %s%s",
			node.Colour(colour.Name),
			i.Name,
		)

	for _, symbol := range i.Symbols {
		node.Text(" %s", symbol)
	}
}

// TODO:
// EnumDeclaration
