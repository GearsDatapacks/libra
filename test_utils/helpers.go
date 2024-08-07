package testutils

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/lowerer"
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

func getIr(t *testing.T, input string) (*ir.Package, []diagnostics.Diagnostic) {
	t.Helper()

	l := lexer.New(text.NewFile("test.lb", input))
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()
	return typechecker.TypeCheck(fakeModule(program), p.Diagnostics)
}

func getLowered(t *testing.T, input string) (*ir.LoweredPackage, []diagnostics.Diagnostic) {
	t.Helper()

	l := lexer.New(text.NewFile("test.lb", input))
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()
	pkg, diags := typechecker.TypeCheck(fakeModule(program), p.Diagnostics)
	return lowerer.Lower(pkg, diags)
}

func fakeModule(program *ast.Program) *module.Module {
	return &module.Module{
		Id:       1,
		Name:     "test",
		Path:     "test",
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

	var name string
	for i := 0; ; i++ {
		pc, _, _, _ := runtime.Caller(i)
		nextName := runtime.FuncForPC(pc).Name()
		if nextName == "testing.tRunner" {
			break
		}
		parts := strings.Split(nextName, "/")
		name = parts[len(parts)-1]
	}

	snaps := snaps.WithConfig(snaps.Filename(name))

	snaps.MatchSnapshot(&namedSnapshot{
		name: fmt.Sprintf("`%s`", strings.ReplaceAll(src, "\n", ";")),
		T:    t,
	}, output)
}
