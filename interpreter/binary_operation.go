package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
)

var binaryOperators = map[string]binOpFn{}

type binOpFn func(values.RuntimeValue, values.RuntimeValue) values.RuntimeValue

func RegisterBinaryOperator(op string, operation binOpFn) {
	binaryOperators[op] = operation
}

func evaluateBinaryOperation(binOp *ast.BinaryOperation, manager *modules.ModuleManager) values.RuntimeValue {
	if binOp.Operator == "||" {
		left := evaluateExpression(binOp.Left, manager)
		if left.Truthy() {
			return left
		}

		return evaluateExpression(binOp.Right, manager)
	}

	if binOp.Operator == "&&" {
		left := evaluateExpression(binOp.Left, manager)
		if !left.Truthy() {
			return left
		}

		return evaluateExpression(binOp.Right, manager)
	}

	left := evaluateExpression(binOp.Left, manager)
	right := evaluateExpression(binOp.Right, manager)

	operation, ok := binaryOperators[binOp.Operator]

	if !ok {
		errors.LogError(errors.DevError(fmt.Sprintf("Operator %q does not exist", binOp.Operator), binOp))
	}

	return operation(left, right)
}
