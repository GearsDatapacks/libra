package lowerer

import "github.com/gearsdatapacks/libra/type_checker/ir"

func (l *lowerer) lowerVariableDeclaration(varDecl *ir.VariableDeclaration, statements *[]ir.Statement) {
	value := l.lowerExpression(varDecl.Value, statements)
	if value == varDecl.Value {
		*statements = append(*statements, varDecl)
		return
	}
	*statements = append(*statements, &ir.VariableDeclaration{
		Symbol: varDecl.Symbol,
		Value:  value,
	})
}

func (l *lowerer) lowerFunctionDeclaration(stmt *ir.FunctionDeclaration, statements *[]ir.Statement) {
	panic("TODO")
}

func (l *lowerer) lowerReturnStatement(ret *ir.ReturnStatement, statements *[]ir.Statement) {
	if ret.Value == nil {
		*statements = append(*statements, ret)
		return
	}

	value := l.lowerExpression(ret.Value, statements)
	if value == ret.Value {
		*statements = append(*statements, ret)
	}
	*statements = append(*statements, &ir.ReturnStatement{
		Value: value,
	})
}

func (l *lowerer) lowerBreakStatement(brk *ir.BreakStatement, statements *[]ir.Statement) {
	context := findContext[loopContext](l)
	if brk.Value != nil {
		value := l.lowerExpression(brk.Value, statements)
		*statements = append(*statements, &ir.Assignment{
			Assignee: &ir.VariableExpression{Symbol: context.breakVariable},
			Value:    value,
		})
	}

	*statements = append(*statements, &ir.Goto{Label: context.breakLabel})
}

func (l *lowerer) lowerContinueStatement(_ *ir.ContinueStatement, statements *[]ir.Statement) {
	context := findContext[loopContext](l)
	*statements = append(*statements, &ir.Goto{Label: context.continueLabel})
}

func (l *lowerer) lowerYieldStatement(yield *ir.YieldStatement, statements *[]ir.Statement) {
	value := l.lowerExpression(yield.Value, statements)
	if value == yield.Value {
		*statements = append(*statements, yield)
		return
	}
	*statements = append(*statements, &ir.YieldStatement{
		Value: value,
	})
}

func (l *lowerer) lowerImportStatement(stmt *ir.ImportStatement, statements *[]ir.Statement) {
	*statements = append(*statements, stmt)
}

func (l *lowerer) lowerTypeDeclaration(stmt *ir.TypeDeclaration, statements *[]ir.Statement) {
	*statements = append(*statements, stmt)
}
