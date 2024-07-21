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
	case *ast.VariableDeclaration:
		return t.typeCheckVariableDeclaration(stmt)
	case *ast.Block:
		return t.typeCheckBlock(stmt, true)
	case *ast.IfExpression:
		return t.typeCheckIfExpression(stmt)
	case *ast.WhileLoop:
		return t.typeCheckWhileLoop(stmt)
	case *ast.ForLoop:
		return t.typeCheckForLoop(stmt)
	case *ast.FunctionDeclaration:
		return t.typeCheckFunctionDeclaration(stmt)
	case *ast.ReturnStatement:
		return t.typeCheckReturn(stmt)
	case *ast.BreakStatement:
		return t.typeCheckBreak(stmt)
	case *ast.YieldStatement:
		return t.typeCheckYield(stmt)
	case *ast.ContinueStatement:
		return t.typeCheckContinue(stmt)
	case *ast.TypeDeclaration,
		*ast.StructDeclaration,
		*ast.InterfaceDeclaration,
		*ast.ImportStatement,
		*ast.UnionDeclaration,
		*ast.TagDeclaration:

		return &ir.Block{
			Statements: []ir.Statement{},
			ResultType: types.Void,
		}
	case ast.Expression:
		return t.typeCheckExpression(stmt)
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
		expectedType = t.typeCheckType(varDec.Type)
	}

	if expectedType != nil {
		conversion := convert(value, expectedType, implicit)
		if conversion == nil {
			t.diagnostics.Report(diagnostics.NotAssignable(varDec.Value.GetLocation(), expectedType, value.Type()))
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
		t.diagnostics.Report(diagnostics.NotConst(varDec.Value.GetLocation()))
	}
	var constVal values.ConstValue
	if !mutable && value.IsConst() {
		constVal = value.ConstValue()
	}

	variable := &symbols.Variable{
		Name:       varDec.Name,
		IsMut:      mutable,
		Type:       varType,
		ConstValue: constVal,
	}
	if !t.symbols.Register(variable) {
		t.diagnostics.Report(diagnostics.VariableDefined(varDec.NameLocation, variable.Name))
	}
	return &ir.VariableDeclaration{
		Name:  variable.Name,
		Value: value,
	}
}

func (t *typeChecker) typeCheckFunctionDeclaration(funcDec *ast.FunctionDeclaration) ir.Statement {
	var fnType *types.Function
	if funcDec.MethodOf != nil {
		fnType = t.symbols.LookupMethod(funcDec.Name, t.typeCheckType(funcDec.MethodOf.Type), false)
	} else if funcDec.MemberOf != nil {
		fnType = t.symbols.LookupMethod(funcDec.Name, t.lookupType(funcDec.MemberOf.Name, funcDec.MemberOf.Location), true)
	} else {
		fnType = t.symbols.Lookup(funcDec.Name).GetType().(*types.Function)
	}

	t.enterScope(symbols.FunctionContext{ReturnType: fnType.ReturnType})
	defer t.exitScope()
	params := []string{}

	for i, param := range funcDec.Parameters {
		if param.Name == nil {
			t.diagnostics.Report(diagnostics.UnnamedParameter(param.Type.GetLocation()))
			continue
		}
		symbol := &symbols.Variable{
			Name:       *param.Name,
			IsMut:      param.Mutable,
			Type:       fnType.Parameters[i],
			ConstValue: nil,
		}
		t.symbols.Register(symbol)
		params = append(params, *param.Name)
	}

	if funcDec.MethodOf != nil {
		symbol := &symbols.Variable{
			Name:       "this",
			IsMut:      funcDec.MethodOf.Mutable,
			Type:       t.typeCheckType(funcDec.MethodOf.Type),
			ConstValue: nil,
		}
		t.symbols.Register(symbol)
	}

	body := t.typeCheckBlock(funcDec.Body, false)

	return &ir.FunctionDeclaration{
		Name:       funcDec.Name,
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
		t.diagnostics.Report(diagnostics.NoReturnOutsideFunction(ret.Location))
		expectedType = types.Invalid
	}

	if ret.Value == nil {
		if expectedType != types.Void {
			t.diagnostics.Report(diagnostics.ExpectedReturnValue(ret.Location))
		}
		return &ir.ReturnStatement{
			Value: nil,
		}
	}

	value := t.typeCheckExpression(ret.Value)
	if conversion := convert(value, expectedType, implicit); conversion != nil {
		value = conversion
	} else {
		t.diagnostics.Report(diagnostics.NotAssignable(ret.Value.GetLocation(), expectedType, value.Type()))
	}

	return &ir.ReturnStatement{
		Value: value,
	}
}

func (t *typeChecker) typeCheckBreak(b *ast.BreakStatement) ir.Statement {
	symbolTable := t.symbols
	var context *symbols.LoopContext
	for symbolTable != nil {
		if ctx, ok := symbolTable.Context.(*symbols.LoopContext); ok {
			context = ctx
			break
		}
		symbolTable = symbolTable.Parent
	}

	var value ir.Expression
	if b.Value != nil {
		value = t.typeCheckExpression(b.Value)
	}
	if context == nil {
		t.diagnostics.Report(diagnostics.CannotUseStatementOutsideLoop(b.Location, "break"))
	} else if value != nil {
		context.ResultType = value.Type()
	}

	return &ir.BreakStatement{
		Value: value,
	}
}

func (t *typeChecker) typeCheckYield(yield *ast.YieldStatement) ir.Statement {
	symbolTable := t.symbols
	var context *symbols.BlockContext
	context, ok := symbolTable.Context.(*symbols.BlockContext)

	value := t.typeCheckExpression(yield.Value)
	if !ok {
		t.diagnostics.Report(diagnostics.CannotUseStatementOutsideBlock(yield.Location, "yield"))
	} else {
		context.ResultType = value.Type()
	}

	return &ir.YieldStatement{
		Value: value,
	}
}

func (t *typeChecker) typeCheckContinue(c *ast.ContinueStatement) ir.Statement {
	symbolTable := t.symbols
	for symbolTable != nil {
		if _, ok := symbolTable.Context.(*symbols.LoopContext); ok {
			break
		}
		symbolTable = symbolTable.Parent
	}

	if symbolTable == nil {
		t.diagnostics.Report(diagnostics.CannotUseStatementOutsideLoop(c.Location, "continue"))
	}

	return &ir.ContinueStatement{}
}
