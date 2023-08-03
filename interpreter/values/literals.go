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

type NullLiteral struct {
	*BaseValue
}

func MakeNull() *NullLiteral {
	return &NullLiteral{}
}

func (nl *NullLiteral) Value() any {
	return nil
}

func (nl *NullLiteral) Type() ValueType {
	return "null"
}

func (nl *NullLiteral) ToString() string {
	return "null"
}

func (nl *NullLiteral) truthy() bool {
	return false
}

func (nl *NullLiteral) equalTo(value RuntimeValue) bool {
	_, ok := value.(*NullLiteral)
	return ok
}

type BooleanLiteral struct {
	*BaseValue
	value bool
}

func MakeBoolean(value bool) *BooleanLiteral {
	return &BooleanLiteral{value: value}
}

func (bl *BooleanLiteral) Value() any {
	return bl.value
}

func (bl *BooleanLiteral) Type() ValueType {
	return "boolean"
}

func (bl *BooleanLiteral) ToString() string {
	return fmt.Sprint(bl.value)
}

func (bl *BooleanLiteral) truthy() bool {
	return bl.value
}

func (bl *BooleanLiteral) equalTo(value RuntimeValue) bool {
	boolean, ok := value.(*BooleanLiteral)

	return ok && boolean.value == bl.value
}