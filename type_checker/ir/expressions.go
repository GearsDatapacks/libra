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
	Union
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

	case Union:
		fallthrough
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

	case Union:
		ty = types.RuntimeType

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

	case Union:
		left := b.Left.ConstValue().(values.TypeValue).Type.(types.Type)
		right := b.Right.ConstValue().(values.TypeValue).Type.(types.Type)
		return values.TypeValue{
			Type: types.MakeUnion(left, right),
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

	if types.Match(c.To, types.Float) {
		return values.FloatValue{
			Value: values.NumericValue(c.Expression.ConstValue()),
		}
	}

	if types.Match(c.To, types.Int) {
		num := values.NumericValue(c.Expression.ConstValue())
		return values.IntValue{
			Value: int64(num),
		}
	}

	if _, ok := c.To.(*types.Union); ok {
		return c.Expression.ConstValue()
	}

	if _, ok := c.To.(*types.Explicit); ok {
		return c.Expression.ConstValue()
	}

	panic("unreachable")
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
	return i.Left.IsConst() && i.Index.IsConst()
}

func (i *IndexExpression) ConstValue() values.ConstValue {
	if !i.IsConst() {
		return nil
	}

	return i.Left.ConstValue().Index(i.Index.ConstValue())
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

type TupleStructExpression struct {
	expression
	Struct *types.TupleStruct
	Fields []Expression
}

func (t *TupleStructExpression) String() string {
	var result bytes.Buffer

	result.WriteString(t.Struct.Name)
	result.WriteString(" {")
	if len(t.Fields) > 0 {
		result.WriteByte(' ')
	}

	for i, value := range t.Fields {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(value.String())
	}
	if len(t.Fields) > 0 {
		result.WriteByte(' ')
	}
	result.WriteByte('}')

	return result.String()
}

func (t *TupleStructExpression) Type() types.Type {
	return t.Struct
}

func (s *TupleStructExpression) IsConst() bool {
	for _, expr := range s.Fields {
		if !expr.IsConst() {
			return false
		}
	}
	return true
}

func (s *TupleStructExpression) ConstValue() values.ConstValue {
	if !s.IsConst() {
		return nil
	}

	members := []values.ConstValue{}

	for _, expr := range s.Fields {
		members = append(members, expr.ConstValue())
	}

	return values.TupleValue{
		Values: members,
	}
}

type MemberExpression struct {
	expression
	Left     Expression
	Member   string
	DataType types.Type
}

func (m *MemberExpression) String() string {
	return fmt.Sprintf("%s.%s", m.Left.String(), m.Member)
}

func (m *MemberExpression) Type() types.Type {
	return m.DataType
}

func (m *MemberExpression) IsConst() bool {
	return m.Left.IsConst() && m.Left.ConstValue().Member(m.Member) != nil
}

func (m *MemberExpression) ConstValue() values.ConstValue {
	if !m.IsConst() {
		return nil
	}
	return m.Left.ConstValue().Member(m.Member)
}

type Block struct {
	expression
	Statements []Statement
	ResultType types.Type
}

func (b *Block) String() string {
	var result bytes.Buffer

	result.WriteByte('{')
	if len(b.Statements) > 0 {
		result.WriteByte('\n')
	}
	for _, stmt := range b.Statements {
		result.WriteString(stmt.String())
		result.WriteByte('\n')
	}
	result.WriteByte('}')

	return result.String()
}

func (b *Block) Type() types.Type {
	return b.ResultType
}

func (*Block) IsConst() bool {
	return false
}

func (*Block) ConstValue() values.ConstValue {
	return nil
}

type IfExpression struct {
	expression
	Condition  Expression
	Body       *Block
	ElseBranch Statement
}

func (i *IfExpression) String() string {
	var result bytes.Buffer
	result.WriteString("if ")
	result.WriteString(i.Condition.String())
	result.WriteByte(' ')
	result.WriteString(i.Body.String())

	if i.ElseBranch != nil {
		result.WriteString("\nelse ")
		result.WriteString(i.ElseBranch.String())
	}

	return result.String()
}

func (i *IfExpression) Type() types.Type {
	return i.Body.Type()
}

func (*IfExpression) IsConst() bool {
	return false
}

func (*IfExpression) ConstValue() values.ConstValue {
	return nil
}

type WhileLoop struct {
	expression
	Condition Expression
	Body      *Block
}

func (w *WhileLoop) String() string {
	var result bytes.Buffer
	result.WriteString("while ")
	result.WriteString(w.Condition.String())
	result.WriteByte(' ')
	result.WriteString(w.Body.String())

	return result.String()
}

func (w *WhileLoop) Type() types.Type {
	return w.Body.Type()
}

func (*WhileLoop) IsConst() bool {
	return false
}

func (*WhileLoop) ConstValue() values.ConstValue {
	return nil
}

type ForLoop struct {
	expression
	Variable symbols.Variable
	Iterator Expression
	Body     *Block
}

func (f *ForLoop) String() string {
	var result bytes.Buffer

	result.WriteString("for ")
	result.WriteString(f.Variable.Name)
	result.WriteString(" in ")
	result.WriteString(f.Iterator.String())
	result.WriteByte(' ')
	result.WriteString(f.Body.String())

	return result.String()
}

func (f *ForLoop) Type() types.Type {
	return f.Body.Type()
}

func (*ForLoop) IsConst() bool {
	return false
}

func (*ForLoop) ConstValue() values.ConstValue {
	return nil
}

type TypeExpression struct {
	expression
	DataType types.Type
}

func (t *TypeExpression) String() string {
	return t.DataType.String()
}

func (t *TypeExpression) Type() types.Type {
	return types.RuntimeType
}

func (t *TypeExpression) IsConst() bool {
	return true
}

func (t *TypeExpression) ConstValue() values.ConstValue {
	return values.TypeValue{Type: t.DataType}
}

type FunctionExpression struct {
	expression
	Parameters []string
	Body       *Block
	DataType   *types.Function
}

func (f *FunctionExpression) String() string {
	var result bytes.Buffer

	result.WriteString("fn(")

	for i, param := range f.Parameters {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(param)
	}

	result.WriteString(") ")
	result.WriteString(f.Body.String())

	return result.String()
}

func (f *FunctionExpression) Type() types.Type {
	return f.DataType
}

func (t *FunctionExpression) IsConst() bool {
	return false
}

func (t *FunctionExpression) ConstValue() values.ConstValue {
	return nil
}

type RefExpression struct {
	expression
	Value   Expression
	Mutable bool
}

func (r *RefExpression) String() string {
	if r.Mutable {
		return fmt.Sprintf("&mut %s", r.Value.String())
	}
	return fmt.Sprintf("&%s", r.Value.String())
}

func (r *RefExpression) Type() types.Type {
	return &types.Pointer{
		Underlying: r.Value.Type(),
		Mutable:    r.Mutable,
	}
}

func (*RefExpression) IsConst() bool {
	return false
}

func (*RefExpression) ConstValue() values.ConstValue {
	return nil
}

type DerefExpression struct {
	expression
	Value Expression
}

func (d *DerefExpression) String() string {
	return fmt.Sprintf("*%s", d.Value.String())
}

func (d *DerefExpression) Type() types.Type {
	if ptr, ok := d.Value.Type().(*types.Pointer); ok {
		return ptr.Underlying
	}
	return types.Invalid
}

func (*DerefExpression) IsConst() bool {
	return false
}

func (*DerefExpression) ConstValue() values.ConstValue {
	return nil
}

// TODO:
// RangeExpression
