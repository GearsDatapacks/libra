package interpreter

import (
	"log"

	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func evaluateExpression(expr ast.Expression) values.RuntimeValue {
	switch expression := expr.(type) {
	case *ast.IntegerLiteral:
		return values.MakeInteger(expression.Value)

	case *ast.BinaryOperation:
		return evaluateBinaryOperation(*expression)

	default:
		log.Fatalf("Unexpected expression type %t", expression)

		return &values.IntegerLiteral{}
	}
}
