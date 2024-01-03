package values

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type UntypedNumber struct {
	BaseValue
	Value float64
}

func MakeUntypedNumber(value float64, isFloat bool) *UntypedNumber {
	isIntAssignable := value == float64(int64(value))
	var defaultType types.ValidType = &types.IntLiteral{}
	if isFloat {
		defaultType = &types.FloatLiteral{}
	}
	return &UntypedNumber{Value: value, BaseValue: BaseValue{DataType: &types.UntypedNumber{
		Default:         defaultType,
		IsIntAssignable: isIntAssignable,
	}}}
}

// func (un *UntypedNumber) Type() ValueType {
// return "number"
// }

func (un *UntypedNumber) ToString() string {
	return fmt.Sprint(un.Value)
}

func (un *UntypedNumber) Truthy() bool {
	return un.Value != 0
}

func (un *UntypedNumber) EqualTo(value RuntimeValue) bool {
	number, ok := value.(*UntypedNumber)

	return ok && number.Value == un.Value
}

func (un *UntypedNumber) AutoCast(ty types.ValidType) RuntimeValue {
	if _, infer := ty.(*types.Infer); infer {
		return un.castTo(un.DataType.(*types.UntypedNumber).Default)
	}

	return un.castTo(ty)
}

func (un *UntypedNumber) castTo(ty types.ValidType) RuntimeValue {
	if _, ok := ty.(*types.FloatLiteral); ok {
		return MakeFloat(un.Value)
	}
	if _, ok := ty.(*types.IntLiteral); ok {
		return MakeInteger(int(un.Value))
	}
	return un
}

func (un *UntypedNumber) Copy() RuntimeValue {
	temp := *un
	return &temp
}

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

func (il *IntegerLiteral) castTo(ty types.ValidType) RuntimeValue {
	if _, ok := ty.(*types.FloatLiteral); ok {
		return MakeFloat(float64(il.Value))
	}
	return il
}

