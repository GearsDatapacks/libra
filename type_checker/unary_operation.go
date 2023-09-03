package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckUnaryOperation(unOp *ast.UnaryOperation, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	valueType, err := typeCheckExpression(unOp.Value, symbolTable)
	if err != nil {
		return nil, err
	}

	checkerFn, exists := registry.UnaryOperators[unOp.Operator]

	if !exists {
		return nil, errors.TypeError(fmt.Sprintf("Operator %q does not exist", unOp.Operator), unOp)
	}

	resultType := checkerFn(valueType)

	if resultType == nil {
		return nil, errors.TypeError(fmt.Sprintf("Operator %q is not defined for type %q", unOp.Operator, valueType), unOp)
	}

	return resultType, nil
}
