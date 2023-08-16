package typechecker

import (
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckExpression(expr ast.Expression) types.Type {
	switch expression := expr.(type) {
	case *ast.IntegerLiteral:
		return types.INT
	case *ast.NullLiteral:
		return types.NULL
	case *ast.BooleanLiteral:
		return types.BOOL
	case *ast.Identifier:
		return types.NULL
	case *ast.BinaryOperation:
		validTypes := valid(typeCheckExpression(expression.Left)) && valid(typeCheckExpression(expression.Right))
		if !validTypes {
			return types.INVALID
		}
		return types.INT
	default:
		return types.INVALID
	}
}
