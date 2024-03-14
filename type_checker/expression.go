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
			Type:    types.Invalid,
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

	leftNumeric := types.Assignable(types.Int, lType) || types.Assignable(types.Float, lType)
	rightNumeric := types.Assignable(types.Int, rType) || types.Assignable(types.Float, rType)
	isFloat := !types.Assignable(types.Int, lType) || !types.Assignable(types.Int, rType)
	var lUntyped bool
	var rUntyped bool
	if v, ok := lType.(types.VariableType); ok {
		lUntyped = v.Untyped
	}
	if v, ok := rType.(types.VariableType); ok {
		rUntyped = v.Untyped
	}
	untyped := lUntyped && rUntyped

	switch op {
	case token.DOUBLE_AMPERSAND:
		if types.Assignable(types.Bool, lType) && types.Assignable(types.Bool, rType) {
			binOp = ir.LogicalAnd
		}
	case token.DOUBLE_PIPE:
		if types.Assignable(types.Bool, lType) && types.Assignable(types.Bool, rType) {
			binOp = ir.LogicalOr
		}
	case token.LEFT_ANGLE:
		if leftNumeric && rightNumeric {
			binOp = ir.Less
			if isFloat {
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			}
		}
	case token.RIGHT_ANGLE:
		if leftNumeric && rightNumeric {
			binOp = ir.Greater
			if isFloat {
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			}
		}
	case token.LEFT_ANGLE_EQUALS:
		if leftNumeric && rightNumeric {
			binOp = ir.LessEq
			if isFloat {
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			}
		}
	case token.RIGHT_ANGLE_EQUALS:
		if leftNumeric && rightNumeric {
			binOp = ir.GreaterEq
			if isFloat {
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			}
		}
	case token.DOUBLE_EQUALS:
		if types.Assignable(rType, lType) || types.Assignable(lType, rType) {
			binOp = ir.Equal
		}
	case token.BANG_EQUALS:
		if types.Assignable(rType, lType) || types.Assignable(lType, rType) {
			binOp = ir.NotEqual
		}
	case token.DOUBLE_LEFT_ANGLE:
		if types.Assignable(types.Int, lType) && types.Assignable(types.Int, rType) {
			binOp = ir.LeftShift
		}
	case token.DOUBLE_RIGHT_ANGLE:
		if types.Assignable(types.Int, lType) && types.Assignable(types.Int, rType) {
			binOp = ir.RightShift
		}
	case token.PLUS:
		if types.Assignable(types.String, lType) && types.Assignable(types.String, rType) {
			binOp = ir.Concat
		}

		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.AddFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.AddInt
			}
		}
	case token.MINUS:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.SubtractFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.SubtractInt
			}
		}
	case token.STAR:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.MultiplyFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.MultiplyInt
			}
		}
	case token.SLASH:
		if leftNumeric && rightNumeric {
			binOp = ir.Divide
			if isFloat {
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			}
		}
	case token.PERCENT:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.ModuloFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.ModuloInt
			}
		}
	case token.DOUBLE_STAR:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.PowerFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.PowerInt
			}
		}
	case token.PIPE:
		if types.Assignable(types.Int, lType) && types.Assignable(types.Int, rType) {
			binOp = ir.BitwiseOr
		}
	case token.AMPERSAND:
		if types.Assignable(types.Int, lType) && types.Assignable(types.Int, rType) {
			binOp = ir.BitwiseAnd
		}
	}

	if untyped {
		binOp = binOp | ir.UntypedBit
	}
	return
}
