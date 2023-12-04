package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
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

	case *ast.TupleExpression:
		return evaluateTuple(expression, env)

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
		if structType, isStruct := env.GetType(ident.Symbol).(*types.TupleStruct); isStruct {
			return evaluateTupleStructExpression(structType, call, env)
		}

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

	if function.This != nil {
		scope.DeclareVariable("this", function.This)
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
	listTypes := []types.ValidType{}

	for _, elem := range list.Elements {
		elemValue := evaluateExpression(elem, env)
		newType := true
		for _, listType := range listTypes {
			if listType.Valid(elemValue.Type()) {
				newType = false
				break
			}
		}

		if newType {
			listTypes = append(listTypes, elemValue.Type())
		}
		evaluatedValues = append(evaluatedValues, elemValue)
	}

	return &values.ListLiteral{
		Elements: evaluatedValues,
		BaseValue: values.BaseValue{
			DataType: &types.ArrayLiteral{
				ElemType: types.MakeUnion(listTypes...),
				Length:   len(list.Elements),
				CanInfer: true,
			},
		},
	}
}

func evaluateMap(maplit *ast.MapLiteral, env *environment.Environment) values.RuntimeValue {
	keyTypes := []types.ValidType{}
	valueTypes := []types.ValidType{}
	evaluatedValues := map[values.RuntimeValue]values.RuntimeValue{}

	for key, value := range maplit.Elements {
		keyValue := evaluateExpression(key, env)
		keyType := keyValue.Type()
		newType := true
		for _, dataType := range keyTypes {
			if dataType.Valid(keyType) {
				newType = false
				break
			}
		}

		if newType {
			keyTypes = append(keyTypes, keyType)
		}

		valueValue := evaluateExpression(value, env)
		evaluatedValues[keyValue] = valueValue

		valueType := valueValue.Type()
		newType = true
		for _, dataType := range valueTypes {
			if dataType.Valid(valueType) {
				newType = false
				break
			}
		}

		if newType {
			valueTypes = append(valueTypes, valueType)
		}
	}

	return &values.MapLiteral{
		Elements: evaluatedValues,
		BaseValue: values.BaseValue{DataType: &types.MapLiteral{
			KeyType:   types.MakeUnion(keyTypes...),
			ValueType: types.MakeUnion(valueTypes...),
		}},
	}
}

func evaluateIndexExpression(indexExpr *ast.IndexExpression, env *environment.Environment) values.RuntimeValue {
	leftValue := evaluateExpression(indexExpr.Left, env)
	indexValue := evaluateExpression(indexExpr.Index, env)

	return leftValue.Index(indexValue)
}

func evaluateMemberExpression(memberExpr ast.MemberExpression, env *environment.Environment) values.RuntimeValue {
	value := evaluateExpression(memberExpr.Left, env)

	method := env.GetMethod(memberExpr.Member, value.Type())
	if method != nil {
		method.This = value
		return method
	}

	memberValue := value.Member(memberExpr.Member)
	if memberValue != nil {
		return memberValue
	}

	return values.MakeNull()
}

func evaluateStructExpression(structExpr ast.StructExpression, env *environment.Environment) values.RuntimeValue {
	members := map[string]values.RuntimeValue{}
	structType := env.GetType(structExpr.Name).(*types.Struct)

	for name, dataType := range structType.Members {
		if value, hasMember := structExpr.Members[name]; hasMember {
			members[name] = evaluateExpression(value, env)
			continue
		}
		members[name] = values.GetZeroValue(dataType.String())
	}

	return &values.StructLiteral{
		Name:      structExpr.Name,
		Members:   members,
		BaseValue: values.BaseValue{DataType: structType},
	}
}

func evaluateTuple(tuple *ast.TupleExpression, env *environment.Environment) values.RuntimeValue {
	members := []values.RuntimeValue{}

	for _, member := range tuple.Members {
		members = append(members, evaluateExpression(member, env))
	}

	return &values.TupleValue{Members: members}
}

func evaluateTupleStructExpression(tupleType *types.TupleStruct, tupleExpr *ast.FunctionCall, env *environment.Environment) values.RuntimeValue {
	members := []values.RuntimeValue{}
	for _, arg := range tupleExpr.Args {
		members = append(members, evaluateExpression(arg, env))
	}

	return &values.TupleStructValue{
		BaseValue: values.BaseValue{DataType: tupleType},
		Members: members,
		Name: tupleExpr.Left.String(),
	}
}
