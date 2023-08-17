package typechecker

import (
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TypeCheck(program ast.Program) bool {
	for _, stmt := range program.Body {
		if !valid(typeCheck(stmt)) {
			return false
		}
	}

	return true
}

func typeCheck(stmt ast.Statement) types.DataType {
	switch statement := stmt.(type) {
	case *ast.VariableDeclaration:
		expressionType := typeCheckExpression(statement.Value)
		dataType := types.FromString(statement.DataType)
		correctType := dataType == expressionType

		// Blank if type to be inferred
		if statement.DataType == "" {
			return expressionType
		}

		if correctType {
			return dataType
		}

		return types.INVALID

	case *ast.ExpressionStatement:
		return typeCheckExpression(statement.Expression)

	default:
		return types.INVALID
	}
}

func valid(t types.DataType) bool {
	return t != types.INVALID
}
