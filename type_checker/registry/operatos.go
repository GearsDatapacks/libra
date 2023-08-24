package registry

import "github.com/gearsdatapacks/libra/type_checker/types"

type operatorChecker func(types.ValidType, types.ValidType) types.ValidType

var Operators = map[string]operatorChecker{}

func registerOperator(operator string, fn operatorChecker) {
	Operators[operator] = fn
}

func registerRegularOperator(name string, left, right, result types.ValidType) {
	fn := func(leftType, rightType types.ValidType) types.ValidType {
		if left.Valid(leftType) && right.Valid(rightType) {
			return result
		}
		return nil
	}

	registerOperator(name, fn)
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

func registerOperators() {
	registerOperator("+", arithmeticOperator)
	registerOperator("-", arithmeticOperator)
	registerOperator("*", arithmeticOperator)
	registerRegularOperator("/", numberType, numberType, floatType)
	registerOperator("%", arithmeticOperator)

	registerRegularOperator(">", numberType, numberType, boolType)
	registerRegularOperator(">=", numberType, numberType, boolType)
	registerRegularOperator("<", numberType, numberType, boolType)
	registerRegularOperator("<=", numberType, numberType, boolType)
	registerRegularOperator("==", &types.Any{}, &types.Any{}, boolType)
	registerRegularOperator("!=", &types.Any{}, &types.Any{}, boolType)

	registerRegularOperator("||", boolType, boolType, boolType)
	registerRegularOperator("&&", boolType, boolType, boolType)
}
