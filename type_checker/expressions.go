package typechecker

import (
	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckExpression(expr ast.Expression, symbolTable *SymbolTable) types.DataType {
	switch expression := expr.(type) {
	case *ast.IntegerLiteral:
		return types.INT
	case *ast.NullLiteral:
		return types.NULL
	case *ast.BooleanLiteral:
		return types.BOOL
	case *ast.Identifier:
		return symbolTable.GetSymbol(expression.Symbol)
	case *ast.BinaryOperation:
		return typeCheckBinaryOperation(expression, symbolTable)
	case *ast.AssignmentExpression:
		return typeCheckAssignmentExpression(expression, symbolTable)
	default:
		return types.INVALID
	}
}

func typeCheckAssignmentExpression(assignment *ast.AssignmentExpression, symbolTable *SymbolTable) types.DataType {
	if assignment.Assignee.Type() != "Identifier" {
		errors.TypeError("Can only assign values to variables")
	}

	symbolName := assignment.Assignee.(*ast.Identifier).Symbol

	if symbolTable.IsConstant(symbolName) {
		return types.INVALID
	}

	dataType := symbolTable.GetSymbol(symbolName)

	expressionType := typeCheckExpression(assignment.Value, symbolTable)
	correctType := dataType == expressionType

	if correctType {
		return dataType
	}

	return types.INVALID
}
