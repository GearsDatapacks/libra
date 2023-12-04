package interpreter

import (
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
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

func registerFunctionDeclaration(funcDec *ast.FunctionDeclaration, env *environment.Environment) values.RuntimeValue {
	params := []string{}
	paramTypes := []types.ValidType{}

	for _, param := range funcDec.Parameters {
		params = append(params, param.Name)
		paramTypes = append(paramTypes, types.FromAst(param.Type, env))
	}
	returnType := types.FromAst(funcDec.ReturnType, env)

	functionType := &types.Function{
		Parameters: paramTypes,
		ReturnType: returnType,
		Name:       funcDec.Name,
	}

	fn := &values.FunctionValue{
		Name:                   funcDec.Name,
		Parameters:             params,
		DeclarationEnvironment: env,
		Body:                   funcDec.Body,
		BaseValue:              values.BaseValue{DataType: functionType},
	}

	if funcDec.MethodOf != nil {
		parentType := types.FromAst(funcDec.MethodOf, env)
		functionType.MethodOf = parentType

		types.AddMethod(funcDec.Name, functionType)
		env.AddMethod(funcDec.Name, fn)
		return fn
	} else {
		return env.DeclareVariable(funcDec.Name, fn)
	}
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
		if ifStatement.Else == nil {
			return values.MakeNull()
		}
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

func evaluateWhileLoop(while *ast.WhileLoop, env *environment.Environment) values.RuntimeValue {
	for evaluateExpression(while.Condition, env).Truthy() {
		newEnv := environment.NewChild(env, environment.GENERIC_SCOPE)

		for _, statement := range while.Body {
			evaluate(statement, newEnv)
		}
	}

	return values.MakeNull()
}

func evaluateForLoop(forLoop *ast.ForLoop, env *environment.Environment) values.RuntimeValue {
	loopEnv := environment.NewChild(env, environment.GENERIC_SCOPE)

	evaluate(forLoop.Initial, loopEnv)

	for evaluateExpression(forLoop.Condition, loopEnv).Truthy() {
		newEnv := environment.NewChild(loopEnv, environment.GENERIC_SCOPE)

		for _, statement := range forLoop.Body {
			evaluate(statement, newEnv)
		}
		evaluate(forLoop.Update, loopEnv)
	}

	return values.MakeNull()
}

func evaluateStructDeclaration(structDecl *ast.StructDeclaration, env *environment.Environment) values.RuntimeValue {
	members := map[string]types.ValidType{}

	for memberName, memberType := range structDecl.Members {
		dataType := types.FromAst(memberType, env)

		members[memberName] = dataType
	}

	structType := &types.Struct{
		Name:    structDecl.Name,
		Members: members,
	}
	env.AddType(structDecl.Name, structType)

	return values.MakeNull()
}

func evaluateTupleStructDeclaration(structDecl *ast.TupleStructDeclaration, env *environment.Environment) values.RuntimeValue {
	members := []types.ValidType{}

	for _, memberType := range structDecl.Members {
		dataType := types.FromAst(memberType, env)

		members = append(members, dataType)
	}

	structType := &types.TupleStruct{
		Name:    structDecl.Name,
		Members: members,
	}
	env.AddType(structDecl.Name, structType)

	return values.MakeNull()
}
