package types

import (
	"bytes"
	"fmt"

	"github.com/gearsdatapacks/libra/type_checker/values"
)

type Type interface {
	String() string
	valid(Type) bool
	indexBy(Type, []values.ConstValue) Type
}

func Assignable(to, from Type) bool {
	if to == Invalid || from == Invalid {
		return true
	}

	return to.valid(from)
}

func Index(left, index Type, constVals ...values.ConstValue) Type {
	if index == Invalid {
		return Invalid
	}
	return left.indexBy(index, constVals)
}

func ToReal(ty Type) Type {
	if pseudo, ok := ty.(Pseudo); ok {
		return pseudo.ToReal()
	}
	return ty
}

func Hashable(ty Type) bool {
	switch ty.(type) {
	case PrimaryType:
		return true
	case VariableType:
		return true
	default:
		return false
	}
}

type PrimaryType int

const (
	Invalid PrimaryType = iota
	Bool
	String
)

var typeNames = map[PrimaryType]string{
	Invalid: "<?>",
	Bool:    "bool",
	String:  "string",
}

func (pt PrimaryType) String() string {
	return typeNames[pt]
}

func (pt PrimaryType) valid(other Type) bool {
	primary, isPrimary := other.(PrimaryType)
	return isPrimary && primary == pt
}

func (pt PrimaryType) indexBy(index Type, _ []values.ConstValue) Type {
	switch pt {
	case String:
		if Assignable(Int, index) {
			return String
		}
	}

	return Invalid
}

type VTKind int

const (
	_ VTKind = iota
	VT_Int
	VT_Float
)

var (
	Int = VariableType{
		Kind:    VT_Int,
		Untyped: false,
	}
	Float = VariableType{
		Kind:    VT_Float,
		Untyped: false,
	}
	UntypedInt = VariableType{
		Kind:    VT_Int,
		Untyped: true,
	}
	UntypedFloat = VariableType{
		Kind:    VT_Float,
		Untyped: true,
	}
)

type VariableType struct {
	Kind         VTKind
	Untyped      bool
	Downcastable bool
}

func (v VariableType) String() string {
	if v.Untyped {
		switch v.Kind {
		case VT_Int:
			return "untyped int"
		case VT_Float:
			return "untyped float"
		default:
			panic("unreachable")
		}
	}

	switch v.Kind {
	case VT_Int:
		return "i32"
	case VT_Float:
		return "f32"
	default:
		panic("unreachable")
	}
}

func (v VariableType) valid(other Type) bool {
	variable, ok := other.(VariableType)
	if !ok {
		return false
	}

	if variable.Untyped {
		return variable.Downcastable || variable.Kind <= v.Kind
	}

	return v.Kind == variable.Kind
}

func (v VariableType) indexBy(index Type, _ []values.ConstValue) Type {
	return Invalid
}

func (v VariableType) ToReal() Type {
	if v.Untyped {
		v.Untyped = false
	}
	return v
}

type ListType struct {
	ElemType Type
}

func (l *ListType) String() string {
	return l.ElemType.String() + "[]"
}

func (l *ListType) valid(other Type) bool {
	if list, ok := other.(*ListType); ok {
		return Assignable(l.ElemType, list.ElemType) && Assignable(list.ElemType, l.ElemType)
	}
	if array, ok := other.(*ArrayType); ok {
		return array.CanInfer && Assignable(l.ElemType, array.ElemType)
	}
	return false
}

func (l *ListType) indexBy(index Type, _ []values.ConstValue) Type {
	if Assignable(Int, index) {
		return l.ElemType
	}
	return Invalid
}

type ArrayType struct {
	ElemType Type
	Length   int
	CanInfer bool
}

func (a *ArrayType) String() string {
	if a.Length == -1 {
		return a.ElemType.String() + "[_]"
	}
	return fmt.Sprintf("%s[%d]", a.ElemType.String(), a.Length)
}

func (a *ArrayType) valid(other Type) bool {
	array, ok := other.(*ArrayType)
	if !ok {
		return false
	}

	lengthsMatch := a.Length == -1 || a.Length == array.Length
	return lengthsMatch && Assignable(a.ElemType, array.ElemType) && Assignable(array.ElemType, a.ElemType)
}

func (a *ArrayType) indexBy(index Type, _ []values.ConstValue) Type {
	if Assignable(Int, index) {
		return a.ElemType
	}
	return Invalid
}

func (a *ArrayType) ToReal() Type {
	a.CanInfer = false
	return a
}

type MapType struct {
	KeyType   Type
	ValueType Type
}

func (m *MapType) String() string {
	return fmt.Sprintf("{%s: %s}", m.KeyType.String(), m.ValueType.String())
}

func (m *MapType) valid(other Type) bool {
	mapType, ok := other.(*MapType)
	if !ok {
		return false
	}

	keysMatch := Assignable(m.KeyType, mapType.KeyType) && Assignable(mapType.KeyType, m.KeyType)
	valuesMatch := Assignable(m.ValueType, mapType.ValueType) && Assignable(mapType.ValueType, m.ValueType)
	return keysMatch && valuesMatch
}

func (m *MapType) indexBy(index Type, _ []values.ConstValue) Type {
	if Assignable(m.KeyType, index) {
		return m.ValueType
	}
	return Invalid
}

type TupleType struct {
	Types []Type
}

func (t *TupleType) String() string {
	var result bytes.Buffer

	result.WriteByte('(')
	for i, ty := range t.Types {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(ty.String())
	}
	result.WriteByte(')')

	return result.String()
}

func (t *TupleType) valid(other Type) bool {
	tuple, ok := other.(*TupleType)
	if !ok {
		return false
	}

	if len(t.Types) != len(tuple.Types) {
		return false
	}

	for i, ty := range t.Types {
		if !Assignable(ty, tuple.Types[i]) && !Assignable(tuple.Types[i], ty) {
			return false
		}
	}

	return true
}

func (a *TupleType) indexBy(t Type, constVals []values.ConstValue) Type {
	if !Assignable(Int, t) {
		return Invalid
	}

	if len(constVals) == 0 {
		return Invalid
	}

	index := constVals[0].(values.IntValue).Value
	return a.Types[index]
}

type Pseudo interface {
	ToReal() Type
}
