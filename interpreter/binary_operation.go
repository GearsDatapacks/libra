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
	if binOp.Operator == "||" {
		left := evaluateExpression(binOp.Left, env)
		if left.Truthy() {
			return left
		}

		return evaluateExpression(binOp.Right, env)
	}

	if binOp.Operator == "&&" {
		left := evaluateExpression(binOp.Left, env)
		if !left.Truthy() {
			return left
		}

		return evaluateExpression(binOp.Right, env)
	}

	left := evaluateExpression(binOp.Left, env)
	right := evaluateExpression(binOp.Right, env)

	operation, ok := binaryOperators[binOp.Operator]

	if !ok {
		errors.LogError(errors.DevError(fmt.Sprintf("Operator %q does not exist", binOp.Operator), binOp))
	}

	return operation(left, right)
}
