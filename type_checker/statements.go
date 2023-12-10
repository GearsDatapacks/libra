package typechecker

import (
	"fmt"
	"log"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckStatement(stmt ast.Statement, manager *modules.ModuleManager) types.ValidType {
	switch statement := stmt.(type) {
	case *ast.VariableDeclaration:
		return typeCheckVariableDeclaration(statement, manager)

	case *ast.ExpressionStatement:
		return typeCheckExpression(statement.Expression, manager)

	case *ast.FunctionDeclaration:
		return typeCheckFunctionDeclaration(statement, manager)

	case *ast.ReturnStatement:
		return typeCheckReturnStatement(statement, manager)

	case *ast.IfStatement:
		return typeCheckIfStatement(statement, manager)

	case *ast.WhileLoop:
		return typeCheckWhileLoop(statement, manager)

	case *ast.ForLoop:
		return typeCheckForLoop(statement, manager)

	case *ast.StructDeclaration:
		// return typeCheckStructDeclaration(statement, manager)
		return &types.Void{}

	case *ast.TupleStructDeclaration:
		// return typeCheckStructDeclaration(statement, manager)
		return &types.Void{}

	case *ast.InterfaceDeclaration:
		// return typeCheckInterfaceDeclaration(statement, manager)
		return &types.Void{}

	case *ast.TypeDeclaration:
		// return typeCheckTypeDeclataion(statement, manager)
		return &types.Void{}

	case *ast.ImportStatement:
		return &types.Void{}

	default:
		log.Fatal(errors.DevError("(Type checker) Unexpected statement type: " + statement.String()))
		return nil
	}
}

func typeCheckGlobalStatement(stmt ast.Statement, manager *modules.ModuleManager) types.ValidType {
	switch statement := stmt.(type) {
	case *ast.FunctionDeclaration:
		return registerFunctionDeclaration(statement, manager)

	case *ast.StructDeclaration:
		return typeCheckStructDeclaration(statement, manager)

	case *ast.TupleStructDeclaration:
		return typeCheckTupleStructDeclaration(statement, manager)

	case *ast.InterfaceDeclaration:
		return typeCheckInterfaceDeclaration(statement, manager)

	case *ast.TypeDeclaration:
		return typeCheckTypeDeclataion(statement, manager)
	default:
		return &types.Void{}
	}
}

func registerTypeStatement(stmt ast.Statement, manager *modules.ModuleManager) types.ValidType {
	switch statement := stmt.(type) {
	// case *ast.ImportStatement:
	// 	return typeCheckImportStatement(statement, manager)

	case *ast.StructDeclaration:
		return registerStructDeclaration(statement, manager)

	case *ast.TupleStructDeclaration:
		return registerTupleStructDeclaration(statement, manager)

	case *ast.InterfaceDeclaration:
		return registerInterfaceDeclaration(statement, manager)

	case *ast.TypeDeclaration:
		return registerTypeDeclataion(statement, manager)
	default:
		return &types.Void{}
	}
}

func typeCheckVariableDeclaration(varDec *ast.VariableDeclaration, manager *modules.ModuleManager) types.ValidType {
	dataType := types.FromAst(varDec.DataType, manager.SymbolTable)
	if dataType.String() == "TypeError" {
		return dataType
	}

	if varDec.Value == nil {
		err := manager.SymbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
		if err != nil {
			err.Line = varDec.Token.Line
			err.Column = varDec.Token.Column
			return err
		}
		return dataType
	}

	expressionType := typeCheckExpression(varDec.Value, manager)
	if expressionType.String() == "TypeError" {
		return expressionType
	}

	if _, ok := expressionType.(*types.Void); ok {
		return types.Error(fmt.Sprintf("Cannot assign void to variable %q", varDec.Name), varDec)
	}

	if dataType.String() == "Infer" {
		err := manager.SymbolTable.RegisterSymbol(varDec.Name, expressionType, varDec.Constant)
		if err != nil {
			err.Line = varDec.Token.Line
			err.Column = varDec.Token.Column
			return err
		}
		return expressionType
	}

	correctType := dataType.Valid(expressionType)

	if partial, ok := dataType.(types.PartialType); ok {
		dataType, correctType = partial.Infer(expressionType)
	}

	if correctType {
		err := manager.SymbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
		if err != nil {
			err.Line = varDec.Token.Line
			err.Column = varDec.Token.Column
			return err
		}
		return dataType
	}

	return types.Error(fmt.Sprintf("Type %q is not assignable to type %q", expressionType, dataType), varDec)
}

func typeCheckFunctionDeclaration(funcDec *ast.FunctionDeclaration, manager *modules.ModuleManager) types.ValidType {
	var fn *types.Function
	if funcDec.MethodOf == nil {
		fn = manager.SymbolTable.GetSymbol(funcDec.Name).(*types.Function)
	} else {
		parentType := types.FromAst(funcDec.MethodOf, manager.SymbolTable)
		fn = types.Member(parentType, funcDec.Name, false).(*types.Function)
	}

	childTable := symbols.NewFunction(manager.SymbolTable, fn.ReturnType)
	manager.EnterScope(childTable)
	for i, param := range funcDec.Parameters {
		paramType := fn.Parameters[i]
		err := childTable.RegisterSymbol(param.Name, paramType, false)
		if err != nil {
			err.Line = funcDec.Token.Line
			err.Column = funcDec.Token.Column
			return err
		}
	}

	if fn.MethodOf != nil {
		childTable.RegisterSymbol("this", fn.MethodOf, true)
	}

	for _, statement := range funcDec.Body {
		err := typeCheckStatement(statement, manager)
		if err.String() == "TypeError" {
			return err
		}
	}

	if !fn.ReturnType.Valid(&types.Void{}) && !childTable.HasReturn() {
		return types.Error(fmt.Sprintf("Missing return from function %q", funcDec.Name), funcDec)
	}

	manager.ExitScope()

	return fn
}

func registerFunctionDeclaration(funcDec *ast.FunctionDeclaration, manager *modules.ModuleManager) types.ValidType {
	params := []types.ValidType{}
	returnType := types.FromAst(funcDec.ReturnType, manager.SymbolTable)
	if returnType.String() == "TypeError" {
		return returnType
	}

	for _, param := range funcDec.Parameters {
		paramType := types.FromAst(param.Type, manager.SymbolTable)
		if paramType.String() == "TypeError" {
			return paramType
		}
		params = append(params, paramType)
	}

	functionType := &types.Function{
		Parameters: params,
		ReturnType: returnType,
		Name:       funcDec.Name,
	}

	if funcDec.MethodOf != nil {
		parentType := types.FromAst(funcDec.MethodOf, manager.SymbolTable)
		if parentType.String() == "TypeError" {
			return parentType
		}

		if types.Member(parentType, funcDec.Name, false) != nil {
			return types.Error(fmt.Sprintf("Type %q already has member %q", parentType.String(), funcDec.Name), funcDec)
		}
		functionType.MethodOf = parentType

		types.AddMethod(funcDec.Name, functionType)
	} else {
		err := manager.SymbolTable.RegisterSymbol(funcDec.Name, functionType, true)
		if err != nil {
			err.Line = funcDec.Token.Line
			err.Column = funcDec.Token.Column
			return err
		}
	}

	if funcDec.IsExport() {
		manager.SymbolTable.GlobalScope().Exports[funcDec.Name] = functionType
	}

	return functionType
}

func typeCheckReturnStatement(ret *ast.ReturnStatement, manager *modules.ModuleManager) types.ValidType {
	expressionType := typeCheckExpression(ret.Value, manager)
	if expressionType.String() == "TypeError" {
		return expressionType
	}
	functionScope := manager.SymbolTable.FindFunctionScope()

	if functionScope == nil {
		return types.Error("Cannot use return statement outside of a function", ret)
	}

	expectedType := manager.SymbolTable.ReturnType()

	if !expectedType.Valid(expressionType) {
		return types.Error(fmt.Sprintf("Invalid return type. Expected type %q, got %q", expectedType, expressionType), ret)
	}

	manager.SymbolTable.AddReturn()
	return expressionType
}

func typeCheckIfStatement(ifStatement *ast.IfStatement, manager *modules.ModuleManager) types.ValidType {
	err := typeCheckExpression(ifStatement.Condition, manager)
	if err.String() == "TypeError" {
		return err
	}

	newScope := symbols.NewChild(manager.SymbolTable, symbols.CONDITIONAL_SCOPE)
	manager.EnterScope(newScope)

	for _, statement := range ifStatement.Body {
		err := typeCheckStatement(statement, manager)
		if err.String() == "TypeError" {
			return err
		}
	}
	manager.ExitScope()

	if nextIf, isIf := ifStatement.Else.(*ast.IfStatement); isIf {
		return typeCheckIfStatement(nextIf, manager)
	}

	if nextElse, isElse := ifStatement.Else.(*ast.ElseStatement); isElse {
		return typeCheckElseStatement(nextElse, manager)
	}

	return &types.Void{}
}

func typeCheckElseStatement(elseStatement *ast.ElseStatement, manager *modules.ModuleManager) types.ValidType {
	newScope := symbols.NewChild(manager.SymbolTable, symbols.FALLBACK_SCOPE)
	manager.EnterScope(newScope)

	for _, statement := range elseStatement.Body {
		err := typeCheckStatement(statement, manager)
		if err.String() == "TypeError" {
			return err
		}
	}
	manager.ExitScope()

	return &types.Void{}
}

func typeCheckWhileLoop(while *ast.WhileLoop, manager *modules.ModuleManager) types.ValidType {
	typeCheckExpression(while.Condition, manager)

	newScope := symbols.NewChild(manager.SymbolTable, symbols.CONDITIONAL_SCOPE)
	manager.EnterScope(newScope)

	for _, statement := range while.Body {
		err := typeCheckStatement(statement, manager)
		if err.String() == "TypeError" {
			return err
		}
	}

	return &types.Void{}
}

func typeCheckForLoop(forLoop *ast.ForLoop, manager *modules.ModuleManager) types.ValidType {
	newScope := symbols.NewChild(manager.SymbolTable, symbols.CONDITIONAL_SCOPE)
	manager.EnterScope(newScope)
	err := typeCheckStatement(forLoop.Initial, manager)
	if err.String() == "TypeError" {
		return err
	}
	err = typeCheckExpression(forLoop.Condition, manager)
	if err.String() == "TypeError" {
		return err
	}
	err = typeCheckStatement(forLoop.Update, manager)
	if err.String() == "TypeError" {
		return err
	}

	for _, statement := range forLoop.Body {
		err := typeCheckStatement(statement, manager)
		if err.String() == "TypeError" {
			return err
		}
	}
	manager.ExitScope()

	return &types.Void{}
}

func typeCheckStructDeclaration(structDecl *ast.StructDeclaration, manager *modules.ModuleManager) types.ValidType {
	structType := manager.SymbolTable.GetType(structDecl.Name).(*types.Struct)

	for memberName, memberType := range structDecl.Members {
		dataType := types.FromAst(memberType, manager.SymbolTable)
		if dataType.String() == "TypeError" {
			return dataType
		}

		structType.Members[memberName] = dataType
	}

	if structDecl.IsExport() {
		manager.SymbolTable.GlobalScope().Exports[structDecl.Name] = &types.Type{DataType: structType}
	}

	return structType
}

func typeCheckTupleStructDeclaration(structDecl *ast.TupleStructDeclaration, manager *modules.ModuleManager) types.ValidType {
	structType := manager.SymbolTable.GetType(structDecl.Name).(*types.TupleStruct)

	for _, memberType := range structDecl.Members {
		dataType := types.FromAst(memberType, manager.SymbolTable)
		if dataType.String() == "TypeError" {
			return dataType
		}

		structType.Members = append(structType.Members, dataType)
	}

	if structDecl.IsExport() {
		manager.SymbolTable.GlobalScope().Exports[structDecl.Name] = &types.Type{DataType: structType}
	}

	return structType
}

func typeCheckInterfaceDeclaration(intDecl *ast.InterfaceDeclaration, manager *modules.ModuleManager) types.ValidType {
	interfaceType := manager.SymbolTable.GetType(intDecl.Name).(*types.Interface)

	for _, member := range intDecl.Members {
		if !member.IsFunction {
			dataType := types.FromAst(member.ResultType, manager.SymbolTable)
			if dataType.String() == "TypeError" {
				return dataType
			}

			interfaceType.Members[member.Name] = dataType
			continue
		}

		fnType := &types.Function{}
		fnType.Name = member.Name

		returnType := types.FromAst(member.ResultType, manager.SymbolTable)
		if returnType.String() == "TypeError" {
			return returnType
		}
		fnType.ReturnType = returnType

		fnType.Parameters = []types.ValidType{}
		for _, param := range member.Parameters {
			paramType := types.FromAst(param, manager.SymbolTable)
			if paramType.String() == "TypeError" {
				return paramType
			}

			fnType.Parameters = append(fnType.Parameters, paramType)
		}

		interfaceType.Members[member.Name] = fnType
	}

	if intDecl.IsExport() {
		manager.SymbolTable.GlobalScope().Exports[intDecl.Name] = &types.Type{DataType: interfaceType}
	}

	return interfaceType
}

func typeCheckTypeDeclataion(typeDecl *ast.TypeDeclaration, manager *modules.ModuleManager) types.ValidType {
	dataType := types.FromAst(typeDecl.DataType, manager.SymbolTable)
	if dataType.String() == "TypeError" {
		return dataType
	}

	err := manager.SymbolTable.UpdateType(typeDecl.Name, dataType)
	if err != nil {
		return err
	}

	if typeDecl.IsExport() {
		manager.SymbolTable.GlobalScope().Exports[typeDecl.Name] = &types.Type{DataType: dataType}
	}

	return dataType
}

func registerStructDeclaration(structDecl *ast.StructDeclaration, manager *modules.ModuleManager) types.ValidType {
	members := map[string]types.ValidType{}

	structType := &types.Struct{
		Name:    structDecl.Name,
		Members: members,
	}

	err := manager.SymbolTable.AddType(structDecl.Name, structType)
	if err != nil {
		return err
	}

	return structType
}

func registerTupleStructDeclaration(structDecl *ast.TupleStructDeclaration, manager *modules.ModuleManager) types.ValidType {
	members := []types.ValidType{}

	structType := &types.TupleStruct{
		Name:    structDecl.Name,
		Members: members,
	}

	err := manager.SymbolTable.AddType(structDecl.Name, structType)
	if err != nil {
		return err
	}

	return structType
}

func registerInterfaceDeclaration(intDecl *ast.InterfaceDeclaration, manager *modules.ModuleManager) types.ValidType {
	members := map[string]types.ValidType{}

	interfaceType := &types.Interface{
		Name:    intDecl.Name,
		Members: members,
	}
	err := manager.SymbolTable.AddType(intDecl.Name, interfaceType)
	if err != nil {
		return err
	}

	return interfaceType
}

func registerTypeDeclataion(typeDecl *ast.TypeDeclaration, manager *modules.ModuleManager) types.ValidType {
	dataType := &types.Void{}
	err := manager.SymbolTable.AddType(typeDecl.Name, dataType)
	if err != nil {
		return err
	}
	return dataType
}

func typeCheckImportStatement(importStatement *ast.ImportStatement, manager *modules.ModuleManager) types.ValidType {
	modPath := importStatement.Module
	mod, exists := manager.Modules[modPath]
	if !exists {
		return types.Error(fmt.Sprintf("Cannot import module %q, it does not exist", modPath), importStatement)
	}

	err := TypeCheck(&mod.Main.Ast, mod)
	if err != nil {
		return err.(*types.TypeError)
	}

	importedMod := &types.Module{
		Name:    mod.Name,
		Exports: mod.SymbolTable.Exports,
	}

	manager.SymbolTable.RegisterSymbol(mod.Name, importedMod, true)
	return importedMod
}
