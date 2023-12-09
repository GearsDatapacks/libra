package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckUnaryOperation(unOp *ast.UnaryOperation, manager *modules.ModuleManager) types.ValidType {
	valueType := typeCheckExpression(unOp.Value, manager)
	if valueType.String() == "TypeError" {
		return valueType
	}

	not_exit_error := types.Error(fmt.Sprintf("Operator %q does not exist", unOp.Operator), unOp)

	if unOp.Operator == "?" {
		if errType, ok := valueType.(*types.ErrorType); ok {
			if !manager.SymbolTable.IsInFunctionScope() {
				return types.Error("Cannot use operator \"?\" outside of a function", unOp)
			}
			if _, ok := manager.SymbolTable.ReturnType().(*types.ErrorType); !ok {
				return types.Error("Can only use operator \"?\" in a function that returns an error type", unOp)
			}
			return errType.ResultType
		}
		return not_exit_error
	}

	checkerFn, exists := registry.UnaryOperators[unOp.Operator]

	if !exists {
		return not_exit_error
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
