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

	// Don't even check for operators with invalid types, to prevent cascading errors
	if left.Type() == types.Invalid || right.Type() == types.Invalid {
		return &ir.BinaryExpression{
			Left:     left,
			Operator: 0,
			Right:    right,
		}
	}

	left, right, operator := getBinaryOperator(binExpr.Operator.Kind, left, right)

	if operator == 0 {
		t.Diagnostics.ReportBinaryOperatorUndefined(binExpr.Operator.Location, binExpr.Operator.Value, left.Type(), right.Type())
	}

	return &ir.BinaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func getBinaryOperator(op token.Kind, left, right ir.Expression) (lhs, rhs ir.Expression, binOp ir.BinaryOperator) {
	lhs = left
	rhs = right
	binOp = 0
	lType := left.Type()
	rType := right.Type()

	leftNumeric := lType == types.Int || lType == types.Float
	rightNumeric := rType == types.Int || rType == types.Float
	isFloat := lType == types.Float || rType == types.Float

	switch op {
	case token.DOUBLE_AMPERSAND:
		if lType == types.Bool && rType == types.Bool {
			binOp = ir.LogicalAnd
		}
	case token.DOUBLE_PIPE:
		if lType == types.Bool && rType == types.Bool {
			binOp = ir.LogicalOr
		}
	case token.LEFT_ANGLE:
		if leftNumeric && rightNumeric {
			binOp = ir.Less
			if isFloat {
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			}
		}
	case token.RIGHT_ANGLE:
		if leftNumeric && rightNumeric {
			binOp = ir.Greater
			if isFloat {
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			}
		}
	case token.LEFT_ANGLE_EQUALS:
		if leftNumeric && rightNumeric {
			binOp = ir.LessEq
			if isFloat {
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			}
		}
	case token.RIGHT_ANGLE_EQUALS:
		if leftNumeric && rightNumeric {
			binOp = ir.GreaterEq
			if isFloat {
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			}
		}
	case token.DOUBLE_EQUALS:
		if lType == rType {
			binOp = ir.Equal
		}
	case token.BANG_EQUALS:
		if lType == rType {
			binOp = ir.NotEqual
		}
	case token.DOUBLE_LEFT_ANGLE:
		if lType == types.Int && rType == types.Int {
			binOp = ir.LeftShift
		}
	case token.DOUBLE_RIGHT_ANGLE:
		if lType == types.Int && rType == types.Int {
			binOp = ir.RightShift
		}
	case token.PLUS:
		if lType == types.String && rType == types.String {
			binOp = ir.Concat
		}

		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.AddFloat
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			} else {
				binOp = ir.AddInt
			}
		}
	case token.MINUS:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.SubtractFloat
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			} else {
				binOp = ir.SubtractInt
			}
		}
	case token.STAR:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.MultiplyFloat
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			} else {
				binOp = ir.MultiplyInt
			}
		}
	case token.SLASH:
		if leftNumeric && rightNumeric {
			binOp = ir.Divide
			if isFloat {
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			}
		}
	case token.PERCENT:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.ModuloFloat
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			} else {
				binOp = ir.ModuloInt
			}
		}
	case token.DOUBLE_STAR:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.PowerFloat
				lhs = convert(lhs, types.Float, false)
				rhs = convert(rhs, types.Float, false)
			} else {
				binOp = ir.PowerInt
			}
		}
	case token.PIPE:
		if lType == types.Int && rType == types.Int {
			binOp = ir.BitwiseOr
		}
	case token.AMPERSAND:
		if lType == types.Int && rType == types.Int {
			binOp = ir.BitwiseAnd
		}
	}

	return
}
