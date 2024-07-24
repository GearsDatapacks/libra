package symbols

import (
	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

type Symbol interface {
	GetType() types.Type
	Value() values.ConstValue
	GetName() string
	Mutable() bool
}

type Variable struct {
	Name       string
	IsMut      bool
	Type       types.Type
	ConstValue values.ConstValue
}

func (v *Variable) Value() values.ConstValue {
	return v.ConstValue
}

func (v *Variable) GetType() types.Type {
	return v.Type
}

func (v *Variable) GetName() string {
	return v.Name
}

func (v *Variable) Mutable() bool {
	return v.IsMut
}

func (v *Variable) Print(node *printer.Node) {
	node.
		Text(
			"%sVAR_SYMBOL %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			v.Name,
		).
		TextIf(v.IsMut, " %smut", node.Colour(colour.Attribute)).
		Node(v.Type).
		OptionalNode(v.ConstValue)
}

type Type struct {
	Name string
	Type types.Type
}

func (t *Type) Value() values.ConstValue {
	return values.TypeValue{Type: t.Type}
}

func (t *Type) GetType() types.Type {
	return types.RuntimeType
}

func (t *Type) GetName() string {
	return t.Name
}

func (*Type) Mutable() bool {
	return false
}

type Method struct {
	MethodOf types.Type
	Static   bool
	Function *types.Function
}
