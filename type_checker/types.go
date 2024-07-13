package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (t *typeChecker) typeCheckType(expression ast.Expression) types.Type {
	return t.typeFromExpr(t.doTypeCheckExpression(expression), expression.GetLocation())
}

func (t *typeChecker) typeFromExpr(expr ir.Expression, location text.Location) types.Type {
	if expr.Type() == types.Invalid {
		return types.Invalid
	}
	if expr.Type() != types.RuntimeType {
		t.diagnostics.Report(diagnostics.ExpressionNotType(location, expr.Type()))
		return types.Invalid
	}
	if !expr.IsConst() {
		t.diagnostics.Report(diagnostics.NotConst(location))
		return types.Invalid
	}
	return expr.ConstValue().(values.TypeValue).Type.(types.Type)
}

func (t *typeChecker) lookupType(name string, location text.Location) types.Type {
	symbol := t.symbols.Lookup(name)
	if symbol == nil {
		t.diagnostics.Report(diagnostics.UndefinedType(location, name))
		return types.Invalid
	}
	if symbol.GetType() != types.RuntimeType {
		t.diagnostics.Report(diagnostics.ExpressionNotType(location, symbol.GetType()))
		return types.Invalid
	}
	if symbol.Value() == nil {
		t.diagnostics.Report(diagnostics.NotConst(location))
		return types.Invalid
	}
	return symbol.Value().(values.TypeValue).Type.(types.Type)
}
