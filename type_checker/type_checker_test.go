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
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func TestIntegerLiteral(t *testing.T) {
	input := "+1_23_456"
	val := 123456

	program := getProgram(t, input)

	integer := getExpr[*ir.IntegerLiteral](t, program)

	utils.AssertEq(t, integer.Value, int64(val))
	utils.AssertEq[types.Type](t, integer.Type(), types.UntypedInt)
}

func TestFloatLiteral(t *testing.T) {
	input := "3.14_15_9"
	val := 3.14159

	program := getProgram(t, input)

	float := getExpr[*ir.FloatLiteral](t, program)

	utils.AssertEq(t, float.Value, val)
	utils.AssertEq[types.Type](t, float.Type(), types.UntypedFloat)
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

func TestVariables(t *testing.T) {
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
		{"const my_float: f32 = 15", "my_float", false, types.Float},
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

func TestBinaryExpression(t *testing.T) {
	tests := []struct {
		src    string
		left   types.Type
		op     ir.BinaryOperator
		right  types.Type
		result types.Type
	}{
		{"true && false", types.Bool, ir.LogicalAnd, types.Bool, types.Bool},
		{"false || false", types.Bool, ir.LogicalOr, types.Bool, types.Bool},
		{"1.5 < 2", types.Float, ir.Less, types.Float, types.Bool},
		{"17 <= 17", types.UntypedInt, ir.LessEq, types.UntypedInt, types.Bool},
		{"3.14 > 2.71", types.Float, ir.Greater, types.Float, types.Bool},
		{"42 >= 69", types.UntypedInt, ir.GreaterEq, types.UntypedInt, types.Bool},
		{"1 == 2", types.UntypedInt, ir.Equal, types.UntypedInt, types.Bool},
		{"true == true", types.Bool, ir.Equal, types.Bool, types.Bool},
		{"1.2 != 7.5", types.UntypedFloat, ir.NotEqual, types.UntypedFloat, types.Bool},
		{"1 << 5", types.UntypedInt, ir.LeftShift, types.UntypedInt, types.UntypedInt},
		{"8362 >> 3", types.UntypedInt, ir.RightShift, types.UntypedInt, types.UntypedInt},
		{"10101 | 1010", types.UntypedInt, ir.BitwiseOr, types.UntypedInt, types.UntypedInt},
		{"73 & 52", types.UntypedInt, ir.BitwiseAnd, types.UntypedInt, types.UntypedInt},
		{"1 + 6", types.UntypedInt, ir.AddInt, types.UntypedInt, types.UntypedInt},
		{"2.3 + 4", types.Float, ir.AddFloat, types.Float, types.UntypedFloat},
		{`"Hello " + "world"`, types.String, ir.Concat, types.String, types.String},
		{"8 - 12", types.UntypedInt, ir.SubtractInt, types.UntypedInt, types.UntypedInt},
		{"3 - 1.3", types.Float, ir.SubtractFloat, types.Float, types.UntypedFloat},
		{"6 * 7", types.UntypedInt, ir.MultiplyInt, types.UntypedInt, types.UntypedInt},
		{"1.3 * 0.4", types.Float, ir.MultiplyFloat, types.Float, types.UntypedFloat},
		{"0.3 / 2", types.Float, ir.Divide, types.Float, types.UntypedFloat},
		{"103 % 2", types.UntypedInt, ir.ModuloInt, types.UntypedInt, types.UntypedInt},
		{"1.4 % 1", types.Float, ir.ModuloFloat, types.Float, types.UntypedFloat},
		{"2 ** 7", types.UntypedInt, ir.PowerInt, types.UntypedInt, types.UntypedInt},
		{"3 ** 0.5", types.Float, ir.PowerFloat, types.Float, types.UntypedFloat},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		binExpr := getExpr[*ir.BinaryExpression](t, program)

		utils.AssertEq(t, binExpr.Left.Type(), test.left)
		op := binExpr.Operator & ^ir.UntypedBit
		utils.AssertEq(t, op, test.op)
		utils.AssertEq(t, binExpr.Operator.Type(), test.result)
		utils.AssertEq(t, binExpr.Right.Type(), test.right)
	}
}

func TestUnaryExpression(t *testing.T) {
	tests := []struct {
		src      string
		operator ir.UnaryOperator
		operand  types.Type
		result   types.Type
	}{
		{"-1", ir.NegateInt, types.UntypedInt, types.UntypedInt},
		{"-2.72", ir.NegateFloat, types.UntypedFloat, types.UntypedFloat},
		{"!true", ir.LogicalNot, types.Bool, types.Bool},
		{"~104", ir.BitwiseNot, types.UntypedInt, types.UntypedInt},
		// TODO:
		// IncrecementInt (Needs a variable to increment)
		// IncrementFloat
		// DecrecementInt
		// DecrementFloat
		// PropagateError
		// CrashError
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		unExpr := getExpr[*ir.UnaryExpression](t, program)

		utils.AssertEq(t, unExpr.Operand.Type(), test.operand)
		op := unExpr.Operator & ^ir.UntypedBit
		utils.AssertEq(t, op, test.operator)
		utils.AssertEq(t, unExpr.Operator.Type(), test.result)
	}
}

func TestCastExpression(t *testing.T) {
	tests := []struct {
		src    string
		result types.Type
	}{
		{"1 -> i32", types.Int},
		{"1 -> f32", types.Float},
		{"1.6 -> f32", types.Float},
		{"true -> bool", types.Bool},
		{"false -> i32", types.Int},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		// note: if the conversion doesn't change the type (true -> bool),
		// the compiler removes the conversion completely, so we can't assume
		// the expression will be a *ir.Conversion
		expr := getExpr[ir.Expression](t, program)

		utils.AssertEq(t, expr.Type(), test.result)
	}
}

func TestCompileTimeValues(t *testing.T) {
	tests := []struct {
		src   string
		value any
	}{
		{"1", 1},
		{"17.5", 17.5},
		{"1.0", 1.0},
		{"1.0 -> i32", 1},
		{"5 -> f32", 5.0},
		{"false", false},
		{"true", true},
		{"-1", -1},
		{"!false", true},
		{"1 + 2 * 3", 7},
		{"1 + 2 / 4", 1.5},
		{`"test" + "123"`, "test123"},
		{"7 == 10", false},
		{"1.5 != 2.3", true},
		{"true || false", true},
		{"true && false", false},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		expr := getExpr[ir.Expression](t, program)

		utils.Assert(t, expr.IsConst(), "Expression was not compile-time known")
		utils.AssertEq(t, expr.ConstValue(), constValue(test.value))
	}
}

func TestArrays(t *testing.T) {
	tests := []struct {
		src    string
		values []any
	}{
		{"[1, 2, 3]", []any{1, 2, 3}},
		{"[true, false, true || true]", []any{true, false, true}},
		{"[1.5 + 2, 6 / 5, 1.2 ** 2]", []any{3.5, 1.2, 1.44}},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		expr := getExpr[*ir.ArrayExpression](t, program)

		utils.Assert(t, expr.IsConst(), "Expression was not compile-time known")
		constVal := expr.ConstValue().(values.ArrayValue)
		utils.AssertEq(t, len(constVal.Elements), len(test.values))
		for i, elem := range constVal.Elements {
			utils.AssertEq(t, elem, constValue(test.values[i]))
		}
	}
}

func constValue(val any) values.ConstValue {
	switch value := val.(type) {
	case int:
		return values.IntValue{
			Value: int64(value),
		}
	case float64:
		return values.FloatValue{
			Value: value,
		}
	case bool:
		return values.BoolValue{
			Value: value,
		}
	case string:
		return values.StringValue{
			Value: value,
		}
	default:
		panic("Unreachable")
	}
}

func TestIndexExpressions(t *testing.T) {
	tests := []struct {
		src      string
		dataType types.Type
	}{
		{"[1, 2, 3][1]", types.Int},
		{"[1.2, 3.4, 1][7]", types.Float},
		{"[7 == 2, 31 > 30.5][0.0]", types.Bool},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		expr := getExpr[*ir.IndexExpression](t, program)

		utils.AssertEq(t, expr.Type(), test.dataType)
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
		{`mut result = 1 [+] "hi"`, []diagnostic{{`Operator "+" is not defined for types "untyped int" and "string"`, diagnostics.Error}}},
		{"const neg_bool = [-]true", []diagnostic{{`Operator "-" is not defined for operand of type "bool"`, diagnostics.Error}}},
		{"let truthy: bool = [1] -> bool", []diagnostic{{`Cannot cast value of type "untyped int" to type "bool"`, diagnostics.Error}}},
		{"let i = 0; [i]++", []diagnostic{{`Variable "i" is immutable`, diagnostics.Error}}},
		{"1 + [2]--", []diagnostic{{`Cannot decrement a non-variable value`, diagnostics.Error}}},
		{"[[1, 2, [true] ]]", []diagnostic{{`Value of type "bool" is not assignable to type "i32"`, diagnostics.Error}}},
		{"mut a = 0; const b = [a + 1]", []diagnostic{{`Value must be known at compile time`, diagnostics.Error}}},
		{`let arr: string[[ [1.5] ]] = [["one", "half"]]`, []diagnostic{{`Array length must be an integer`, diagnostics.Error}}},
		{`[[1, 2, 3]][[ [3.14] ]]`, []diagnostic{{`Cannot index value of type "i32[3]" with value of type "untyped float"`, diagnostics.Error}}},
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
	escaped := false
	for i, c := range sourceText {
		if c == '[' && !escaped {
			if i+1 < len(sourceText) && sourceText[i+1] == '[' {
				escaped = true
			} else {
				spans = append(spans, text.NewSpan(line, line, col, 0))
			}
			continue
		}
		if c == ']' && !escaped {
			if i+1 < len(sourceText) && sourceText[i+1] == ']' {
				escaped = true
			} else {
				spans[len(spans)-1].End = col
				spans[len(spans)-1].EndLine = line
			}
			continue
		}

		escaped = false
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
		"expression is not %T (is %T)", struct{ t T }{}.t, exprStmt.Expression))
	return expr
}
