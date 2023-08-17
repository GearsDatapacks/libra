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
		if expression.Assignee.Type() != "Identifier" {
			errors.TypeError("Can only assign values to variables")
		}

		dataType := symbolTable.GetSymbol(expression.Assignee.(*ast.Identifier).Symbol)

		expressionType := typeCheckExpression(expression.Value, symbolTable)
		correctType := dataType == expressionType

		if correctType {
			return dataType
		}

		return types.INVALID

	default:
		return types.INVALID
	}
}
