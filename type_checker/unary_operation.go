package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckUnaryOperation(unOp *ast.UnaryOperation, symbolTable *symbols.SymbolTable) types.ValidType {
	valueType := typeCheckExpression(unOp.Value, symbolTable)
	if valueType.String() == "TypeError" {
		return valueType
	}

	checkerFn, exists := registry.UnaryOperators[unOp.Operator]

	if !exists {
		return types.Error(fmt.Sprintf("Operator %q does not exist", unOp.Operator), unOp)
	}

	resultType := checkerFn(valueType, unOp.Postfix)

	if resultType == nil {
		return types.Error(fmt.Sprintf("Operator %q is not defined for type %q", unOp.Operator, valueType), unOp)
	}

	if resultType.String() == "TypeError" {
		resultType.(*types.TypeError).Line = unOp.Token.Line
		resultType.(*types.TypeError).Column = unOp.Token.Column
	}

	return resultType
}
