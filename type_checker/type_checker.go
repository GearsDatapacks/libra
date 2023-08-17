package typechecker

import (
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TypeCheck(program ast.Program) bool {
	symbolTable := NewSymbolTable()

	for _, stmt := range program.Body {
		if !valid(typeCheck(stmt, symbolTable)) {
			return false
		}
	}

	return true
}

func typeCheck(stmt ast.Statement, symbolTable *SymbolTable) types.DataType {
	switch statement := stmt.(type) {
	case *ast.VariableDeclaration:
		expressionType := typeCheckExpression(statement.Value, symbolTable)
		dataType := types.FromString(statement.DataType)
		correctType := dataType == expressionType

		// Blank if type to be inferred
		if statement.DataType == "" {
			symbolTable.RegisterSymbol(statement.Name, expressionType)
			return expressionType
		}

		if correctType {
			symbolTable.RegisterSymbol(statement.Name, dataType)
			return dataType
		}

		return types.INVALID

	case *ast.ExpressionStatement:
		return typeCheckExpression(statement.Expression, symbolTable)

	default:
		return types.INVALID
	}
}

func valid(t types.DataType) bool {
	return t != types.INVALID
}
