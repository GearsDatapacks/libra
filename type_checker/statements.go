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
		symbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
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
			return err
		}
		return expressionType
	}

	correctType := dataType.Valid(expressionType)

	if partial, ok := dataType.(types.PartialType); ok {
		dataType, correctType = partial.Infer(expressionType)
	}

	if correctType {
		symbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
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
			return err
		}
	}

	functionType := &types.Function{
		Parameters: params,
		ReturnType: returnType,
	}

	err := symbolTable.RegisterSymbol(funcDec.Name, functionType, true)
	if err != nil {
		return err
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
