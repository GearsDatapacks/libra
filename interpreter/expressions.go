package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func evaluateExpression(expr ast.Expression, env *environment.Environment) values.RuntimeValue {
	switch expression := expr.(type) {
	case *ast.IntegerLiteral:
		return values.MakeInteger(expression.Value)

	case *ast.Identifier:
		return env.GetVariable(expression.Symbol)

	case *ast.AssignmentExpression:
		return evaluateAssignmentExpression(*expression, env)

	case *ast.BinaryOperation:
		return evaluateBinaryOperation(*expression, env)

	default:
		errors.DevError(fmt.Sprintf("Unexpected expression type %t", expression), expr)

		return &values.IntegerLiteral{}
	}
}
