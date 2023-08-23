package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)



func typeCheckBinaryOperation(binOp *ast.BinaryOperation, symbolTable *symbols.SymbolTable) types.ValidType {
	leftType := typeCheckExpression(binOp.Left, symbolTable)
	rightType := typeCheckExpression(binOp.Right, symbolTable)

	checkerFn, exists := registry.Operators[binOp.Operator]

	if !exists {
		errors.TypeError(fmt.Sprintf("Operator %q does not exist", binOp.Operator), binOp)
	}

	resultType := checkerFn(leftType, rightType)

	if resultType == nil {
		errors.TypeError(fmt.Sprintf("Operator %q is not defined for types %q and %q", binOp.Operator, leftType, rightType), binOp)
	}

	return resultType
}
