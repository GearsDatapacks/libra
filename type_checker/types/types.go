package types

import (
	"bytes"
	"fmt"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

type Type interface {
	String() string
	valid(Type) bool
}

var Context interface {
	LookupMethod(string, Type, bool) *Function
	Id() uint
}

func Assignable(to, from Type) bool {
	if to == Invalid || from == Invalid {
		return true
	}

	if alias, ok := to.(*Alias); ok {
		return Assignable(alias.Type, from)
	}
	if alias, ok := from.(*Alias); ok {
		return Assignable(to, alias.Type)
	}

	return to.valid(from)
}

func Match(a, b Type) bool {
	return Assignable(a, b) && Assignable(b, a)
}

func Index(left, index Type, constVals ...values.ConstValue) (Type, *diagnostics.Partial) {
	if index == Invalid {
		return Invalid, nil
	}
	if indexable, ok := left.(indexable); ok {
		return indexable.indexBy(index, constVals)
	}
	return Invalid, diagnostics.CannotIndex(left, index)
}

func Member(left Type, member string, constVal ...values.ConstValue) (Type, *diagnostics.Partial) {
	if left == Invalid {
		return Invalid, nil
	}

	if left == RuntimeType {
		ty := constVal[0].(values.TypeValue).Type.(Type)
		if method := Context.LookupMethod(member, ty, true); method != nil {
			return method, nil
		}
	} else if method := Context.LookupMethod(member, left, false); method != nil {
		return method, nil
	}

	if hasMember, ok := left.(hasMembers); ok {
		ty, diag := hasMember.member(member)
		if diag == nil {
			return ty, nil
		}
		return Invalid, diag
	}

	return Invalid, diagnostics.NoMember(left, member)
}

func ToReal(ty Type) Type {
	if pseudo, ok := ty.(pseudo); ok {
		return pseudo.toReal()
	}
	return ty
}

func Hashable(ty Type) bool {
	switch ty.(type) {
	case PrimaryType:
		return true
	case VariableType:
		return true
	default:
		return false
	}
}

type PrimaryType int

const (
	Invalid PrimaryType = iota
	Void
	Bool
	String
	RuntimeType
)

var typeNames = map[PrimaryType]string{
	Invalid:     "<?>",
	Void:        "void",
	Bool:        "bool",
	String:      "string",
	RuntimeType: "Type",
}

func (pt PrimaryType) String() string {
	return typeNames[pt]
}

func (pt PrimaryType) valid(other Type) bool {
	primary, isPrimary := other.(PrimaryType)
	return isPrimary && primary == pt
}

func (pt PrimaryType) indexBy(index Type, _ []values.ConstValue) (Type, *diagnostics.Partial) {
	switch pt {
	case String:
		if Assignable(Int, index) {
			return String, nil
		}
	case Invalid:
		return Invalid, nil
	}

	return Invalid, diagnostics.CannotIndex(pt, index)
}

type VTKind int

const (
	_ VTKind = iota
	VT_Int
	VT_Float
)

var (
	Int = VariableType{
		Kind:    VT_Int,
		Untyped: false,
	}
	Float = VariableType{
		Kind:    VT_Float,
		Untyped: false,
	}
	UntypedInt = VariableType{
		Kind:    VT_Int,
		Untyped: true,
	}
	UntypedFloat = VariableType{
		Kind:    VT_Float,
		Untyped: true,
	}
)

type VariableType struct {
	Kind         VTKind
	Untyped      bool
	Downcastable bool
}

func (v VariableType) String() string {
	if v.Untyped {
		switch v.Kind {
		case VT_Int:
			return "untyped int"
		case VT_Float:
			return "untyped float"
		default:
			panic("unreachable")
		}
	}

	switch v.Kind {
	case VT_Int:
		return "i32"
	case VT_Float:
		return "f32"
	default:
		panic("unreachable")
	}
}

func (v VariableType) valid(other Type) bool {
	variable, ok := other.(VariableType)
	if !ok {
		return false
	}

	if variable.Untyped {
		return variable.Downcastable || variable.Kind <= v.Kind
	}

	return v.Kind == variable.Kind
}

func (v VariableType) toReal() Type {
	if v.Untyped {
		v.Untyped = false
	}
	return v
}

type ListType struct {
	ElemType Type
}

func (l *ListType) String() string {
	return l.ElemType.String() + "[]"
}

func (l *ListType) valid(other Type) bool {
	if list, ok := other.(*ListType); ok {
		return Match(l.ElemType, list.ElemType)
	}
	if array, ok := other.(*ArrayType); ok {
		return array.CanInfer && Assignable(l.ElemType, array.ElemType)
	}
	return false
}

func (l *ListType) indexBy(index Type, _ []values.ConstValue) (Type, *diagnostics.Partial) {
	if Assignable(Int, index) {
		return l.ElemType, nil
	}
	return Invalid, diagnostics.CannotIndex(l, index)
}

func (l *ListType) Item() Type {
	return l.ElemType
}

type ArrayType struct {
	ElemType Type
	Length   int
	CanInfer bool
}

func (a *ArrayType) String() string {
	if a.Length == -1 {
		return a.ElemType.String() + "[_]"
	}
	return fmt.Sprintf("%s[%d]", a.ElemType.String(), a.Length)
}

func (a *ArrayType) valid(other Type) bool {
	array, ok := other.(*ArrayType)
	if !ok {
		return false
	}

	lengthsMatch := a.Length == -1 || a.Length == array.Length
	return lengthsMatch && Match(a.ElemType, array.ElemType)
}

func (a *ArrayType) indexBy(index Type, constVals []values.ConstValue) (Type, *diagnostics.Partial) {
	if !Assignable(Int, index) {
		return Invalid, diagnostics.CannotIndex(a, index)
	}
	if len(constVals) > 0 && a.Length != -1 {
		intIndex := int64(values.NumericValue(constVals[0]))
		if intIndex < 0 || intIndex >= int64(a.Length) {
			return Invalid, diagnostics.IndexOutOfBounds(intIndex, int64(a.Length))
		}
	}

	return a.ElemType, nil
}

func (a *ArrayType) toReal() Type {
	a.CanInfer = false
	return a
}

func (a *ArrayType) Item() Type {
	return a.ElemType
}

type MapType struct {
	KeyType   Type
	ValueType Type
}

func (m *MapType) String() string {
	return fmt.Sprintf("{%s: %s}", m.KeyType.String(), m.ValueType.String())
}

func (m *MapType) valid(other Type) bool {
	mapType, ok := other.(*MapType)
	if !ok {
		return false
	}

	keysMatch := Match(mapType.KeyType, m.KeyType)
	valuesMatch := Match(mapType.ValueType, m.ValueType)
	return keysMatch && valuesMatch
}

func (m *MapType) indexBy(index Type, _ []values.ConstValue) (Type, *diagnostics.Partial) {
	if Assignable(m.KeyType, index) {
		return m.ValueType, nil
	}
	return Invalid, diagnostics.CannotIndex(m, index)
}

func (m *MapType) Item() Type {
	return &TupleType{Types: []Type{m.KeyType, m.ValueType}}
}

type TupleType struct {
	Types []Type
}

func (t *TupleType) String() string {
	var result bytes.Buffer

	result.WriteByte('(')
	for i, ty := range t.Types {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(ty.String())
	}
	result.WriteByte(')')

	return result.String()
}

func (t *TupleType) valid(other Type) bool {
	tuple, ok := other.(*TupleType)
	if !ok {
		return false
	}

	if len(t.Types) != len(tuple.Types) {
		return false
	}

	for i, ty := range t.Types {
		if !Assignable(ty, tuple.Types[i]) && !Assignable(tuple.Types[i], ty) {
			return false
		}
	}

	return true
}

func (a *TupleType) indexBy(t Type, constVals []values.ConstValue) (Type, *diagnostics.Partial) {
	if !Assignable(Int, t) {
		return Invalid, diagnostics.CannotIndex(a, t)
	}

	if len(constVals) == 0 {
		return Invalid, diagnostics.NotConstPartial
	}

	index := constVals[0].(values.IntValue).Value
	return a.Types[index], nil
}

type Function struct {
	Parameters []Type
	ReturnType Type
}

func (fn *Function) String() string {
	var result bytes.Buffer

	result.WriteString("fn(")
	for i, ty := range fn.Parameters {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(ty.String())
	}
	result.WriteByte(')')

	if fn.ReturnType != Void {
		result.WriteString(": ")
		result.WriteString(fn.ReturnType.String())
	}

	return result.String()
}

func (fn *Function) valid(other Type) bool {
	function, ok := other.(*Function)
	if !ok {
		return false
	}

	if len(fn.Parameters) != len(function.Parameters) {
		return false
	}

	for i, ty := range fn.Parameters {
		if !Match(ty, function.Parameters[i]) {
			return false
		}
	}

	return Match(fn.ReturnType, function.ReturnType)
}

type Alias struct {
	Type
}

func (a *Alias) toReal() Type {
	return ToReal(a.Type)
}

type StructField struct {
	Type     Type
	Exported bool
}

type Struct struct {
	Name     string
	ModuleId uint
	Fields   map[string]StructField
}

func (s *Struct) String() string {
	return s.Name
}

func (s *Struct) valid(other Type) bool {
	struc, ok := other.(*Struct)
	if !ok {
		return false
	}

	if s.Name != struc.Name {
		return false
	}

	if len(s.Fields) != len(struc.Fields) {
		return false
	}

	for name, field := range s.Fields {
		ty := field.Type
		if field, ok := struc.Fields[name]; ok {
			if !Match(ty, field.Type) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (s *Struct) member(member string) (Type, *diagnostics.Partial) {
	if field, ok := s.Fields[member]; ok {
		if s.ModuleId != Context.Id() && !field.Exported {
			return Invalid, diagnostics.FieldPrivate(s, member)
		}
		return field.Type, nil
	}
	return Invalid, diagnostics.NoMember(s, member)
}

type Interface struct {
	Name    string
	Methods map[string]*Function
}

func (i *Interface) String() string {
	return i.Name
}

func (i *Interface) valid(other Type) bool {
	for name, ty := range i.Methods {
		member, diag := Member(other, name)
		if diag != nil {
			return false
		}
		if !Assignable(ty, member) {
			return false
		}
	}

	return true
}

func (i *Interface) member(member string) (Type, *diagnostics.Partial) {
	if ty, ok := i.Methods[member]; ok {
		return ty, nil
	}
	return Invalid, diagnostics.NoMember(i, member)
}

type Union struct {
	Types []Type
}

func (u *Union) String() string {
	var result bytes.Buffer

	for i, ty := range u.Types {
		if i != 0 {
			result.WriteString(" | ")
		}
		result.WriteString(ty.String())
	}

	return result.String()
}

func (u *Union) valid(other Type) bool {
	if union, ok := other.(*Union); ok {
		for _, ty := range union.Types {
			if !Assignable(u, ty) {
				return false
			}
		}
		return true
	} else {
		for _, ty := range u.Types {
			if Assignable(ty, other) {
				return true
			}
		}
	}

	return false
}

func MakeUnion(a, b Type) Type {
	types := []Type{}
	if union, ok := a.(*Union); ok {
		types = append(types, union.Types...)
	} else {
		types = append(types, a)
	}

	if union, ok := b.(*Union); ok {
		types = append(types, union.Types...)
	} else {
		types = append(types, b)
	}

	return &Union{
		Types: types,
	}
}

type Module struct {
	Name   string
	Module interface {
		LookupExportType(string) Type
	}
}

func (m *Module) String() string {
	return m.Name
}

func (*Module) valid(Type) bool {
	return false
}

func (m *Module) member(member string) (Type, *diagnostics.Partial) {
	if ty := m.Module.LookupExportType(member); ty != nil {
		return ty, nil
	}
	return Invalid, diagnostics.NoMember(m, member)
}

type Pointer struct {
	Underlying Type
	Mutable    bool
}

func (p *Pointer) String() string {
	if p.Mutable {
		return fmt.Sprintf("*mut %s", p.Underlying.String())
	}
	return fmt.Sprintf("*%s", p.Underlying.String())
}

func (p *Pointer) valid(other Type) bool {
	ptr, ok := other.(*Pointer)
	if !ok {
		return false
	}

	if !Assignable(p.Underlying, ptr.Underlying) {
		return false
	}

	return !p.Mutable || ptr.Mutable
}

func (p *Pointer) member(member string) (Type, *diagnostics.Partial) {
	return Member(p.Underlying, member)
}

type pseudo interface {
	toReal() Type
}

type indexable interface {
	indexBy(Type, []values.ConstValue) (Type, *diagnostics.Partial)
}

type hasMembers interface {
	member(string) (Type, *diagnostics.Partial)
}

type Iterator interface {
	Item() Type
}

// TODO:
// ErrorType
// OptionType
