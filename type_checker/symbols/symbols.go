package symbols

import (
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
