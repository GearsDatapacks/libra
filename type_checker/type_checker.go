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

func TypeCheck(mod *module.Module, manager diagnostics.Manager) (*ir.Package, diagnostics.Manager) {
	t := new(mod, &manager)

	pkg := &ir.Package{
		Modules: map[string]*ir.Module{},
	}

	t.registerDeclarations()
	t.typeCheckImports(pkg)
	t.typeCheckDeclarations(pkg)
	t.typeCheckFunctions()

	t.typeCheckStatements(pkg)

	return pkg, *t.diagnostics
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

func (t *typeChecker) typeCheckImports(pkg *ir.Package) {
	if t.stage >= tcImports {
		return
	}
	t.stage = tcImports

	for _, subMod := range t.subModules {
		subMod.typeCheckImports(pkg)
	}

	t.updateContext()
	if _, ok := pkg.Modules[t.module.Path]; !ok {
		pkg.Modules[t.module.Path] = &ir.Module{
			Name:       t.module.Name,
			Statements: []ir.Statement{},
		}
	}
	module := pkg.Modules[t.module.Path]

	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			if importStmt, ok := stmt.(*ast.ImportStatement); ok {
				stmt := t.typeCheckImport(importStmt)
				if stmt != nil {
					module.Statements = append(module.Statements, stmt)
				}
			}
		}
	}
}

func (t *typeChecker) typeCheckDeclarations(pkg *ir.Package) {
	if t.stage >= tcDecls {
		return
	}
	t.stage = tcDecls

	for _, subMod := range t.subModules {
		subMod.typeCheckDeclarations(pkg)
	}

	t.updateContext()
	module := pkg.Modules[t.module.Path]

	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			decl := t.typeCheckDeclaration(stmt)

			if decl != nil {
				module.Statements = append(module.Statements, decl)
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

func (t *typeChecker) typeCheckStatements(pkg *ir.Package) {
	if t.stage >= tcStmts {
		return
	}
	t.stage = tcStmts

	for _, subMod := range t.subModules {
		subMod.typeCheckStatements(pkg)
	}

	t.updateContext()
	module := pkg.Modules[t.module.Path]

	for _, file := range t.module.Files {
		for _, stmt := range file.Ast.Statements {
			nextStatement := t.typeCheckStatement(stmt)
			if nextStatement != nil {
				module.Statements = append(module.Statements, nextStatement)
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
