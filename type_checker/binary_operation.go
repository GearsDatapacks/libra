package typechecker

import (
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

var operators = map[operation]types.DataType{}

type operation struct {
	operator string
	left     types.DataType
	right    types.DataType
}

func RegisterOperator(operator string, left, right, result types.DataType) {
	operation := operation{
		operator: operator,
		left:     left,
		right:    right,
	}
	operators[operation] = result
}

func typeCheckBinaryOperation(binOp *ast.BinaryOperation, symbolTable *SymbolTable) types.DataType {
	leftType := typeCheckExpression(binOp.Left, symbolTable)
	rightType := typeCheckExpression(binOp.Right, symbolTable)
	operation := operation{
		operator: binOp.Operator,
		left:     leftType,
		right:    rightType,
	}

	resultType, validTypes := operators[operation]

	if !validTypes {
		return types.INVALID
	}

	return resultType
}

func RegisterOperators() {
	RegisterOperator("+", types.INT, types.INT, types.INT)
	RegisterOperator("-", types.INT, types.INT, types.INT)
	RegisterOperator("*", types.INT, types.INT, types.INT)
	RegisterOperator("/", types.INT, types.INT, types.INT)
	RegisterOperator("%", types.INT, types.INT, types.INT)

	RegisterOperator(">", types.INT, types.INT, types.BOOL)
	RegisterOperator(">=", types.INT, types.INT, types.BOOL)
	RegisterOperator("<", types.INT, types.INT, types.BOOL)
	RegisterOperator("<=", types.INT, types.INT, types.BOOL)

	RegisterOperator("||", types.BOOL, types.BOOL, types.BOOL)
	RegisterOperator("&&", types.BOOL, types.BOOL, types.BOOL)
}
