package typechecker

import (
	"fmt"

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

	default:
		errors.DevError("(Type checker) Unexpected statment type: " + statement.String())
		return &types.IntLiteral{}
	}
}

func typeCheckVariableDeclaration(varDec *ast.VariableDeclaration, symbolTable *symbols.SymbolTable) types.ValidType {
	expressionType := typeCheckExpression(varDec.Value, symbolTable)

	if _, ok := expressionType.(*types.Void); ok {
		errors.TypeError(fmt.Sprintf("Cannot assign void to variable %q", varDec.Name), varDec)
	}

	// Blank if type to be inferred
	if varDec.DataType.Type() == "Infer" {
		symbolTable.RegisterSymbol(varDec.Name, expressionType, varDec.Constant)
		return expressionType
	}

	dataType := types.FromAst(varDec.DataType)
	correctType := dataType.Valid(expressionType)
	
	if correctType {
		symbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
		return dataType
	}

	errors.TypeError(fmt.Sprintf("Type %q is not assignable to type %q", expressionType, dataType), varDec)
	return &types.IntLiteral{}
}

func typeCheckFunctionDeclaration(funcDec *ast.FunctionDeclaration, symbolTable *symbols.SymbolTable) types.ValidType {
	params := []types.ValidType{}
	returnType := types.FromAst(funcDec.ReturnType)
	
	childTable := symbols.NewFunction(symbolTable, returnType)
	for _, param := range funcDec.Parameters {
		paramType := types.FromAst(param.Type)
		params = append(params, paramType)
		childTable.RegisterSymbol(param.Name, paramType, false)
	}

	functionType := &types.Function{
		Parameters: params,
		ReturnType: returnType,
	}

	symbolTable.RegisterSymbol(funcDec.Name, functionType, true)

	for _, statement := range funcDec.Body {
		typeCheckStatement(statement, childTable)
	}

	return functionType
}

func typeCheckReturnStatement(ret *ast.ReturnStatement, symbolTable *symbols.SymbolTable) types.ValidType {
	expressionType := typeCheckExpression(ret.Value, symbolTable)
	functionScope := symbolTable.FindFunctionScope()

	if functionScope == nil {
		errors.TypeError("Cannot use return statement outside of a function", ret)
	}

	expectedType := functionScope.ReturnType()

	if !expectedType.Valid(expressionType) {
		errors.TypeError(fmt.Sprintf("Invalid return type. Expected type %q, got %q", expectedType, expressionType), ret)
	}

	return expressionType
}

func typeCheckIfStatement(ifStatement *ast.IfStatement, symbolTable *symbols.SymbolTable) types.ValidType {
	typeCheckExpression(ifStatement.Condition, symbolTable)

	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)

	for _, statement := range ifStatement.Body {
		typeCheckStatement(statement, newScope)
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
		typeCheckStatement(statement, newScope)
	}

	return &types.Void{}
}

func typeCheckWhileLoop(while *ast.WhileLoop, symbolTable *symbols.SymbolTable) types.ValidType {
	typeCheckExpression(while.Condition, symbolTable)
	
	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)

	for _, statement := range while.Body {
		typeCheckStatement(statement, newScope)
	}

	return &types.Void{}
}

func typeCheckForLoop(forLoop *ast.ForLoop, symbolTable *symbols.SymbolTable) types.ValidType {
	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)
	typeCheckStatement(forLoop.Initial, newScope)
	typeCheckExpression(forLoop.Condition, newScope)
	typeCheckStatement(forLoop.Update, newScope)

	for _, statement := range forLoop.Body {
		typeCheckStatement(statement, newScope)
	}

	return &types.Void{}
}
