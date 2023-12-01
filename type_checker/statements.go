package typechecker

import (
	"fmt"
	"log"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckStatement(stmt ast.Statement, symbolTable *symbols.SymbolTable) types.ValidType {
	switch statement := stmt.(type) {
	case *ast.VariableDeclaration:
		return typeCheckVariableDeclaration(statement, symbolTable)

	case *ast.ExpressionStatement:
		return typeCheckExpression(statement.Expression, symbolTable)

	case *ast.FunctionDeclaration:
		return typeCheckFunctionDeclaration(statement, symbolTable)

	case *ast.ReturnStatement:
		return typeCheckReturnStatement(statement, symbolTable)

	case *ast.IfStatement:
		return typeCheckIfStatement(statement, symbolTable)

	case *ast.WhileLoop:
		return typeCheckWhileLoop(statement, symbolTable)

	case *ast.ForLoop:
		return typeCheckForLoop(statement, symbolTable)

	case *ast.StructDeclaration:
		return typeCheckStructDeclaration(statement, symbolTable)

	case *ast.InterfaceDeclaration:
		return typeCheckInterfaceDeclaration(statement, symbolTable)

	case *ast.TypeDeclaration:
		return typeCheckTypeDeclataion(statement, symbolTable)

	default:
		log.Fatal(errors.DevError("(Type checker) Unexpected statment type: " + statement.String()))
		return nil
	}
}

func typeCheckVariableDeclaration(varDec *ast.VariableDeclaration, symbolTable *symbols.SymbolTable) types.ValidType {
	dataType := types.FromAst(varDec.DataType, symbolTable)
	if dataType.String() == "TypeError" {
		return dataType
	}

	if varDec.Value == nil {
		err := symbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
		if err != nil {
			err.Line = varDec.Token.Line
			err.Column = varDec.Token.Column
			return err
		}
		return dataType
	}

	expressionType := typeCheckExpression(varDec.Value, symbolTable)
	if expressionType.String() == "TypeError" {
		return expressionType
	}

	if _, ok := expressionType.(*types.Void); ok {
		return types.Error(fmt.Sprintf("Cannot assign void to variable %q", varDec.Name), varDec)
	}

	if dataType.String() == "Infer" {
		err := symbolTable.RegisterSymbol(varDec.Name, expressionType, varDec.Constant)
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
		err := symbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
		if err != nil {
			err.Line = varDec.Token.Line
			err.Column = varDec.Token.Column
			return err
		}
		return dataType
	}

	return types.Error(fmt.Sprintf("Type %q is not assignable to type %q", expressionType, dataType), varDec)
}

func typeCheckFunctionDeclaration(funcDec *ast.FunctionDeclaration, symbolTable *symbols.SymbolTable) types.ValidType {
	params := []types.ValidType{}
	returnType := types.FromAst(funcDec.ReturnType, symbolTable)
	if returnType.String() == "TypeError" {
		return returnType
	}

	childTable := symbols.NewFunction(symbolTable, returnType)
	for _, param := range funcDec.Parameters {
		paramType := types.FromAst(param.Type, symbolTable)
		if paramType.String() == "TypeError" {
			return paramType
		}
		params = append(params, paramType)
		err := childTable.RegisterSymbol(param.Name, paramType, false)
		if err != nil {
			err.Line = funcDec.Token.Line
			err.Column = funcDec.Token.Column
			return err
		}
	}

	functionType := &types.Function{
		Parameters: params,
		ReturnType: returnType,
		Name: funcDec.Name,
	}

	if funcDec.MethodOf != nil {
		parentType := types.FromAst(funcDec.MethodOf, childTable)
		if parentType.String() == "TypeError" {
			return parentType
		}
		functionType.MethodOf = parentType

		childTable.RegisterSymbol("this", parentType, true)
		types.AddMethod(funcDec.Name, functionType)
	} else {
		err := symbolTable.RegisterSymbol(funcDec.Name, functionType, true)
		if err != nil {
			err.Line = funcDec.Token.Line
			err.Column = funcDec.Token.Column
			return err
		}
	}
	
	for _, statement := range funcDec.Body {
		err := typeCheckStatement(statement, childTable)
		if err.String() == "TypeError" {
			return err
		}
	}

	return functionType
}

func typeCheckReturnStatement(ret *ast.ReturnStatement, symbolTable *symbols.SymbolTable) types.ValidType {
	expressionType := typeCheckExpression(ret.Value, symbolTable)
	if expressionType.String() == "TypeError" {
		return expressionType
	}
	functionScope := symbolTable.FindFunctionScope()

	if functionScope == nil {
		return types.Error("Cannot use return statement outside of a function", ret)
	}

	expectedType := functionScope.ReturnType()

	if !expectedType.Valid(expressionType) {
		return types.Error(fmt.Sprintf("Invalid return type. Expected type %q, got %q", expectedType, expressionType), ret)
	}

	return expressionType
}

func typeCheckIfStatement(ifStatement *ast.IfStatement, symbolTable *symbols.SymbolTable) types.ValidType {
	err := typeCheckExpression(ifStatement.Condition, symbolTable)
	if err.String() == "TypeError" {
		return err
	}

	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)

	for _, statement := range ifStatement.Body {
		err := typeCheckStatement(statement, newScope)
		if err.String() == "TypeError" {
			return err
		}
	}

	if nextIf, isIf := ifStatement.Else.(*ast.IfStatement); isIf {
		return typeCheckIfStatement(nextIf, symbolTable)
	}

	if nextElse, isElse := ifStatement.Else.(*ast.ElseStatement); isElse {
		return typeCheckElseStatement(nextElse, symbolTable)
	}

	return &types.Void{}
}

