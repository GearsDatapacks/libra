package lowerer

import (
	"fmt"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type lowerer struct {
	currentModule *ir.LoweredModule
	diagnostics   diagnostics.Manager
	labelId       int
	varId         int
	scope         *scope
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
	endLabel      string
	yieldVariable symbols.Variable
}

func makeMain() *ir.FunctionDeclaration {
	return &ir.FunctionDeclaration{
		Name:       "main",
		Parameters: []string{},
		Body: &ir.Block{
			Statements: []ir.Statement{},
			ResultType: types.Void,
		},
		Type: &types.Function{
			Parameters: []types.Type{},
			ReturnType: types.Void,
		},
		Exported: false,
	}
}

func Lower(pkg *ir.Package, diagnostics diagnostics.Manager) (*ir.LoweredPackage, diagnostics.Manager) {
	lowerer := lowerer{
		diagnostics: diagnostics,
	}

	lowered := &ir.LoweredPackage{
		Modules: map[string]*ir.LoweredModule{},
	}

	for name, module := range pkg.Modules {
		mainFunction, definesMain := getMain(module)

		mod := &ir.LoweredModule{
			Name:         name,
			MainFunction: mainFunction,
			Types:        []*ir.TypeDeclaration{},
			Functions:    []*ir.FunctionDeclaration{mainFunction},
			Globals:      []*ir.VariableDeclaration{},
		}
		lowered.Modules[name] = mod
		lowerer.currentModule = mod

		for _, stmt := range module.Statements {
			lowerer.lowerGlobal(stmt, mod, !definesMain)
		}
		// Only compile main if there are any statements there
		if definesMain || len(mainFunction.Body.Statements) == 0 {
			mod.Functions = mod.Functions[1:]
		} else {
			mainFunction.Body.Statements = lowerer.cfa(mainFunction.Body.Statements, nil, false)
		}
	}
	fixAbi(lowered)
	return lowered, lowerer.diagnostics
}

func getMain(module *ir.Module) (*ir.FunctionDeclaration, bool) {
	for _, statement := range module.Statements {
		if funcDecl, ok := statement.(*ir.FunctionDeclaration); ok && funcDecl.Name == "main" {
			return funcDecl, true
		}
	}

	return makeMain(), false
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

func (l *lowerer) lowerGlobal(
	statement ir.Statement,
	mod *ir.LoweredModule,
	allowArbitraryStatements bool,
) {
	switch stmt := statement.(type) {
	case *ir.FunctionDeclaration:
		mod.Functions = append(mod.Functions, l.lowerFunctionDeclaration(stmt))
	case *ir.TypeDeclaration:
		mod.Types = append(mod.Types, l.lowerTypeDeclaration(stmt))
	case *ir.ImportStatement:
		mod.Imports = append(mod.Imports, l.lowerImportStatement(stmt))
	default:
		// TODO: Allow global variable/constant declarations here
		if allowArbitraryStatements {
			l.lower(statement, &mod.MainFunction.Body.Statements)
		} else {
			l.diagnostics.Report(diagnostics.NonDeclOutsideMain(statement.GetLocation()))
		}
	}
}

func (l *lowerer) lower(statement ir.Statement, statements *[]ir.Statement) {
	switch stmt := statement.(type) {
	case *ir.VariableDeclaration:
		l.lowerVariableDeclaration(stmt, statements)
	case *ir.ReturnStatement:
		l.lowerReturnStatement(stmt, statements)
	case *ir.BreakStatement:
		l.lowerBreakStatement(stmt, statements)
	case *ir.ContinueStatement:
		l.lowerContinueStatement(stmt, statements)
	case *ir.YieldStatement:
		l.lowerYieldStatement(stmt, statements)

	case *ir.TypeDeclaration, *ir.FunctionDeclaration, *ir.ImportStatement:
		panic("Declarations not allowed here")

	case ir.Expression:
		expr := l.lowerExpression(stmt, statements, false)
		if expr != nil {
			*statements = append(*statements, expr)
		}

	default:
		panic(fmt.Sprintf("TODO: lower %T", stmt))
	}
}

func (l *lowerer) lowerExpression(expression ir.Expression, statements *[]ir.Statement, used bool) ir.Expression {
	var lowered ir.Expression
	switch expr := expression.(type) {
	case *ir.IntegerLiteral:
		if !used {
			return nil
		}
		lowered = l.lowerIntegerLiteral(expr, statements)
	case *ir.FloatLiteral:
		if !used {
			return nil
		}
		lowered = l.lowerFloatLiteral(expr, statements)
	case *ir.BooleanLiteral:
		if !used {
			return nil
		}
		lowered = l.lowerBooleanLiteral(expr, statements)
	case *ir.StringLiteral:
		if !used {
			return nil
		}
		lowered = l.lowerStringLiteral(expr, statements)
	case *ir.VariableExpression:
		if !used {
			return nil
		}
		lowered = l.lowerVariableExpression(expr, statements)
	case *ir.BinaryExpression:
		lowered = l.lowerBinaryExpression(expr, statements, used)
	case *ir.UnaryExpression:
		lowered = l.lowerUnaryExpression(expr, statements, used)
	case *ir.Conversion:
		lowered = l.lowerConversion(expr, statements, used)
	case *ir.InvalidExpression:
		if !used {
			return nil
		}
		lowered = l.lowerInvalidExpression(expr, statements)
	case *ir.ArrayExpression:
		lowered = l.lowerArrayExpression(expr, statements)
	case *ir.IndexExpression:
		lowered = l.lowerIndexExpression(expr, statements, used)
	case *ir.MapExpression:
		lowered = l.lowerMapExpression(expr, statements)
	case *ir.Assignment:
		lowered = l.lowerAssignment(expr, statements)
	case *ir.TupleExpression:
		lowered = l.lowerTupleExpression(expr, statements)
	case *ir.TypeCheck:
		lowered = l.lowerTypeCheck(expr, statements, used)
	case *ir.FunctionCall:
		lowered = l.lowerFunctionCall(expr, statements)
	case *ir.StructExpression:
		lowered = l.lowerStructExpression(expr, statements)
	case *ir.TupleStructExpression:
		lowered = l.lowerTupleStructExpression(expr, statements)
	case *ir.MemberExpression:
		lowered = l.lowerMemberExpression(expr, statements, used)
	case *ir.Block:
		lowered = l.lowerBlock(expr, statements, used)
	case *ir.IfExpression:
		lowered = l.lowerIfExpression(expr, statements, nil, used)
	case *ir.WhileLoop:
		lowered = l.lowerWhileLoop(expr, statements, used)
	case *ir.ForLoop:
		lowered = l.lowerForLoop(expr, statements)
	case *ir.TypeExpression:
		if !used {
			return nil
		}
		lowered = l.lowerTypeExpression(expr, statements)
	case *ir.FunctionExpression:
		if !used {
			return nil
		}
		lowered = l.lowerFunctionExpression(expr, statements)
	case *ir.RefExpression:
		lowered = l.lowerRefExpression(expr, statements, used)
	case *ir.DerefExpression:
		lowered = l.lowerDerefExpression(expr, statements, used)

	default:
		panic(fmt.Sprintf("TODO: lower %T", expr))
	}

	if lowered == nil {
		return nil
	}
	return optimiseExpression(lowered)
}
