package ir

import (
	"bytes"
	"os"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

type Statement interface {
	printer.Printable
}

type Expression interface {
	Statement
	Type() types.Type
	IsConst() bool
	ConstValue() values.ConstValue
	irExpr()
}

type Package struct {
	Modules map[string]*Module
}

func (p *Package) String() string {
	var text bytes.Buffer

	irPrinter := printer.New(&text, false)
	for _, kv := range printer.SortMap(p.Modules) {
		irPrinter.Node(kv.Value)
	}
	irPrinter.Print()

	return text.String()
}

func (p *Package) Print() {
	irPrinter := printer.New(os.Stdout, true)
	for _, kv := range printer.SortMap(p.Modules) {
		irPrinter.Node(kv.Value)
	}
	irPrinter.Print()
}

type Module struct {
	Name       string
	Statements []Statement
}

func (m *Module) Print(node *printer.Node) {
	node.Text(
		"%sMODULE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		m.Name,
	)

	printer.Nodes(node, m.Statements)
}

func AssignableExpr(expr Expression) bool {
	switch e := expr.(type) {
	case *VariableExpression:
		return true
	case *IndexExpression:
		return AssignableExpr(e.Left)
	case *MemberExpression:
		return AssignableExpr(e.Left)
	case *DerefExpression:
		return AssignableExpr(e.Value)
	case *InvalidExpression:
		return true
	default:
		return false
	}
}

func MutableExpr(expr Expression) bool {
	switch e := expr.(type) {
	// TODO: Currently, when trying to assign to a field of a referenced struct, this will check if the pointer variable itself is mutable,
	// not whether the pointer is a mutable pointer, which means it will often be incorrect. This needs to be fixed somehow
	case *VariableExpression:
		return e.Symbol.IsMut
	case *IndexExpression:
		return MutableExpr(e.Left)
	case *MemberExpression:
		return MutableExpr(e.Left)
	case *DerefExpression:
		return e.Value.Type().(*types.Pointer).Mutable
	case *InvalidExpression:
		return true
	default:
		return false
	}
}

func Index(left, index Expression) (types.Type, *diagnostics.Partial) {
	if index.IsConst() {
		if left.IsConst() {
			return types.Index(left.Type(), index.Type(), index.ConstValue(), left.ConstValue())
		}
		return types.Index(left.Type(), index.Type(), index.ConstValue())
	}
	return types.Index(left.Type(), index.Type())
}

func Member(left Expression, member string) (types.Type, *diagnostics.Partial) {
	if left.IsConst() {
		return types.Member(left.Type(), member, left.ConstValue())
	}
	return types.Member(left.Type(), member)
}
