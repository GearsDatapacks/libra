package typechecker

import (
	"fmt"
	"log"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckExpression(expr ast.Expression, symbolTable *symbols.SymbolTable) types.ValidType {
	switch expression := expr.(type) {
	case *ast.IntegerLiteral:
		return &types.IntLiteral{}
	case *ast.FloatLiteral:
		return &types.FloatLiteral{}
	case *ast.StringLiteral:
		return &types.StringLiteral{}
	case *ast.NullLiteral:
		return &types.NullLiteral{}
	case *ast.BooleanLiteral:
		return &types.BoolLiteral{}
	case *ast.VoidValue:
		return &types.Void{}

	case *ast.Identifier:
		dataType := symbolTable.GetSymbol(expression.Symbol)
		if err, isErr := dataType.(*types.TypeError); isErr {
			err.Line = expression.Token.Line
			err.Column = expression.Token.Column
		}
		return dataType

	case *ast.BinaryOperation:
		return typeCheckBinaryOperation(expression, symbolTable)

	case *ast.UnaryOperation:
		return typeCheckUnaryOperation(expression, symbolTable)

	case *ast.AssignmentExpression:
		return typeCheckAssignmentExpression(expression, symbolTable)

	case *ast.FunctionCall:
		return typeCheckFunctionCall(expression, symbolTable)

	case *ast.ListLiteral:
		return typeCheckList(expression, symbolTable)

	case *ast.MapLiteral:
		return typeCheckMap(expression, symbolTable)

	case *ast.IndexExpression:
		return typeCheckIndexExpression(expression, symbolTable)
	
	case *ast.MemberExpression:
		return typeCheckMemberExpression(expression, symbolTable)

	case *ast.StructExpression:
		return typeCheckStructExpression(expression, symbolTable)
	
	case *ast.TupleExpression:
		return typeCheckTuple(expression, symbolTable)

	default:
		log.Fatal(errors.DevError("(Type checker) Unexpected expression type: " + expr.String()))
		return nil
	}
}

func typeCheckAssignmentExpression(assignment *ast.AssignmentExpression, symbolTable *symbols.SymbolTable) types.ValidType {
	var dataType types.ValidType
	if assignment.Assignee.Type() == "Identifier" {
		symbolName := assignment.Assignee.(*ast.Identifier).Symbol
	
		dataType = symbolTable.GetSymbol(symbolName)
		
	} else if assignment.Assignee.Type() == "IndexExpression" {
		index := assignment.Assignee.(*ast.IndexExpression)
		leftType := typeCheckExpression(index.Left, symbolTable)
		if leftType.String() == "TypeError" {
			return leftType
		}
		indexType := typeCheckExpression(index.Index, symbolTable)
		if indexType.String() == "TypeError" {
			return indexType
		}

		dataType = leftType.IndexBy(indexType)
	}  else if assignment.Assignee.Type() == "MemberExpression" {
		member := assignment.Assignee.(*ast.MemberExpression)
		leftType := typeCheckExpression(member.Left, symbolTable)
		if leftType.String() == "TypeError" {
			return leftType
		}

		dataType = types.Member(leftType, member.Member, member.IsNumberMember)
	} else {
		return types.Error("Can only assign values to variables", assignment)
	}

	if dataType.String() == "TypeError" {
		return dataType
	}

	if dataType.Constant() {
		return types.Error("Cannot assign data to constant value", assignment)
	}

	expressionType := typeCheckExpression(assignment.Value, symbolTable)
	if expressionType.String() == "TypeError" {
		return expressionType
	}
	correctType := dataType.Valid(expressionType)

	if correctType {
		return dataType
	}

	return types.Error(fmt.Sprintf("Type %q is not assignable to type %q", expressionType, dataType), assignment)
}

func typeCheckFunctionCall(call *ast.FunctionCall, symbolTable *symbols.SymbolTable) types.ValidType {
	if ident, ok := call.Left.(*ast.Identifier); ok {
		name := ident.Symbol

		if structType, isStruct := symbolTable.GetType(name).(*types.TupleStruct); isStruct {
			return typeCheckTupleStructExpression(structType, call, symbolTable)
		}

		if builtin, ok := registry.Builtins[name]; ok {
			if len(builtin.Parameters) != len(call.Args) {
				if len(call.Args) < len(builtin.Parameters) {

					return types.Error(fmt.Sprintf("Missing argument for function %q", name), call)
				}
				return types.Error(fmt.Sprintf("Extra argument passed to function %q", name), call)
			}

			for i, param := range builtin.Parameters {
				arg := typeCheckExpression(call.Args[i], symbolTable)
				if arg.String() == "TypeError" {
					return arg
				}
				if !param.Valid(arg) {
					return types.Error(fmt.Sprintf("Invalid arguments passed to function %q: Type %q is not a valid argument for parameter of type %q", name, arg, param), call)
				}
			}

			return builtin.ReturnType
		}
	}

	callVar := typeCheckExpression(call.Left, symbolTable)
	if callVar.String() == "TypeError" {
		return callVar
	}

	function, ok := callVar.(*types.Function)

	if !ok {
		return types.Error(fmt.Sprintf("%q is not a function", call.Left.String()), call)
	}

	name := function.Name

	if len(function.Parameters) != len(call.Args) {
		if len(call.Args) < len(function.Parameters) {

			return types.Error(fmt.Sprintf("Missing argument for function %q", name), call)
		}
		return types.Error(fmt.Sprintf("Extra argument passed to function %q", name), call)
	}

	for i, param := range function.Parameters {
		arg := typeCheckExpression(call.Args[i], symbolTable)
		if arg.String() == "TypeError" {
			return arg
		}
		if !param.Valid(arg) {
			return types.Error(fmt.Sprintf("Invalid arguments passed to function %q: Type %q is not a valid argument for parameter of type %q", name, arg, param), call)
		}
	}

	return function.ReturnType
}

func typeCheckTupleStructExpression(tuple *types.TupleStruct, instance *ast.FunctionCall, symbolTable *symbols.SymbolTable) types.ValidType {
	if len(tuple.Members) != len(instance.Args) {
		return types.Error("Tuple struct expression incompatible with type", instance)
	}
	for i, arg := range instance.Args {
		argType := typeCheckExpression(arg, symbolTable)
		if argType.String() == "TypeError" {
			return argType
		}
		if !tuple.Members[i].Valid(argType) {
			return types.Error("Tuple struct expression incompatible with type", instance)
		}
	}

	return tuple
}

func typeCheckList(list *ast.ListLiteral, symbolTable *symbols.SymbolTable) types.ValidType {
	listTypes := []types.ValidType{}

	for _, elem := range list.Elements {
		elemType := typeCheckExpression(elem, symbolTable)
		if elemType.String() == "TypeError" {
			return elemType
		}
		newType := true
		for _, listType := range listTypes {
			if listType.Valid(elemType) {
				newType = false
				break
			}
		}

		if newType {
			listTypes = append(listTypes, elemType)
		}
	}

	return &types.ArrayLiteral{
		ElemType: types.MakeUnion(listTypes...),
		Length:   len(list.Elements),
		CanInfer: true,
	}
}

func typeCheckMap(maplit *ast.MapLiteral, symbolTable *symbols.SymbolTable) types.ValidType {
	keyTypes := []types.ValidType{}
	valueTypes := []types.ValidType{}

	for key, value := range maplit.Elements {
		keyType := typeCheckExpression(key, symbolTable)
		if keyType.String() == "TypeError" {
			return keyType
		}
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

		valueType := typeCheckExpression(value, symbolTable)
		if valueType.String() == "TypeError" {
			return valueType
		}

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

	return &types.MapLiteral{
		KeyType:   types.MakeUnion(keyTypes...),
		ValueType: types.MakeUnion(valueTypes...),
	}
}

func typeCheckIndexExpression(indexExpr *ast.IndexExpression, symbolTable *symbols.SymbolTable) types.ValidType {
	leftType := typeCheckExpression(indexExpr.Left, symbolTable)
	if leftType.String() == "TypeError" {
		return leftType
	}

	indexType := typeCheckExpression(indexExpr.Index, symbolTable)
	if indexType.String() == "TypeError" {
		return indexType
	}

	resultType := leftType.IndexBy(indexType)
	if resultType == nil {
		return types.Error(fmt.Sprintf("Type %q is not indexable with type %q", leftType.String(), indexType.String()), indexExpr)
	}

	return resultType
}

func typeCheckMemberExpression(memberExpr *ast.MemberExpression, symbolTable *symbols.SymbolTable) types.ValidType {
	leftType := typeCheckExpression(memberExpr.Left, symbolTable)
	if leftType.String() == "TypeError" {
		return leftType
	}

	resultType := types.Member(leftType, memberExpr.Member, memberExpr.IsNumberMember)
	if resultType == nil {
		return types.Error(fmt.Sprintf("Type %q does not have member %q", leftType.String(), memberExpr.Member), memberExpr)
	}

	return resultType
}

func typeCheckStructExpression(structExpr *ast.StructExpression, symbolTable *symbols.SymbolTable) types.ValidType {
	definedType := symbolTable.GetType(structExpr.Name)
	if definedType.String() == "TypeError" {
		return types.Error(fmt.Sprintf("Struct %q is undefined", structExpr.Name), structExpr)
	}

	members := map[string]types.ValidType{}

	for name, member := range structExpr.Members {
		dataType := typeCheckExpression(member, symbolTable)
		if dataType.String() == "TypeError" {
			return dataType
		}

		members[name] = dataType
	}

	structType := &types.Struct{
		Name:    structExpr.Name,
		Members: members,
	}

	if !definedType.Valid(structType) {
		return types.Error("Struct expression incompatiable with type", structExpr)
	}

	return structType
}

func typeCheckTuple(tuple *ast.TupleExpression, symbolTable *symbols.SymbolTable) types.ValidType {
	members := []types.ValidType{}
	for _, member := range tuple.Members {
		memberType := typeCheckExpression(member, symbolTable)
		if memberType.String() == "TypeError" {
			return memberType
		}
		members = append(members, memberType)
	}

	return &types.Tuple{Members: members}
}
