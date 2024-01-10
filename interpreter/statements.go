package interpreter

import (
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
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

	return manager.Env.DeclareVariable(varDec.Name, varDec.DataType.GetType(), value)
}

func registerFunctionDeclaration(funcDec *ast.FunctionDeclaration, manager *modules.ModuleManager) values.RuntimeValue {
	params := []values.Parameter{}

	for _, param := range funcDec.Parameters {
		params = append(params, values.Parameter{
			Name: param.Name,
			Type: param.Type.GetType(),
		})
	}

	functionType := funcDec.GetType().(*types.Function)

	fn := &values.FunctionValue{
		Name:       funcDec.Name,
		Parameters: params,
		Env:        manager.Env,
		Manager:    manager,
		Body:       funcDec.Body,
		BaseValue:  values.BaseValue{DataType: functionType},
	}

	if funcDec.MethodOf != nil {
		parentType := funcDec.MethodOf.GetType()
		functionType.MethodOf = parentType

		types.AddMethod(funcDec.Name, functionType)
		environment.AddMethod(funcDec.Name, fn)
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

	if importStatement.ImportAll {
		for name, variable := range mod.Env.Exports {
			manager.Env.DeclareVariable(name, variable.Type(), variable)
		}
		return values.MakeNull()
	}

	if importStatement.ImportedSymbols != nil {
		for _, symbol := range importStatement.ImportedSymbols {
			value, ok := mod.Env.Exports[symbol]
			if ok {
				manager.Env.DeclareVariable(symbol, value.Type(), value)
			}
		}
	}

	name := mod.Name
	if importStatement.Alias != "" {
		name = importStatement.Alias
	}

	importedMod := &values.Module{
		Name:    name,
		Exports: mod.Env.Exports,
		BaseValue: values.BaseValue{
			DataType: importStatement.GetType(),
		},
	}

	manager.Env.DeclareVariable(name, importedMod.DataType, importedMod)
	return importedMod
}

func evaluateUnitStructDeclaration(structDecl *ast.UnitStructDeclaration, manager *modules.ModuleManager) values.RuntimeValue {
	unit := values.MakeUnitStruct(structDecl.Name, structDecl.GetType().(*types.UnitStruct))
	manager.Env.DeclareVariable(structDecl.Name, unit.DataType, unit)

	if structDecl.IsExport() {
		manager.Env.Exports[structDecl.Name] = unit
	}

	return values.MakeNull()
}

func evaluateEnumDeclaration(enumDec *ast.EnumDeclaration, manager *modules.ModuleManager) values.RuntimeValue {
	members := map[string]values.RuntimeValue{}
	for name, member := range enumDec.GetType().(*types.Enum).Types {
		if unitType, ok := member.DataType.(*types.Type).DataType.(*types.UnitStruct); ok {
			unit := values.MakeUnitStruct(name, unitType)
			members[name] = unit
		}
	}

	enum := &values.Enum{
		Name:      enumDec.Name,
		Members:   members,
		BaseValue: values.BaseValue{DataType: enumDec.GetType()},
	}

	manager.Env.DeclareVariable(enum.Name, enumDec.GetType(), enum)
	if enumDec.IsExport() {
		manager.Env.Exports[enumDec.Name] = enum
	}

	return enum
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
