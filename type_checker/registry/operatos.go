package registry

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type binaryOperatorChecker func(types.ValidType, types.ValidType) types.ValidType
type unaryOperatorChecker func(types.ValidType) types.ValidType

var BinaryOperators = map[string]binaryOperatorChecker{}
var UnaryOperators = map[string]unaryOperatorChecker{}

func registerBinaryOperator(operator string, fn binaryOperatorChecker) {
	BinaryOperators[operator] = fn
}

func registerUnaryOperator(operator string, fn unaryOperatorChecker) {
	UnaryOperators[operator] = fn
}

func registerRegularBinaryOperator(name string, left, right, result types.ValidType) {
	fn := func(leftType, rightType types.ValidType) types.ValidType {
		if left.Valid(leftType) && right.Valid(rightType) {
			return result
		}
		return nil
	}

	registerBinaryOperator(name, fn)
}

func registerRegularUnaryOperator(name string, value, result types.ValidType) {
	fn := func(valueType types.ValidType) types.ValidType {
		if value.Valid(valueType) {
			return result
		}
		return nil
	}

	registerUnaryOperator(name, fn)
}

func arithmeticOperator(leftType, rightType types.ValidType) types.ValidType {
	if !numberType.Valid(leftType) || !numberType.Valid(rightType) {
		return nil
	}

	if leftType.Valid(floatType) || rightType.Valid(floatType) {
		return floatType
	}

	return intType
}

func plusOperator(leftType, rightType types.ValidType) types.ValidType {
	if !stringType.Valid(leftType) || !stringType.Valid(rightType) {
		return arithmeticOperator(leftType, rightType)
	}

	return stringType
}

func powerOperator(leftType, rightType types.ValidType) types.ValidType {
	if !numberType.Valid(leftType) || !intType.Valid(rightType) {
		return nil
	}

	if leftType.Valid(floatType) {
		return floatType
	}

	return intType
}

func logicalOperator(leftType, rightType types.ValidType) types.ValidType {
	if leftType.Valid(rightType) {
		return leftType
	}
	return types.MakeUnion(leftType, rightType)
}

func incDecOperator(dataType types.ValidType, op string) types.ValidType {
	if !numberType.Valid(dataType) {
		return nil
	}

	if !dataType.WasVariable() {
		errors.TypeError(fmt.Sprintf("Operator %q is not defined for non-variable values", op))
	}

	return dataType
}

func registerOperators() {
	registerBinaryOperator("+", plusOperator)
	registerBinaryOperator("-", arithmeticOperator)
	registerBinaryOperator("*", arithmeticOperator)
	registerRegularBinaryOperator("/", numberType, numberType, floatType)
	registerBinaryOperator("%", arithmeticOperator)
	registerBinaryOperator("**", powerOperator)

	registerRegularBinaryOperator(">", numberType, numberType, boolType)
	registerRegularBinaryOperator(">=", numberType, numberType, boolType)
	registerRegularBinaryOperator("<", numberType, numberType, boolType)
	registerRegularBinaryOperator("<=", numberType, numberType, boolType)
	registerRegularBinaryOperator("==", &types.Any{}, &types.Any{}, boolType)
	registerRegularBinaryOperator("!=", &types.Any{}, &types.Any{}, boolType)

	registerBinaryOperator("||", logicalOperator)
	registerBinaryOperator("&&", logicalOperator)

	registerUnaryOperator("++", func(v types.ValidType) types.ValidType { return incDecOperator(v, "++") })
	registerUnaryOperator("--", func(v types.ValidType) types.ValidType { return incDecOperator(v, "--") })
	registerRegularUnaryOperator("!", &types.Any{}, boolType)
}
