package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckExpression(expr ast.Expression, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	switch expression := expr.(type) {
	case *ast.IntegerLiteral:
		return &types.IntLiteral{}, nil
	case *ast.FloatLiteral:
		return &types.FloatLiteral{}, nil
	case *ast.StringLiteral:
		return &types.StringLiteral{}, nil
	case *ast.NullLiteral:
		return &types.NullLiteral{}, nil
	case *ast.BooleanLiteral:
		return &types.BoolLiteral{}, nil
	case *ast.VoidValue:
		return &types.Void{}, nil

	case *ast.Identifier:
		return symbolTable.GetSymbol(expression.Symbol)

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

	case *ast.IndexExpression:
		return typeCheckIndexExpression(expression, symbolTable)

	default:
		return nil, errors.DevError("(Type checker) Unexpected expression type: " + expr.String())
	}
}

func typeCheckAssignmentExpression(assignment *ast.AssignmentExpression, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	if assignment.Assignee.Type() != "Identifier" {
		return nil, errors.TypeError("Can only assign values to variables", assignment)
	}

	symbolName := assignment.Assignee.(*ast.Identifier).Symbol

	if symbolTable.IsConstant(symbolName) {
		return nil, errors.TypeError("Cannot reassign constant "+symbolName, assignment)
	}

	dataType, err := symbolTable.GetSymbol(symbolName)
	if err != nil {
		return nil, err
	}

	expressionType, err := typeCheckExpression(assignment.Value, symbolTable)
	if err != nil {
		return nil, err
	}
	correctType := dataType.Valid(expressionType)

	if correctType {
		return dataType, nil
	}

	return nil, errors.TypeError(fmt.Sprintf("Type %q is not assignable to type %q", expressionType, dataType), assignment)
}

func typeCheckFunctionCall(call *ast.FunctionCall, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	if builtin, ok := registry.Builtins[call.Name]; ok {
		if len(builtin.Parameters) != len(call.Args) {
			return nil, errors.TypeError(fmt.Sprintf("Invalid arguments passed to function %q", call.Name), call)
		}

		for i, param := range builtin.Parameters {
			arg, err := typeCheckExpression(call.Args[i], symbolTable)
			if err != nil {
				return nil, err
			}
			if !param.Valid(arg) {
				return nil, errors.TypeError(fmt.Sprintf("Invalid arguments passed to function %q", call.Name), call)
			}
		}

		return builtin.ReturnType, nil
	}

	if !symbolTable.Exists(call.Name) {
		return nil, errors.TypeError(fmt.Sprintf("Function %q is undefined", call.Name), call)
	}

	callVar, err := symbolTable.GetSymbol(call.Name)
	if err != nil {
		return nil, err
	}

	function, ok := callVar.(*types.Function)

	if !ok {
		return nil, errors.TypeError(fmt.Sprintf("Variable %q is not a function", call.Name), call)
	}

	if len(function.Parameters) != len(call.Args) {
		return nil, errors.TypeError(fmt.Sprintf("Invalid arguments passed to function %q", call.Name), call)
	}

	for i, param := range function.Parameters {
		arg, err := typeCheckExpression(call.Args[i], symbolTable)
		if err != nil {
			return nil, err
		}
		if !param.Valid(arg) {
			return nil, errors.TypeError(fmt.Sprintf("Invalid arguments passed to function %q", call.Name), call)
		}
	}

	return function.ReturnType, nil
}

func typeCheckList(list *ast.ListLiteral, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	listTypes := []types.ValidType{}

	for _, elem := range list.Elements {
		elemType, err := typeCheckExpression(elem, symbolTable)
		if err != nil {
			return nil, err
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
	}, nil
}

func typeCheckIndexExpression(indexExpr *ast.IndexExpression, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	leftType, err := typeCheckExpression(indexExpr.Left, symbolTable)
	if err != nil {
		return nil, err
	}

	indexType, err := typeCheckExpression(indexExpr.Index, symbolTable)
	if err != nil {
		return nil, err
	}

	resultType, indexable := leftType.Indexable(indexType)
	if !indexable {
		return nil, errors.TypeError(fmt.Sprintf("Type %q is not indexable with type %q", leftType.String(), indexType.String()), indexExpr)
	}

	return resultType, nil
}
