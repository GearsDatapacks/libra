package interpreter

import (
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func evaluateExpressionStatement(exprStmt ast.ExpressionStatement, env *environment.Environment) values.RuntimeValue {
	return evaluateExpression(exprStmt.Expression, env)
}
