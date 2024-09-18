package ir

import (
	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type VariableDeclaration struct {
	Location text.Location
	Symbol   *symbols.Variable
	Value    Expression
}

func (v *VariableDeclaration) GetLocation() text.Location {
	return v.Location
}

func (v *VariableDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sVAR_DECL",
			node.Colour(colour.NodeName),
		).
		Node(v.Symbol).
		OptionalNode(v.Value)
}

type FunctionDeclaration struct {
	Location   text.Location
	Name       string
	Parameters []string
	Body       *Block
	Type       *types.Function
	Exported   bool
	Extern     *string
}

func (f *FunctionDeclaration) GetLocation() text.Location {
	return f.Location
}

func (f *FunctionDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sFUNC_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			f.Name,
		).
		TextIf(
			f.Exported,
			" %spub",
			node.Colour(colour.Attribute),
		).
		Node(f.Type)

	for _, param := range f.Parameters {
		node.Text(" %s%s", node.Colour(colour.Name), param)
	}

	if f.Extern != nil {
		node.Text(
			" %sextern %s%s",
			node.Colour(colour.Attribute),
			node.Colour(colour.Name),
			*f.Extern,
		)
	}

	node.OptionalNode(f.Body)
}

type ReturnStatement struct {
	Location text.Location
	Value    Expression
}

func (r *ReturnStatement) GetLocation() text.Location {
	return r.Location
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
	Location text.Location
	Value    Expression
}

func (b *BreakStatement) GetLocation() text.Location {
	return b.Location
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
	Location text.Location
	Value    Expression
}

func (y *YieldStatement) GetLocation() text.Location {
	return y.Location
}

func (y *YieldStatement) Print(node *printer.Node) {
	node.
		Text(
			"%sYIELD",
			node.Colour(colour.NodeName),
		).
		Node(y.Value)
}

type ContinueStatement struct {
	Location text.Location
}

func (c *ContinueStatement) GetLocation() text.Location {
	return c.Location
}

func (*ContinueStatement) Print(node *printer.Node) {
	node.Text(
		"%sCONTINUE",
		node.Colour(colour.NodeName),
	)
}

type ImportStatement struct {
	Location  text.Location
	Module    string
	Name      string
	Symbols   []string
	ImportAll bool
}

func (i *ImportStatement) GetLocation() text.Location {
	return i.Location
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

type TypeDeclaration struct {
	Location text.Location
	Name     string
	Exported bool
	Type     types.Type
}

func (t *TypeDeclaration) GetLocation() text.Location {
	return t.Location
}

func (t *TypeDeclaration) Print(node *printer.Node) {
	node.
		Text(
			"%sTYPE_DECL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			t.Name,
		).
		TextIf(
			t.Exported,
			" %spub",
			node.Colour(colour.Attribute),
		).
		Node(t.Type)
}
