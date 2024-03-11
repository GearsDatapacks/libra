package typechecker_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/module"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/parser/ast"
	utils "github.com/gearsdatapacks/libra/test_utils"
	"github.com/gearsdatapacks/libra/text"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TestIntegerLiteral(t *testing.T) {
	input := "1_23_456"
	val := 123456

	program := getProgram(t, input)

	integer := getExpr[*ir.IntegerLiteral](t, program)

	utils.AssertEq(t, integer.Value, int64(val))
	utils.AssertEq[types.Type](t, integer.Type(), types.Int)
}

func TestFloatLiteral(t *testing.T) {
	input := "3.14_15_9"
	val := 3.14159

	program := getProgram(t, input)

	float := getExpr[*ir.FloatLiteral](t, program)

	utils.AssertEq(t, float.Value, val)
	utils.AssertEq[types.Type](t, float.Type(), types.Float)
}

func TestBooleanLiteral(t *testing.T) {
	input := "true"

	program := getProgram(t, input)

	boolean := getExpr[*ir.BooleanLiteral](t, program)

	utils.AssertEq(t, boolean.Value, true)
	utils.AssertEq[types.Type](t, boolean.Type(), types.Bool)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hi \"there\\"`
	val := "Hi \"there\\"

	program := getProgram(t, input)

	str := getExpr[*ir.StringLiteral](t, program)

	utils.AssertEq(t, str.Value, val)
	utils.AssertEq[types.Type](t, str.Type(), types.String)
}

func TesVariables(t *testing.T) {
	tests := []struct {
		src      string
		varName  string
		mutable  bool
		dataType types.Type
	}{
		{"let x = 1", "x", false, types.Int},
		{"mut foo: f32 = 1.4", "foo", true, types.Float},
		{`const greeting: string = "Hi!"`, "greeting", false, types.String},
		{"mut is_awesome = true", "is_awesome", true, types.Bool},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		varDec := getStmt[*ir.VariableDeclaration](t, program)
		utils.AssertEq(t, varDec.Name, test.varName)

		program = getProgram(t, test.src+"\n"+test.varName)
		utils.AssertEq(t, len(program.Statements), 2)
		exprStmt, ok := program.Statements[1].(*ir.ExpressionStatement)
		utils.Assert(t, ok, "Statement is not an expressions statement")
		variable, ok := exprStmt.Expression.(*ir.VariableExpression)
		utils.Assert(t, ok, "Expression is not a variable")

		utils.AssertEq(t, variable.Symbol.Name, test.varName)
		utils.AssertEq(t, variable.Symbol.Mutable, test.mutable)
		utils.AssertEq(t, variable.Symbol.Type, test.dataType)
	}
}

type diagnostic struct {
	message string
	kind    diagnostics.DiagnosticKind
}

func TestTCDiagnostics(t *testing.T) {
	tests := []struct {
		src         string
		diagnostics []diagnostic
	}{
		{"let x: [foo] = 1", []diagnostic{{`Type "foo" is not defined`, diagnostics.Error}}},
		{"const text: string = [false]", []diagnostic{{`Value of type "bool" is not assignable to type "string"`, diagnostics.Error}}},
		{"let foo = 1; let [foo] = 2", []diagnostic{{`Variable "foo" is already defined`, diagnostics.Error}}},
		{"let a = [b]", []diagnostic{{`Variable "b" is not defined`, diagnostics.Error}}},
	}

	for _, test := range tests {
		src, spans := getSpans(test.src)
		utils.AssertEq(t, len(spans), len(test.diagnostics), "Mismatch of spans to diagnostic messages")

		lexer := lexer.New(text.NewFile("test.lb", src))
		tokens := lexer.Tokenise()

		utils.AssertEq(t, len(lexer.Diagnostics), 0, "Expected no lexer diagnostics")

		p := parser.New(tokens, lexer.Diagnostics)
		program := p.Parse()
		utils.AssertEq(t, len(p.Diagnostics), 0, "Expected no parser diagnostics")
		tc := typechecker.New(p.Diagnostics)
		tc.TypeCheck(fakeModule(program))

		utils.AssertEq(t, len(tc.Diagnostics), len(test.diagnostics),
			fmt.Sprintf("Incorrect number of diagnostics (expected %d, got %d)", len(test.diagnostics), len(tc.Diagnostics)))

		for i, diag := range test.diagnostics {
			// FIXME: Do this in a better way
			span := spans[len(spans)-i-1]
			testDiagnostic(t, tc.Diagnostics[i], diag.kind, diag.message, span)
		}
	}
}

func fakeModule(program *ast.Program) *module.Module {
	return &module.Module{
		Name: "test",
		Files: []module.File{{
			Path: "test.lb",
			Ast:  program,
		}},
		Imported: map[string]*module.Module{},
	}
}

func getSpans(sourceText string) (string, []text.Span) {
	var resultText bytes.Buffer
	spans := []text.Span{}
	line := 0
	col := 0
	for _, c := range sourceText {
		if c == '[' {
			spans = append(spans, text.NewSpan(line, col, 0))
			continue
		}
		if c == ']' {
			spans[len(spans)-1].End = col
			continue
		}

		col++
		if c == '\n' {
			line++
			col = 0
		}
		resultText.WriteRune(c)
	}

	return resultText.String(), spans
}

func testDiagnostic(t *testing.T,
	diagnostic diagnostics.Diagnostic,
	kind diagnostics.DiagnosticKind,
	msg string,
	span text.Span) {
	utils.AssertEq(t, diagnostic.Kind, kind)
	utils.AssertEq(t, diagnostic.Message, msg)
	utils.AssertEq(t, diagnostic.Location.Span, span)
}

func getProgram(t *testing.T, input string) *ir.Program {
	t.Helper()

	l := lexer.New(text.NewFile("test.lb", input))
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()
	tc := typechecker.New(p.Diagnostics)
	irProgram := tc.TypeCheckProgram(program)
	utils.AssertEq(t, len(tc.Diagnostics), 0,
		fmt.Sprintf("Expected no diagnostics (got %d)", len(tc.Diagnostics)))

	return irProgram
}

func getStmt[T ir.Statement](t *testing.T, program *ir.Program) T {
	t.Helper()

	stmt, ok := utils.AssertSingle(t, program.Statements).(T)
	utils.Assert(t, ok, fmt.Sprintf(
		"Statement is not %T (is %T)", struct{ t T }{}.t, stmt))
	return stmt
}

func getExpr[T ir.Expression](t *testing.T, program *ir.Program) T {
	t.Helper()

	exprStmt := getStmt[*ir.ExpressionStatement](t, program)
	expr, ok := exprStmt.Expression.(T)
	utils.Assert(t, ok, fmt.Sprintf(
		"expression is not %T (is %T)", struct{ t T }{}.t, expr))
	return expr
}
