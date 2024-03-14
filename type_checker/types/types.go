package types

import (
	"fmt"

	"github.com/gearsdatapacks/libra/parser/ast"
)

type Type interface {
	String() string
	valid(Type) bool
}

func Assignable(to, from Type) bool {
	if to == Invalid || from == Invalid {
		return true
	}

	return to.valid(from)
}

func ToReal(ty Type) Type {
	if pseudo, ok := ty.(Pseudo); ok {
		return pseudo.ToReal()
	}
	return ty
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

func (v VariableType) ToReal() Type {
	if v.Untyped {
		v.Untyped = false
	}
	return v
}

type Pseudo interface {
	ToReal() Type
}

func FromAst(node ast.TypeExpression) Type {
	switch ty := node.(type) {
	case *ast.TypeName:
		return lookupType(ty.Name.Value)
	default:
		panic(fmt.Sprintf("TODO: Types from %T", ty))
	}
}

func lookupType(name string) Type {
	switch name {
	case "i32":
		return Int
	case "f32":
		return Float
	case "bool":
		return Bool
	case "string":
		return String
	default:
		return nil
	}
}