func typeCheckElseStatement(elseStatement *ast.ElseStatement, symbolTable *symbols.SymbolTable) types.ValidType {
	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)

	for _, statement := range elseStatement.Body {
		err := typeCheckStatement(statement, newScope)
		if err.String() == "TypeError" {
			return err
		}
	}

	return &types.Void{}
}

func typeCheckWhileLoop(while *ast.WhileLoop, symbolTable *symbols.SymbolTable) types.ValidType {
	typeCheckExpression(while.Condition, symbolTable)

	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)

	for _, statement := range while.Body {
		err := typeCheckStatement(statement, newScope)
		if err.String() == "TypeError" {
			return err
		}
	}

	return &types.Void{}
}

func typeCheckForLoop(forLoop *ast.ForLoop, symbolTable *symbols.SymbolTable) types.ValidType {
	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)
	err := typeCheckStatement(forLoop.Initial, newScope)
	if err.String() == "TypeError" {
		return err
	}
	err = typeCheckExpression(forLoop.Condition, newScope)
	if err.String() == "TypeError" {
		return err
	}
	err = typeCheckStatement(forLoop.Update, newScope)
	if err.String() == "TypeError" {
		return err
	}

	for _, statement := range forLoop.Body {
		err := typeCheckStatement(statement, newScope)
		if err.String() == "TypeError" {
			return err
		}
	}

	return &types.Void{}
}

func typeCheckStructDeclaration(structDecl *ast.StructDeclaration, symbolTable *symbols.SymbolTable) types.ValidType {
	members := map[string]types.ValidType{}

	for memberName, memberType := range structDecl.Members {
		dataType := types.FromAst(memberType, symbolTable)
		if dataType.String() == "TypeError" {
			return dataType
		}

		members[memberName] = dataType
	}

	structType := &types.Struct{
		Name:    structDecl.Name,
		Members: members,
	}
	err := symbolTable.AddType(structDecl.Name, structType)
	if err != nil {
		return err
	}

	return structType
}

func typeCheckInterfaceDeclaration(intDecl *ast.InterfaceDeclaration, symbolTable *symbols.SymbolTable) types.ValidType {
	members := map[string]types.ValidType{}

	for _, member := range intDecl.Members {
		if !member.IsFunction {
			dataType := types.FromAst(member.ResultType, symbolTable)
			if dataType.String() == "TypeError" {
				return dataType
			}

			members[member.Name] = dataType
			continue
		}

		fnType := &types.Function{}
		fnType.Name = member.Name

		returnType := types.FromAst(member.ResultType, symbolTable)
		if returnType.String() == "TypeError" {
			return returnType
		}
		fnType.ReturnType = returnType

		fnType.Parameters = []types.ValidType{}
		for _, param := range member.Parameters {
			paramType := types.FromAst(param, symbolTable)
			if paramType.String() == "TypeError" {
				return paramType
			}

			fnType.Parameters = append(fnType.Parameters, paramType)
		}

		members[member.Name] = fnType
	}

	interfaceType := &types.Interface{
		Name:    intDecl.Name,
		Members: members,
	}
	err := symbolTable.AddType(intDecl.Name, interfaceType)
	if err != nil {
		return err
	}

	return interfaceType
}

func typeCheckTypeDeclataion(typeDecl *ast.TypeDeclaration, symbolTable *symbols.SymbolTable) types.ValidType {
	dataType := types.FromAst(typeDecl.DataType, symbolTable)
	if dataType.String() == "TypeError" {
		return dataType
	}

	err := symbolTable.AddType(typeDecl.Name, dataType)
	if err != nil {
		return err
	}
	return dataType
}
