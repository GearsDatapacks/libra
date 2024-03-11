package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
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
	value := t.typeCheckExpression(varDec.Value)
	var expectedType types.Type = nil
	if varDec.Type != nil {
		expectedType = types.FromAst(varDec.Type.Type)
		if expectedType == nil {
			name := varDec.Type.Type.(*ast.TypeName)
			t.Diagnostics.ReportUndefinedType(name.Name.Location, name.Name.Value)
		}
	}

	if expectedType != nil && expectedType != value.Type() {
		t.Diagnostics.ReportNotAssignable(varDec.Value.Tokens()[0].Location, expectedType, value.Type())
	}
	varType := expectedType
	if expectedType == nil {
		varType = value.Type()
	}

	variable := symbols.Variable{
		Name:    varDec.Identifier.Value,
		Mutable: mutable,
		Type:    varType,
	}
	if !t.symbols.DeclareVariable(variable) {
		t.Diagnostics.ReportVariableDefined(varDec.Identifier.Location, variable.Name)
	}
	return &ir.VariableDeclaration{
		Name:  variable.Name,
		Value: value,
	}
}
