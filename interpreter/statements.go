package interpreter

import (
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func evaluateExpressionStatement(exprStmt *ast.ExpressionStatement, manager *modules.ModuleManager) values.RuntimeValue {
	return evaluateExpression(exprStmt.Expression, manager)
}

func evaluateVariableDeclaration(varDec *ast.VariableDeclaration, manager *modules.ModuleManager) values.RuntimeValue {
	var value values.RuntimeValue

	if varDec.Value == nil {
		value = values.GetZeroValue(varDec.DataType.String())
	} else {
		value = evaluateExpression(varDec.Value, manager)
	}

	return manager.Env.DeclareVariable(varDec.Name, typechecker.TypeCheckType(varDec.DataType, manager), value)
}

func registerFunctionDeclaration(funcDec *ast.FunctionDeclaration, manager *modules.ModuleManager) values.RuntimeValue {
	params := []values.Parameter{}
	paramTypes := []types.ValidType{}

	for _, param := range funcDec.Parameters {
		params = append(params, values.Parameter{
			Name:  param.Name,
			Type:  typechecker.TypeCheckType(param.Type, manager),
		})
	}
	returnType := typechecker.TypeCheckType(funcDec.ReturnType, manager)

	functionType := &types.Function{
		Parameters: paramTypes,
		ReturnType: returnType,
		Name:       funcDec.Name,
	}

	fn := &values.FunctionValue{
		Name:       funcDec.Name,
		Parameters: params,
		Env:        manager.Env,
		Manager:    manager,
		Body:       funcDec.Body,
		BaseValue:  values.BaseValue{DataType: functionType},
	}

	if funcDec.MethodOf != nil {
		parentType := typechecker.TypeCheckType(funcDec.MethodOf, manager)
		functionType.MethodOf = parentType

		types.AddMethod(funcDec.Name, functionType)
		manager.Env.AddMethod(funcDec.Name, fn)
		return fn
	} else {

		if funcDec.IsExport() {
			manager.Env.GlobalScope().Exports[fn.Name] = fn
		}
		return manager.Env.DeclareVariable(funcDec.Name, fn.Type(), fn)
	}
}

func evaluateReturnStatement(ret *ast.ReturnStatement, manager *modules.ModuleManager) values.RuntimeValue {
	value := evaluateExpression(ret.Value, manager)
	functionScope := manager.Env.FindFunctionScope()
	functionScope.ReturnValue = value
	return value
}

func evaluateIfStatement(ifStatement *ast.IfStatement, manager *modules.ModuleManager) values.RuntimeValue {
	condition := evaluateExpression(ifStatement.Condition, manager)

	if !condition.Truthy() {
		if ifStatement.Else == nil {
			return values.MakeNull()
		}
		if elseStatement, isElse := ifStatement.Else.(*ast.ElseStatement); isElse {
			return evaluateElseStatement(elseStatement, manager)
		}
		if nextIf, isIf := ifStatement.Else.(*ast.IfStatement); isIf {
			return evaluateIfStatement(nextIf, manager)
		}
	}

	newScope := environment.NewChild(manager.Env, environment.GENERIC_SCOPE)
	manager.EnterEnv(newScope)

	for _, statement := range ifStatement.Body {
		evaluate(statement, manager)
	}

	manager.ExitEnv()

	return values.MakeNull()
}

func evaluateElseStatement(elseStatement *ast.ElseStatement, manager *modules.ModuleManager) values.RuntimeValue {
	newScope := environment.NewChild(manager.Env, environment.GENERIC_SCOPE)
	manager.EnterEnv(newScope)

	for _, statement := range elseStatement.Body {
		evaluate(statement, manager)
	}
	manager.ExitEnv()

	return values.MakeNull()
}

func evaluateWhileLoop(while *ast.WhileLoop, manager *modules.ModuleManager) values.RuntimeValue {
	for evaluateExpression(while.Condition, manager).Truthy() {
		newEnv := environment.NewChild(manager.Env, environment.GENERIC_SCOPE)
		manager.EnterEnv(newEnv)

		for _, statement := range while.Body {
			evaluate(statement, manager)
		}
		manager.ExitEnv()
	}

	return values.MakeNull()
}

func evaluateForLoop(forLoop *ast.ForLoop, manager *modules.ModuleManager) values.RuntimeValue {
	loopEnv := environment.NewChild(manager.Env, environment.GENERIC_SCOPE)
	manager.EnterEnv(loopEnv)

	evaluate(forLoop.Initial, manager)

	for evaluateExpression(forLoop.Condition, manager).Truthy() {
		newEnv := environment.NewChild(loopEnv, environment.GENERIC_SCOPE)
		manager.EnterEnv(newEnv)

		for _, statement := range forLoop.Body {
			evaluate(statement, manager)
		}
		manager.ExitEnv()
		evaluate(forLoop.Update, manager)
	}

	manager.ExitEnv()

	return values.MakeNull()
}

func evaluateImportStatement(importStatement *ast.ImportStatement, manager *modules.ModuleManager) values.RuntimeValue {
	modPath := importStatement.Module
	mod := manager.Imported[modPath]

	importedMod := &values.Module{
		Name:    mod.Name,
		Exports: mod.Env.Exports,
		BaseValue: values.BaseValue{
			DataType: &types.Module{
				Name: mod.Name,
				Exports: mod.SymbolTable.Exports,
			},
		},
	}

	manager.Env.DeclareVariable(mod.Name, importedMod.DataType, importedMod)
	return importedMod
}

/*
func evaluateStructDeclaration(structDecl *ast.StructDeclaration, manager *modules.ModuleManager) values.RuntimeValue {
	members := map[string]types.ValidType{}

	for memberName, memberType := range structDecl.Members {
		dataType := typechecker.TypeCheckType(memberType, manager)

		members[memberName] = dataType
	}

	structType := &types.Struct{
		Name:    structDecl.Name,
		Members: members,
	}
	manager.Env.AddType(structDecl.Name, structType)

	return values.MakeNull()
}

func evaluateInterfaceDeclaration(intDecl *ast.InterfaceDeclaration, manager *modules.ModuleManager) values.RuntimeValue {
	members := map[string]types.ValidType{}

	for _, member := range intDecl.Members {
		if !member.IsFunction {
			dataType := typechecker.TypeCheckType(member.ResultType, manager)

			members[member.Name] = dataType
			continue
		}

		fnType := &types.Function{}
		fnType.Name = member.Name

		returnType := typechecker.TypeCheckType(member.ResultType, manager)
		fnType.ReturnType = returnType

		fnType.Parameters = []types.ValidType{}
		for _, param := range member.Parameters {
			paramType := typechecker.TypeCheckType(param, manager)

			fnType.Parameters = append(fnType.Parameters, paramType)
		}

		members[member.Name] = fnType
	}

	interfaceType := &types.Interface{
		Name:    intDecl.Name,
		Members: members,
	}

	manager.Env.AddType(intDecl.Name, interfaceType)
	return values.MakeNull()
}

func evaluateTupleStructDeclaration(structDecl *ast.TupleStructDeclaration, manager *modules.ModuleManager) values.RuntimeValue {
	members := []types.ValidType{}

	for _, memberType := range structDecl.Members {
		dataType := typechecker.TypeCheckType(memberType, manager)

		members = append(members, dataType)
	}

	structType := &types.TupleStruct{
		Name:    structDecl.Name,
		Members: members,
	}
	manager.Env.AddType(structDecl.Name, structType)

	return values.MakeNull()
}
*/