func (il *IntegerLiteral) Copy() RuntimeValue {
	temp := *il
	return &temp
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

func (fl *FloatLiteral) castTo(ty types.ValidType) RuntimeValue {
	if _, ok := ty.(*types.IntLiteral); ok {
		return MakeInteger(int(fl.Value))
	}
	return fl
}

func (fl *FloatLiteral) Copy() RuntimeValue {
	temp := *fl
	return &temp
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

func (str *StringLiteral) Copy() RuntimeValue {
	temp := *str
	return &temp
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

func (nl *NullLiteral) Copy() RuntimeValue {
	temp := *nl
	return &temp
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

func (bl *BooleanLiteral) Copy() RuntimeValue {
	temp := *bl
	return &temp
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

	index := Expect(indexValue, &types.IntLiteral{}).(*IntegerLiteral).Value
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

func (list *ListLiteral) Copy() RuntimeValue {
	// Lists are passed by reference, arrays are passed by value
	if _, isList := list.DataType.(*types.ListLiteral); isList {
		return list
	}

	elements := []RuntimeValue{}
	for _, elem := range list.Elements {
		elements = append(elements, elem.Copy())
	}
	return &ListLiteral{
		Elements:  elements,
		BaseValue: list.BaseValue,
	}
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
	for key := range maplit.Elements {
		if key.EqualTo(indexValue) {
			maplit.Elements[key] = value
			return value
		}
	}
	maplit.Elements[indexValue] = value
	return value
}

func (maplit *MapLiteral) Copy() RuntimeValue {
	return maplit
	// elements := map[RuntimeValue]RuntimeValue{}
	// for key, value := range maplit.Elements {
	// 	elements[key] = value.Copy()
	// }
	// return &MapLiteral{
	// 	Elements:  elements,
	// 	BaseValue: maplit.BaseValue,
	// }
}

type Parameter struct {
	Name string
	Type types.ValidType
}

type FunctionValue struct {
	BaseValue
	Name       string
	Parameters []Parameter
	Env        any
	Manager    any
	Body       []ast.Statement
	This       RuntimeValue
}

// func (fn *FunctionValue) Type() ValueType {
// return "function"
// }

func (fn *FunctionValue) ToString() string {
	result := "fn ("

	for i, param := range fn.Parameters {
		if i != 0 {
			result += ", "
		}
		result += param.Name
	}
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

func (fn *FunctionValue) Copy() RuntimeValue {
	return fn
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

func (sl *StructLiteral) Copy() RuntimeValue {
	members := map[string]RuntimeValue{}
	for name, member := range sl.Members {
		members[name] = member.Copy()
	}

	return &StructLiteral{
		BaseValue: sl.BaseValue,
		Name:      sl.Name,
		Members:   members,
	}
}

type TupleValue struct {
	BaseValue
	Members []RuntimeValue
}

func (tv *TupleValue) ToString() string {
	result := "("

	for i, expr := range tv.Members {
		if i != 0 {
			result += ", "
		}
		result += expr.ToString()
	}

	result += ")"

	return result
}

func (tv *TupleValue) Truthy() bool {
	return len(tv.Members) != 0
}

func (tv *TupleValue) EqualTo(value RuntimeValue) bool {
	tuple, ok := value.(*TupleValue)
	if !ok {
		return false
	}

	if len(tv.Members) != len(tuple.Members) {
		return false
	}

	for i, member := range tv.Members {
		value := tuple.Members[i]

		if !member.EqualTo(value) {
			return false
		}
	}

	return true
}

func (tv *TupleValue) Member(member string) RuntimeValue {
	number, _ := strconv.ParseInt(member, 10, 32)

	return tv.Members[number]
}

func (tv *TupleValue) SetMember(member string, value RuntimeValue) RuntimeValue {
	number, _ := strconv.ParseInt(member, 10, 32)
	tv.Members[number] = value

	return value
}

func (tv *TupleValue) Copy() RuntimeValue {
	members := []RuntimeValue{}
	for _, member := range tv.Members {
		members = append(members, member.Copy())
	}
	return &TupleValue{
		Members:   members,
		BaseValue: tv.BaseValue,
	}
}

type TupleStructValue struct {
	BaseValue
	Name    string
	Members []RuntimeValue
}

func (sl *TupleStructValue) ToString() string {
	result := "("

	for i, expr := range sl.Members {
		if i != 0 {
			result += ", "
		}
		result += expr.ToString()
	}

	result += ")"

	return result
}

func (sl *TupleStructValue) Truthy() bool {
	return true
}

func (ts *TupleStructValue) EqualTo(value RuntimeValue) bool {
	struc, ok := value.(*TupleStructValue)
	if !ok || ts.Name != struc.Name {
		return false
	}

	if len(ts.Members) != len(struc.Members) {
		return false
	}

	for i, member := range ts.Members {
		value := struc.Members[i]

		if !member.EqualTo(value) {
			return false
		}
	}

	return true
}

func (tv *TupleStructValue) Member(member string) RuntimeValue {
	number, _ := strconv.ParseInt(member, 10, 32)

	return tv.Members[number]
}

func (tv *TupleStructValue) SetMember(member string, value RuntimeValue) RuntimeValue {
	number, _ := strconv.ParseInt(member, 10, 32)
	tv.Members[number] = value

	return value
}

func (tv *TupleStructValue) Copy() RuntimeValue {
	members := []RuntimeValue{}
	for _, member := range tv.Members {
		members = append(members, member.Copy())
	}
	return &TupleStructValue{
		BaseValue: tv.BaseValue,
		Name:      tv.Name,
		Members:   members,
	}
}

type Pointer struct {
	BaseValue
	Value RuntimeValue
}

func (p *Pointer) Truthy() bool {
	return p.Value.Truthy()
}

func (p *Pointer) EqualTo(other RuntimeValue) bool {
	ptr, ok := other.(*Pointer)
	if !ok {
		return false
	}
	return p.Value == ptr.Value
}

func (p *Pointer) Copy() RuntimeValue {
	return p
}

func (p *Pointer) ToString() string {
	return "&" + p.Value.ToString()
}

func MakePointer(value RuntimeValue) *Pointer {
	return &Pointer{Value: value}
}

type Error struct {
	BaseValue
	Msg string
}

func (err *Error) ToString() string {
	return err.Msg
}

func (sl *Error) Truthy() bool {
	return true
}

func (*Error) EqualTo(value RuntimeValue) bool {
	return true
}

func (e *Error) Copy() RuntimeValue {
	temp := *e
	return &temp
}

func MakeError(msg string) RuntimeValue {
	return &Error{
		Msg: msg,
	}
}

type Module struct {
	BaseValue
	Name    string
	Exports map[string]RuntimeValue
}

func (m *Module) ToString() string {
	return m.Name
}

func (*Module) Truthy() bool {
	return true
}

func (*Module) EqualTo(value RuntimeValue) bool {
	return true
}

func (m *Module) Member(member string) RuntimeValue {
	return m.Exports[member]
}

func (m *Module) Copy() RuntimeValue {
	return m
}

type UnitStruct struct {
	BaseValue
	Name string
}

func (u *UnitStruct) ToString() string {
	return u.Name
}

func (*UnitStruct) Truthy() bool {
	return true
}

func (u *UnitStruct) EqualTo(value RuntimeValue) bool {
	unit, isUnit := value.(*UnitStruct)
	if !isUnit {
		return false
	}
	return u.Name == unit.Name
}

func (u *UnitStruct) Copy() RuntimeValue {
	return u
}
