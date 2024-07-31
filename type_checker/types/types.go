package types

import (
	"bytes"
	"fmt"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

type Type interface {
	printer.Printable
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
	if indexable, ok := Unwrap(left).(indexable); ok {
		return indexable.indexBy(index, constVals)
	}
	return Invalid, diagnostics.CannotIndex(left, index)
}

func Member(left Type, member string, constVal ...values.ConstValue) (Type, *diagnostics.Partial) {
	if left == Invalid {
		return Invalid, nil
	}

	if left == RuntimeType && len(constVal) > 0 {
		ty := constVal[0].(values.TypeValue).Type.(Type)
		if method := Context.LookupMethod(member, ty, true); method != nil {
			return method, nil
		}
		if sm, ok := ty.(staticMember); ok {
			static, diag := sm.staticMember(member)
			if diag == nil {
				return static, nil
			} else {
				return Invalid, diag
			}
		}
	} else if method := Context.LookupMethod(member, left, false); method != nil {
		return method, nil
	}

	if hasMember, ok := Unwrap(left).(hasMembers); ok {
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

func Unwrap(ty Type) Type {
	if container, ok := ty.(container); ok {
		return container.unwrap()
	}
	return ty
}

type PrimaryType int

const (
	Invalid PrimaryType = iota
	Bool
	String
	RuntimeType
)

var typeNames = map[PrimaryType]string{
	Invalid:     "<?>",
	Bool:        "bool",
	String:      "string",
	RuntimeType: "Type",
}

var Void = &UnitStruct{"void"}

func (pt PrimaryType) String() string {
	return typeNames[pt]
}

func (pt PrimaryType) Print(node *printer.Node) {
	node.Text(
		"%sPRIMARY_TYPE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		pt.String(),
	)
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

func (pt PrimaryType) GetEnumValue(
	_ []values.ConstValue,
	name string,
) (values.ConstValue, *diagnostics.Partial) {
	if pt != String {
		return nil, diagnostics.CannotEnumPartial(pt)
	}

	return values.StringValue{Value: name}, nil
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

func (v VariableType) Print(node *printer.Node) {
	node.
		Text(
			"%sVARIABLE_TYPE %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			v.String(),
		).
		TextIf(
			v.Downcastable,
			"%sdowncastable",
			node.Colour(colour.Attribute),
		)
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

func (v VariableType) GetEnumValue(
	prevValues []values.ConstValue,
	_name string,
) (values.ConstValue, *diagnostics.Partial) {
	if v.Kind != VT_Int {
		return nil, diagnostics.CannotEnumPartial(v)
	}

	if len(prevValues) == 0 {
		return values.IntValue{Value: 0}, nil
	}
	last := prevValues[len(prevValues)-1].(values.IntValue)
	return values.IntValue{Value: last.Value + 1}, nil
}

type ListType struct {
	ElemType Type
}

func (l *ListType) String() string {
	return l.ElemType.String() + "[]"
}

func (l *ListType) Print(node *printer.Node) {
	node.
		Text("%sLIST_TYPE", node.Colour(colour.NodeName)).
		Node(l.ElemType)
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

func (a *ArrayType) Print(node *printer.Node) {
	node.
		Text(
			"%sARRAY_TYPE %s%d",
			node.Colour(colour.NodeName),
			node.Colour(colour.Literal),
			a.Length,
		).
		TextIf(
			a.CanInfer,
			" %scan_infer",
			node.Colour(colour.Attribute),
		).
		Node(a.ElemType)
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

func (m *MapType) Print(node *printer.Node) {
	node.
		Text("%sMAP_TYPE", node.Colour(colour.NodeName)).
		Node(m.KeyType).
		Node(m.ValueType)
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

func (t *TupleType) Print(node *printer.Node) {
	node.Text("%sTUPLE_TYPE", node.Colour(colour.NodeName))

	printer.Nodes(node, t.Types)
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

func (fn *Function) Print(node *printer.Node) {
	node.
		Text("%sFUNCTION_TYPE", node.Colour(colour.NodeName)).
		Node(fn.ReturnType)

	printer.Nodes(node, fn.Parameters)
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

func (a *Alias) unwrap() Type {
	return Unwrap(a.Type)
}

type StructField struct {
	Name     string
	Type     Type
	Exported bool
}

func (s StructField) Print(node *printer.Node) {
	node.
		Text(
			"%sSTRUCT_FIELD %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			s.Name,
		).
		TextIf(
			s.Exported,
			" %spub",
			node.Colour(colour.Attribute),
		).
		Node(s.Type)
}

type Struct struct {
	Name     string
	ModuleId uint
	Fields   map[string]StructField
}

func (s *Struct) String() string {
	return s.Name
}

func (s *Struct) Print(node *printer.Node) {
	node.Text(
		"%sSTRUCT_TYPE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		s.Name,
	)

	printer.Map(node, s.Fields)
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

type TupleStruct struct {
	Name  string
	Types []Type
}

func (t *TupleStruct) String() string {
	return t.Name
}

func (t *TupleStruct) Print(node *printer.Node) {
	node.Text(
		"%sTUPLE_STRUCT_TYPE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		t.Name,
	)

	printer.Nodes(node, t.Types)
}

func (t *TupleStruct) valid(other Type) bool {
	tuple, ok := other.(*TupleStruct)
	if !ok {
		return false
	}

	if t.Name != tuple.Name {
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

func (a *TupleStruct) indexBy(t Type, constVals []values.ConstValue) (Type, *diagnostics.Partial) {
	if !Assignable(Int, t) {
		return Invalid, diagnostics.CannotIndex(a, t)
	}

	if len(constVals) == 0 {
		return Invalid, diagnostics.NotConstPartial
	}

	index := constVals[0].(values.IntValue).Value
	return a.Types[index], nil
}

type Interface struct {
	Name    string
	Methods map[string]*Function
}

func (i *Interface) String() string {
	return i.Name
}

func (i *Interface) Print(node *printer.Node) {
	node.Text(
		"%sINTERFACE_TYPE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		i.Name,
	)

	for name, ty := range i.Methods {
		node.FakeNode(
			"%sINTERFACE_MEMBER %s%s",
			func(n *printer.Node) { n.Node(ty) },
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			name,
		)
	}
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
	Name    string
	Members map[string]Type
}

func (u *Union) String() string {
	return u.Name
}

func (u *Union) Print(node *printer.Node) {
	node.Text(
		"%sUNION_TYPE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		u.Name,
	)

	printer.Map(node, u.Members)
}

func (u *Union) valid(other Type) bool {
	if union, ok := other.(*Union); ok {
		// FIXME: compare more than just the name
		return u.Name == union.Name
	} else {
		canAssign := false
		for _, ty := range u.Members {
			if expl, ok := ty.(*Explicit); (ok && Assignable(expl.Type, other)) || Assignable(ty, other) {
				// TODO: Prevent ambiguity of untyped numbers when assigned to a union
				if canAssign {
					// TODO: Make a proper error message for this
					return false
				}
				canAssign = true
			}
		}
		return canAssign
	}
}

func (u *Union) member(member string) (Type, *diagnostics.Partial) {
	if ty, ok := u.Members[member]; ok {
		return ty, nil
	}
	return nil, diagnostics.NoVariant(u.Name, member)
}

func (u *Union) staticMember(member string) (Type, *diagnostics.Partial) {
	if _, ok := u.Members[member]; ok {
		return RuntimeType, nil
	}
	return nil, diagnostics.NoVariant(u.Name, member)
}

func (u *Union) StaticMemberValue(member string) values.ConstValue {
	if ty, ok := u.Members[member]; ok {
		return values.TypeValue{Type: ty}
	}
	return nil
}

type UnionVariant struct {
	Union *Union
	Name  string
	Type
}

func (v *UnionVariant) String() string {
	return fmt.Sprintf("%s.%s", v.Union.Name, v.Name)
}

func (v *UnionVariant) Print(node *printer.Node) {
	node.
		Text(
			"%sUNION_VARIANT %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			v.Name,
		).
		Node(v.Type)
}

func (v *UnionVariant) valid(other Type) bool {
	if expl, ok := other.(*UnionVariant); ok {
		return expl.Name == v.Name && Assignable(expl.Type, v.Type)
	}
	return Assignable(v.Type, other)
}

func (v *UnionVariant) toReal() Type {
	return v.Union
}

func (v *UnionVariant) unwrap() Type {
	return Unwrap(v.Type)
}

type InlineUnion struct {
	Types []Type
}

func (u *InlineUnion) String() string {
	var result bytes.Buffer

	for i, ty := range u.Types {
		if i != 0 {
			result.WriteString(" | ")
		}
		result.WriteString(ty.String())
	}

	return result.String()
}

func (u *InlineUnion) Print(node *printer.Node) {
	node.Text("%sINLINE_UNION_TYPE", node.Colour(colour.NodeName))

	printer.Nodes(node, u.Types)
}

func (u *InlineUnion) valid(other Type) bool {
	if union, ok := other.(*InlineUnion); ok {
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
	if union, ok := a.(*InlineUnion); ok {
		types = append(types, union.Types...)
	} else {
		types = append(types, a)
	}

	if union, ok := b.(*InlineUnion); ok {
		types = append(types, union.Types...)
	} else {
		types = append(types, b)
	}

	return &InlineUnion{
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

func (m *Module) Print(node *printer.Node) {
	node.Text(
		"%sMODULE_TYPE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		m.Name,
	)
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

func (p *Pointer) Print(node *printer.Node) {
	node.
		Text("%sPOINTER_TYPE", node.Colour(colour.NodeName)).
		TextIf(
			p.Mutable,
			" %smut",
			node.Colour(colour.Attribute),
		).
		Node(p.Underlying)
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

type Explicit struct {
	Name string
	Type
}

func (e *Explicit) String() string {
	return e.Name
}

func (e *Explicit) Print(node *printer.Node) {
	node.
		Text(
			"%sEXPLICIT_TYPE %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			e.Name,
		).
		Node(e.Type)
}

func (e *Explicit) valid(other Type) bool {
	if expl, ok := other.(*Explicit); ok {
		return expl.Name == e.Name
	}
	return Assignable(e.Type, other)
}

func (e *Explicit) unwrap() Type {
	return Unwrap(e.Type)
}

type UnitStruct struct {
	Name string
}

func (u *UnitStruct) String() string {
	return u.Name
}

func (u *UnitStruct) Print(node *printer.Node) {
	node.
		Text(
			"%sUNIT_STRUCT %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			u.Name,
		)
}

func (u *UnitStruct) valid(other Type) bool {
	// FIXME: compare more than just the name
	if unit, ok := other.(*UnitStruct); ok {
		return unit.Name == u.Name
	}
	return false
}

type Tag struct {
	Name  string
	Types []Type
}

func (t *Tag) String() string {
	return t.Name
}

func (t *Tag) Print(node *printer.Node) {
	node.Text(
		"%sTAG_TYPE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		t.Name,
	)

	printer.Nodes(node, t.Types)
}

func (t *Tag) valid(other Type) bool {
	// FIXME: compare more than just the name
	if tag, ok := other.(*Tag); ok {
		return tag.Name == t.Name
	}
	for _, ty := range t.Types {
		if Assignable(ty, other) {
			return true
		}
	}
	return false
}

var ErrorTag = Tag{
	Name:  "Error",
	Types: []Type{},
}

type Result struct {
	OkType Type
}

func (r *Result) String() string {
	return "!" + r.OkType.String()
}

func (r *Result) Print(node *printer.Node) {
	node.
		Text(
			"%sRESULT_TYPE",
			node.Colour(colour.NodeName),
		).
		Node(r.OkType)
}

func (r *Result) valid(other Type) bool {
	if result, ok := other.(*Result); ok && Assignable(r.OkType, result.OkType) {
		return true
	}
	if Assignable(&ErrorTag, other) {
		return true
	}
	return Assignable(r.OkType, other)
}

type Option struct {
	SomeType Type
}

func (r *Option) String() string {
	return "?" + r.SomeType.String()
}

func (r *Option) Print(node *printer.Node) {
	node.
		Text(
			"%sOPTION_TYPE",
			node.Colour(colour.NodeName),
		).
		Node(r.SomeType)
}

func (r *Option) valid(other Type) bool {
	if option, ok := other.(*Option); ok && Assignable(r.SomeType, option.SomeType) {
		return true
	}
	if other == Void {
		return true
	}
	return Assignable(r.SomeType, other)
}

type Enum struct {
	Name       string
	Underlying Type
	Members    map[string]values.ConstValue
}

func (e *Enum) String() string {
	return e.Name
}

func (e *Enum) Print(node *printer.Node) {
	node.
		Text(
			"%sENUM_TYPE %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			e.Name,
		).
		Node(e.Underlying)

		for _, kv := range printer.SortMap(e.Members) {
			node.FakeNode(
				"%sENUM_MEMBER %s%s",
				func(n *printer.Node) {n.Node(kv.Value)},
				node.Colour(colour.NodeName),
				node.Colour(colour.Name),
				kv.Key,
			)
		}
}

func (e *Enum) valid(other Type) bool {
	if enum, ok := other.(*Enum); ok && enum.Name == e.Name {
		return true
	}
	return false
}

func (e *Enum) staticMember(member string) (Type, *diagnostics.Partial) {
	if _, ok := e.Members[member]; ok {
		return e, nil
	}
	return nil, diagnostics.NoEnumMember(e.Name, member)
}

func (e *Enum) StaticMemberValue(member string) values.ConstValue {
	return e.Members[member]
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

type staticMember interface {
	staticMember(string) (Type, *diagnostics.Partial)
}

type container interface {
	unwrap() Type
}

type Iterator interface {
	Item() Type
}

type HasEnumValue interface {
	GetEnumValue([]values.ConstValue, string) (values.ConstValue, *diagnostics.Partial)
}
