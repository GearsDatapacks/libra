package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
)

type operatorChecker func(types.ValidType, types.ValidType)types.ValidType

var operators = map[string]operatorChecker{}

func RegisterOperator(operator string, fn operatorChecker) {
	operators[operator] = fn
}

func typeCheckBinaryOperation(binOp *ast.BinaryOperation, symbolTable *symbols.SymbolTable) types.ValidType {
	leftType := typeCheckExpression(binOp.Left, symbolTable)
	rightType := typeCheckExpression(binOp.Right, symbolTable)

	checkerFn, exists := operators[binOp.Operator]

	if !exists {
		errors.TypeError(fmt.Sprintf("Operator %q does not exist", binOp.Operator), binOp)
	}

	resultType := checkerFn(leftType, rightType)

	if resultType == nil {
		errors.TypeError(fmt.Sprintf("Operator %q is not defined for types %q and %q", binOp.Operator, leftType, rightType), binOp)
	}

	return resultType
}

func registerRegularOperator(name string, left, right, result types.ValidType) {
	fn := func(leftType, rightType types.ValidType) types.ValidType {
		if left.Valid(leftType) && right.Valid(rightType) {
			return result
		}
		return nil
	}

	RegisterOperator(name, fn)
}

var boolType = types.MakeLiteral(types.BOOL)
var floatType = types.MakeLiteral(types.FLOAT)
var intType = types.MakeLiteral(types.INT)
var numberType = types.MakeUnion(intType, floatType)

func arithmeticOperator(leftType, rightType types.ValidType) types.ValidType {
	if !numberType.Valid(leftType) || !numberType.Valid(rightType) {
		return nil
	}
	
	if leftType.Valid(floatType) || rightType.Valid(floatType) {
		return floatType
	}

	return intType
}

func RegisterOperators() {
	RegisterOperator("+", arithmeticOperator)
	RegisterOperator("-", arithmeticOperator)
	RegisterOperator("*", arithmeticOperator)
	registerRegularOperator("/", numberType, numberType, floatType)
	RegisterOperator("%", arithmeticOperator)

	registerRegularOperator(">", numberType, numberType, boolType)
	registerRegularOperator(">=", numberType, numberType, boolType)
	registerRegularOperator("<", numberType, numberType, boolType)
	registerRegularOperator("<=", numberType, numberType, boolType)

	registerRegularOperator("||", boolType, boolType, boolType)
	registerRegularOperator("&&", boolType, boolType, boolType)
}
