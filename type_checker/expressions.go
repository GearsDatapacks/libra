package typechecker

import (
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
	default:
		return types.INVALID
	}
}
