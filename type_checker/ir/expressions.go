package ir

import (
	"bytes"
	"fmt"
	"math"

	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

type expression struct{}

func (expression) irExpr() {}

type IntegerLiteral struct {
	expression
	Value int64
}

func (i *IntegerLiteral) String() string {
	return fmt.Sprint(i.Value)
}

func (IntegerLiteral) Type() types.Type {
	return types.UntypedInt
}

func (IntegerLiteral) IsConst() bool {
	return true
}

func (i *IntegerLiteral) ConstValue() values.ConstValue {
	return values.IntValue{
		Value: i.Value,
	}
}

type FloatLiteral struct {
	expression
	Value float64
}

func (f *FloatLiteral) String() string {
	return fmt.Sprint(f.Value)
}

func (f *FloatLiteral) Type() types.Type {
	uf := types.UntypedFloat
	if float64(int64(f.Value)) == f.Value {
		uf.Downcastable = true
	}
	return uf
}

func (FloatLiteral) IsConst() bool {
	return true
}

func (f *FloatLiteral) ConstValue() values.ConstValue {
	return values.FloatValue{
		Value: f.Value,
	}
}

type BooleanLiteral struct {
	expression
	Value bool
}

func (b *BooleanLiteral) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (BooleanLiteral) Type() types.Type {
	return types.Bool
}

func (BooleanLiteral) IsConst() bool {
	return true
}

func (b *BooleanLiteral) ConstValue() values.ConstValue {
	return values.BoolValue{
		Value: b.Value,
	}
}

type StringLiteral struct {
	expression
	Value string
}

func (b *StringLiteral) String() string {
	return "\"" + b.Value + "\""
}

func (StringLiteral) Type() types.Type {
	return types.String
}

func (StringLiteral) IsConst() bool {
	return true
}

func (s *StringLiteral) ConstValue() values.ConstValue {
	return values.StringValue{
		Value: s.Value,
	}
}

type VariableExpression struct {
	expression
	Symbol symbols.Variable
}

func (v *VariableExpression) String() string {
	return v.Symbol.Name
}

func (v *VariableExpression) Type() types.Type {
	return v.Symbol.Type
}

func (v *VariableExpression) IsConst() bool {
	return v.Symbol.ConstValue != nil
}

func (v *VariableExpression) ConstValue() values.ConstValue {
	return v.Symbol.ConstValue
}

type BinaryOperator int

const (
	_ BinaryOperator = iota
	LogicalAnd
	LogicalOr
	Less
	LessEq
	Greater
	GreaterEq
	Equal
	NotEqual
	LeftShift
	RightShift
	BitwiseOr
	BitwiseAnd
	AddInt
	AddFloat
	Concat
	SubtractInt
	SubtractFloat
	MultiplyInt
	MultiplyFloat
	Divide
	ModuloInt
	ModuloFloat
	PowerInt
	PowerFloat
	UntypedBit = 1 << 8
)

func (b BinaryOperator) String() string {
	b = b & ^UntypedBit

	switch b {
	case LogicalAnd:
		return "&&"
	case LogicalOr:
		return "||"

	case Less:
		return "<"

	case LessEq:
		return "<="

	case Greater:
		return "<"

	case GreaterEq:
		return ">="

	case Equal:
		return "=="

	case NotEqual:
		return "!="

	case LeftShift:
		return "<<"

	case RightShift:
		return ">>"

	case BitwiseOr:
		return "|"

	case BitwiseAnd:
		return "&"

	case AddInt:
		fallthrough
	case AddFloat:
		fallthrough
	case Concat:
		return "+"

	case SubtractInt:
		fallthrough
	case SubtractFloat:
		return "-"

	case MultiplyInt:
		fallthrough
	case MultiplyFloat:
		return "*"

	case Divide:
		return "/"

	case ModuloInt:
		fallthrough
	case ModuloFloat:
		return "%"

	case PowerInt:
		fallthrough
	case PowerFloat:
		return "**"

	default:
		return "<?>"
	}
}

func (b BinaryOperator) Type() types.Type {
	untyped := b&UntypedBit != 0
	b = b & ^UntypedBit
	var ty types.Type = types.Invalid

	switch b {
	case LogicalAnd:
		ty = types.Bool

	case LogicalOr:
		ty = types.Bool

	case Less:
		ty = types.Bool

	case LessEq:
		ty = types.Bool

	case Greater:
		ty = types.Bool

	case GreaterEq:
		ty = types.Bool

	case Equal:
		ty = types.Bool

	case NotEqual:
		ty = types.Bool

	case LeftShift:
		ty = types.Int

	case RightShift:
		ty = types.Int

	case BitwiseOr:
		ty = types.Int

	case BitwiseAnd:
		ty = types.Int

	case AddInt:
		ty = types.Int

	case AddFloat:
		ty = types.Float

	case Concat:
		ty = types.String

	case SubtractInt:
		ty = types.Int

	case SubtractFloat:
		ty = types.Float

	case MultiplyInt:
		ty = types.Int

	case MultiplyFloat:
		ty = types.Float

	case Divide:
		ty = types.Float

	case ModuloInt:
		ty = types.Int

	case ModuloFloat:
		ty = types.Float

	case PowerInt:
		ty = types.Int

	case PowerFloat:
		ty = types.Float
	}

	if untyped {
		if variable, ok := ty.(types.VariableType); ok {
			variable.Untyped = true
			return variable
		}
	}

	return ty
}

type BinaryExpression struct {
	expression
	Left     Expression
	Operator BinaryOperator
	Right    Expression
}

func (b *BinaryExpression) String() string {
	var result bytes.Buffer

	result.WriteString(b.Left.String())
	result.WriteByte(' ')
	result.WriteString(b.Operator.String())
	result.WriteByte(' ')
	result.WriteString(b.Right.String())

	return result.String()
}

func (b *BinaryExpression) Type() types.Type {
	return b.Operator.Type()
}

func (b *BinaryExpression) IsConst() bool {
	return b.Left.IsConst() && b.Right.IsConst()
}

func (b *BinaryExpression) ConstValue() values.ConstValue {
	if !b.IsConst() {
		return nil
	}

	switch b.Operator & ^UntypedBit {
	case LogicalAnd:
		left := b.Left.ConstValue().(values.BoolValue)
		right := b.Right.ConstValue().(values.BoolValue)
		return values.BoolValue{
			Value: left.Value && right.Value,
		}

	case LogicalOr:
		left := b.Left.ConstValue().(values.BoolValue)
		right := b.Right.ConstValue().(values.BoolValue)
		return values.BoolValue{
			Value: left.Value || right.Value,
		}

	case Less:
		left := values.NumericValue(b.Left.ConstValue())
		right := values.NumericValue(b.Right.ConstValue())
		return values.BoolValue{
			Value: left < right,
		}

	case LessEq:
		left := values.NumericValue(b.Left.ConstValue())
		right := values.NumericValue(b.Right.ConstValue())
		return values.BoolValue{
			Value: left <= right,
		}

	case Greater:
		left := values.NumericValue(b.Left.ConstValue())
		right := values.NumericValue(b.Right.ConstValue())
		return values.BoolValue{
			Value: left > right,
		}

	case GreaterEq:
		left := values.NumericValue(b.Left.ConstValue())
		right := values.NumericValue(b.Right.ConstValue())
		return values.BoolValue{
			Value: left >= right,
		}

	case Equal:
		return values.BoolValue{
			Value: b.Left.ConstValue() == b.Right.ConstValue(),
		}

	case NotEqual:
		return values.BoolValue{
			Value: b.Left.ConstValue() != b.Right.ConstValue(),
		}

	case LeftShift:
		left := b.Left.ConstValue().(values.IntValue)
		right := b.Right.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: left.Value << right.Value,
		}

	case RightShift:
		left := b.Left.ConstValue().(values.IntValue)
		right := b.Right.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: left.Value >> right.Value,
		}

	case BitwiseOr:
		left := b.Left.ConstValue().(values.IntValue)
		right := b.Right.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: left.Value | right.Value,
		}

	case BitwiseAnd:
		left := b.Left.ConstValue().(values.IntValue)
		right := b.Right.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: left.Value & right.Value,
		}

	case AddInt:
		left := b.Left.ConstValue().(values.IntValue)
		right := b.Right.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: left.Value + right.Value,
		}
	case AddFloat:
		left := b.Left.ConstValue().(values.FloatValue)
		right := b.Right.ConstValue().(values.FloatValue)
		return values.FloatValue{
			Value: left.Value + right.Value,
		}
	case Concat:
		left := b.Left.ConstValue().(values.StringValue)
		right := b.Right.ConstValue().(values.StringValue)
		return values.StringValue{
			Value: left.Value + right.Value,
		}

	case SubtractInt:
		left := b.Left.ConstValue().(values.IntValue)
		right := b.Right.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: left.Value - right.Value,
		}
	case SubtractFloat:
		left := b.Left.ConstValue().(values.FloatValue)
		right := b.Right.ConstValue().(values.FloatValue)
		return values.FloatValue{
			Value: left.Value - right.Value,
		}

	case MultiplyInt:
		left := b.Left.ConstValue().(values.IntValue)
		right := b.Right.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: left.Value * right.Value,
		}
	case MultiplyFloat:
		left := b.Left.ConstValue().(values.FloatValue)
		right := b.Right.ConstValue().(values.FloatValue)
		return values.FloatValue{
			Value: left.Value * right.Value,
		}

	case Divide:
		left := values.NumericValue(b.Left.ConstValue())
		right := values.NumericValue(b.Right.ConstValue())
		return values.FloatValue{
			Value: left / right,
		}

	case ModuloInt:
		left := b.Left.ConstValue().(values.IntValue)
		right := b.Right.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: left.Value % right.Value,
		}
	case ModuloFloat:
		left := b.Left.ConstValue().(values.FloatValue)
		right := b.Right.ConstValue().(values.FloatValue)

		return values.FloatValue{
			Value: math.Mod(left.Value, right.Value),
		}

	case PowerInt:
		left := b.Left.ConstValue().(values.IntValue)
		right := b.Right.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: int64(math.Pow(float64(left.Value), float64(right.Value))),
		}
	case PowerFloat:
		left := b.Left.ConstValue().(values.FloatValue)
		right := b.Right.ConstValue().(values.FloatValue)

		return values.FloatValue{
			Value: math.Pow(left.Value, right.Value),
		}

	default:
		panic("unreachable")
	}
}

