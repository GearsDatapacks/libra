package registry

import (
	"fmt"

	"github.com/gearsdatapacks/libra/type_checker/types"
)

type binaryOperatorChecker func(types.ValidType, types.ValidType) types.ValidType
type unaryOperatorChecker func(types.ValidType, bool) types.ValidType

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

func registerRegularUnaryOperator(name string, postfix bool, value, result types.ValidType) {
	fn := func(valueType types.ValidType, post bool) types.ValidType {
		if post == postfix && value.Valid(valueType) {
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

func leftShift(leftType, rightType types.ValidType) types.ValidType {
	if intType.Valid(leftType) && intType.Valid(rightType) {
		return intType
	}

	list, isList := leftType.(*types.ListLiteral)
	if isList {
		if leftType.Constant() {
			return types.Error("Cannot append value to constant list")
		}

		if !list.ElemType.Valid(rightType) {
			return types.Error(fmt.Sprintf("Cannot append value of type %q to list of type %q", rightType, leftType))
		}
		return leftType
	}

	_, isArray := leftType.(*types.ArrayLiteral)
	if isArray {
		return types.Error("Cannot append items to an array of fixed length")
	}

	return nil
}

func rightShift(leftType, rightType types.ValidType) types.ValidType {
	if intType.Valid(leftType) && intType.Valid(rightType) {
		return intType
	}

	list, isList := rightType.(*types.ListLiteral)
	if isList {
		if rightType.Constant() {
			return types.Error("Cannot prepend value to constant list")
		}

		if !list.ElemType.Valid(leftType) {
			return types.Error(fmt.Sprintf("Cannot prepend value of type %q to list of type %q", rightType, leftType))
		}
		return rightType
	}

	_, isArray := rightType.(*types.ArrayLiteral)
	if isArray {
		return types.Error("Cannot prepend items to an array of fixed length")
	}

	return nil
}

func incDecOperator(dataType types.ValidType, op string) types.ValidType {
	if !numberType.Valid(dataType) {
		return nil
	}

	if !dataType.WasVariable() {
		return types.Error(fmt.Sprintf("Operator %q is not defined for non-variable values", op))
	}

	if dataType.Constant() {
		incDecStr := "increment"
		if op == "--" {
			incDecStr = "decrement"
		}
		return types.Error(fmt.Sprintf("Cannot %s a constant", incDecStr))
	}

	return dataType
}

func negateOperator(dataType types.ValidType, _ bool) types.ValidType {
	if !numberType.Valid(dataType) {
		return nil
	}

	return dataType
}

func notOperator(dataType types.ValidType, postfix bool) types.ValidType {
	if !postfix {
		if boolType.Valid(dataType) {
			return boolType
		}
		return nil
	}

	if errType, ok := dataType.(*types.ErrorType); ok {
		return errType.ResultType
	}
	return nil
}

func unwrapOperator(dataType types.ValidType, _ bool) types.ValidType {
	if errType, ok := dataType.(*types.ErrorType); ok {
		return errType.ResultType
	}
	return nil
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

	registerBinaryOperator("<<", leftShift)
	registerBinaryOperator(">>", rightShift)

	registerBinaryOperator("||", logicalOperator)
	registerBinaryOperator("&&", logicalOperator)

	registerUnaryOperator("++", func(v types.ValidType, _ bool) types.ValidType { return incDecOperator(v, "++") })
	registerUnaryOperator("--", func(v types.ValidType, _ bool) types.ValidType { return incDecOperator(v, "--") })
	registerUnaryOperator("!", notOperator)
	registerUnaryOperator("?", unwrapOperator)
	registerUnaryOperator("-", negateOperator)
}
