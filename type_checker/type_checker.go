package typechecker

import (
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TypeCheck(program *ast.Program, manager *modules.ModuleManager) error {
	if manager.TypeChecked {
		return nil
	}
	manager.TypeChecked = true

	err := typeCheckGlobalScope(program, manager)
	if err != nil {
		return err
	}

	for _, stmt := range program.Body {
		nextType := typeCheckStatement(stmt, manager)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}
	return nil
}

func typeCheckGlobalScope(program *ast.Program, manager *modules.ModuleManager) error {
	for _, stmt := range program.Body {
		nextType := registerTypeStatement(stmt, manager)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}

	for _, stmt := range program.Body {
		nextType := typeCheckGlobalStatement(stmt, manager)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}

	for _, stmt := range program.Body {
		if importStmt, ok := stmt.(*ast.ImportStatement); ok {
			nextType := typeCheckImportStatement(importStmt, manager)
			if nextType.String() == "TypeError" {
				return nextType.(*types.TypeError)
			}
		}
	}


	return nil
}
