package testutils

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/module"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gkampitakis/go-snaps/snaps"
)

func getAst(t *testing.T, input string) (*ast.Program, []diagnostics.Diagnostic) {
	t.Helper()

	l := lexer.New(text.NewFile("test.lb", input))
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()

	return program, p.Diagnostics
}

func getIr(t *testing.T, input string) (*ir.Program, []diagnostics.Diagnostic) {
	t.Helper()

	l := lexer.New(text.NewFile("test.lb", input))
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()
	return typechecker.TypeCheck(fakeModule(program), p.Diagnostics)
}

func fakeModule(program *ast.Program) *module.Module {
	return &module.Module{
		Id:       1,
		Name:     "test",
		Files:    []module.File{{Path: "test.lb", Ast: program}},
		Imported: map[string]*module.Module{},
	}
}

type namedSnapshot struct {
	name string
	*testing.T
}

func (t *namedSnapshot) Name() string {
	return t.name
}

func matchSnap(t *testing.T, src, output string) {
	t.Helper()

	pc, _, _, _ := runtime.Caller(2)
	name := runtime.FuncForPC(pc).Name()
	parts := strings.Split(name, "/")
	name = parts[len(parts)-1]
	snaps := snaps.WithConfig(snaps.Filename(name))

	snaps.MatchSnapshot(&namedSnapshot{name: fmt.Sprintf("`%s`", strings.ReplaceAll(src, "\n", " ")), T: t}, output)
}
