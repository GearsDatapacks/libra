package interpreter

import (
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func evaluateExpressionStatement(exprStmt ast.ExpressionStatement) values.RuntimeValue {
	return evaluateExpression(exprStmt.Expression)
}
