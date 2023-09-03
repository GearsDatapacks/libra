package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func typeCheckStatement(stmt ast.Statement, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
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
		return nil, errors.DevError("(Type checker) Unexpected statment type: " + statement.String())
	}
}

func typeCheckVariableDeclaration(varDec *ast.VariableDeclaration, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	expressionType, err := typeCheckExpression(varDec.Value, symbolTable)
	if err != nil {
		return nil, err
	}

	if _, ok := expressionType.(*types.Void); ok {
		return nil, errors.TypeError(fmt.Sprintf("Cannot assign void to variable %q", varDec.Name), varDec)
	}

	// Blank if type to be inferred
	if varDec.DataType.Type() == "Infer" {
		err := symbolTable.RegisterSymbol(varDec.Name, expressionType, varDec.Constant)
		if err != nil {
			return nil, err
		}
		return expressionType, nil
	}

	dataType, err := types.FromAst(varDec.DataType)
	if err != nil {
		return nil, err
	}
	correctType := dataType.Valid(expressionType)
	
	if correctType {
		symbolTable.RegisterSymbol(varDec.Name, dataType, varDec.Constant)
		return dataType, nil
	}

	return nil, errors.TypeError(fmt.Sprintf("Type %q is not assignable to type %q", expressionType, dataType), varDec)
}

func typeCheckFunctionDeclaration(funcDec *ast.FunctionDeclaration, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	params := []types.ValidType{}
	returnType, err := types.FromAst(funcDec.ReturnType)
	if err != nil {
		return nil, err
	}
	
	childTable := symbols.NewFunction(symbolTable, returnType)
	for _, param := range funcDec.Parameters {
		paramType, err := types.FromAst(param.Type)
		if err != nil {
			return nil, err
		}
		params = append(params, paramType)
		err = childTable.RegisterSymbol(param.Name, paramType, false)
		if err != nil {
			return nil, err
		}
	}

	functionType := &types.Function{
		Parameters: params,
		ReturnType: returnType,
	}

	err = symbolTable.RegisterSymbol(funcDec.Name, functionType, true)
	if err != nil {
		return nil, err
	}

	for _, statement := range funcDec.Body {
		_, err := typeCheckStatement(statement, childTable)
		if err != nil {
			return nil, err
		}
	}

	return functionType, nil
}

func typeCheckReturnStatement(ret *ast.ReturnStatement, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	expressionType, err := typeCheckExpression(ret.Value, symbolTable)
	if err != nil {
		return nil, err
	}
	functionScope := symbolTable.FindFunctionScope()

	if functionScope == nil {
		return nil, errors.TypeError("Cannot use return statement outside of a function", ret)
	}

	expectedType := functionScope.ReturnType()

	if !expectedType.Valid(expressionType) {
		return nil, errors.TypeError(fmt.Sprintf("Invalid return type. Expected type %q, got %q", expectedType, expressionType), ret)
	}

	return expressionType, nil
}

func typeCheckIfStatement(ifStatement *ast.IfStatement, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	_, err := typeCheckExpression(ifStatement.Condition, symbolTable)
	if err != nil {
		return nil, err
	}

	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)

	for _, statement := range ifStatement.Body {
		_, err := typeCheckStatement(statement, newScope)
		if err != nil {
			return nil, err
		}
	}

	if nextIf, isIf := ifStatement.Else.(*ast.IfStatement); isIf {
		return typeCheckIfStatement(nextIf, symbolTable)
	}

	if nextElse, isElse := ifStatement.Else.(*ast.ElseStatement); isElse {
		return typeCheckElseStatement(nextElse, symbolTable)
	}

	return &types.Void{}, nil
}

func typeCheckElseStatement(elseStatement *ast.ElseStatement, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)

	for _, statement := range elseStatement.Body {
		_, err := typeCheckStatement(statement, newScope)
		if err != nil {
			return nil, err
		}
	}

	return &types.Void{}, nil
}

func typeCheckWhileLoop(while *ast.WhileLoop, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	typeCheckExpression(while.Condition, symbolTable)
	
	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)

	for _, statement := range while.Body {
		_, err := typeCheckStatement(statement, newScope)
		if err != nil {
			return nil, err
		}
	}

	return &types.Void{}, nil
}

func typeCheckForLoop(forLoop *ast.ForLoop, symbolTable *symbols.SymbolTable) (types.ValidType, error) {
	newScope := symbols.NewChild(symbolTable, symbols.GENERIC_SCOPE)
	_, err := typeCheckStatement(forLoop.Initial, newScope)
	if err != nil {
		return nil, err
	}
	_, err = typeCheckExpression(forLoop.Condition, newScope)
	if err != nil {
		return nil, err
	}
	_, err = typeCheckStatement(forLoop.Update, newScope)
	if err != nil {
		return nil, err
	}

	for _, statement := range forLoop.Body {
		_, err := typeCheckStatement(statement, newScope)
		if err != nil {
			return nil, err
		}
	}

	return &types.Void{}, nil
}
