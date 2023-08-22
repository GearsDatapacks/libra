package types

import (
	"strings"

	"github.com/gearsdatapacks/libra/utils"
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
	Types []DataType
}

func MakeUnion(types ...DataType) *Union {
	return &Union{Types: types}
}

func (u *Union) Valid(dataType ValidType) bool {
	for _, unionType := range u.Types {
		if dataType.valid(unionType) {
			return true
		}
	}

	return false
}

func (u *Union) valid(dataType DataType) bool {
	return utils.Contains(u.Types, dataType)
}

func (u *Union) String() string {
	return strings.Join(u.Types, " | ")
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
