package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckBinaryOperation(binOp *ast.BinaryOperation, manager *modules.ModuleManager) types.ValidType {
	leftType := typeCheckExpression(binOp.Left, manager)
	if leftType.String() == "TypeError" {
		return leftType
	}

	rightType := typeCheckExpression(binOp.Right, manager)
	if rightType.String() == "TypeError" {
		return rightType
	}

	checkerFn, exists := registry.BinaryOperators[binOp.Operator]

	if !exists {
		return types.Error(fmt.Sprintf("Operator %q does not exist", binOp.Operator), binOp)
	}

	resultType := checkerFn(leftType, rightType)

	if resultType == nil {
		return types.Error(fmt.Sprintf("Operator %q is not defined for types %q and %q", binOp.Operator, leftType, rightType), binOp)
	}

	if resultType.String() == "TypeError" {
		resultType.(*types.TypeError).Line = binOp.Token.Line
		resultType.(*types.TypeError).Column = binOp.Token.Column
	}

	return resultType
}
