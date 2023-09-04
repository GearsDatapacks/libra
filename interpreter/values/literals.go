package values

import (
	"fmt"
	"strings"

	"github.com/gearsdatapacks/libra/parser/ast"
)

type IntegerLiteral struct {
	BaseValue
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

func (il *IntegerLiteral) Truthy() bool {
	return il.value != 0
}

func (il *IntegerLiteral) EqualTo(value RuntimeValue) bool {
	integer, ok := value.(*IntegerLiteral)

	return ok && integer.value == il.value
}

type FloatLiteral struct {
	BaseValue
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

func (fl *FloatLiteral) Truthy() bool {
	return fl.value != 0
}

func (fl *FloatLiteral) EqualTo(value RuntimeValue) bool {
	float, ok := value.(*FloatLiteral)

	return ok && float.value == fl.value
}

type StringLiteral struct {
	BaseValue
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
	return "\"" + str.value + "\""
}

func (str *StringLiteral) Truthy() bool {
	return len(str.value) != 0
}

func (str *StringLiteral) EqualTo(value RuntimeValue) bool {
	s, ok := value.(*StringLiteral)

	return ok && s.value == str.value
}

type NullLiteral struct {
	BaseValue
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

func (nl *NullLiteral) Truthy() bool {
	return false
}

func (nl *NullLiteral) EqualTo(value RuntimeValue) bool {
	_, ok := value.(*NullLiteral)
	return ok
}

type BooleanLiteral struct {
	BaseValue
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

func (bl *BooleanLiteral) Truthy() bool {
	return bl.value
}

func (bl *BooleanLiteral) EqualTo(value RuntimeValue) bool {
	boolean, ok := value.(*BooleanLiteral)

	return ok && boolean.value == bl.value
}

type ListLiteral struct {
	BaseValue
	Elements []RuntimeValue
}

func (list *ListLiteral) Value() any {
	return list.Elements
}

func (list *ListLiteral) Type() ValueType {
	return "list"
}

func (list *ListLiteral) ToString() string {
	result := "["

	elemStrings := []string{}

	for _, elem := range list.Elements {
		elemStrings = append(elemStrings, elem.ToString())
	}

	result += strings.Join(elemStrings, ", ")
	result += "]"
	return result
}

func (list *ListLiteral) EqualTo(other RuntimeValue) bool {
	otherList, ok := other.(*ListLiteral)
	if !ok {
		return false
	}

	if len(otherList.Elements) != len(list.Elements) {
		return false
	}

	for i, elem := range list.Elements {
		if !elem.EqualTo(otherList.Elements[i]) {
			return false
		}
	}

	return true
}

func (list *ListLiteral) Truthy() bool {
	return len(list.Elements) != 0
}

type FunctionValue struct {
	BaseValue
	Name                   string
	Parameters             []string
	DeclarationEnvironment any
	Body                   []ast.Statement
}

func (fn *FunctionValue) Value() any {
	return fn.Name
}

func (fn *FunctionValue) Type() ValueType {
	return "function"
}

func (fn *FunctionValue) ToString() string {
	result := "fn ("

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

func (fn *FunctionValue) Truthy() bool {
	return true
}

func (fn *FunctionValue) EqualTo(value RuntimeValue) bool {
	function, ok := value.(*FunctionValue)

	return ok && function.Name == fn.Name
}
