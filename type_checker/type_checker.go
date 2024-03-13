package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/module"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type typeChecker struct {
	Diagnostics diagnostics.Manager
	symbols     *symbols.Table
}

func New(diagnostics diagnostics.Manager) *typeChecker {
	return &typeChecker{
		Diagnostics: diagnostics,
		symbols:     symbols.New(),
	}
}

func (t *typeChecker) TypeCheckProgram(program *ast.Program) *ir.Program {
	stmts := []ir.Statement{}

	for _, stmt := range program.Statements {
		stmts = append(stmts, t.typeCheckStatement(stmt))
	}

	return &ir.Program{
		Statements: stmts,
	}
}

func (t *typeChecker) TypeCheck(mod *module.Module) *ir.Program {
	stmts := []ir.Statement{}

	for _, file := range mod.Files {
		for _, stmt := range file.Ast.Statements {
			stmts = append(stmts, t.typeCheckStatement(stmt))
		}
	}

	return &ir.Program{
		Statements: stmts,
	}
}

func canConvert(from, to types.Type) (exists, explicit bool) {
	if types.Assignable(to, from) {
		return false, true
	}

	if from == types.Int && to == types.Float {
		return true, false
	}

	if from == types.Float && to == types.Int {
		return true, true
	}

	if from == types.Bool && to == types.Int {
		return true, true
	}

	return false, false
}

func convert(from ir.Expression, to types.Type, allowExplicit bool) ir.Expression {
	exists, explicit := canConvert(from.Type(), to)

	if !exists && explicit {
		return from
	}
	if exists && (!explicit || allowExplicit) {
		return &ir.Conversion{
			Expression: from,
			To:         to,
		}
	}

	return nil
}
