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

	case *ast.FloatLiteral:
		return values.MakeFloat(expression.Value)
		
	case *ast.BooleanLiteral:
		return values.MakeBoolean(expression.Value)

	case *ast.NullLiteral:
		return values.MakeNull()

	case *ast.VoidValue:
		return values.MakeNull()

	case *ast.Identifier:
		return env.GetVariable(expression.Symbol)

	case *ast.AssignmentExpression:
		return evaluateAssignmentExpression(expression, env)

	case *ast.BinaryOperation:
		return evaluateBinaryOperation(expression, env)
	
	case *ast.FunctionCall:
		return evaluateFunctionCall(expression, env)

	default:
		errors.DevError(fmt.Sprintf("Unexpected expression type %q", expression.String()), expr)

		return &values.IntegerLiteral{}
	}
}

func evaluateAssignmentExpression(assignment *ast.AssignmentExpression, env *environment.Environment) values.RuntimeValue {
	varName := assignment.Assignee.(*ast.Identifier).Symbol
	value := evaluateExpression(assignment.Value, env)

	return env.AssignVariable(varName, value)
}

func evaluateFunctionCall(call *ast.FunctionCall, env *environment.Environment) values.RuntimeValue {
	function := env.GetVariable(call.Name).(*values.FunctionValue)
	declarationEnvironment := function.DeclarationEnvironment.(*environment.Environment)
	scope := environment.NewChild(declarationEnvironment, environment.FUNCTION_SCOPE)

	for i, param := range function.Parameters {
		arg := evaluateExpression(call.Args[i], env)
		scope.DeclareVariable(param, arg)
	}

	for _, statement := range function.Body {
		evaluate(statement, scope)

		if scope.ReturnValue != nil {
			return scope.ReturnValue
		}
	}

	return values.MakeNull()
}
