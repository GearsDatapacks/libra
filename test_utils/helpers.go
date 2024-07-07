package testutils

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/module"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gkampitakis/go-snaps/snaps"
)

func getAst(t *testing.T, input string) *ast.Program {
	t.Helper()

	l := lexer.New(text.NewFile("test.lb", input))
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()
	AssertEq(t, len(p.Diagnostics), 0,
		fmt.Sprintf("Expected no diagnostics (got %d)", len(p.Diagnostics)))

	return program
}

func getIr(t *testing.T, input string) *ir.Program {
	t.Helper()

	l := lexer.New(text.NewFile("test.lb", input))
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()
	irProgram, diags := typechecker.TypeCheck(fakeModule(program), p.Diagnostics)
	AssertEq(t, len(diags), 0,
		fmt.Sprintf("Expected no diagnostics (got %d)", len(diags)))

	return irProgram
}

func fakeModule(program *ast.Program) *module.Module {
	return &module.Module{
		Id:       1,
		Name:     "test",
		Files:    []module.File{{Path: "test.lb", Ast: program}},
		Imported: map[string]*module.Module{},
	}
}

func matchSnap(t *testing.T, values ...any) {
	t.Helper()

	pc, _, _, _ := runtime.Caller(2)
	name := runtime.FuncForPC(pc).Name()
	parts := strings.Split(name, "/")
	name = parts[len(parts)-1]
	snaps := snaps.WithConfig(snaps.Filename(name))

	snaps.MatchSnapshot(t, values...)
}