type UnaryOperator int

const (
	_ UnaryOperator = iota
	NegateInt
	NegateFloat
	Identity
	LogicalNot
	BitwiseNot
)

const (
	postfix UnaryOperator = iota | (1 << 16)
	IncrecementInt
	IncrementFloat
	DecrecementInt
	DecrementFloat
	PropagateError
	CrashError
	// TODO: Deref and ref
)

func (u UnaryOperator) String() string {
	u = u & ^UntypedBit

	switch u {
	case NegateInt:
		fallthrough
	case NegateFloat:
		return "-"
	case Identity:
		return "+"
	case LogicalNot:
		return "!"
	case BitwiseNot:
		return "~"
	case IncrecementInt:
		fallthrough
	case IncrementFloat:
		return "++"
	case DecrecementInt:
		fallthrough
	case DecrementFloat:
		return "--"
	case PropagateError:
		return "?"
	case CrashError:
		return "!"
	default:
		return "<?>"
	}
}

func (b UnaryOperator) Type() types.Type {
	untyped := b&UntypedBit != 0
	b = b & ^UntypedBit
	var ty types.Type = types.Invalid

	switch b {
	case NegateInt:
		ty = types.Int
	case NegateFloat:
		ty = types.Float
	case Identity:
		ty = types.Int
	case LogicalNot:
		ty = types.Bool
	case BitwiseNot:
		ty = types.Int
	case IncrecementInt:
		ty = types.Int
	case IncrementFloat:
		ty = types.Float
	case DecrecementInt:
		ty = types.Int
	case DecrementFloat:
		ty = types.Float
	case PropagateError:
		panic("TODO: Type for PropagateError unary operator")
	case CrashError:
		panic("TODO: Type for CrashError unary operator")
	}

	if untyped {
		if variable, ok := ty.(types.VariableType); ok {
			variable.Untyped = true
			return variable
		}
	}

	return ty
}

