package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TypeCheck(program ast.Program) {
	symbolTable := NewSymbolTable()

	for _, stmt := range program.Body {
		typeCheck(stmt, symbolTable)
	}
}

func typeCheck(stmt ast.Statement, symbolTable *SymbolTable) types.DataType {
	switch statement := stmt.(type) {
	case *ast.VariableDeclaration:
		return typeCheckVariableDeclaration(statement, symbolTable)

	case *ast.ExpressionStatement:
		return typeCheckExpression(statement.Expression, symbolTable)

	default:
		errors.DevError("Unexpected statment type")
		return types.INT
	}
}

func typeCheckVariableDeclaration(varDec *ast.VariableDeclaration, symbolTable *SymbolTable) types.DataType {
	expressionType := typeCheckExpression(varDec.Value, symbolTable)
	
	// Blank if type to be inferred
	if varDec.DataType == "" {
		symbolTable.RegisterSymbol(varDec.Name, expressionType, varDec.Constant)
		return expressionType
	}
	
	dataType := types.FromString(varDec.DataType)
	correctType := dataType == expressionType

	if correctType {
		symbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
		return dataType
	}

	errors.TypeError(fmt.Sprintf("Type %q is not assignable to type %q", expressionType, dataType))
	return types.INT
}
