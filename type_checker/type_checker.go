package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/module"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
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
