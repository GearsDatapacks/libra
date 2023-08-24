package types

import (
	"strings"
)

type Literal struct {
	DataType DataType
}

func MakeLiteral(dataType DataType) *Literal {
	return &Literal{DataType: dataType}
}

func (lit *Literal) Valid(dataType ValidType) bool {
	return dataType.valid(lit.DataType)
}

func (lit *Literal) valid(dataType DataType) bool {
	return lit.DataType == dataType
}

func (lit *Literal) String() string {
	return string(lit.DataType)
}

type Union struct {
	Types []ValidType
}

func MakeUnion(types ...ValidType) *Union {
	return &Union{Types: types}
}

func (u *Union) Valid(dataType ValidType) bool {
	for _, unionType := range u.Types {
		if dataType.Valid(unionType) {
			return true
		}
	}

	return false
}

func (u *Union) valid(dataType DataType) bool {
	for _, unionType := range u.Types {
		if unionType.valid(dataType) {
			return true
		}
	}

	return false
}

func (u *Union) String() string {
	typeStrings := []string{}

	for _, dataType := range u.Types {
		typeStrings = append(typeStrings, dataType.String())
	}

	return strings.Join(typeStrings, " | ")
}

type Function struct {
	Parameters []ValidType
	ReturnType ValidType
}

func (fn *Function) Valid(dataType ValidType) bool {
	return dataType.valid(FUNCTION)
}

func (fn *Function) valid(dataType DataType) bool {
	return dataType == FUNCTION
}

func (fn *Function) String() string {
	return FUNCTION
}

type Void struct {}

func (v *Void) Valid(dataType ValidType) bool {
	_, isVoid := dataType.(*Void)
	return isVoid
}

func (v *Void) valid(dataType DataType) bool {
	return false
}

func (v *Void) String() string {
	return "void"
}

type Any struct {}

func (a *Any) Valid(dataType ValidType) bool { return true }
func (a *Any) valid(dataType DataType) bool { return true }

func (a *Any) String() string {
	return "any"
}
