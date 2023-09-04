package typechecker

import (
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TypeCheck(program ast.Program) error {
	symbolTable := symbols.New()

	for _, stmt := range program.Body {
		nextType := typeCheckStatement(stmt, symbolTable)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}
	return nil
}
