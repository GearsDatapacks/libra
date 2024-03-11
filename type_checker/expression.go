package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/ir"
)

func (t *typeChecker) typeCheckExpression(expression ast.Expression) ir.Expression {
	switch expr := expression.(type) {
	case *ast.IntegerLiteral:
		return &ir.IntegerLiteral{
			Value: expr.Value,
		}
	case *ast.FloatLiteral:
		return &ir.FloatLiteral{
			Value: expr.Value,
		}
	case *ast.BooleanLiteral:
		return &ir.BooleanLiteral{
			Value: expr.Value,
		}
	case *ast.StringLiteral:
		return &ir.StringLiteral{
			Value: expr.Value,
		}
	default:
		panic(fmt.Sprintf("TODO: Type-check %T", expr))
	}
}