type UnaryExpression struct {
	expression
	Operator UnaryOperator
	Operand  Expression
}

func (u *UnaryExpression) String() string {
	var result bytes.Buffer
	isPost := (u.Operator & postfix) != 0

	if isPost {
		result.WriteString(u.Operand.String())
		result.WriteString(u.Operator.String())
	} else {
		result.WriteString(u.Operator.String())
		result.WriteString(u.Operand.String())
	}

	return result.String()
}

func (u *UnaryExpression) Type() types.Type {
	return u.Operator.Type()
}

func (u *UnaryExpression) IsConst() bool {
	return u.Operand.IsConst()
}

func (u *UnaryExpression) ConstValue() values.ConstValue {
	if !u.IsConst() {
		return nil
	}

	switch u.Operator & ^UntypedBit {
	case NegateInt:
		value := u.Operand.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: -value.Value,
		}
	case NegateFloat:
		value := u.Operand.ConstValue().(values.FloatValue)
		return values.FloatValue{
			Value: -value.Value,
		}
	case Identity:
		return u.Operand.ConstValue()
	case LogicalNot:
		value := u.Operand.ConstValue().(values.BoolValue)
		return values.BoolValue{
			Value: !value.Value,
		}
	case BitwiseNot:
		value := u.Operand.ConstValue().(values.IntValue)
		return values.IntValue{
			Value: ^value.Value,
		}
	case IncrecementInt:
		fallthrough
	case IncrementFloat:
		fallthrough
	case DecrecementInt:
		fallthrough
	case DecrementFloat:
		fallthrough
	case PropagateError:
		fallthrough
	case CrashError:
		return nil
	default:
		return nil
	}
}

