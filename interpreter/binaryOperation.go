package interpreter

import (
	"log"

	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

var operators = map[[3]string]opFn{}

type opFn func(values.RuntimeValue, values.RuntimeValue) values.RuntimeValue

func RegisterOperator(op string, left string, right string, operation opFn) {
	operators[[3]string{op, left, right}] = operation
}

func evaluateBinaryOperation(binOp ast.BinaryOperation, env *environment.Environment) values.RuntimeValue {
	left := evaluateExpression(binOp.Left, env)
	right := evaluateExpression(binOp.Right, env)

	operation, ok := operators[[3]string{binOp.Operator, string(left.Type()), string(right.Type())}]

	if !ok {
		log.Fatalf("Operator %q does not exist or does not support operands of type %q and %q", binOp.Operator, left.Type(), right.Type())
	}

	return operation(left, right)
}
