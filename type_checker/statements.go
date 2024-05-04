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
