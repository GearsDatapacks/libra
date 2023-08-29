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

	case *ast.StringLiteral:
		return values.MakeString(expression.Value)
		
	case *ast.BooleanLiteral:
		return values.MakeBoolean(expression.Value)

	case *ast.NullLiteral:
		return values.MakeNull()

	case *ast.VoidValue:
		return values.MakeNull()

	case *ast.Identifier:
		return env.GetVariable(expression.Symbol)
	
	case *ast.ListLiteral:
		return evaluateList(expression, env)

	case *ast.AssignmentExpression:
		return evaluateAssignmentExpression(expression, env)

	case *ast.BinaryOperation:
		return evaluateBinaryOperation(expression, env)
	
	case *ast.UnaryOperation:
		return evaluateUnaryOperation(expression, env)
	
	case *ast.FunctionCall:
		return evaluateFunctionCall(expression, env)

	default:
		errors.DevError(fmt.Sprintf("(Interpreter) Unexpected expression type %q", expression.String()), expr)

		return &values.IntegerLiteral{}
	}
}

func evaluateAssignmentExpression(assignment *ast.AssignmentExpression, env *environment.Environment) values.RuntimeValue {
	varName := assignment.Assignee.(*ast.Identifier).Symbol

	if assignment.Operation != "=" {
		operator := assignment.Operation[:len(assignment.Operation)-1]
		newValue := evaluateBinaryOperation(&ast.BinaryOperation{
			Left: assignment.Assignee,
			Right: assignment.Value,
			Operator: operator,
		}, env)

		return env.AssignVariable(varName, newValue)
	}

	value := evaluateExpression(assignment.Value, env)

	return env.AssignVariable(varName, value)
}

func evaluateFunctionCall(call *ast.FunctionCall, env *environment.Environment) values.RuntimeValue {
	if builtin, ok := builtins[call.Name]; ok {
		args := []values.RuntimeValue{}

		for _, arg := range call.Args {
			args = append(args, evaluateExpression(arg, env))
		}

		return builtin(args, env)
	}

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

func evaluateList(list *ast.ListLiteral, env *environment.Environment) values.RuntimeValue {
	evaluatedValues := []values.RuntimeValue{}

	for _, elem := range list.Elements {
		evaluatedValues = append(evaluatedValues, evaluateExpression(elem, env))
	}

	return &values.ListLiteral{
		Elements: evaluatedValues,
	}
}
