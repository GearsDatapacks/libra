package values

import "log"

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
		log.Fatalf("Cannot create runtime value of type %v", v)
		return MakeInteger(0)
	}
}
