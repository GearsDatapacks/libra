package values

import (
	"fmt"
	"strings"

	"github.com/gearsdatapacks/libra/parser/ast"
)

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

type FloatLiteral struct {
	*BaseValue
	value float64
}

func MakeFloat(value float64) *FloatLiteral {
	return &FloatLiteral{value: value}
}

func (fl *FloatLiteral) Value() any {
	return fl.value
}

func (fl *FloatLiteral) Type() ValueType {
	return "float"
}

func (fl *FloatLiteral) ToString() string {
	return fmt.Sprint(fl.value)
}

func (fl *FloatLiteral) truthy() bool {
	return fl.value != 0
}

func (fl *FloatLiteral) equalTo(value RuntimeValue) bool {
	float, ok := value.(*FloatLiteral)

	return ok && float.value == fl.value
}

type StringLiteral struct {
	*BaseValue
	value string
}

func MakeString(value string) *StringLiteral {
	return &StringLiteral{value: value}
}

func (str *StringLiteral) Value() any {
	return str.value
}

func (str *StringLiteral) Type() ValueType {
	return "string"
}

func (str *StringLiteral) ToString() string {
	return  "\"" + str.value + "\""
}

func (str *StringLiteral) truthy() bool {
	return len(str.value) != 0
}

func (str *StringLiteral) equalTo(value RuntimeValue) bool {
	s, ok := value.(*StringLiteral)

	return ok && s.value == str.value
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

type FunctionValue struct {
	*BaseValue
	Name string
	Parameters []string
	DeclarationEnvironment any
	Body []ast.Statement
}

func (fn *FunctionValue) Value() any {
	return fn.Name
}

func (fn *FunctionValue) Type() ValueType {
	return "function"
}

func (fn *FunctionValue) ToString() string {
	result := "function ("

	result += strings.Join(fn.Parameters, ", ")
	result += ") {"

	for _, statement := range fn.Body {
		result += "  "
		result += statement.String()
		result += "\n"
	}

	result += "}"

	return result
}

func (fn *FunctionValue) truthy() bool {
	return true
}

func (fn *FunctionValue) equalTo(value RuntimeValue) bool {
	function, ok := value.(*FunctionValue)

	return ok && function.Name == fn.Name
}
