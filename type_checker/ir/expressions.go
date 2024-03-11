package ir

import (
	"fmt"

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

func (i *IntegerLiteral) Type() types.Type {
	return types.Int
}

// TODO:
// FloatLiteral
// BooleanLiteral
// StringLiteral
// Identifier
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
