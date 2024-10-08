package types

import (
	"bytes"
	"fmt"
	"math"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/type_checker/values"
	"tinygo.org/x/go-llvm"
)

type Type interface {
	printer.Printable
	String() string
	ToLlvm(llvm.Context) llvm.Type
	valid(Type) bool
	byteSize() int
}

var Context interface {
	LookupMethod(string, Type, bool) *Function
	Id() uint
}

func Assignable(to, from Type) bool {
	if to == Invalid || from == Invalid {
		return true
	}

	if from == Never {
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
	case Numeric:
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

func ByteSize(ty Type) int {
	return ty.byteSize()
}

func BitSize(ty Type) int {
	return ty.byteSize() * 8
}

func bitsToBytes(bits int) int {
	// This makes sure we don't discard any bits, e.g. 9 bits goes to 2 bytes
	return (bits + 7) / 8
}

type CastKind int

const (
	NoCast CastKind = iota
	IdentityCast
	ImplicitCast
	OperatorCast
	ExplicitCast
)

func Cast(from, to Type, maxKind CastKind) CastKind {
	if castable, ok := from.(castTo); ok {
		kind := castable.castTo(to)
		if kind != NoCast && kind <= maxKind {
			return kind
		}
	}

	if castable, ok := to.(castFrom); ok {
		kind := castable.castFrom(from)
		if kind != NoCast && kind <= maxKind {
			return kind
		}
	}

	if Match(to, from) {
		return IdentityCast
	}

	if Assignable(to, from) {
		return ImplicitCast
	}

	return NoCast
}

type PrimaryType int

const (
	Invalid PrimaryType = iota
	Bool
	String
	RuntimeType
	Never
)

var typeNames = map[PrimaryType]string{
	Invalid:     "<?>",
	Bool:        "bool",
	String:      "string",
	RuntimeType: "Type",
	Never:       "never",
}

var Void = NewUnit("void")

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
		if Assignable(I32, index) {
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

func (pt PrimaryType) ToLlvm(context llvm.Context) llvm.Type {
	switch pt {
	case Invalid:
		panic("Type should not be invalid at this point")
	case Bool:
		return context.Int1Type()
	case String:
		// TODO: Use proper strings, not cstrings
		return llvm.PointerType(context.Int8Type(), 0)
	case RuntimeType:
		panic("TODO: Runtime types")
	case Never:
		panic("TODO: Never types")
	default:
		panic("Unreachable")
	}
}

func (pt PrimaryType) byteSize() int {
	switch pt {
	case Bool:
		return 1
	case Invalid:
		return 0
	case Never:
		return 0
	case RuntimeType:
		panic("TODO: Size of RuntimeType")
	case String:
		// TODO: Make this not a cstring
		return 8
	default:
		panic("Unreachable")
	}
}

type NumKind int

const (
	_ NumKind = iota
	NumUint
	NumInt
	NumFloat
)

var (
	I32 = Int(32)
	F32 = Float(32)
)

func Int(width int) Numeric {
	return Numeric{
		Kind:     NumInt,
		BitWidth: width,
	}
}

func Uint(width int) Numeric {
	return Numeric{
		Kind:     NumUint,
		BitWidth: width,
	}
}

func Float(width int) Numeric {
	return Numeric{
		Kind:     NumFloat,
		BitWidth: width,
	}
}

func minIntWidth(i int64) int {
	if i <= math.MaxInt8 && i >= math.MinInt8 {
		return 8
	}
	if i <= math.MaxInt16 && i >= math.MinInt16 {
		return 16
	}
	if i <= math.MaxInt32 && i >= math.MinInt32 {
		return 32
	}
	return 64
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minUintWidth(u uint64) int {
	if u <= math.MaxUint8 {
		return 8
	}
	if u <= math.MaxUint16 {
		return 16
	}
	if u <= math.MaxUint32 {
		return 32
	}
	return 64
}

func minFloatWidth(f float64) int {
	if f <= MaxFloat16 && f >= -MaxFloat16 {
		return 16
	}
	if f <= math.MaxFloat32 && f >= -math.MaxFloat32 {
		return 32
	}
	return 64
}

func IntType(value int64) Numeric {
	ui := Numeric{
		Kind:         NumInt,
		BitWidth:     maxInt(minIntWidth(value), 32),
		Downcastable: &Downcastable{},
	}
	ui.Downcastable.MinFloatWidth = minFloatWidth(float64(value))
	ui.Downcastable.MinIntWidth = minIntWidth(value)
	ui.Downcastable.MinUintWidth = minUintWidth(uint64(value))

	return ui
}

func UintType(value uint64) Numeric {
	ui := Numeric{
		Kind:         NumInt,
		BitWidth:     maxInt(minUintWidth(value), 32),
		Downcastable: &Downcastable{},
	}
	ui.Downcastable.MinFloatWidth = minFloatWidth(float64(value))
	ui.Downcastable.MinIntWidth = minIntWidth(int64(value))
	ui.Downcastable.MinUintWidth = minUintWidth(value)

	return ui
}

func FloatType(value float64) Numeric {
	uf := Numeric{
		Kind:         NumFloat,
		BitWidth:     64,
		Downcastable: &Downcastable{},
	}
	uf.Downcastable.MinFloatWidth = minFloatWidth(value)
	if float64(int64(value)) == value {
		uf.Downcastable.MinIntWidth = minIntWidth(int64(value))
		uf.Downcastable.MinUintWidth = minUintWidth(uint64(value))
	}

	return uf
}

type Numeric struct {
	Kind         NumKind
	BitWidth     int
	Downcastable *Downcastable
}

type Downcastable struct {
	MinIntWidth,
	MinUintWidth,
	MinFloatWidth int
}

func (n Numeric) String() string {
	if n.Untyped() {
		switch n.Kind {
		case NumInt:
			return "untyped int"
		case NumUint:
			return "untyped uint"
		case NumFloat:
			return "untyped float"
		default:
			panic("unreachable")
		}
	}

	switch n.Kind {
	case NumInt:
		return fmt.Sprintf("i%d", n.BitWidth)
	case NumUint:
		return fmt.Sprintf("u%d", n.BitWidth)
	case NumFloat:
		return fmt.Sprintf("f%d", n.BitWidth)
	default:
		panic("unreachable")
	}
}

func (n Numeric) Untyped() bool {
	return n.Downcastable != nil
}

func (v Numeric) Print(node *printer.Node) {
	node.
		Text(
			"%sVARIABLE_TYPE %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			v.String(),
		)
}

func (v Numeric) valid(other Type) bool {
	variable, ok := other.(Numeric)
	if !ok {
		return false
	}

	if v.Kind == variable.Kind && v.BitWidth == variable.BitWidth {
		return true
	}

	if variable.Untyped() {
		switch v.Kind {
		case NumInt:
			if variable.Downcastable.MinIntWidth != 0 &&
				variable.Downcastable.MinIntWidth <= v.BitWidth {
				return true
			}
		case NumUint:
			if variable.Downcastable.MinUintWidth != 0 &&
				variable.Downcastable.MinUintWidth <= v.BitWidth {
				return true
			}
		case NumFloat:
			if variable.Downcastable.MinFloatWidth != 0 &&
				variable.Downcastable.MinFloatWidth <= v.BitWidth {
				return true
			}
		}
	}

	// TODO: Add a proper error message for this
	return false
}

func (v Numeric) toReal() Type {
	v.Downcastable = nil
	return v
}

const MaxFloat16 = 6.5504e+4

func (n Numeric) MaxValue() float64 {
	switch n.Kind {
	case NumUint:
		return math.Exp2(float64(n.BitWidth)) - 1
	case NumInt:
		return math.Exp2(float64(n.BitWidth-1)) - 1
	case NumFloat:
		switch n.BitWidth {
		case 16:
			return MaxFloat16
		case 32:
			return math.MaxFloat32
		case 64:
			return math.MaxFloat64
		default:
			panic("Unreachable")
		}
	default:
		panic("unreachable")
	}
}

func (from Numeric) castTo(t Type) CastKind {
	to, ok := t.(Numeric)
	if !ok {
		return NoCast
	}

	if Assignable(to, from) {
		kind := IdentityCast
		if from.Untyped() && !to.Untyped() {
			kind = ImplicitCast
		} else if from.Untyped() && to.Untyped() && from.Kind != to.Kind {
			kind = ImplicitCast
		}

		return kind
	}

	if to.Kind >= from.Kind && to.MaxValue() >= from.MaxValue() {
		return OperatorCast
	}

	return ExplicitCast
}

func (to Numeric) castFrom(from Type) CastKind {
	if from == Bool {
		if to.Kind == NumInt || to.Kind == NumUint {
			return ExplicitCast
		}
	}
	return NoCast
}

func (v Numeric) GetEnumValue(
	prevValues []values.ConstValue,
	_name string,
) (values.ConstValue, *diagnostics.Partial) {
	if v.Kind == NumFloat {
		return nil, diagnostics.CannotEnumPartial(v)
	}

	if len(prevValues) == 0 {
		return values.IntValue{Value: 0}, nil
	}
	last := prevValues[len(prevValues)-1].(values.IntValue)
	return values.IntValue{Value: last.Value + 1}, nil
}

func (n Numeric) ToLlvm(context llvm.Context) llvm.Type {
	switch n.Kind {
	case NumFloat:
		switch n.BitWidth {
		case 16:
			panic("TODO: F16 types")
		case 32:
			return context.FloatType()
		case 64:
			return context.DoubleType()
		default:
			panic("Invalid float bit-width")
		}
	case NumInt:
		switch n.BitWidth {
		case 8:
			return context.Int8Type()
		case 16:
			return context.Int16Type()
		case 32:
			return context.Int32Type()
		case 64:
			return context.Int64Type()
		default:
			panic("Invalid int bit-width")
		}
	case NumUint:
		switch n.BitWidth {
		case 8:
			return context.Int8Type()
		case 16:
			return context.Int16Type()
		case 32:
			return context.Int32Type()
		case 64:
			return context.Int64Type()
		default:
			panic("Invalid uint bit-width")
		}
	default:
		panic("unreachable")
	}
}

func (n Numeric) byteSize() int {
	return bitsToBytes(n.BitWidth)
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
	if Assignable(I32, index) {
		return l.ElemType, nil
	}
	return Invalid, diagnostics.CannotIndex(l, index)
}

func (l *ListType) Item() Type {
	return l.ElemType
}

func (*ListType) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (l *ListType) byteSize() int {
	// len + cap + ptr
	return 24
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
	if !Assignable(I32, index) {
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

func (*ArrayType) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (a *ArrayType) byteSize() int {
	return a.Length * a.ElemType.byteSize()
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

func (*MapType) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (m *MapType) byteSize() int {
	// len + cap + ptr
	return 24
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
	if !Assignable(I32, t) {
		return Invalid, diagnostics.CannotIndex(a, t)
	}

	if len(constVals) == 0 {
		return Invalid, diagnostics.NotConstPartial
	}

	index := constVals[0].(values.IntValue).Value
	return a.Types[index], nil
}

func (*TupleType) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (t *TupleType) byteSize() int {
	size := 0
	for _, ty := range t.Types {
		size += ty.byteSize()
	}
	return size
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

func (*Function) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (*Function) byteSize() int {
	// Just a pointer
	return 8
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
	Name       string
	ModuleId   uint
	Fields     map[string]StructField
	FieldOrder []string
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

func (s *Struct) ToLlvm(context llvm.Context) llvm.Type {
	types := make([]llvm.Type, 0, len(s.FieldOrder))
	for _, name := range s.FieldOrder {
		types = append(types, s.Fields[name].Type.ToLlvm(context))
	}
	return llvm.StructType(types, false)
}

func (s *Struct) byteSize() int {
	size := 0
	for _, field := range s.Fields {
		size += field.Type.byteSize()
	}
	return size
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
	if !Assignable(I32, t) {
		return Invalid, diagnostics.CannotIndex(a, t)
	}

	if len(constVals) == 0 {
		return Invalid, diagnostics.NotConstPartial
	}

	index := constVals[0].(values.IntValue).Value
	return a.Types[index], nil
}

func (*TupleStruct) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (t *TupleStruct) byteSize() int {
	size := 0
	for _, ty := range t.Types {
		size += ty.byteSize()
	}
	return size
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

func (*Interface) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (*Interface) byteSize() int {
	panic("TODO")
}

// TODO: Untagged unions
type Union struct {
	Name    string
	Id      int
	Members map[string]Type
}

var unionId = 0

func NewUnion(name string) *Union {
	id := unionId
	unionId++
	return &Union{
		Name:    name,
		Id:      id,
		Members: map[string]Type{},
	}
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
		return u.Id == union.Id
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

func (to *Union) castFrom(from Type) CastKind {
	if Assignable(to, from) {
		return ImplicitCast
	}
	return NoCast
}

func (*Union) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (u *Union) byteSize() int {
	if len(u.Members) > 255 {
		panic("TODO: More than 1-bit tags")
	}
	size := 0
	for _, member := range u.Members {
		size = maxInt(size, member.byteSize())
	}
	// Add one for the tag size
	return size + 1
}

type UnionVariant struct {
	Union *Union
	Name  string
	Id    int
	Type
}

var variantId = 0

func NewVariant(union *Union, name string, ty Type) *UnionVariant {
	id := variantId
	variantId++
	return &UnionVariant{
		Union: union,
		Name:  name,
		Id:    id,
		Type:  ty,
	}
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
		return expl.Id == v.Id
	}
	return Assignable(v.Type, other)
}

func (v *UnionVariant) toReal() Type {
	return v.Union
}

func (v *UnionVariant) unwrap() Type {
	return Unwrap(v.Type)
}

func (to *UnionVariant) castFrom(from Type) CastKind {
	if Assignable(to, from) {
		return ImplicitCast
	}
	return NoCast
}

func (*UnionVariant) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
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

func (*InlineUnion) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (u *InlineUnion) byteSize() int {
	if len(u.Types) > 255 {
		panic("TODO: More than 1-bit tags")
	}
	size := 0
	for _, member := range u.Types {
		size = maxInt(size, member.byteSize())
	}
	// Add one for the tag size
	return size + 1
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

func (*Module) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (m *Module) byteSize() int {
	panic("TODO")
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

func (p *Pointer) ToLlvm(context llvm.Context) llvm.Type {
	return llvm.PointerType(p.Underlying.ToLlvm(context), 0)
}

func (p *Pointer) byteSize() int {
	// TODO: target-specific pointer size
	return 8
}

type Explicit struct {
	Name string
	Id   int
	Type
}

var explicitId = 0

func NewExplicit(name string, ty Type) *Explicit {
	id := explicitId
	explicitId++
	return &Explicit{
		Name: name,
		Id:   id,
		Type: ty,
	}
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
		return expl.Id == e.Id
	}
	return Assignable(e.Type, other)
}

func (e *Explicit) unwrap() Type {
	return Unwrap(e.Type)
}

func (from *Explicit) castTo(to Type) CastKind {
	if Assignable(to, from.Type) {
		return ExplicitCast
	}

	return NoCast
}

func (to *Explicit) castFrom(from Type) CastKind {
	if Assignable(to, from) {
		return ImplicitCast
	}
	return NoCast
}

type UnitStruct struct {
	Name string
	Id   int
}

var unitId = 0

func NewUnit(name string) *UnitStruct {
	id := unitId
	unitId++
	return &UnitStruct{name, id}
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
	if unit, ok := other.(*UnitStruct); ok {
		return unit.Id == u.Id
	}
	return false
}

func (*UnitStruct) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (*UnitStruct) byteSize() int {
	return 0
}

type Tag struct {
	Name  string
	Id    int
	Types []Type
}

var tagId = 0

func NewTag(name string) *Tag {
	id := tagId
	tagId++
	return &Tag{
		Name:  name,
		Id:    id,
		Types: []Type{},
	}
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
	if tag, ok := other.(*Tag); ok {
		return tag.Id == t.Id
	}
	for _, ty := range t.Types {
		if Assignable(ty, other) {
			return true
		}
	}
	return false
}

func (*Tag) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (to *Tag) castFrom(from Type) CastKind {
	if Assignable(to, from) {
		return IdentityCast
	}
	return NoCast
}

func (*Tag) byteSize() int {
	panic("TODO")
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

func (*Result) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (r *Result) byteSize() int {
	return maxInt(ErrorTag.byteSize(), r.OkType.byteSize()) + 1
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

func (*Option) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (o *Option) byteSize() int {
	// void is zero-size so it's always the size of the some type
	return o.byteSize() + 1
}

type Enum struct {
	Name       string
	Id         int
	Underlying Type
	Members    map[string]values.ConstValue
}

var enumId = 0

func NewEnum(name string, underlying Type) *Enum {
	id := enumId
	enumId++
	return &Enum{
		Name:       name,
		Id:         id,
		Underlying: underlying,
		Members:    map[string]values.ConstValue{},
	}
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
			func(n *printer.Node) { n.Node(kv.Value) },
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			kv.Key,
		)
	}
}

func (e *Enum) valid(other Type) bool {
	if enum, ok := other.(*Enum); ok && enum.Id == e.Id {
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

func (from *Enum) castTo(to Type) CastKind {
	if Assignable(to, from.Underlying) {
		return ExplicitCast
	}
	return NoCast
}

func (to *Enum) castFrom(from Type) CastKind {
	if Assignable(to.Underlying, from) {
		return ExplicitCast
	}
	return NoCast
}

func (*Enum) ToLlvm(llvm.Context) llvm.Type {
	panic("TODO")
}

func (*Enum) byteSize() int {
	panic("TODO")
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

type castTo interface {
	castTo(Type) CastKind
}

type castFrom interface {
	castFrom(Type) CastKind
}
