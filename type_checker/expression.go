package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func (t *typeChecker) typeCheckExpression(expression ast.Expression) ir.Expression {
	switch expr := expression.(type) {
	case *ast.IntegerLiteral:
		return &ir.IntegerLiteral{
			Value: expr.Value,
		}
	case *ast.FloatLiteral:
		return &ir.FloatLiteral{
			Value: expr.Value,
		}
	case *ast.BooleanLiteral:
		return &ir.BooleanLiteral{
			Value: expr.Value,
		}
	case *ast.StringLiteral:
		return &ir.StringLiteral{
			Value: expr.Value,
		}
	case *ast.Identifier:
		return t.typeCheckIdentifier(expr)
	case *ast.BinaryExpression:
		return t.typeCheckBinaryExpression(expr)
	default:
		panic(fmt.Sprintf("TODO: Type-check %T", expr))
	}
}

func (t *typeChecker) typeCheckIdentifier(ident *ast.Identifier) ir.Expression {
	variable := t.symbols.LookupVariable(ident.Name)
	if variable == nil {
		t.Diagnostics.ReportVariableUndefined(ident.Token.Location, ident.Name)
		variable = &symbols.Variable{
			Name:    ident.Name,
			Mutable: true,
			Type:    nil,
		}
	}
	return &ir.VariableExpression{
		Symbol: *variable,
	}
}

func (t *typeChecker) typeCheckBinaryExpression(binExpr *ast.BinaryExpression) ir.Expression {
	left := t.typeCheckExpression(binExpr.Left)
	right := t.typeCheckExpression(binExpr.Right)

	operator := getBinaryOperator(binExpr.Operator.Kind, left.Type(), right.Type())

	if operator == 0 {
		t.Diagnostics.ReportBinaryOperatorUndefined(binExpr.Operator.Location, binExpr.Operator.Value, left.Type(), right.Type())
	}

	return &ir.BinaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func getBinaryOperator(op token.Kind, left, right types.Type) ir.BinaryOperator {
	switch op {
	case token.DOUBLE_AMPERSAND:
		if left == types.Bool && right == types.Bool {
			return ir.LogicalAnd
		}
	case token.DOUBLE_PIPE:
		if left == types.Bool && right == types.Bool {
			return ir.LogicalOr
		}
	case token.LEFT_ANGLE:
		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		if leftNumeric && rightNumeric {
			return ir.Less
		}
	case token.RIGHT_ANGLE:
		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		if leftNumeric && rightNumeric {
			return ir.Greater
		}
	case token.LEFT_ANGLE_EQUALS:
		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		if leftNumeric && rightNumeric {
			return ir.LessEq
		}
	case token.RIGHT_ANGLE_EQUALS:
		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		if leftNumeric && rightNumeric {
			return ir.GreaterEq
		}
	case token.DOUBLE_EQUALS:
		if left == right {
			return ir.Equal
		}
	case token.BANG_EQUALS:
		if left == right {
			return ir.NotEqual
		}
	case token.DOUBLE_LEFT_ANGLE:
		if left == types.Int && right == types.Int {
			return ir.LeftShift
		}
	case token.DOUBLE_RIGHT_ANGLE:
		if left == types.Int && right == types.Int {
			return ir.RightShift
		}
	case token.PLUS:
		if left == types.String && right == types.String {
			return ir.Concat
		}

		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		isFloat := left == types.Float || right == types.Float
		if leftNumeric && rightNumeric {
			if isFloat {
				return ir.AddFloat
			} else {
				return ir.AddInt
			}
		}
	case token.MINUS:
		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		isFloat := left == types.Float || right == types.Float
		if leftNumeric && rightNumeric {
			if isFloat {
				return ir.SubtractFloat
			} else {
				return ir.SubtractInt
			}
		}
	case token.STAR:
		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		isFloat := left == types.Float || right == types.Float
		if leftNumeric && rightNumeric {
			if isFloat {
				return ir.MultiplyFloat
			} else {
				return ir.MultiplyInt
			}
		}
	case token.SLASH:
		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		if leftNumeric && rightNumeric {
			return ir.Divide
		}
	case token.PERCENT:
		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		isFloat := left == types.Float || right == types.Float
		if leftNumeric && rightNumeric {
			if isFloat {
				return ir.ModuloFloat
			} else {
				return ir.ModuloInt
			}
		}
	case token.DOUBLE_STAR:
		leftNumeric := left == types.Int || left == types.Float
		rightNumeric := right == types.Int || right == types.Float
		isFloat := left == types.Float || right == types.Float
		if leftNumeric && rightNumeric {
			if isFloat {
				return ir.PowerFloat
			} else {
				return ir.PowerInt
			}
		}
	case token.PIPE:
		if left == types.Int && right == types.Int {
			return ir.BitwiseOr
		}
	case token.AMPERSAND:
		if left == types.Int && right == types.Int {
			return ir.BitwiseAnd
		}
	}

	return 0
}
