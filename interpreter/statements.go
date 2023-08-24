package interpreter

import (
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func evaluateExpressionStatement(exprStmt *ast.ExpressionStatement, env *environment.Environment) values.RuntimeValue {
	return evaluateExpression(exprStmt.Expression, env)
}

func evaluateVariableDeclaration(varDec *ast.VariableDeclaration, env *environment.Environment) values.RuntimeValue {
	var value values.RuntimeValue

	if varDec.Value == nil {
		value = values.GetZeroValue(varDec.DataType.String())
	} else {
		value = evaluateExpression(varDec.Value, env)
	}

	return env.DeclareVariable(varDec.Name, value)
}

func evaluateFunctionDeclaration(funcDec *ast.FunctionDeclaration, env *environment.Environment) values.RuntimeValue {
	params := []string{}

	// Only need names
	for _, param := range funcDec.Parameters {
		params = append(params, param.Name)
	}

	fn := &values.FunctionValue{
		Name:                   funcDec.Name,
		Parameters:             params,
		DeclarationEnvironment: env,
		Body:                   funcDec.Body,
	}

	return env.DeclareVariable(funcDec.Name, fn)
}

func evaluateReturnStatement(ret *ast.ReturnStatement, env *environment.Environment) values.RuntimeValue {
	value := evaluateExpression(ret.Value, env)
	functionScope := env.FindFunctionScope()
	functionScope.ReturnValue = value
	return value
}

func evaluateIfStatement(ifStatement *ast.IfStatement, env *environment.Environment) values.RuntimeValue {
	condition := evaluateExpression(ifStatement.Condition, env)

	if !condition.Truthy() {
		if ifStatement.Condition == nil { return values.MakeNull() }
		if elseStatement, isElse := ifStatement.Else.(*ast.ElseStatement); isElse {
			return evaluateElseStatement(elseStatement, env)
		}
		if nextIf, isIf := ifStatement.Else.(*ast.IfStatement); isIf {
			return evaluateIfStatement(nextIf, env)
		}
	}

	newScope := environment.NewChild(env, environment.GENERIC_SCOPE)

	for _, statement := range ifStatement.Body {
		evaluate(statement, newScope)
	}

	return values.MakeNull()
}

func evaluateElseStatement(elseStatement *ast.ElseStatement, env *environment.Environment) values.RuntimeValue {
	newScope := environment.NewChild(env, environment.GENERIC_SCOPE)

	for _, statement := range elseStatement.Body {
		evaluate(statement, newScope)
	}

	return values.MakeNull()
}
