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

	case *ast.MapLiteral:
		return evaluateMap(expression, env)

	case *ast.AssignmentExpression:
		return evaluateAssignmentExpression(expression, env)

	case *ast.BinaryOperation:
		return evaluateBinaryOperation(expression, env)

	case *ast.UnaryOperation:
		return evaluateUnaryOperation(expression, env)

	case *ast.FunctionCall:
		return evaluateFunctionCall(expression, env)

	case *ast.IndexExpression:
		return evaluateIndexExpression(expression, env)
	
	case *ast.MemberExpression:
		return evaluateMemberExpression(*expression, env)

	case *ast.StructExpression:
		return evaluateStructExpression(*expression, env)

	default:
		errors.LogError(errors.DevError(fmt.Sprintf("(Interpreter) Unexpected expression type %q", expression.String()), expr))

		return &values.IntegerLiteral{}
	}
}

func evaluateAssignmentExpression(assignment *ast.AssignmentExpression, env *environment.Environment) values.RuntimeValue {
	var value values.RuntimeValue

	if assignment.Operation != "=" {
		operator := assignment.Operation[:len(assignment.Operation)-1]
		value = evaluateBinaryOperation(&ast.BinaryOperation{
			Left:     assignment.Assignee,
			Right:    assignment.Value,
			Operator: operator,
		}, env)
	} else {
		value = evaluateExpression(assignment.Value, env)
	}

	switch assignee := assignment.Assignee.(type) {
	case *ast.Identifier:
		return env.AssignVariable(assignee.Symbol, value)
	
	case *ast.IndexExpression:
		leftValue := evaluateExpression(assignee.Left, env)
		indexValue := evaluateExpression(assignee.Index, env)
		return leftValue.SetIndex(indexValue, value)
	
	case *ast.MemberExpression:
		leftValue := evaluateExpression(assignee.Left, env)
		return leftValue.SetMember(assignee.Member, value)
	}

	return value
}

func evaluateFunctionCall(call *ast.FunctionCall, env *environment.Environment) values.RuntimeValue {
	if ident, ok := call.Left.(*ast.Identifier); ok {
		if builtin, ok := builtins[ident.Symbol]; ok {
			args := []values.RuntimeValue{}

			for _, arg := range call.Args {
				args = append(args, evaluateExpression(arg, env))
			}

			return builtin(args, env)
		}
	}

	function := evaluateExpression(call.Left, env).(*values.FunctionValue)
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

func evaluateMap(maplit *ast.MapLiteral, env *environment.Environment) values.RuntimeValue {
	evaluatedValues := map[values.RuntimeValue]values.RuntimeValue{}

	for key, value := range maplit.Elements {
		keyValue := evaluateExpression(key, env)
		valueValue := evaluateExpression(value, env)
		evaluatedValues[keyValue] = valueValue
	}

	return &values.MapLiteral{
		Elements: evaluatedValues,
	}
}

func evaluateIndexExpression(indexExpr *ast.IndexExpression, env *environment.Environment) values.RuntimeValue {
	leftValue := evaluateExpression(indexExpr.Left, env)
	indexValue := evaluateExpression(indexExpr.Index, env)

	return leftValue.Index(indexValue)
}

func evaluateMemberExpression(memberExpr ast.MemberExpression, env *environment.Environment) values.RuntimeValue {
	value := evaluateExpression(memberExpr.Left, env)
	return value.Member(memberExpr.Member)
}

func evaluateStructExpression(structExpr ast.StructExpression, env *environment.Environment) values.RuntimeValue {
	members := map[string]values.RuntimeValue{}
	structType := env.GetStruct(structExpr.Name)

	for name, dataType := range structType.Members {
		if value, hasMember := structExpr.Members[name]; hasMember {
			members[name] = evaluateExpression(value, env)
			continue
		}
		members[name] = values.GetZeroValue(dataType.String())
	}

	return &values.StructLiteral{
		Name:    structExpr.Name,
		Members: members,
	}
}
