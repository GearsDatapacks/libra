package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func evaluateExpression(expr ast.Expression, manager *modules.ModuleManager) values.RuntimeValue {
	switch expression := expr.(type) {
	case *ast.IntegerLiteral:
		return values.MakeUntypedNumber(float64(expression.Value), false)

	case *ast.FloatLiteral:
		return values.MakeUntypedNumber(expression.Value, true)

	case *ast.StringLiteral:
		return values.MakeString(expression.Value)

	case *ast.BooleanLiteral:
		return values.MakeBoolean(expression.Value)

	case *ast.NullLiteral:
		return values.MakeNull()

	case *ast.VoidValue:
		return values.MakeNull()

	case *ast.Identifier:
		return manager.Env.GetVariable(expression.Symbol)

	case *ast.ListLiteral:
		return evaluateList(expression, manager)

	case *ast.MapLiteral:
		return evaluateMap(expression, manager)

	case *ast.AssignmentExpression:
		return evaluateAssignmentExpression(expression, manager)

	case *ast.BinaryOperation:
		return evaluateBinaryOperation(expression, manager)

	case *ast.UnaryOperation:
		return evaluateUnaryOperation(expression, manager)

	case *ast.FunctionCall:
		return evaluateFunctionCall(expression, manager)

	case *ast.IndexExpression:
		return evaluateIndexExpression(expression, manager)

	case *ast.MemberExpression:
		return evaluateMemberExpression(*expression, manager)

	case *ast.StructExpression:
		return evaluateStructExpression(*expression, manager)

	case *ast.TupleExpression:
		return evaluateTuple(expression, manager)

	case *ast.CastExpression:
		return evaluateCastExpression(expression, manager)

	case *ast.TypeCheckExpression:
		return evaluateTypeCheckExpression(expression, manager)

	default:
		errors.LogError(errors.DevError(fmt.Sprintf("(Interpreter) Unexpected expression type %q", expression.String()), expr))

		return &values.IntegerLiteral{}
	}
}

func evaluateAssignmentExpression(assignment *ast.AssignmentExpression, manager *modules.ModuleManager) values.RuntimeValue {
	var value values.RuntimeValue

	if assignment.Operation != "=" {
		operator := assignment.Operation[:len(assignment.Operation)-1]
		value = evaluateBinaryOperation(&ast.BinaryOperation{
			Left:     assignment.Assignee,
			Right:    assignment.Value,
			Operator: operator,
		}, manager)
	} else {
		value = evaluateExpression(assignment.Value, manager)
	}

	switch assignee := assignment.Assignee.(type) {
	case *ast.Identifier:
		leftValue := manager.Env.GetVariable(assignee.Symbol)
		return manager.Env.AssignVariable(assignee.Symbol, leftValue.Type(), value)

	case *ast.IndexExpression:
		leftValue := evaluateExpression(assignee.Left, manager)
		indexValue := evaluateExpression(assignee.Index, manager)
		return leftValue.SetIndex(indexValue, value)

	case *ast.MemberExpression:
		leftValue := evaluateExpression(assignee.Left, manager)
		return leftValue.SetMember(assignee.Member, value)
	}

	return value
}

func evaluateFunctionCall(call *ast.FunctionCall, manager *modules.ModuleManager) values.RuntimeValue {
	if ident, ok := call.Left.(*ast.Identifier); ok {
		if structType, isStruct := manager.SymbolTable.GetType(ident.Symbol).(*types.TupleStruct); isStruct {
			return evaluateTupleStructExpression(structType, call, manager)
		}

		if builtin, ok := builtins[ident.Symbol]; ok {
			args := []values.RuntimeValue{}

			for _, arg := range call.Args {
				args = append(args, evaluateExpression(arg, manager))
			}

			return builtin(args, manager.Env)
		}
	}

	ty := typechecker.TypeCheckTypeExpression(call.Left, manager)
	if structType, isStruct := ty.(*types.TupleStruct); isStruct {
		return evaluateTupleStructExpression(structType, call, manager)
	}

	function := evaluateExpression(call.Left, manager).(*values.FunctionValue)
	declarationEnv := function.Env.(*environment.Environment)
	mod := function.Manager.(*modules.ModuleManager)
	scope := environment.NewChild(declarationEnv, environment.FUNCTION_SCOPE)

	for i, param := range function.Parameters {
		arg := evaluateExpression(call.Args[i], manager)
		scope.DeclareVariable(param.Name, param.Type, arg)
	}

	if function.This != nil {
		scope.DeclareVariable("this", function.This.Type(), function.This)
	}

	mod.EnterEnv(scope)

	for _, statement := range function.Body {
		evaluate(statement, mod)

		if scope.ReturnValue != nil {
			return scope.ReturnValue
		}
	}

	mod.ExitEnv()

	return values.MakeNull()
}

