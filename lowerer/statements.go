package lowerer

import (
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func (l *lowerer) lowerVariableDeclaration(varDecl *ir.VariableDeclaration, statements *[]ir.Statement) {
	value := l.lowerExpression(varDecl.Value, statements, true)
	if value == varDecl.Value {
		*statements = append(*statements, varDecl)
		return
	}
	*statements = append(*statements, &ir.VariableDeclaration{
		Symbol: varDecl.Symbol,
		Value:  value,
	})
}

func (l *lowerer) lowerFunctionDeclaration(funcDecl *ir.FunctionDeclaration) *ir.FunctionDeclaration {
	var body *ir.Block
	if funcDecl.Body != nil {
		statements := []ir.Statement{}
		for _, stmt := range funcDecl.Body.Statements {
			l.lower(stmt, &statements)
		}
		statements = l.cfa(statements, &funcDecl.Location, funcDecl.Type.ReturnType != types.Void)
		body = &ir.Block{Statements: statements, ResultType: funcDecl.Body.ResultType}
	}

	fn := &ir.FunctionDeclaration{
		Name:       funcDecl.Name,
		Parameters: funcDecl.Parameters,
		Body:       body,
		Type:       funcDecl.Type,
		Exported:   funcDecl.Exported,
		Extern:     funcDecl.Extern,
		Location:   funcDecl.Location,
	}
	return fn
}

func (l *lowerer) lowerReturnStatement(ret *ir.ReturnStatement, statements *[]ir.Statement) {
	if ret.Value == nil {
		*statements = append(*statements, ret)
		return
	}

	value := l.lowerExpression(ret.Value, statements, true)
	if value == ret.Value {
		*statements = append(*statements, ret)
		return
	}
	*statements = append(*statements, &ir.ReturnStatement{
		Value: value,
	})
}

func (l *lowerer) lowerBreakStatement(brk *ir.BreakStatement, statements *[]ir.Statement) {
	context := findContext[loopContext](l)
	if brk.Value != nil {
		value := l.lowerExpression(brk.Value, statements, true)
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
	context := findContext[blockContext](l)
	value := l.lowerExpression(yield.Value, statements, true)
	*statements = append(*statements, &ir.Assignment{
		Assignee: &ir.VariableExpression{Symbol: context.yieldVariable},
		Value:    value,
	})
	*statements = append(*statements, &ir.Goto{Label: context.endLabel})
}

func (l *lowerer) lowerImportStatement(stmt *ir.ImportStatement) *ir.ImportStatement {
	return stmt
}

func (l *lowerer) lowerTypeDeclaration(stmt *ir.TypeDeclaration) *ir.TypeDeclaration {
	return stmt
}
