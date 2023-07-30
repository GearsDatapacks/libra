package interpreter

import (
	"log"

	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

var operators = map[[3]string]opFn{}

type opFn func(values.RuntimeValue, values.RuntimeValue) values.RuntimeValue
func RegisterOperator(op string, left string, right string, operation opFn)  {
	operators[[3]string{op, left, right}] = operation
}

func evaluateBinaryOperation(binOp ast.BinaryOperation) values.RuntimeValue {
	left := evaluateExpression(binOp.Left)
	right := evaluateExpression(binOp.Right)

	operation, ok := operators[[3]string{binOp.Op, string(left.Type()), string(right.Type())}]

	if !ok {
		log.Fatalf("Operator %q does not exist or does not support operands of type %q and %q", binOp.Op, left.Type(), right.Type())
	}

	return operation(left, right)
}
