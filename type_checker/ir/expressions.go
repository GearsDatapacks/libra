package ir

import "fmt"

type expression struct{}

func (expression) irExpr() {}

type IntegerLiteral struct {
	expression
	Value int64
}

func (i *IntegerLiteral) String() string {
	return fmt.Sprint(i.Value)
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
