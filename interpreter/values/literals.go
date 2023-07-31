package values

import "fmt"

type IntegerLiteral struct {
	*BaseValue
	value int
}

func MakeInteger(value int) *IntegerLiteral {
	return &IntegerLiteral{value: value}
}

func (il *IntegerLiteral) Value() any {
	return il.value
}

func (il *IntegerLiteral) Type() ValueType {
	return "integer"
}

func (il *IntegerLiteral) ToString() string {
	return fmt.Sprint(il.value)
}

func (il *IntegerLiteral) truthy() bool {
	return il.value != 0
}

func (il *IntegerLiteral) equalTo(value RuntimeValue) bool {
	integer, ok := value.(*IntegerLiteral)

	return ok && integer.value == il.value
}
