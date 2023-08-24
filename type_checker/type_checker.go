package typechecker

import (
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
)

func TypeCheck(program ast.Program) {
	symbolTable := symbols.New()

	for _, stmt := range program.Body {
		typeCheckStatement(stmt, symbolTable)
	}
}
