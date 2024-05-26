package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (t *typeChecker) typeCheckType(expression ast.Expression) types.Type {
	return t.typeFromExpr(t.typeCheckExpression(expression), expression.Location())
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

func (t *typeChecker) lookupType(tok token.Token) types.Type {
	symbol := t.symbols.Lookup(tok.Value)
	if symbol == nil {
		t.diagnostics.Report(diagnostics.UndefinedType(tok.Location, tok.Value))
		return types.Invalid
	}
	if symbol.GetType() != types.RuntimeType {
		t.diagnostics.Report(diagnostics.ExpressionNotType(tok.Location, symbol.GetType()))
		return types.Invalid
	}
	if symbol.Value() == nil {
		t.diagnostics.Report(diagnostics.NotConst(tok.Location))
		return types.Invalid
	}
	return symbol.Value().(values.TypeValue).Type.(types.Type)
}
