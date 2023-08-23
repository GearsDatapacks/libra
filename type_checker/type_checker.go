package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TypeCheck(program ast.Program) {
	symbolTable := symbols.New()

	for _, stmt := range program.Body {
		typeCheck(stmt, symbolTable)
	}
}

func typeCheck(stmt ast.Statement, symbolTable *symbols.SymbolTable) types.ValidType {
	switch statement := stmt.(type) {
	case *ast.VariableDeclaration:
		return typeCheckVariableDeclaration(statement, symbolTable)

	case *ast.ExpressionStatement:
		return typeCheckExpression(statement.Expression, symbolTable)

	case *ast.FunctionDeclaration:
		return typeCheckFunctionDeclaration(statement, symbolTable)
	
	case *ast.ReturnStatement:
		return typeCheckReturnStatement(statement, symbolTable)

	default:
		errors.DevError("Unexpected statment type: " + statement.String())
		return &types.Literal{}
	}
}

func typeCheckVariableDeclaration(varDec *ast.VariableDeclaration, symbolTable *symbols.SymbolTable) types.ValidType {
	expressionType := typeCheckExpression(varDec.Value, symbolTable)

	// Blank if type to be inferred
	if varDec.DataType.Type() == "Infer" {
		symbolTable.RegisterSymbol(varDec.Name, expressionType, varDec.Constant)
		return expressionType
	}

	dataType := types.FromAst(varDec.DataType)
	correctType := expressionType.Valid(dataType)
	
	if correctType {
		symbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
		return dataType
	}

	errors.TypeError(fmt.Sprintf("Type %q is not assignable to type %q", expressionType, dataType), varDec)
	return &types.Literal{}
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
		typeCheck(statement, childTable)
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
