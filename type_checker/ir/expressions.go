package ir

import (
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

// TODO:
// BinaryExpression
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
