package interpreter

import (
	"log"

	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func Evaluate(program ast.Program, env *environment.Environment) values.RuntimeValue {
	Register()

	var lastValue values.RuntimeValue

	for _, statement := range program.Body {
		lastValue = evaluate(statement, env)
	}

	return lastValue
}

func evaluate(astNode ast.Statement, env *environment.Environment) values.RuntimeValue {
	switch statement := astNode.(type) {
	case *ast.ExpressionStatement:
		return evaluateExpressionStatement(*statement, env)
	
	case *ast.VariableDeclaration:
		return evaluateVariableDeclaration(*statement, env)

	default:
		log.Fatalf("Unreconised AST node: %s", astNode.String())
		return &values.IntegerLiteral{}
	}
}
