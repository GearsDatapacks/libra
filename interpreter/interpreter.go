package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func Evaluate(program ast.Program, env *environment.Environment) values.RuntimeValue {
	var lastValue values.RuntimeValue

	for _, statement := range program.Body {
		register(statement, env)
	}

	for _, statement := range program.Body {
		lastValue = evaluate(statement, env)
	}

	return lastValue
}

func evaluate(astNode ast.Statement, env *environment.Environment) values.RuntimeValue {
	switch statement := astNode.(type) {
	case *ast.ExpressionStatement:
		return evaluateExpressionStatement(statement, env)

	case *ast.VariableDeclaration:
		return evaluateVariableDeclaration(statement, env)

	case *ast.FunctionDeclaration:
		// return registerFunctionDeclaration(statement, env)
		return values.MakeNull()

	case *ast.ReturnStatement:
		return evaluateReturnStatement(statement, env)

	case *ast.IfStatement:
		return evaluateIfStatement(statement, env)

	case *ast.WhileLoop:
		return evaluateWhileLoop(statement, env)

	case *ast.ForLoop:
		return evaluateForLoop(statement, env)

	case *ast.StructDeclaration:
		// return evaluateStructDeclaration(statement, env)
		return values.MakeNull()
	
	case *ast.TupleStructDeclaration:
		return values.MakeNull()

	case *ast.InterfaceDeclaration:
		return values.MakeNull()

	case *ast.TypeDeclaration:
		// env.AddType(statement.Name, types.FromAst(statement.DataType, env))
		return values.MakeNull()

	default:
		errors.LogError(errors.DevError(fmt.Sprintf("(Interpreter) Unreconised AST node: %s", astNode.String()), astNode))
		return nil
	}
}

func register(astNode ast.Statement, env *environment.Environment) {
	switch statement := astNode.(type) {
	case *ast.ExpressionStatement:
	case *ast.FunctionDeclaration:
		registerFunctionDeclaration(statement, env)

	case *ast.StructDeclaration:
		evaluateStructDeclaration(statement, env)

	case *ast.InterfaceDeclaration:
		evaluateInterfaceDeclaration(statement, env)

	case *ast.TupleStructDeclaration:
		evaluateTupleStructDeclaration(statement, env)

	case *ast.TypeDeclaration:
		env.AddType(statement.Name, types.FromAst(statement.DataType, env))
		// return values.MakeNull()
	}
}
