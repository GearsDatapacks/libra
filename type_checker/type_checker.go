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
		return typeCheckVariableDeclaration(statement, symbolTable)

	case *ast.ExpressionStatement:
		return typeCheckExpression(statement.Expression, symbolTable)

	default:
		return types.INVALID
	}
}

func typeCheckVariableDeclaration(varDex *ast.VariableDeclaration, symbolTable *SymbolTable) types.DataType {
	expressionType := typeCheckExpression(varDex.Value, symbolTable)
	dataType := types.FromString(varDex.DataType)
	correctType := dataType == expressionType

	// Blank if type to be inferred
	if varDex.DataType == "" {
		symbolTable.RegisterSymbol(varDex.Name, expressionType, varDex.Constant)
		return expressionType
	}

	if correctType {
		symbolTable.RegisterSymbol(varDex.Name, dataType, varDex.Constant)
		return dataType
	}

	return types.INVALID
}

func valid(t types.DataType) bool {
	return t != types.INVALID
}
