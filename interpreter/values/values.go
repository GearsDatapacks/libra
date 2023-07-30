package values

type ValueType string

type BaseValue struct {}

type RuntimeValue interface {
	Value() any
	Type() ValueType
	ToString() string
	truthy() bool
	equalTo(RuntimeValue) bool
}
