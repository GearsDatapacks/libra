package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/module"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type tcStage int

const (
	tcNone tcStage = iota
	tcRegister
	tcImports
	tcDecls
	tcFns
	tcStmts
)

type typeChecker struct {
	diagnostics *diagnostics.Manager
	module      *module.Module
	symbols     *symbols.Table
	subModules  map[string]*typeChecker
	stage       tcStage
}

var mods = map[string]*typeChecker{}

func new(mod *module.Module, diagnostics *diagnostics.Manager) *typeChecker {
	t := &typeChecker{
		diagnostics: diagnostics,
		module:      mod,
		symbols:     symbols.New(),
		subModules:  map[string]*typeChecker{},
		stage:       tcNone,
	}
	mods[mod.Path] = t

	for name, subMod := range mod.Imported {
		if mod, ok := mods[subMod.Path]; ok {
			t.subModules[name] = mod
			continue
		}
		mods[subMod.Path] = new(subMod, diagnostics)
		t.subModules[name] = mods[subMod.Path]
	}
	return t
}

type typeContext struct {
	*symbols.Table
	id uint
}

func (t *typeContext) Id() uint {
	return t.id
}

func TypeCheck(mod *module.Module, manager diagnostics.Manager) (*ir.Program, diagnostics.Manager) {
	t := new(mod, &manager)

	statements := []ir.Statement{}

	t.registerDeclarations()
	t.typeCheckImports(&statements)
	t.typeCheckDeclarations(&statements)
	t.typeCheckFunctions()

	t.typeCheckStatements(&statements)

	return &ir.Program{
		Statements: statements,
	}, *t.diagnostics
}

func (t *typeChecker) updateContext() {
	types.Context = &typeContext{
		Table: t.symbols,
		id:    t.module.Id,
	}
}

func (t *typeChecker) registerDeclarations() {
	if t.stage >= tcRegister {
		return
	}
	t.stage = tcRegister

	for _, subMod := range t.subModules {
		subMod.registerDeclarations()
	}

	t.updateContext()
	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			t.registerDeclaration(stmt)
		}
	}
}

func (t *typeChecker) typeCheckImports(statements *[]ir.Statement) {
	if t.stage >= tcImports {
		return
	}
	t.stage = tcImports

	for _, subMod := range t.subModules {
		subMod.typeCheckImports(&[]ir.Statement{})
	}

	t.updateContext()
	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			if importStmt, ok := stmt.(*ast.ImportStatement); ok {
				stmt := t.typeCheckImport(importStmt)
				if stmt != nil {
					*statements = append(*statements, stmt)
				}
			}
		}
	}
}

func (t *typeChecker) typeCheckDeclarations(statements *[]ir.Statement) {
	if t.stage >= tcDecls {
		return
	}
	t.stage = tcDecls

	for _, subMod := range t.subModules {
		subMod.typeCheckDeclarations(&[]ir.Statement{})
	}

	t.updateContext()
	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			decl := t.typeCheckDeclaration(stmt)

			if decl != nil {
				*statements = append(*statements, decl)
			}
		}
	}
}

func (t *typeChecker) typeCheckFunctions() {
	if t.stage >= tcFns {
		return
	}
	t.stage = tcFns

	for _, subMod := range t.subModules {
		subMod.typeCheckFunctions()
	}

	t.updateContext()
	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			if fn, ok := stmt.(*ast.FunctionDeclaration); ok {
				t.typeCheckFunctionType(fn)
			}
		}
	}
}

func (t *typeChecker) typeCheckStatements(statements *[]ir.Statement) {
	if t.stage >= tcStmts {
		return
	}
	t.stage = tcStmts

	for _, subMod := range t.subModules {
		// TODO: return ir for other modules too
		subMod.typeCheckStatements(&[]ir.Statement{})
	}

	t.updateContext()
	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			nextStatement := t.typeCheckStatement(stmt)
			if nextStatement != nil {
				*statements = append(*statements, nextStatement)
			}
		}
	}
}

func (t *typeChecker) enterScope(context ...any) {
	if len(context) > 0 {
		t.symbols = t.symbols.ChildWithContext(context[0])
	} else {
		t.symbols = t.symbols.Child()
	}
}

func (t *typeChecker) exitScope() {
	t.symbols = t.symbols.Parent
}

type conversionKind int

const (
	none conversionKind = iota
	identity
	implicit
	operator
	explicit
)

func canConvert(from, to types.Type) conversionKind {
	kind := none

	if expl, ok := to.(*types.Explicit); ok && types.Assignable(expl.Type, from) {
		return implicit
	}
	if expl, ok := from.(*types.Explicit); ok {
		if canConvert(expl.Type, to) != none {
			return explicit
		}
	}

	if types.Assignable(to, from) {
		kind = identity
	} else if from == types.Int && to == types.Float {
		kind = operator
	} else if from == types.Float && to == types.Int {
		kind = explicit
	} else if from == types.Bool && to == types.Int {
		kind = explicit
	}

	if v, ok := from.(types.VariableType); ok &&
		v.Untyped && kind == identity {
		return implicit
	}

	if _, ok := to.(*types.Union); ok && kind == identity {
		return implicit
	}

	return kind
}

func convert(from ir.Expression, to types.Type, maxKind conversionKind) ir.Expression {
	kind := canConvert(from.Type(), to)

	if kind == identity {
		return from
	}
	if kind == none {
		return nil
	}

	if kind <= maxKind {
		return &ir.Conversion{
			Expression: from,
			To:         to,
		}
	}

	return nil
}
