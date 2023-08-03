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
	case bool:
		return MakeBoolean(value)
	case nil:
		return MakeNull()
	default:
		errors.DevError(fmt.Sprintf("Cannot create runtime value of type %T", v))
		return MakeInteger(0)
	}
}

var typeStringMap = map[string]string {
	"int": "integer",
	"bool": "boolean",
	"null": "null",
}

func TypeToString[T any]() string {
	tValue := struct{t T}{}.t

	typeString, ok := typeStringMap[fmt.Sprintf("%T", tValue)]

	if !ok {
		errors.DevError(fmt.Sprintf("Invalid runtime value type: %T", tValue))
	}

	return typeString
}