type Conversion struct {
	expression
	Expression Expression
	To         types.Type
}

func (c *Conversion) String() string {
	var result bytes.Buffer

	result.WriteString(c.Expression.String())
	result.WriteString(" -> ")
	result.WriteString(c.To.String())

	return result.String()
}

func (c *Conversion) Type() types.Type {
	return c.To
}

func (c *Conversion) IsConst() bool {
	return c.Expression.IsConst()
}

func (c *Conversion) ConstValue() values.ConstValue {
	if !c.IsConst() {
		return nil
	}

	switch {
	case types.Match(c.To, types.Float):
		return values.FloatValue{
			Value: values.NumericValue(c.Expression.ConstValue()),
		}
	case types.Match(c.To, types.Int):
		num := values.NumericValue(c.Expression.ConstValue())
		return values.IntValue{
			Value: int64(num),
		}
	default:
		panic("unreachable")
	}
}

type InvalidExpression struct {
	Expression
}

func (i *InvalidExpression) Type() types.Type {
	return types.Invalid
}

func (i *InvalidExpression) IsConst() bool {
	return false
}

type ArrayExpression struct {
	expression
	DataType *types.ArrayType
	Elements []Expression
}

func (a *ArrayExpression) String() string {
	var result bytes.Buffer

	result.WriteByte('[')
	for i, elem := range a.Elements {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(elem.String())
	}

	result.WriteByte(']')
	return result.String()
}

func (a *ArrayExpression) Type() types.Type {
	return a.DataType
}

func (a *ArrayExpression) IsConst() bool {
	for _, elem := range a.Elements {
		if !elem.IsConst() {
			return false
		}
	}
	return true
}
func (a *ArrayExpression) ConstValue() values.ConstValue {
	if !a.IsConst() {
		return nil
	}

	elems := []values.ConstValue{}
	for _, elem := range a.Elements {
		elems = append(elems, elem.ConstValue())
	}

	return values.ArrayValue{
		Elements: elems,
	}
}

type IndexExpression struct {
	expression
	Left     Expression
	Index    Expression
	DataType types.Type
}

func (i *IndexExpression) String() string {
	return fmt.Sprintf("%s[%s]", i.Left.String(), i.Index.String())
}

func (i *IndexExpression) Type() types.Type {
	return i.DataType
}

func (i *IndexExpression) IsConst() bool {
	return false
}

func (i *IndexExpression) ConstValue() values.ConstValue {
	return nil
}

type KeyValue struct {
	Key   Expression
	Value Expression
}

type MapExpression struct {
	expression
	KeyValues []KeyValue
	DataType  *types.MapType
}

func (m *MapExpression) String() string {
	var result bytes.Buffer

	result.WriteByte('{')
	for i, kv := range m.KeyValues {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(kv.Key.String())
		result.WriteString(": ")
		result.WriteString(kv.Value.String())
	}
	result.WriteByte('}')

	return result.String()
}

func (m *MapExpression) Type() types.Type {
	return m.DataType
}

func (m *MapExpression) IsConst() bool {
	for _, kv := range m.KeyValues {
		if !kv.Key.IsConst() || !kv.Value.IsConst() {
			return false
		}
	}
	return true
}

