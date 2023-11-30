package values

import (
	"fmt"
	"strings"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type IntegerLiteral struct {
	BaseValue
	Value int
}

func MakeInteger(value int) *IntegerLiteral {
	return &IntegerLiteral{Value: value, BaseValue: BaseValue{DataType: &types.IntLiteral{}}}
}

// func (il *IntegerLiteral) Type() ValueType {
// return "integer"
// }

func (il *IntegerLiteral) ToString() string {
	return fmt.Sprint(il.Value)
}

func (il *IntegerLiteral) Truthy() bool {
	return il.Value != 0
}

func (il *IntegerLiteral) EqualTo(value RuntimeValue) bool {
	integer, ok := value.(*IntegerLiteral)

	return ok && integer.Value == il.Value
}

type FloatLiteral struct {
	BaseValue
	Value float64
}

func MakeFloat(value float64) *FloatLiteral {
	return &FloatLiteral{Value: value, BaseValue: BaseValue{DataType: &types.FloatLiteral{}}}
}

// func (fl *FloatLiteral) Type() ValueType {
// return "float"
// }

func (fl *FloatLiteral) ToString() string {
	return fmt.Sprint(fl.Value)
}

func (fl *FloatLiteral) Truthy() bool {
	return fl.Value != 0
}

func (fl *FloatLiteral) EqualTo(value RuntimeValue) bool {
	float, ok := value.(*FloatLiteral)

	return ok && float.Value == fl.Value
}

type StringLiteral struct {
	BaseValue
	Value string
}

func MakeString(value string) *StringLiteral {
	return &StringLiteral{Value: value, BaseValue: BaseValue{DataType: &types.StringLiteral{}}}
}

// func (str *StringLiteral) Type() ValueType {
// return "string"
// }

func (str *StringLiteral) ToString() string {
	return "\"" + str.Value + "\""
}

func (str *StringLiteral) Truthy() bool {
	return len(str.Value) != 0
}

func (str *StringLiteral) EqualTo(value RuntimeValue) bool {
	s, ok := value.(*StringLiteral)

	return ok && s.Value == str.Value
}

type NullLiteral struct {
	BaseValue
}

func MakeNull() *NullLiteral {
	return &NullLiteral{BaseValue: BaseValue{DataType: &types.NullLiteral{}}}
}

// func (nl *NullLiteral) Type() ValueType {
// return "null"
// }

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
	Value bool
}

func MakeBoolean(value bool) *BooleanLiteral {
	return &BooleanLiteral{Value: value, BaseValue: BaseValue{DataType: &types.BoolLiteral{}}}
}

// func (bl *BooleanLiteral) Type() ValueType {
// return "boolean"
// }

func (bl *BooleanLiteral) ToString() string {
	return fmt.Sprint(bl.Value)
}

func (bl *BooleanLiteral) Truthy() bool {
	return bl.Value
}

func (bl *BooleanLiteral) EqualTo(value RuntimeValue) bool {
	boolean, ok := value.(*BooleanLiteral)

	return ok && boolean.Value == bl.Value
}

type ListLiteral struct {
	BaseValue
	Elements []RuntimeValue
}

// func (list *ListLiteral) Type() ValueType {
// return "list"
// }

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

func (list *ListLiteral) Index(indexValue RuntimeValue) RuntimeValue {
	index := indexValue.(*IntegerLiteral).Value
	indexSize := index
	// negative indexing
	if index < 0 {
		indexSize = -indexSize - 1
	}

	if indexSize >= len(list.Elements) {
		errors.LogError(fmt.Sprintf("Index out of range [%d] with length %d", index, len(list.Elements)))
	}

	if index < 0 {
		return list.Elements[len(list.Elements)+index]
	}
	return list.Elements[index]
}

func (list *ListLiteral) SetIndex(indexValue RuntimeValue, value RuntimeValue) RuntimeValue {
	index := indexValue.(*IntegerLiteral).Value
	indexSize := index
	// negative indexing
	if index < 0 {
		indexSize = -indexSize - 1
	}

	if indexSize >= len(list.Elements) {
		errors.LogError(fmt.Sprintf("Index out of range [%d] with length %d", index, len(list.Elements)))
	}

	if index < 0 {
		list.Elements[len(list.Elements)+index] = value
	}
	list.Elements[index] = value
	return value
}

type MapLiteral struct {
	BaseValue
	Elements map[RuntimeValue]RuntimeValue
}

// func (maplit *MapLiteral) Type() ValueType {
// return "map"
// }

func (maplit *MapLiteral) ToString() string {
	result := "{"

	elemStrings := []string{}

	for key, value := range maplit.Elements {
		elemStrings = append(elemStrings, key.ToString())
		elemStrings[len(elemStrings)-1] += ": "
		elemStrings[len(elemStrings)-1] += value.ToString()
	}

	result += strings.Join(elemStrings, ", ")
	result += "}"
	return result
}

func (maplit *MapLiteral) EqualTo(other RuntimeValue) bool {
	otherMap, ok := other.(*MapLiteral)
	if !ok {
		return false
	}

	if len(otherMap.Elements) != len(maplit.Elements) {
		return false
	}

	for key, value := range maplit.Elements {
		if !value.EqualTo(otherMap.Elements[key]) {
			return false
		}
	}

	return true
}

func (maplit *MapLiteral) Truthy() bool {
	return len(maplit.Elements) != 0
}

func (maplit *MapLiteral) Index(indexValue RuntimeValue) RuntimeValue {
	for key, value := range maplit.Elements {
		if key.EqualTo(indexValue) {
			return value
		}
	}

	return MakeNull()
}

func (maplit *MapLiteral) SetIndex(indexValue RuntimeValue, value RuntimeValue) RuntimeValue {
	maplit.Elements[indexValue] = value
	return value
}

type FunctionValue struct {
	BaseValue
	Name                   string
	Parameters             []string
	DeclarationEnvironment any
	Body                   []ast.Statement
	This                   RuntimeValue
}

// func (fn *FunctionValue) Type() ValueType {
// return "function"
// }

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

type StructLiteral struct {
	BaseValue
	Name    string
	Members map[string]RuntimeValue
}

// func (sl *StructLiteral) Type() ValueType {
// return "struct"
// }

func (sl *StructLiteral) ToString() string {
	result := "{ "

	for name, expr := range sl.Members {
		result += name
		result += ": "
		result += expr.ToString()
		result += ", "
	}

	result += "}"

	return result
}

func (sl *StructLiteral) Truthy() bool {
	return len(sl.Members) != 0
}

func (sl *StructLiteral) EqualTo(value RuntimeValue) bool {
	struc, ok := value.(*StructLiteral)
	if !ok || sl.Name != struc.Name {
		return false
	}

	for name, member := range sl.Members {
		value, ok := struc.Members[name]
		if !ok {
			return false
		}

		if !member.EqualTo(value) {
			return false
		}
	}

	return true
}

func (sl *StructLiteral) Member(member string) RuntimeValue {
	value, ok := sl.Members[member]
	if !ok {
		return nil
	}

	return value
}

func (sl *StructLiteral) SetMember(member string, value RuntimeValue) RuntimeValue {
	sl.Members[member] = value

	return value
}
