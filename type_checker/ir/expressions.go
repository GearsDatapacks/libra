package ir

import (
	"math"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

type expression struct{}

func (expression) irExpr() {}

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

type IntegerLiteral struct {
	expression
	Value int64
}

func (i *IntegerLiteral) Print(node *printer.Node) {
	node.Text(
		"%sINT_LIT %s%v",
		node.Colour(colour.NodeName),
		node.Colour(colour.Literal),
		i.Value,
	)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (i IntegerLiteral) Type() types.Type {
	ui := types.Numeric{
		Kind:         types.NumInt,
		BitWidth:     maxInt(minIntWidth(i.Value), 32),
		Downcastable: &types.Downcastable{},
	}
	ui.Downcastable.MinFloatWidth = minFloatWidth(float64(i.Value))
	ui.Downcastable.MinIntWidth = minIntWidth(i.Value)
	ui.Downcastable.MinUintWidth = minUintWidth(uint64(i.Value))

	return ui
}

func (IntegerLiteral) IsConst() bool {
	return true
}

func (i *IntegerLiteral) ConstValue() values.ConstValue {
	return values.IntValue{
		Value: i.Value,
	}
}

type UintLiteral struct {
	expression
	Value uint64
}

func (i UintLiteral) Print(node *printer.Node) {
	node.Text(
		"%sUINT_LIT %s%v",
		node.Colour(colour.NodeName),
		node.Colour(colour.Literal),
		i.Value,
	)
}

func (i UintLiteral) Type() types.Type {
	ui := types.Numeric{
		Kind:         types.NumInt,
		BitWidth:     maxInt(minUintWidth(i.Value), 32),
		Downcastable: &types.Downcastable{},
	}
	ui.Downcastable.MinFloatWidth = minFloatWidth(float64(i.Value))
	ui.Downcastable.MinIntWidth = minIntWidth(int64(i.Value))
	ui.Downcastable.MinUintWidth = minUintWidth(i.Value)

	return ui
}

func (UintLiteral) IsConst() bool {
	return true
}

func (i UintLiteral) ConstValue() values.ConstValue {
	return values.UintValue{
		Value: i.Value,
	}
}

type FloatLiteral struct {
	expression
	Value float64
}

func (f *FloatLiteral) Print(node *printer.Node) {
	node.Text(
		"%sFLOAT_LIT %s%v",
		node.Colour(colour.NodeName),
		node.Colour(colour.Literal),
		f.Value,
	)
}

func minFloatWidth(f float64) int {
	if f <= types.MaxFloat16 && f >= -types.MaxFloat16 {
		return 16
	}
	if f <= math.MaxFloat32 && f >= -math.MaxFloat32 {
		return 32
	}
	return 64
}

func (f *FloatLiteral) Type() types.Type {
	uf := types.Numeric{
		Kind:         types.NumFloat,
		BitWidth:     64,
		Downcastable: &types.Downcastable{},
	}
	uf.Downcastable.MinFloatWidth = minFloatWidth(f.Value)
	if float64(int64(f.Value)) == f.Value {
		uf.Downcastable.MinIntWidth = minIntWidth(int64(f.Value))
		uf.Downcastable.MinUintWidth = minUintWidth(uint64(f.Value))
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

func (b *BooleanLiteral) Print(node *printer.Node) {
	node.Text(
		"%sBOOL_LIT %s%t",
		node.Colour(colour.NodeName),
		node.Colour(colour.Literal),
		b.Value,
	)
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

func (s *StringLiteral) Print(node *printer.Node) {
	node.Text(
		"%sSTRING_LIT %s%q",
		node.Colour(colour.NodeName),
		node.Colour(colour.Literal),
		s.Value,
	)
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

func (v *VariableExpression) Print(node *printer.Node) {
	v.Symbol.Print(node)
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

type BinOpId int

const (
	_ BinOpId = iota
	LogicalAnd
	LogicalOr
	Less
	LessEq
	Greater
	GreaterEq
	Equal
	NotEqual
	LeftShift
	ArithmeticRightShift
	LogicalRightShift
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

type BinaryOperator struct {
	Id       BinOpId
	DataType types.Type
}

func (bo BinaryOperator) String() string {
	b := bo.Id & ^UntypedBit

	switch b {
	case LogicalAnd:
		return "LogicalAnd"
	case LogicalOr:
		return "LogicalOr"
	case Less:
		return "Less"
	case LessEq:
		return "LessEq"
	case Greater:
		return "Greater"
	case GreaterEq:
		return "GreaterEq"
	case Equal:
		return "Equal"
	case NotEqual:
		return "NotEqual"
	case LeftShift:
		return "LeftShift"
	case ArithmeticRightShift:
		return "ArithmeticRightShift"
	case LogicalRightShift:
		return "LogicalRightShift"
	case Union:
		return "Union"
	case BitwiseOr:
		return "BitwiseOr"
	case BitwiseAnd:
		return "BitwiseAnd"
	case AddInt:
		return "AddInt"
	case AddFloat:
		return "AddFloat"
	case Concat:
		return "Concat"
	case SubtractInt:
		return "SubtractInt"
	case SubtractFloat:
		return "SubtractFloat"
	case MultiplyInt:
		return "MultiplyInt"
	case MultiplyFloat:
		return "MultiplyFloat"
	case Divide:
		return "Divide"
	case ModuloInt:
		return "ModuloInt"
	case ModuloFloat:
		return "ModuloFloat"
	case PowerInt:
		return "PowerInt"
	case PowerFloat:
		return "PowerFloat"
	default:
		return "<?>"
	}
}

func (b BinaryOperator) Type() types.Type {
	untyped := b.Id&UntypedBit != 0
	id := b.Id & ^UntypedBit
	var ty types.Type = types.Invalid

	switch id {
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
		ty = types.I32

	case ArithmeticRightShift:
		ty = types.I32

	case LogicalRightShift:
		ty = types.I32

	case BitwiseOr:
		ty = types.I32

	case Union:
		ty = types.RuntimeType

	case BitwiseAnd:
		ty = types.I32

	case AddInt:
		ty = b.DataType

	case AddFloat:
		ty = b.DataType

	case Concat:
		ty = types.String

	case SubtractInt:
		ty = b.DataType

	case SubtractFloat:
		ty = b.DataType

	case MultiplyInt:
		ty = b.DataType

	case MultiplyFloat:
		ty = b.DataType

	case Divide:
		ty = types.F32

	case ModuloInt:
		ty = b.DataType

	case ModuloFloat:
		ty = b.DataType

	case PowerInt:
		ty = b.DataType

	case PowerFloat:
		ty = b.DataType
	}

	if untyped {
		if variable, ok := ty.(types.Numeric); ok {
			variable.Downcastable = &types.Downcastable{}

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

func (b *BinaryExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sBINARY_EXPR %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Attribute),
			b.Operator.String(),
		).
		Node(b.Left).
		Node(b.Right).
		Node(b.Type()).
		OptionalNode(b.ConstValue())
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

	switch b.Operator.Id & ^UntypedBit {
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

	case ArithmeticRightShift:
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

type UnaryOperator struct {
	Id       UnOpId
	DataType types.Type
}

type UnOpId int

const (
	_ UnOpId = iota
	NegateInt
	NegateFloat
	Identity
	LogicalNot
	BitwiseNot
)

const (
	postfix UnOpId = iota | (1 << 16)
	IncrecementInt
	IncrementFloat
	DecrecementInt
	DecrementFloat
	PropagateError
	CrashError
)

func (u UnaryOperator) String() string {
	id := u.Id & ^UntypedBit

	switch id {
	case NegateInt:
		return "NegateInt"
	case NegateFloat:
		return "NegateFloat"
	case Identity:
		return "Identity"
	case LogicalNot:
		return "LogicalNot"
	case BitwiseNot:
		return "BitwiseNot"
	case IncrecementInt:
		return "IncrecementInt"
	case IncrementFloat:
		return "IncrementFloat"
	case DecrecementInt:
		return "DecrecementInt"
	case DecrementFloat:
		return "DecrementFloat"
	case PropagateError:
		return "PropagateError"
	case CrashError:
		return "CrashError"
	default:
		return "<?>"
	}
}

func (u UnaryOperator) Type() types.Type {
	untyped := u.Id&UntypedBit != 0
	id := u.Id & ^UntypedBit
	var ty types.Type = types.Invalid

	switch id {
	case NegateInt:
		ty = types.I32
	case NegateFloat:
		ty = types.F32
	case Identity:
		ty = types.I32
	case LogicalNot:
		ty = types.Bool
	case BitwiseNot:
		ty = types.I32
	case IncrecementInt:
		ty = types.I32
	case IncrementFloat:
		ty = types.F32
	case DecrecementInt:
		ty = types.I32
	case DecrementFloat:
		ty = types.F32
	case PropagateError:
		return u.DataType
	case CrashError:
		return u.DataType
	}

	if untyped {
		if variable, ok := ty.(types.Numeric); ok {
			variable.Downcastable = &types.Downcastable{}
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

func (u *UnaryExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sUNARY_EXPR %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Attribute),
			u.Operator.String(),
		).
		Node(u.Operand).
		Node(u.Type()).
		OptionalNode(u.ConstValue())
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

	switch u.Operator.Id & ^UntypedBit {
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

func (c *Conversion) Print(node *printer.Node) {
	node.
		Text(
			"%sCONVERSION",
			node.Colour(colour.NodeName),
		).
		Node(c.Expression).
		Node(c.To).
		OptionalNode(c.ConstValue())
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

	if n, ok := c.To.(types.Numeric); ok {
		num := values.NumericValue(c.Expression.ConstValue())
		if n.Kind == types.NumFloat {
			return values.FloatValue{
				Value: num,
			}
		} else if n.Kind == types.NumInt {
			return values.IntValue{
				Value: int64(num),
			}
		} else if n.Kind == types.NumUint {
			return values.UintValue{
				Value: uint64(num),
			}
		}
	}

	return c.Expression.ConstValue()
}

type InvalidExpression struct {
	Expression
}

func (i *InvalidExpression) Print(node *printer.Node) {
	node.
		Text("%sINVALID_EXPR", node.Colour(colour.NodeName)).
		OptionalNode(i.Expression)
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

func (a *ArrayExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sARRAY_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(a.DataType).
		OptionalNode(a.ConstValue())

	printer.Nodes(node, a.Elements)
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

func (i *IndexExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sINDEX_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(i.Left).
		Node(i.Index).
		Node(i.DataType).
		OptionalNode(i.ConstValue())
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

func (kv KeyValue) Print(node *printer.Node) {
	node.
		Text("%sKEY_VALUE", node.Colour(colour.NodeName)).
		Node(kv.Key).
		Node(kv.Value)
}

type MapExpression struct {
	expression
	KeyValues []KeyValue
	DataType  *types.MapType
}

func (m *MapExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sMAP_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(m.DataType).
		OptionalNode(m.ConstValue())

	printer.Nodes(node, m.KeyValues)
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

	value := map[uint64]values.KeyValue{}
	for _, kv := range m.KeyValues {
		keyValue := kv.Key.ConstValue()
		valueValue := kv.Value.ConstValue()
		hash := keyValue.Hash()
		value[hash] = values.KeyValue{
			Key:   keyValue,
			Value: valueValue,
		}
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

func (a *Assignment) Print(node *printer.Node) {
	node.
		Text(
			"%sASSIGNMENT",
			node.Colour(colour.NodeName),
		).
		Node(a.Assignee).
		Node(a.Value)
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

func (t *TupleExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sTUPLE_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(t.Type()).
		OptionalNode(t.ConstValue())

	printer.Nodes(node, t.Values)
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

func (t *TypeCheck) Print(node *printer.Node) {
	node.
		Text(
			"%sTYPE_CHECK_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(t.Value).
		Node(t.DataType).
		OptionalNode(t.ConstValue())
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

func (f *FunctionCall) Print(node *printer.Node) {
	node.
		Text(
			"%sFUNCTION_CALL",
			node.Colour(colour.NodeName),
		).
		Node(f.Function).
		Node(f.ReturnType)

	printer.Nodes(node, f.Arguments)
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
	Struct types.Type
	Fields map[string]Expression
}

func (s *StructExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sSTRUCT_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(s.Struct).
		OptionalNode(s.ConstValue())

	for _, keyValue := range printer.SortMap(s.Fields) {
		node.FakeNode(
			"%sSTRUCT_FIELD %s%s",
			func(n *printer.Node) { n.Node(keyValue.Value) },
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			keyValue.Key,
		)
	}
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
	Struct types.Type
	Fields []Expression
}

func (t *TupleStructExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sTUPLE_STRUCT_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(t.Struct).
		OptionalNode(t.ConstValue())

	printer.Nodes(node, t.Fields)
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

func (m *MemberExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sMEMBER_EXPR %s%s",
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			m.Member,
		).
		Node(m.Left).
		Node(m.DataType).
		OptionalNode(m.ConstValue())
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

func (b *Block) Print(node *printer.Node) {
	node.
		Text(
			"%sBLOCK",
			node.Colour(colour.NodeName),
		).
		Node(b.ResultType)

	printer.Nodes(node, b.Statements)
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
	ResultType types.Type
	Body       *Block
	ElseBranch Statement
}

func (i *IfExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sIF_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(i.Condition).
		Node(i.Body)

	if i.ElseBranch != nil {
		node.FakeNode(
			"%sELSE_BRANCH",
			func(n *printer.Node) { n.Node(i.ElseBranch) },
			node.Colour(colour.NodeName),
		)
	}
}

func (i *IfExpression) Type() types.Type {
	return i.ResultType
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

func (w *WhileLoop) Print(node *printer.Node) {
	node.
		Text(
			"%sWHILE_LOOP",
			node.Colour(colour.NodeName),
		).
		Node(w.Condition).
		Node(w.Body)
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

func (f *ForLoop) Print(node *printer.Node) {
	node.
		Text(
			"%sFOR_LOOP",
			node.Colour(colour.NodeName),
		).
		Node(&f.Variable).
		Node(f.Iterator).
		Node(f.Body)
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

func (t *TypeExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sTYPE_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(t.DataType)
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
	Location   text.Location
}

func (f *FunctionExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sFUNC_EXPR",
			node.Colour(colour.NodeName),
		)

	for _, param := range f.Parameters {
		node.Text(
			" %s%s",
			node.Colour(colour.Name),
			param,
		)
	}

	node.
		Node(f.Body).
		Node(f.DataType)
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

func (r *RefExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sREF_EXPR",
			node.Colour(colour.NodeName),
		).
		TextIf(
			r.Mutable,
			" %smut",
			node.Colour(colour.Attribute),
		).
		Node(r.Value)
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

func (d *DerefExpression) Print(node *printer.Node) {
	node.
		Text(
			"%sDEREF_EXPR",
			node.Colour(colour.NodeName),
		).
		Node(d.Value)
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
