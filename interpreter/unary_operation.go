package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
)

var unaryOperators = map[string]unOpFn{}

type unOpFn func(values.RuntimeValue, bool, *environment.Environment) values.RuntimeValue

func RegisterUnaryOperator(op string, operation unOpFn) {
	unaryOperators[op] = operation
}

func evaluateUnaryOperation(unOp *ast.UnaryOperation, manager *modules.ModuleManager) values.RuntimeValue {
	value := evaluateExpression(unOp.Value, manager)

	operation, ok := unaryOperators[unOp.Operator]

	if !ok {
		errors.LogError(errors.DevError(fmt.Sprintf("Operator %q does not exist", unOp.Operator), unOp))
	}

	return operation(value, unOp.Postfix, manager.Env)
}
