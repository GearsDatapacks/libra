package values

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
)

type ValueType string

type BaseValue struct{}

type RuntimeValue interface {
	Value() any
	Type() ValueType
	ToString() string
	Truthy() bool
	EqualTo(RuntimeValue) bool
}

func MakeValue(v any) RuntimeValue {
	switch value := v.(type) {
	case int:
		return MakeInteger(value)
	case float64:
		return MakeFloat(value)
	case bool:
		return MakeBoolean(value)
	case string:
		return MakeString(value)
	case nil:
		return MakeNull()
	default:
		errors.DevError(fmt.Sprintf("Cannot create runtime value of type %T", v))
		return MakeInteger(0)
	}
}

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
		errors.DevError(fmt.Sprintf("Cannot create runtime value of type %s", dataType))
		return MakeInteger(0)
	}
}
