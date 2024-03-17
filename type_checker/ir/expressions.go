package ir

import (
	"bytes"
	"fmt"

	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
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

// TODO:
// ListLiteral
// MapLiteral
// FunctionCall
// IndexExpression
// AssignmentExpression
// TupleExpression
// MemberExpression
// StructMember
// StructExpression
// CastExpression
// TypeCheckExpression
// RangeExpression
