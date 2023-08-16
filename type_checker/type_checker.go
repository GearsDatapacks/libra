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

func typeCheck(stmt ast.Statement) types.Type {
	switch statement := stmt.(type) {
	case *ast.VariableDeclaration:
		expressionType := typeCheckExpression(statement.Value)
		dataType := types.FromString(statement.DataType)
		correctType := dataType == expressionType
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

func valid(t types.Type) bool {
	return t != types.INVALID
}
