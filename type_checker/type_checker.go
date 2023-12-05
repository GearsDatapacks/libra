package typechecker

import (
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TypeCheck(program *ast.Program, symbolTable *symbols.SymbolTable) error {
	err := typeCheckGlobalScope(program, symbolTable)
	if err != nil {
		return err
	}

	for _, stmt := range program.Body {
		nextType := typeCheckStatement(stmt, symbolTable)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}
	return nil
}

func typeCheckGlobalScope(program *ast.Program, symbolTable *symbols.SymbolTable) error {
	for _, stmt := range program.Body {
		nextType := registerTypeStatement(stmt, symbolTable)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}

	for _, stmt := range program.Body {
		nextType := typeCheckGlobalStatement(stmt, symbolTable)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}

	return nil
}
