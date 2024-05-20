package parser_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/parser/ast"
	utils "github.com/gearsdatapacks/libra/test_utils"
	"github.com/gearsdatapacks/libra/text"
)

func TestIdentifierExpression(t *testing.T) {
	input := "hello_World123"
	program := getProgram(t, input)

	ident := getExpr[*ast.Identifier](t, program)

	utils.AssertEq(t, ident.Name, input)
	utils.AssertEq(t, ident.Token.Kind, token.IDENTIFIER)
	utils.AssertEq(t, ident.Token.Value, input)
	utils.AssertEq(t, ident.Token.Location.Span, text.NewSpan(0, 0, 0, len(input)))
}

func TestIntegerExpression(t *testing.T) {
	input := "156_098"
	val := 156098

	program := getProgram(t, input)

	integer := getExpr[*ast.IntegerLiteral](t, program)

	utils.AssertEq(t, integer.Value, int64(val))

	utils.AssertEq(t, integer.Token.Kind, token.INTEGER)
	utils.AssertEq(t, integer.Token.Value, "156098")
	utils.AssertEq(t, integer.Token.Location.Span, text.NewSpan(0, 0, 0, len(input)))
}

func TestFloatExpression(t *testing.T) {
	input := "3.1415"
	val := 3.1415

	program := getProgram(t, input)

	float := getExpr[*ast.FloatLiteral](t, program)

	utils.AssertEq(t, float.Value, val)

	utils.AssertEq(t, float.Token.Kind, token.FLOAT)
	utils.AssertEq(t, float.Token.Value, input)
	utils.AssertEq(t, float.Token.Location.Span, text.NewSpan(0, 0, 0, len(input)))
}

func TestBooleanExpression(t *testing.T) {
	input := "true"

	program := getProgram(t, input)

	boolean := getExpr[*ast.BooleanLiteral](t, program)

	utils.AssertEq(t, boolean.Value, true)
	utils.AssertEq(t, boolean.Token.Kind, token.IDENTIFIER)
	utils.AssertEq(t, boolean.Token.Value, input)
	utils.AssertEq(t, boolean.Token.Location.Span, text.NewSpan(0, 0, 0, len(input)))
}

func TestListExpression(t *testing.T) {
	tests := []struct {
		src    string
		values []any
	}{
		{"[]", []any{}},
		{"[1 ,2, 3,]", []any{1, 2, 3}},
		{"[true, false, 5]", []any{true, false, 5}},
		{`[a,b,"c!"]`, []any{"$a", "$b", "c!"}},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		list := getExpr[*ast.ListLiteral](t, program)
		testLiteral(t, list, tt.values)
	}
}

func TestMapExpression(t *testing.T) {
	tests := []struct {
		src       string
		keyValues [][2]any
	}{
		{"({})", [][2]any{}},
		{"({1: 2, 2: 3, 3:4})", [][2]any{{1, 2}, {2, 3}, {3, 4}}},
		{`({"foo": "bar", "hello": "world"})`, [][2]any{{"foo", "bar"}, {"hello", "world"}}},
		{`({hi: "there", "x": computed})`, [][2]any{{"$hi", "there"}, {"x", "$computed"}}},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		parenExpr := getExpr[*ast.ParenthesisedExpression](t, program)
		mapLit := parenExpr.Expression.(*ast.MapLiteral)
		testLiteral(t, mapLit, tt.keyValues)
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		src  string
		name string
		args []any
	}{
		{"add(1, 2)", "add", []any{1, 2}},
		{`print("Hello, world!" ,)`, "print", []any{"Hello, world!"}},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		call := getExpr[*ast.FunctionCall](t, program)
		ident, ok := call.Callee.(*ast.Identifier)
		utils.Assert(t, ok, "Callee is not an identifier")
		utils.AssertEq(t, tt.name, ident.Name)
		for i, arg := range call.Arguments {
			testLiteral(t, arg, tt.args[i])
		}
	}
}