func evaluateList(list *ast.ListLiteral, manager *modules.ModuleManager) values.RuntimeValue {
	evaluatedValues := []values.RuntimeValue{}
	listTypes := []types.ValidType{}

	for _, elem := range list.Elements {
		elemValue := evaluateExpression(elem, manager)
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

func evaluateMap(maplit *ast.MapLiteral, manager *modules.ModuleManager) values.RuntimeValue {
	keyTypes := []types.ValidType{}
	valueTypes := []types.ValidType{}
	evaluatedValues := map[values.RuntimeValue]values.RuntimeValue{}

	for key, value := range maplit.Elements {
		keyValue := evaluateExpression(key, manager)
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

		valueValue := evaluateExpression(value, manager)
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

func evaluateIndexExpression(indexExpr *ast.IndexExpression, manager *modules.ModuleManager) values.RuntimeValue {
	leftValue := evaluateExpression(indexExpr.Left, manager)
	indexValue := evaluateExpression(indexExpr.Index, manager)

	return leftValue.Index(indexValue)
}

func evaluateMemberExpression(memberExpr ast.MemberExpression, manager *modules.ModuleManager) values.RuntimeValue {
	value := evaluateExpression(memberExpr.Left, manager)

	method := environment.GetMethod(memberExpr.Member, value.Type())
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

func evaluateStructExpression(structExpr ast.StructExpression, manager *modules.ModuleManager) values.RuntimeValue {
	members := map[string]values.RuntimeValue{}
	structType := typechecker.TypeCheckTypeExpression(structExpr.InstanceOf, manager).(*types.Struct)

	for name, dataType := range structType.Members {
		if value, hasMember := structExpr.Members[name]; hasMember {
			members[name] = evaluateExpression(value, manager)
			continue
		}
		members[name] = values.GetZeroValue(dataType.Type.String())
	}

	return &values.StructLiteral{
		Name:      structExpr.InstanceOf.String(),
		Members:   members,
		BaseValue: values.BaseValue{DataType: structType},
	}
}

func evaluateTuple(tuple *ast.TupleExpression, manager *modules.ModuleManager) values.RuntimeValue {
	members := []values.RuntimeValue{}

	for _, member := range tuple.Members {
		members = append(members, evaluateExpression(member, manager))
	}

	return &values.TupleValue{Members: members}
}

func evaluateTupleStructExpression(tupleType *types.TupleStruct, tupleExpr *ast.FunctionCall, manager *modules.ModuleManager) values.RuntimeValue {
	members := []values.RuntimeValue{}
	for _, arg := range tupleExpr.Args {
		members = append(members, evaluateExpression(arg, manager))
	}

	return &values.TupleStructValue{
		BaseValue: values.BaseValue{DataType: tupleType},
		Members:   members,
		Name:      tupleExpr.Left.String(),
	}
}

func evaluateCastExpression(cast *ast.CastExpression, manager *modules.ModuleManager) values.RuntimeValue {
	left := evaluateExpression(cast.Left, manager)
	ty := typechecker.TypeCheckType(cast.DataType, manager)
	castable, ok := ty.(types.CustomCastable)
	if !ty.Valid(left.Type()) && !(ok && castable.CanCast(left.Type())) {
		errors.LogError(fmt.Sprintf("%q is type %q, not %q", left.ToString(), left.Type(), ty))
	}

	return values.Cast(left, ty)
}

func evaluateTypeCheckExpression(expr *ast.TypeCheckExpression, manager *modules.ModuleManager) values.RuntimeValue {
	left := evaluateExpression(expr.Left, manager)
	ty := typechecker.TypeCheckType(expr.DataType, manager)
	return values.MakeBoolean(ty.Valid(left.Type()))
}
