package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/module"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
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
	t.registerDeclarations(mod)
	t.typeCheckFunctions(mod)
	stmts := t.typeCheckStatements(mod)

	return &ir.Program{
		Statements: stmts,
	}
}

func (t *typeChecker) registerDeclarations(mod *module.Module) {
	for _, file := range mod.Files {
		for _, stmt := range file.Ast.Statements {
			 t.registerDeclaration(stmt)
		}
	}
}

func (t *typeChecker) typeCheckFunctions(mod *module.Module) {
	for _, file := range mod.Files {
		for _, stmt := range file.Ast.Statements {
			 if fn, ok := stmt.(*ast.FunctionDeclaration); ok {
				t.typeCheckFunctionType(fn)
			 }
		}
	}
}

func (t *typeChecker) typeCheckStatements(mod *module.Module) []ir.Statement {
	stmts := []ir.Statement{}
	for _, file := range mod.Files {
		for _, stmt := range file.Ast.Statements {
			stmts = append(stmts, t.typeCheckStatement(stmt))
		}
	}

	return stmts
}

func (t *typeChecker) enterScope(context ...any) {
	if len(context) > 0 {
		t.symbols = t.symbols.ChildWithContext(context[0])
	} else {
		t.symbols = t.symbols.Child()
	}
}

func (t *typeChecker) exitScope() {
	t.symbols = t.symbols.Parent
}

type conversionKind int

const (
	none conversionKind = iota
	identity
	implicit
	operator
	explicit
)

func canConvert(from, to types.Type) conversionKind {
	kind := none

	if types.Assignable(to, from) {
		kind = identity
	} else if from == types.Int && to == types.Float {
		kind = operator
	} else if from == types.Float && to == types.Int {
		kind = explicit
	} else if from == types.Bool && to == types.Int {
		kind = explicit
	}

	if v, ok := from.(types.VariableType); ok &&
		v.Untyped && kind == identity {
		return implicit
	}

	return kind
}

func convert(from ir.Expression, to types.Type, maxKind conversionKind) ir.Expression {
	kind := canConvert(from.Type(), to)

	if kind == identity {
		return from
	}
	if kind == none {
		return nil
	}

	if kind <= maxKind {
		return &ir.Conversion{
			Expression: from,
			To:         to,
		}
	}

	return nil
}