func TestIndexExpression(t *testing.T) {
	tests := []struct {
		src   string
		left  any
		index any
	}{
		{"(arr[7])", "$arr", 7},
		{`({"a": 1}["b"])`, [][2]any{{"a", 1}}, "b"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		parenExpr := getExpr[*ast.ParenthesisedExpression](t, program)
		index := parenExpr.Expression.(*ast.IndexExpression)
		testLiteral(t, index.Left, tt.left)
		testLiteral(t, index.Index, tt.index)
	}
}

func TestMemberExpression(t *testing.T) {
	tests := []struct {
		src    string
		left   any
		member string
	}{
		{"foo.bar", "$foo", "bar"},
		{"1.to_string", 1, "to_string"},
		{"a\n.b", "$a", "b"},
		{".None", nil, "None"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		member := getExpr[*ast.MemberExpression](t, program)
		testLiteral(t, member.Left, tt.left)
		utils.AssertEq(t, tt.member, member.Member.Value)
	}
}

type structField struct {
	name  string
	value any
}

func TestStructExpression(t *testing.T) {
	tests := []struct {
		src     string
		name    string
		members []structField
	}{
		{"foo {bar: 1, baz: 2}", "foo", []structField{{"bar", 1}, {"baz", 2}}},
		{"rect {width: 9, height: 7.8}", "rect", []structField{{"width", 9}, {"height", 7.8}}},
		{`message {greeting: "Hello", name: name,}`, "message", []structField{{"greeting", "Hello"}, {"name", "$name"}}},
		{".{a:1, b:2}", ".", []structField{{"a", 1}, {"b", 2}}},
		// FIXME: Make this parse the expression somehow
		// {`struct {field: "value"}`, "struct", []structField{{"field", "value"}}},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		structExpr := getExpr[*ast.StructExpression](t, program)
		if tt.name == "." {
			_, ok := structExpr.Struct.(*ast.InferredExpression)
			utils.Assert(t, ok, "Struct's type should be inferred")
		} else {
			ident, ok := structExpr.Struct.(*ast.Identifier)
			utils.Assert(t, ok, "Struct is not an identifier")
			utils.AssertEq(t, tt.name, ident.Name)
		}

		for i, member := range structExpr.Members {
			tMember := tt.members[i]
			utils.AssertEq(t, tMember.name, member.Name.Value)
			testLiteral(t, member.Value, tMember.value)
		}
	}
}

func TestCastExpression(t *testing.T) {
	tests := []struct {
		src  string
		left any
		to   string
	}{
		{"1->f32", 1, "f32"},
		{"foo -> bar", "$foo", "bar"},
		{`"_" -> u8`, "_", "u8"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		cast := getExpr[*ast.CastExpression](t, program)

		testLiteral(t, cast.Left, tt.left)
		ident, ok := cast.Type.(*ast.TypeName)
		utils.Assert(t, ok, "Didn't cast to a type name")
		utils.AssertEq(t, tt.to, ident.Name.Value)
	}
}

func TestTypeCheckExpression(t *testing.T) {
	tests := []struct {
		src      string
		left     any
		typeName string
	}{
		{"1 is i32", 1, "i32"},
		{`"Hello" is string`, "Hello", "string"},
		{"thing is bool", "$thing", "bool"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		typeCheck := getExpr[*ast.TypeCheckExpression](t, program)

		testLiteral(t, typeCheck.Left, tt.left)
		ident, ok := typeCheck.Type.(*ast.TypeName)
		utils.Assert(t, ok, "Didn't check for a type name")
		utils.AssertEq(t, tt.typeName, ident.Name.Value)
	}
}

func TestRangeExpression(t *testing.T) {
	tests := []struct {
		src   string
		start any
		end   any
	}{
		{"1..10", 1, 10},
		{"1.5..78.03", 1.5, 78.03},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		rangeExpr := getExpr[*ast.RangeExpression](t, program)
		testLiteral(t, rangeExpr.Start, tt.start)
		testLiteral(t, rangeExpr.End, tt.end)
	}
}

func TestBinaryExpressions(t *testing.T) {
	tests := []struct {
		src   string
		left  any
		op    string
		right any
	}{
		{"1 + 2", 1, "+", 2},
		{`"Hello" + "world"`, "Hello", "+", "world"},
		{"foo - bar", "$foo", "-", "$bar"},
		{"19 / 27", 19, "/", 27},
		{"1 << 2", 1, "<<", 2},
		{"7 &19", 7, "&", 19},
		{"15.04* 1_2_3", 15.04, "*", 123},
		{"true||false", true, "||", false},
		{"[1,2,3]<< 4", []any{1, 2, 3}, "<<", 4},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		expr := getExpr[*ast.BinaryExpression](t, program)

		testLiteral(t, expr.Left, tt.left)
		utils.AssertEq(t, expr.Operator.Value, tt.op)
		testLiteral(t, expr.Right, tt.right)
	}
}

func TestAssignmentExpressions(t *testing.T) {
	tests := []struct {
		src      string
		assignee any
		op       string
		value    any
	}{
		{"a = b", "$a", "=", "$b"},
		{"foo -= 1", "$foo", "-=", 1},
		{`msg += "Hello"`, "$msg", "+=", "Hello"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		expr := getExpr[*ast.AssignmentExpression](t, program)

		testLiteral(t, expr.Assignee, tt.assignee)
		utils.AssertEq(t, expr.Operator.Value, tt.op)
		testLiteral(t, expr.Value, tt.value)
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		src      string
		operator string
		operand  any
	}{
		{"-2", "-", 2},
		{"!true", "!", true},
		{"+foo", "+", "$foo"},
		{"~123", "~", 123},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		expr := getExpr[*ast.PrefixExpression](t, program)

		utils.AssertEq(t, expr.Operator.Value, tt.operator)
		testLiteral(t, expr.Operand, tt.operand)
	}
}

func TestPostfixExpressions(t *testing.T) {
	tests := []struct {
		src      string
		operand  any
		operator string
	}{
		{"a?", "$a", "?"},
		{"foo++", "$foo", "++"},
		{"5!", 5, "!"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		expr := getExpr[*ast.PostfixExpression](t, program)

		testLiteral(t, expr.Operand, tt.operand)
		utils.AssertEq(t, expr.Operator.Value, tt.operator)
	}
}

func TestParenthesisedExpressions(t *testing.T) {
	tests := []struct {
		src   string
		left  any
		op    string
		right any
	}{
		{"(1 + 2)", 1, "+", 2},
		{"(true && false)", true, "&&", false},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		expr := getExpr[*ast.ParenthesisedExpression](t, program)
		binExpr, ok := expr.Expression.(*ast.BinaryExpression)
		utils.Assert(t, ok, fmt.Sprintf(
			"Expression was not binary expression (was %T)", expr.Expression))

		testLiteral(t, binExpr.Left, tt.left)
		utils.AssertEq(t, binExpr.Operator.Value, tt.op)
		testLiteral(t, binExpr.Right, tt.right)
	}
}

func TestTupleExpressions(t *testing.T) {
	tests := []struct {
		src    string
		values []any
	}{
		{"()", []any{}},
		{"(1, 2, 3)", []any{1, 2, 3}},
		{`(1, "hi", false, thing)`, []any{1, "hi", false, "$thing"}},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		expr := getExpr[*ast.TupleExpression](t, program)
		utils.AssertEq(t, len(tt.values), len(expr.Values))

		for i, value := range expr.Values {
			testLiteral(t, value, tt.values[i])
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		src string
		res string
	}{
		{"1 + 2", "(1 + 2)"},
		{"1 + 2 + 3", "((1 + 2) + 3)"},
		{"1 + 2 * 3", "(1 + (2 * 3))"},
		{"1 * 2 + 3", "((1 * 2) + 3)"},
		{"foo + bar * baz ** qux", "(foo + (bar * (baz ** qux)))"},
		{"a **b** c", "(a ** (b ** c))"},
		{"1 << 2 & 3", "((1 << 2) & 3)"},
		{"true || false == true", "(true || (false == true))"},
		{"1 + (2 + 3)", "(1 + (2 + 3))"},
		{"( 2**2 ) ** 2", "((2 ** 2) ** 2)"},
		{"-1 + 2", "(-(1) + 2)"},
		{"foo + -(bar * baz)", "(foo + -((bar * baz)))"},
		{"1 - foo++", "(1 - (foo)++)"},
		{"hi + (a || b)!", "(hi + ((a || b))!)"},
		{"foo++-- + 1", "(((foo)++)-- + 1)"},
		{"-a! / 4", "(-((a)!) / 4)"},
		{"!foo() / 79", "(!(foo()) / 79)"},
		{"-a[b] + 4", "(-(a[b]) + 4)"},
		{"fns[1]() * 3", "(fns[1]() * 3)"},
		{"a = 1 + 2", "(a = (1 + 2))"},
		{"foo = bar = baz", "(foo = (bar = baz))"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		expr := getExpr[ast.HasPrecedence](t, program)

		utils.AssertEq(t, expr.PrecedenceString(), tt.res)
	}
}

type diagnostic struct {
	message string
	kind    diagnostics.DiagnosticKind
}

func TestParserDiagnostics(t *testing.T) {
	tests := []struct {
		src         string
		diagnostics []diagnostic
	}{
		{"let a = [;]", []diagnostic{{"Expected expression, found `;`", diagnostics.Error}}},
		{"1 [2]", []diagnostic{{"Expected newline after statement, found integer", diagnostics.Error}}},
		{"(1 + 2[]", []diagnostic{{"Expected `)`, found <Eof>", diagnostics.Error}}},
		{"[else]\n {}", []diagnostic{{"Else statement not allowed without preceding if", diagnostics.Error}}},
		{"for i [42] {}", []diagnostic{{`Expected "in" keyword, found integer`, diagnostics.Error}}},
		{"let [in] = 1\nfor i [in] 20 {}", []diagnostic{
			{`Expected "in" keyword, but it has been overwritten by a variable`, diagnostics.Error},
			{"Try removing or renaming this variable", diagnostics.Info},
		}},
		{"fn add(a: i32, b, [c]): f32 {}", []diagnostic{{"The last parameter of a function must have a type annotation", diagnostics.Error}}},
		{"fn [foo](\n[bar]\n): baz {}", []diagnostic{
			{"The last parameter of a function must have a type annotation", diagnostics.Error},
			{"Parameter of this function", diagnostics.Info},
		}},
		{"struct Rect { w: i32, [h] }", []diagnostic{{"The last field of a struct must have a type annotation", diagnostics.Error}}},
		{"struct [Wrapper] {\n[value]\n}", []diagnostic{
			{"The last field of a struct must have a type annotation", diagnostics.Error},
			{"Field in this struct", diagnostics.Info},
		}},
		{"fn (string) [bool].maybe() {}", []diagnostic{{"Functions cannot be both methods and static members", diagnostics.Error}}},
		{`import * from "io" [as] in_out`, []diagnostic{{"Only one import modifier is allowed", diagnostics.Error}}},
		{`import {read, write} from [*] from "io"`, []diagnostic{{"Only one import modifier is allowed", diagnostics.Error}}},
		{`if true { [fn] a() {} }`, []diagnostic{{"Function declaration not allowed here", diagnostics.Error}}},
		{`type T = [;]`, []diagnostic{{"Expected type, found `;`", diagnostics.Error}}},
	}

	for _, test := range tests {
		src, spans := getSpans(test.src)
		utils.AssertEq(t, len(spans), len(test.diagnostics), "Mismatch of spans to diagnostic messages")

		lexer := lexer.New(text.NewFile("test.lb", src))
		tokens := lexer.Tokenise()

		utils.AssertEq(t, len(lexer.Diagnostics), 0, "Expected no lexer diagnostics")

		p := parser.New(tokens, lexer.Diagnostics)
		p.Parse()
		diagnostics := p.Diagnostics
		utils.AssertEq(t, len(diagnostics), len(test.diagnostics),
			fmt.Sprintf("Incorrect number of diagnostics (expected %d, got %d)", len(test.diagnostics), len(diagnostics)))

		for i, diag := range test.diagnostics {
			// FIXME: Do this in a better way
			span := spans[len(spans)-i-1]
			testDiagnostic(t, diagnostics[i], diag.kind, diag.message, span)
		}
	}
}

func TestErrorExpression(t *testing.T) {
	input := ")"

	l := lexer.New(text.NewFile("test.lb", input))
	tokens := l.Tokenise()
	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()

	diag := utils.AssertSingle(t, p.Diagnostics)
	msg := "Expected expression, found `)`"
	testDiagnostic(t, diag, diagnostics.Error, msg, text.NewSpan(0, 0, 0, 1))

	getExpr[*ast.ErrorNode](t, program)
}

func getSpans(sourceText string) (string, []text.Span) {
	var resultText bytes.Buffer
	spans := []text.Span{}
	line := 0
	col := 0
	for _, c := range sourceText {
		if c == '[' {
			spans = append(spans, text.NewSpan(line, line, col, 0))
			continue
		}
		if c == ']' {
			spans[len(spans)-1].End = col
			spans[len(spans)-1].EndLine = line
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

func getProgram(t *testing.T, input string) *ast.Program {
	t.Helper()

	l := lexer.New(text.NewFile("test.lb", input))
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()
	utils.AssertEq(t, len(p.Diagnostics), 0,
		fmt.Sprintf("Expected no diagnostics (got %d)", len(p.Diagnostics)))

	return program
}

func getExpr[T ast.Expression](t *testing.T, program *ast.Program) T {
	t.Helper()

	utils.AssertEq(t, len(program.Statements), 1,
		fmt.Sprintf("Program does not contain one statement. (has %d)",
			len(program.Statements)))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	utils.Assert(t, ok, fmt.Sprintf(
		"Statement is not an expression statement (is %T)", program.Statements[0]))

	expr, ok := stmt.Expression.(T)
	utils.Assert(t, ok, fmt.Sprintf("Expression is not %T (is %T)",
		struct{ t T }{}.t, stmt.Expression))
	return expr
}

func testLiteral(t *testing.T, expr ast.Expression, expected any) {
	t.Helper()

	switch val := expected.(type) {
	case nil:
		utils.Assert(t, expr == nil, "Expected value to be nil")
	case int:
		integer, ok := expr.(*ast.IntegerLiteral)
		utils.Assert(t, ok, fmt.Sprintf("Value was not an integer (was %T)", expr))
		utils.AssertEq(t, integer.Value, int64(val))

	case float64:
		float, ok := expr.(*ast.FloatLiteral)
		utils.Assert(t, ok, fmt.Sprintf("Value was not an float (was %T)", expr))
		utils.AssertEq(t, float.Value, val)

	case bool:
		boolean, ok := expr.(*ast.BooleanLiteral)
		utils.Assert(t, ok, fmt.Sprintf("Value was not a bool (was %T)", expr))
		utils.AssertEq(t, boolean.Value, val)

	case string:
		if val[0] == '$' {
			ident, ok := expr.(*ast.Identifier)
			utils.Assert(t, ok, fmt.Sprintf("Value was not an identifier (was %T)", expr))
			utils.AssertEq(t, ident.Name, val[1:])
		} else {
			str, ok := expr.(*ast.StringLiteral)
			utils.Assert(t, ok, fmt.Sprintf("Value was not a string (was %T)", expr))
			utils.AssertEq(t, str.Value, val)
		}

	case []any:
		list, ok := expr.(*ast.ListLiteral)
		utils.Assert(t, ok, fmt.Sprintf("Value was not a list (was %T)", expr))
		utils.AssertEq(t, len(list.Values), len(val),
			fmt.Sprintf("Lists' lengths do not match. Expected %d elements, got %d",
				len(val), len(list.Values)))

		for i, value := range list.Values {
			testLiteral(t, value, val[i])
		}

	case [][2]any:
		mapLit, ok := expr.(*ast.MapLiteral)
		utils.Assert(t, ok, fmt.Sprintf("Value was not a map (was %T)", expr))
		utils.AssertEq(t, len(mapLit.KeyValues), len(val),
			fmt.Sprintf("Maps' lengths do not match. Expected %d key-value pairs, got %d",
				len(val), len(mapLit.KeyValues)))

		for i, kv := range mapLit.KeyValues {
			testLiteral(t, kv.Key, val[i][0])
			testLiteral(t, kv.Value, val[i][1])
		}

	default:
		panic("Invalid literal kind")
	}
}
