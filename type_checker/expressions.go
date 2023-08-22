package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckExpression(expr ast.Expression, symbolTable *symbols.SymbolTable) types.ValidType {
	switch expression := expr.(type) {
	case *ast.IntegerLiteral:
		return types.MakeLiteral(types.INT)
	case *ast.FloatLiteral:
		return types.MakeLiteral(types.FLOAT)
	case *ast.NullLiteral:
		return types.MakeLiteral(types.NULL)
	case *ast.BooleanLiteral:
		return types.MakeLiteral(types.BOOL)
	case *ast.Identifier:
		return symbolTable.GetSymbol(expression.Symbol)
	case *ast.BinaryOperation:
		return typeCheckBinaryOperation(expression, symbolTable)
	case *ast.AssignmentExpression:
		return typeCheckAssignmentExpression(expression, symbolTable)
	case *ast.FunctionCall:
		return typeCheckFunctionCall(expression, symbolTable)
	default:
		errors.DevError("Unexpected expression type: " + expr.String())
		return &types.Literal{}
	}
}

func typeCheckAssignmentExpression(assignment *ast.AssignmentExpression, symbolTable *symbols.SymbolTable) types.ValidType {
	if assignment.Assignee.Type() != "Identifier" {
		errors.TypeError("Can only assign values to variables", assignment)
	}

	symbolName := assignment.Assignee.(*ast.Identifier).Symbol

	if symbolTable.IsConstant(symbolName) {
		errors.TypeError("Cannot reassign constant " + symbolName, assignment)
	}

	dataType := symbolTable.GetSymbol(symbolName)

	expressionType := typeCheckExpression(assignment.Value, symbolTable)
	correctType := dataType == expressionType

	if correctType {
		return dataType
	}

	errors.TypeError(fmt.Sprintf("Type %q is not assignable to type %q", expressionType, dataType), assignment)
	return &types.Literal{}
}

func typeCheckFunctionCall(call *ast.FunctionCall, symbolTable *symbols.SymbolTable) types.ValidType {
	if !symbolTable.Exists(call.Name) {
		errors.TypeError(fmt.Sprintf("Function %q is undefined", call.Name), call)
	}

	callVar := symbolTable.GetSymbol(call.Name)

	if !callVar.Valid(&types.Function{}) {
		errors.TypeError(fmt.Sprintf("Variable %q is not a function", call.Name), call)
	}

	function := callVar.(*types.Function)

	if len(function.Parameters) != len(call.Args) {
		errors.TypeError(fmt.Sprintf("Invalid arguments passed to function %q", call.Name), call)
	}

	for i, param := range function.Parameters {
		arg := typeCheckExpression(call.Args[i], symbolTable)
		if !param.Valid(arg) {
			errors.TypeError(fmt.Sprintf("Invalid arguments passed to function %q", call.Name), call)
		}
	}

	return function.ReturnType
}
