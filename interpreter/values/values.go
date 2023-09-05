package values

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
)

type ValueType string

type RuntimeValue interface {
	Type() ValueType
	ToString() string
	Truthy() bool
	EqualTo(RuntimeValue) bool
	Varname() string
	SetVarname(string)
	Index(RuntimeValue) RuntimeValue
}

type BaseValue struct {
	varname string
}

func (b *BaseValue) Varname() string {
	return b.varname
}

func (b *BaseValue) SetVarname(name string) {
	b.varname = name
}

func (b *BaseValue) Index(v RuntimeValue) RuntimeValue {
	return MakeNull()
}

// func MakeValue(v any) RuntimeValue {
// 	switch value := v.(type) {
// 	case int:
// 		return MakeInteger(value)
// 	case float64:
// 		return MakeFloat(value)
// 	case bool:
// 		return MakeBoolean(value)
// 	case string:
// 		return MakeString(value)
// 	case nil:
// 		return MakeNull()
// 	default:
// 		errors.DevError(fmt.Sprintf("Cannot create runtime value of type %T", v))
// 		return MakeInteger(0)
// 	}
// }

func GetZeroValue(dataType string) RuntimeValue {
	switch dataType {
	case "int":
		return MakeInteger(0)
	case "float":
		return MakeFloat(0)
	case "boolean":
		return MakeBoolean(false)
	case "null":
		return MakeNull()
	case "string":
		return MakeString("")
	default:
		errors.LogError(errors.DevError(fmt.Sprintf("Cannot create runtime value of type %s", dataType)))
		return MakeInteger(0)
	}
}
