package registry

import "github.com/gearsdatapacks/libra/type_checker/types"

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

func incDecOperator(dataType types.ValidType) types.ValidType {
	if !numberType.Valid(dataType) {
		return nil
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

	registerRegularBinaryOperator("||", boolType, boolType, boolType)
	registerRegularBinaryOperator("&&", boolType, boolType, boolType)

	registerUnaryOperator("++", incDecOperator)
	registerUnaryOperator("--", incDecOperator)
	registerRegularUnaryOperator("!", boolType, boolType)
}
