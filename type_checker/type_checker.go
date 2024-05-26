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
	stage        tcStage
}

var mods = map[string]*typeChecker{}

func new(mod *module.Module, diagnostics *diagnostics.Manager) *typeChecker {
	t := &typeChecker{
		diagnostics: diagnostics,
		module:      mod,
		symbols:     symbols.New(),
		subModules:  map[string]*typeChecker{},
		stage:        tcNone,
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

	t.registerDeclarations()
	t.typeCheckImports()
	t.typeCheckDeclarations()
	t.typeCheckFunctions()

	stmts := t.typeCheckStatements()

	return &ir.Program{
		Statements: stmts,
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

func (t *typeChecker) typeCheckImports() {
	if t.stage >= tcImports {
		return
	}
	t.stage = tcImports

	for _, subMod := range t.subModules {
		subMod.typeCheckImports()
	}

	t.updateContext()
	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			if importStmt, ok := stmt.(*ast.ImportStatement); ok {
				t.typeCheckImport(importStmt)
			}
		}
	}
}

func (t *typeChecker) typeCheckDeclarations() {
	if t.stage >= tcDecls {
		return
	}
	t.stage = tcDecls

	for _, subMod := range t.subModules {
		subMod.typeCheckDeclarations()
	}

	t.updateContext()
	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			t.typeCheckDeclaration(stmt)
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

// TODO: return ir for other modules too
func (t *typeChecker) typeCheckStatements() []ir.Statement {
	if t.stage >= tcStmts {
		return []ir.Statement{}
	}
	t.stage = tcStmts
	
	for _, subMod := range t.subModules {
		subMod.typeCheckStatements()
	}
	
	t.updateContext()
	stmts := []ir.Statement{}
	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			stmts = append(stmts, t.typeCheckStatement(stmt))
		}
	}

	return stmts
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
