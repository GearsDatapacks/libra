package values

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
)

type ValueType string

type BaseValue struct {}

type RuntimeValue interface {
	Value() any
	Type() ValueType
	ToString() string
	truthy() bool
	equalTo(RuntimeValue) bool
}

func MakeValue(v any) RuntimeValue {
	switch value := v.(type) {
	case int:
		return MakeInteger(value)
	case float64:
		return MakeFloat(value)
	case bool:
		return MakeBoolean(value)
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
	default:
		errors.DevError(fmt.Sprintf("Cannot create runtime value of type %s", dataType))
		return MakeInteger(0)
	}
}

var typeStringMap = map[string]string {
	"int": "integer",
	"float64": "float",
	"bool": "boolean",
	"null": "null",
	"func": "function",
}

func TypeToString[T any]() string {
	tValue := struct{t T}{}.t

	typeString, ok := typeStringMap[fmt.Sprintf("%T", tValue)]

	if !ok {
		errors.DevError(fmt.Sprintf("Invalid runtime value type: %T", tValue))
	}

	return typeString
}
