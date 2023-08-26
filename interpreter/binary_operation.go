package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

var binaryOperators = map[string]binOpFn{}

type binOpFn func(values.RuntimeValue, values.RuntimeValue) values.RuntimeValue

func RegisterBinaryOperator(op string, operation binOpFn) {
	binaryOperators[op] = operation
}

func evaluateBinaryOperation(binOp *ast.BinaryOperation, env *environment.Environment) values.RuntimeValue {
	left := evaluateExpression(binOp.Left, env)
	right := evaluateExpression(binOp.Right, env)

	operation, ok := binaryOperators[binOp.Operator]

	if !ok {
		errors.DevError(fmt.Sprintf("Operator %q does not exist", binOp.Operator), binOp)
	}

	return operation(left, right)
}
