package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

var unaryOperators = map[string]unOpFn{}

type unOpFn func(values.RuntimeValue, *environment.Environment) values.RuntimeValue

func RegisterUnaryOperator(op string, operation unOpFn) {
	unaryOperators[op] = operation
}

func evaluateUnaryOperation(unOp *ast.UnaryOperation, env *environment.Environment) values.RuntimeValue {
	value := evaluateExpression(unOp.Value, env)

	operation, ok := unaryOperators[unOp.Operator]

	if !ok {
		errors.DevError(fmt.Sprintf("Operator %q does not exist", unOp.Operator), unOp)
	}

	return operation(value, env)
}
