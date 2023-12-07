package modules

import (
	"os"

	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/parser/ast"
)

type Module struct {
	Ast ast.Program
	Exports []ast.Statement
}

func Get(file string) (*Module, error) {
	code, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	lexer := lexer.New(code)
	tokens, err := lexer.Tokenise()
	if err != nil {
		return nil, err
	}

	parser := parser.New()
	program, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	exports := []ast.Statement{}
	for _, stmt := range program.Body {
		if stmt.IsExport() {
			exports = append(exports, stmt)
		}
	}

	return &Module{
		Ast: program,
		Exports: exports,
	}, nil
}
