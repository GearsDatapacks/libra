package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (t *typeChecker) typeCheckStatement(statement ast.Statement) ir.Statement {
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		return &ir.ExpressionStatement{
			Expression: t.typeCheckExpression(stmt.Expression),
		}
	case *ast.VariableDeclaration:
		return t.typeCheckVariableDeclaration(stmt)
	case *ast.BlockStatement:
		return t.typeCheckBlock(stmt, true)
	case *ast.IfStatement:
		return t.typeCheckIfStatement(stmt)
	case *ast.WhileLoop:
		return t.typeCheckWhileLoop(stmt)
	case *ast.ForLoop:
		return t.typeCheckForLoop(stmt)
	case *ast.FunctionDeclaration:
		return t.typeCheckFunctionDeclaration(stmt)
	case *ast.ReturnStatement:
		return t.typeCheckReturn(stmt)
	default:
		panic(fmt.Sprintf("TODO: Type-check %T", statement))
	}
}

func (t *typeChecker) typeCheckVariableDeclaration(varDec *ast.VariableDeclaration) ir.Statement {
	mutable := varDec.Keyword.Value == "mut"
	constant := varDec.Keyword.Value == "const"
	value := t.typeCheckExpression(varDec.Value)
	var expectedType types.Type = nil
	if varDec.Type != nil {
		expectedType = t.typeFromAst(varDec.Type.Type)
		if expectedType == nil {
			name := varDec.Type.Type.(*ast.TypeName)
			t.Diagnostics.Report(diagnostics.UndefinedType(name.Location(), name.Name.Value))
		}
	}

	if expectedType != nil {
		conversion := convert(value, expectedType, implicit)
		if conversion == nil {
			t.Diagnostics.Report(diagnostics.NotAssignable(varDec.Value.Location(), expectedType, value.Type()))
		} else {
			value = conversion
		}
	} else {
		value = convert(value, types.ToReal(value.Type()), implicit)
	}

	varType := expectedType
	if expectedType == nil {
		varType = value.Type()
	}

	if constant && !value.IsConst() {
		t.Diagnostics.Report(diagnostics.NotConst(varDec.Value.Location()))
	}
	var constVal values.ConstValue
	if !mutable && value.IsConst() {
		constVal = value.ConstValue()
	}

	variable := symbols.Variable{
		Name:       varDec.Identifier.Value,
		Mutable:    mutable,
		Type:       varType,
		ConstValue: constVal,
	}
	if !t.symbols.DeclareVariable(variable) {
		t.Diagnostics.Report(diagnostics.VariableDefined(varDec.Identifier.Location, variable.Name))
	}
	return &ir.VariableDeclaration{
		Name:  variable.Name,
		Value: value,
	}
}

func (t *typeChecker) typeCheckBlock(block *ast.BlockStatement, createScope bool) *ir.Block {
	if createScope {
		t.enterScope()
		defer t.exitScope()
	}

	stmts := []ir.Statement{}
	for _, stmt := range block.Statements {
		stmts = append(stmts, t.typeCheckStatement(stmt))
	}
	return &ir.Block{
		Statements: stmts,
	}
}

func (t *typeChecker) typeCheckIfStatement(ifStmt *ast.IfStatement) ir.Statement {
	condition := t.typeCheckExpression(ifStmt.Condition)
	if !types.Assignable(types.Bool, condition.Type()) {
		t.Diagnostics.Report(diagnostics.ConditionMustBeBool(ifStmt.Condition.Location()))
	}

	body := t.typeCheckBlock(ifStmt.Body, true)
	var elseBranch ir.Statement
	if ifStmt.ElseBranch != nil {
		elseBranch = t.typeCheckStatement(ifStmt.ElseBranch.Statement)
	}
	return &ir.IfStatement{
		Condition:  condition,
		Body:       body,
		ElseBranch: elseBranch,
	}
}

func (t *typeChecker) typeCheckWhileLoop(loop *ast.WhileLoop) ir.Statement {
	condition := t.typeCheckExpression(loop.Condition)
	if !types.Assignable(types.Bool, condition.Type()) {
		t.Diagnostics.Report(diagnostics.ConditionMustBeBool(loop.Condition.Location()))
	}

	body := t.typeCheckBlock(loop.Body, true)
	return &ir.WhileLoop{
		Condition: condition,
		Body:      body,
	}
}

func (t *typeChecker) typeCheckForLoop(loop *ast.ForLoop) ir.Statement {
	iter := t.typeCheckExpression(loop.Iterator)
	var itemType types.Type = types.Invalid
	if iterator, ok := iter.Type().(types.Iterator); ok {
		itemType = iterator.Item()
	} else {
		t.Diagnostics.Report(diagnostics.NotIterable(loop.Iterator.Location()))
	}
	variable := symbols.Variable{
		Name:       loop.Variable.Value,
		Mutable:    false,
		Type:       itemType,
		ConstValue: nil,
	}
	t.enterScope()
	defer t.exitScope()
	t.symbols.DeclareVariable(variable)
	body := t.typeCheckBlock(loop.Body, false)

	return &ir.ForLoop{
		Variable: variable,
		Iterator: iter,
		Body:     body,
	}
}

func (t *typeChecker) typeCheckFunctionDeclaration(funcDec *ast.FunctionDeclaration) ir.Statement {
	fnType := t.symbols.LookupVariable(funcDec.Name.Value).Type.(*types.Function)

	t.enterScope(symbols.FunctionContext{ReturnType: fnType.ReturnType})
	defer t.exitScope()
	params := []string{}

	for i, param := range funcDec.Parameters {
		symbol := symbols.Variable{
			Name:       param.Name.Value,
			Mutable:    param.Mutable != nil,
			Type:       fnType.Parameters[i],
			ConstValue: nil,
		}
		t.symbols.DeclareVariable(symbol)
		params = append(params, param.Name.Value)
	}

	body := t.typeCheckBlock(funcDec.Body, false)

	return &ir.FunctionDeclaration{
		Name:       funcDec.Name.Value,
		Parameters: params,
		Body:       body,
	}
}

func (t *typeChecker) typeCheckReturn(ret *ast.ReturnStatement) ir.Statement {
	var expectedType types.Type = nil
	symbolTable := t.symbols
	for symbolTable != nil {
		if fnContext, ok := symbolTable.Context.(symbols.FunctionContext); ok {
			expectedType = fnContext.ReturnType
			break
		}
		symbolTable = symbolTable.Parent
	}
	if expectedType == nil {
		t.Diagnostics.Report(diagnostics.NoReturnOutsideFunction(ret.Keyword.Location))
		expectedType = types.Invalid
	}

	if ret.Value == nil {
		if expectedType != types.Void {
			t.Diagnostics.Report(diagnostics.ExpectedReturnValue(ret.Keyword.Location))
		}
		return &ir.ReturnStatement{
			Value: nil,
		}
	}

	value := t.typeCheckExpression(ret.Value)
	if conversion := convert(value, expectedType, implicit); conversion != nil {
		value = conversion
	} else {
		t.Diagnostics.Report(diagnostics.NotAssignable(ret.Value.Location(), expectedType, value.Type()))
	}

	return &ir.ReturnStatement{
		Value: value,
	}
}
