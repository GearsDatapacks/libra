package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
)

var operators = map[string][3]types.ValidType{}

func RegisterOperator(operator string, left, right, result types.ValidType) {
	operators[operator] = [3]types.ValidType{left, right, result}
}

func typeCheckBinaryOperation(binOp *ast.BinaryOperation, symbolTable *symbols.SymbolTable) types.ValidType {
	leftType := typeCheckExpression(binOp.Left, symbolTable)
	rightType := typeCheckExpression(binOp.Right, symbolTable)

	types, exists := operators[binOp.Operator]

	if !exists {
		errors.TypeError(fmt.Sprintf("Operator %q does not exist", binOp.Operator), binOp)
	}

	validTypes := types[0].Valid(leftType) && types[1].Valid(rightType)

	if !validTypes {
		errors.TypeError(fmt.Sprintf("Operator %q is not defined for types %q and %q", binOp.Operator, leftType, rightType), binOp)
	}

	return types[2]
}

func RegisterOperators() {
	numberType := types.MakeUnion(types.INT, types.FLOAT)
	boolType := types.MakeLiteral(types.BOOL)

	RegisterOperator("+", numberType, numberType, numberType)
	RegisterOperator("-", numberType, numberType, numberType)
	RegisterOperator("*", numberType, numberType, numberType)
	RegisterOperator("/", numberType, numberType, numberType)
	RegisterOperator("%", numberType, numberType, numberType)

	RegisterOperator(">", numberType, numberType, boolType)
	RegisterOperator(">=", numberType, numberType, boolType)
	RegisterOperator("<", numberType, numberType, boolType)
	RegisterOperator("<=", numberType, numberType, boolType)

	RegisterOperator("||", boolType, boolType, boolType)
	RegisterOperator("&&", boolType, boolType, boolType)
}
