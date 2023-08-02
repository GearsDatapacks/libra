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
	default:
		errors.DevError(fmt.Sprintf("Cannot create runtime value of type %T", v))
		return MakeInteger(0)
	}
}
