package lowerer

import "github.com/gearsdatapacks/libra/type_checker/ir"

func (l *lowerer) lowerVariableDeclaration(varDecl *ir.VariableDeclaration) ir.Statement {
	value := l.lowerExpression(varDecl.Value)
	if value == varDecl.Value {
		return varDecl
	}
	return &ir.VariableDeclaration{
		Symbol: varDecl.Symbol,
		Value:  value,
	}
}

func (l *lowerer) lowerFunctionDeclaration(stmt *ir.FunctionDeclaration) ir.Statement {
	return stmt
}

func (l *lowerer) lowerReturnStatement(ret *ir.ReturnStatement) ir.Statement {
	if ret.Value == nil {
		return ret
	}

	value := l.lowerExpression(ret.Value)
	if value == ret.Value {
		return ret
	}
	return &ir.ReturnStatement{
		Value: value,
	}
}

func (l *lowerer) lowerBreakStatement(brk *ir.BreakStatement) ir.Statement {
	if brk.Value == nil {
		return brk
	}

	value := l.lowerExpression(brk.Value)
	if value == brk.Value {
		return brk
	}
	return &ir.ReturnStatement{
		Value: value,
	}
}

func (l *lowerer) lowerYieldStatement(yield *ir.YieldStatement) ir.Statement {
	value := l.lowerExpression(yield.Value)
	if value == yield.Value {
		return yield
	}
	return &ir.ReturnStatement{
		Value: value,
	}
}

func (l *lowerer) lowerImportStatement(stmt *ir.ImportStatement) ir.Statement {
	return stmt
}

func (l *lowerer) lowerTypeDeclaration(stmt *ir.TypeDeclaration) ir.Statement {
	return stmt
}