func (m *MapExpression) ConstValue() values.ConstValue {
	if !m.IsConst() {
		return nil
	}

	value := map[uint64]values.ConstValue{}
	for _, kv := range m.KeyValues {
		keyValue := kv.Key.ConstValue()
		valueValue := kv.Value.ConstValue()
		value[keyValue.Hash()] = valueValue
	}

	return values.MapValue{
		Values: value,
	}
}

type Assignment struct {
	expression
	Assignee Expression
	Value    Expression
}

func (a *Assignment) String() string {
	return fmt.Sprintf("%s = %s", a.Assignee, a.Value)
}

func (a *Assignment) Type() types.Type {
	return a.Value.Type()
}

func (a *Assignment) IsConst() bool {
	return a.Value.IsConst()
}

func (a *Assignment) ConstValue() values.ConstValue {
	return a.Value.ConstValue()
}

type TupleExpression struct {
	expression
	Values   []Expression
	DataType *types.TupleType
}

func (t *TupleExpression) String() string {
	var result bytes.Buffer

	result.WriteByte('(')
	for i, val := range t.Values {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(val.String())
	}

	result.WriteByte(')')
	return result.String()
}

func (t *TupleExpression) Type() types.Type {
	return t.DataType
}

func (t *TupleExpression) IsConst() bool {
	for _, val := range t.Values {
		if !val.IsConst() {
			return false
		}
	}
	return true
}

func (t *TupleExpression) ConstValue() values.ConstValue {
	if !t.IsConst() {
		return nil
	}

	vals := []values.ConstValue{}
	for _, val := range t.Values {
		vals = append(vals, val.ConstValue())
	}

	return values.TupleValue{
		Values: vals,
	}
}

type TypeCheck struct {
	expression
	Value    Expression
	DataType types.Type
}

func (t *TypeCheck) String() string {
	return fmt.Sprintf("%s is %s", t.Value.String(), t.DataType.String())
}

func (t *TypeCheck) Type() types.Type {
	return types.Bool
}

func (t *TypeCheck) IsConst() bool {
	return types.Assignable(t.DataType, t.Value.Type()) ||
		!types.Assignable(t.Value.Type(), t.DataType)
}

func (t *TypeCheck) ConstValue() values.ConstValue {
	if types.Assignable(t.DataType, t.Value.Type()) {
		return values.BoolValue{Value: true}
	}
	if !types.Assignable(t.Value.Type(), t.DataType) {
		return values.BoolValue{Value: false}
	}
	return nil
}

type FunctionCall struct {
	expression
	Function   Expression
	Arguments  []Expression
	ReturnType types.Type
}

func (f *FunctionCall) String() string {
	var result bytes.Buffer

	result.WriteString(f.Function.String())
	result.WriteByte('(')
	for i, arg := range f.Arguments {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(arg.String())
	}
	result.WriteByte(')')

	return result.String()
}

func (f *FunctionCall) Type() types.Type {
	return f.ReturnType
}

func (f *FunctionCall) IsConst() bool {
	return false
}

func (f *FunctionCall) ConstValue() values.ConstValue {
	return nil
}

type StructExpression struct {
	expression
	Struct *types.Struct
	Fields map[string]Expression
}

func (s *StructExpression) String() string {
	var result bytes.Buffer

	result.WriteString(s.Struct.Name)
	result.WriteString(" {")
	if len(s.Fields) > 0 {
		result.WriteByte(' ')
	}
	isFirst := true
	for name, value := range s.Fields {
		if !isFirst {
			result.WriteString(", ")
		} else {
			isFirst = false
		}
		result.WriteString(name)
		result.WriteString(": ")
		result.WriteString(value.String())
	}
	if len(s.Fields) > 0 {
		result.WriteByte(' ')
	}
	result.WriteByte('}')

	return result.String()
}

func (s *StructExpression) Type() types.Type {
	return s.Struct
}

func (s *StructExpression) IsConst() bool {
	for _, expr := range s.Fields {
		if !expr.IsConst() {
			return false
		}
	}
	return true
}

func (s *StructExpression) ConstValue() values.ConstValue {
	if !s.IsConst() {
		return nil
	}

	members := map[string]values.ConstValue{}

	for name, expr := range s.Fields {
		members[name] = expr.ConstValue()
	}

	return values.StructValue{
		Members: members,
	}
}

// TODO:
// MemberExpression
// RangeExpression
