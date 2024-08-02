package lowerer

import (
	"fmt"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
)

type lowerer struct {
	diagnostics diagnostics.Manager
	labelId     int
	varId       int
	scope       *scope
}

type scope struct {
	parent  *scope
	context any
}

type loopContext struct {
	breakLabel,
	continueLabel string
	breakVariable symbols.Variable
}

type blockContext struct {
	endLabel string
	yieldVariable symbols.Variable
}

func Lower(pkg *ir.Package, diagnostics diagnostics.Manager) (*ir.Package, diagnostics.Manager) {
	lowerer := lowerer{
		diagnostics: diagnostics,
	}

	lowered := &ir.Package{
		Modules: map[string]*ir.Module{},
	}

	for name, module := range pkg.Modules {
		mod := &ir.Module{
			Name:       name,
			Statements: []ir.Statement{},
		}
		lowered.Modules[name] = mod

		for _, stmt := range module.Statements {
			lowerer.lower(stmt, &mod.Statements)
		}
	}
	return lowered, lowerer.diagnostics
}

func (l *lowerer) genLabel() string {
	id := l.labelId
	l.labelId++
	return fmt.Sprintf("label%d", id)
}

func (l *lowerer) genVar() string {
	id := l.varId
	l.varId++
	return fmt.Sprintf("var%d", id)
}

func (l *lowerer) beginScope(context any) any {
	l.scope = &scope{
		parent:  l.scope,
		context: context,
	}
	return true
}

func (l *lowerer) endScope(_ any) {
	l.scope = l.scope.parent
}

func findContext[Context any](l *lowerer) Context {
	scope := l.scope
	for scope != nil {
		if context, ok := scope.context.(Context); ok {
			return context
		}
		scope = scope.parent
	}

	panic("Should find context")
}

func (l *lowerer) lower(statement ir.Statement, statements *[]ir.Statement) {
	switch stmt := statement.(type) {
	case *ir.VariableDeclaration:
		l.lowerVariableDeclaration(stmt, statements)
	case *ir.FunctionDeclaration:
		l.lowerFunctionDeclaration(stmt, statements)
	case *ir.ReturnStatement:
		l.lowerReturnStatement(stmt, statements)
	case *ir.BreakStatement:
		l.lowerBreakStatement(stmt, statements)
	case *ir.ContinueStatement:
		l.lowerContinueStatement(stmt, statements)
	case *ir.YieldStatement:
		l.lowerYieldStatement(stmt, statements)
	case *ir.ImportStatement:
		l.lowerImportStatement(stmt, statements)
	case *ir.TypeDeclaration:
		l.lowerTypeDeclaration(stmt, statements)

	case ir.Expression:
		*statements = append(*statements, l.lowerExpression(stmt, statements))

	default:
		panic(fmt.Sprintf("TODO: lower %T", stmt))
	}
}

func (l *lowerer) lowerExpression(expression ir.Expression, statements *[]ir.Statement) ir.Expression {
	switch expr := expression.(type) {
	case *ir.IntegerLiteral:
		return l.lowerIntegerLiteral(expr, statements)
	case *ir.FloatLiteral:
		return l.lowerFloatLiteral(expr, statements)
	case *ir.BooleanLiteral:
		return l.lowerBooleanLiteral(expr, statements)
	case *ir.StringLiteral:
		return l.lowerStringLiteral(expr, statements)
	case *ir.VariableExpression:
		return l.lowerVariableExpression(expr, statements)
	case *ir.BinaryExpression:
		return l.lowerBinaryExpression(expr, statements)
	case *ir.UnaryExpression:
		return l.lowerUnaryExpression(expr, statements)
	case *ir.Conversion:
		return l.lowerConversion(expr, statements)
	case *ir.InvalidExpression:
		return l.lowerInvalidExpression(expr, statements)
	case *ir.ArrayExpression:
		return l.lowerArrayExpression(expr, statements)
	case *ir.IndexExpression:
		return l.lowerIndexExpression(expr, statements)
	case *ir.MapExpression:
		return l.lowerMapExpression(expr, statements)
	case *ir.Assignment:
		return l.lowerAssignment(expr, statements)
	case *ir.TupleExpression:
		return l.lowerTupleExpression(expr, statements)
	case *ir.TypeCheck:
		return l.lowerTypeCheck(expr, statements)
	case *ir.FunctionCall:
		return l.lowerFunctionCall(expr, statements)
	case *ir.StructExpression:
		return l.lowerStructExpression(expr, statements)
	case *ir.TupleStructExpression:
		return l.lowerTupleStructExpression(expr, statements)
	case *ir.MemberExpression:
		return l.lowerMemberExpression(expr, statements)
	case *ir.Block:
		return l.lowerBlock(expr, statements)
	case *ir.IfExpression:
		return l.lowerIfExpression(expr, statements, nil)
	case *ir.WhileLoop:
		return l.lowerWhileLoop(expr, statements)
	case *ir.ForLoop:
		return l.lowerForLoop(expr, statements)
	case *ir.TypeExpression:
		return l.lowerTypeExpression(expr, statements)
	case *ir.FunctionExpression:
		return l.lowerFunctionExpression(expr, statements)
	case *ir.RefExpression:
		return l.lowerRefExpression(expr, statements)
	case *ir.DerefExpression:
		return l.lowerDerefExpression(expr, statements)

	default:
		panic(fmt.Sprintf("TODO: lower %T", expr))
	}
}
