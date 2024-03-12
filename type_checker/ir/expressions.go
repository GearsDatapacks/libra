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
	return types.Int
}

type FloatLiteral struct {
	expression
	Value float64
}

func (f *FloatLiteral) String() string {
	return fmt.Sprint(f.Value)
}

func (FloatLiteral) Type() types.Type {
	return types.Float
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
)

func (b BinaryOperator) String() string {
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
	switch b {
	case LogicalAnd:
		return types.Bool

	case LogicalOr:
		return types.Bool

	case Less:
		return types.Bool

	case LessEq:
		return types.Bool

	case Greater:
		return types.Bool

	case GreaterEq:
		return types.Bool

	case Equal:
		return types.Bool

	case NotEqual:
		return types.Bool

	case LeftShift:
		return types.Int

	case RightShift:
		return types.Int

	case BitwiseOr:
		return types.Int

	case BitwiseAnd:
		return types.Int

	case AddInt:
		return types.Int

	case AddFloat:
		return types.Float

	case Concat:
		return types.String

	case SubtractInt:
		return types.Int

	case SubtractFloat:
		return types.Float

	case MultiplyInt:
		return types.Int

	case MultiplyFloat:
		return types.Float

	case Divide:
		return types.Float

	case ModuloInt:
		return types.Int

	case ModuloFloat:
		return types.Float

	case PowerInt:
		return types.Int

	case PowerFloat:
		return types.Float

	default:
		return types.Invalid
	}
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
// UnaryExpression
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
