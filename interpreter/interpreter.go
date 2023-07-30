package interpreter

import (
	"log"

	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func Evaluate(program ast.Program) values.RuntimeValue {
	Register()

	var lastValue values.RuntimeValue

	for _, statement := range program.Body {
		lastValue = evaluate(statement)
	}

	return lastValue
}

func evaluate(astNode ast.Statement) values.RuntimeValue {
	switch statement := astNode.(type) {
	case *ast.ExpressionStatement:
		return evaluateExpressionStatement(*statement)

	default:
		log.Fatalf("Unreconised AST node: %s", astNode.String())
		return &values.IntegerLiteral{}
	}
}
