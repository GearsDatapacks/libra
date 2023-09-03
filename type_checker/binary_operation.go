package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckBinaryOperation(binOp *ast.BinaryOperation, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	leftType, err := typeCheckExpression(binOp.Left, symbolTable)
	if err != nil {
		return nil, err
	}

	rightType, err := typeCheckExpression(binOp.Right, symbolTable)
	if err != nil {
		return nil, err
	}

	checkerFn, exists := registry.BinaryOperators[binOp.Operator]

	if !exists {
		return nil, errors.TypeError(fmt.Sprintf("Operator %q does not exist", binOp.Operator), binOp)
	}

	resultType := checkerFn(leftType, rightType)

	if resultType == nil {
		return nil, errors.TypeError(fmt.Sprintf("Operator %q is not defined for types %q and %q", binOp.Operator, leftType, rightType), binOp)
	}

	return resultType, nil
}
