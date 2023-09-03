package typechecker

import (
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
)

func TypeCheck(program ast.Program) error {
	symbolTable := symbols.New()

	for _, stmt := range program.Body {
		_, err := typeCheckStatement(stmt, symbolTable)
		if err != nil {
			return err
		}
	}
	return nil
}
