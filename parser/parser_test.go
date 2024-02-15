package parser_test

import (
	"fmt"
	"testing"

	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/parser/ast"
	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestIdentifierExpression(t *testing.T) {
	input := "hello_World123"
	program := getProgram(t, input)

	ident := getExpr[*ast.Identifier](t, program)

	utils.AssertEq(t, ident.Name, input)
	utils.AssertEq(t, ident.Token,
		token.New(token.IDENTIFIER, input,
			token.NewSpan(0, 0, len(input)),
		),
	)
}

func TestIntegerExpression(t *testing.T) {
	input := "156_098"
	val := 156098

	program := getProgram(t, input)

	integer := getExpr[*ast.IntegerLiteral](t, program)

	utils.AssertEq(t, integer.Value, int64(val))
	utils.AssertEq(t, integer.Token,
		token.New(token.INTEGER, "156098",
			token.NewSpan(0, 0, len(input)),
		),
	)
}

func TestFloatExpression(t *testing.T) {
	input := "3.1415"
	val := 3.1415

	program := getProgram(t, input)

	float := getExpr[*ast.FloatLiteral](t, program)

	utils.AssertEq(t, float.Value, val)
	utils.AssertEq(t, float.Token,
		token.New(token.FLOAT, input,
			token.NewSpan(0, 0, len(input)),
		),
	)
}

func TestBooleanExpression(t *testing.T) {
	input := "true"

	program := getProgram(t, input)

	boolean := getExpr[*ast.BooleanLiteral](t, program)

	utils.AssertEq(t, boolean.Value, true)
	utils.AssertEq(t, boolean.Token,
		token.New(token.IDENTIFIER, input,
			token.NewSpan(0, 0, len(input)),
		),
	)
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
		{"{}", [][2]any{}},
		{"{1: 2, 2: 3, 3:4}", [][2]any{{1, 2}, {2, 3}, {3, 4}}},
		{`{"foo": "bar", "hello": "world"}`, [][2]any{{"foo", "bar"}, {"hello", "world"}}},
		{`{hi: "there", "x": computed}`, [][2]any{{"$hi", "there"}, {"x", "$computed"}}},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		mapLit := getExpr[*ast.MapLiteral](t, program)
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
		{"arr[7]", "$arr", 7},
		{`{"a": 1}["b"]`, [][2]any{{"a", 1}}, "b"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		call := getExpr[*ast.IndexExpression](t, program)
		testLiteral(t, call.Left, tt.left)
		testLiteral(t, call.Index, tt.index)
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
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		member := getExpr[*ast.MemberExpression](t, program)
		testLiteral(t, member.Left, tt.left)
		utils.AssertEq(t, tt.member, member.Member.Value)
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

func TestErrorExpression(t *testing.T) {
	input := ")"

	l := lexer.New(input, "test.lb")
	tokens := l.Tokenise()
	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()

	utils.AssertEq(t, len(p.Diagnostics.Diagnostics), 1)
	diag := p.Diagnostics.Diagnostics[0]
	utils.AssertEq(t, diag.Message, "Expected expression, got `)`")
	utils.AssertEq(t, diag.Span, token.NewSpan(0, 0, 1))

	getExpr[*ast.ErrorExpression](t, program)
}

func TestMissingNewlineError(t *testing.T) {
	input := "1 2"

	l := lexer.New(input, "test.lb")
	tokens := l.Tokenise()
	p := parser.New(tokens, l.Diagnostics)
	p.Parse()

	utils.AssertEq(t, len(p.Diagnostics.Diagnostics), 1)
	diag := p.Diagnostics.Diagnostics[0]
	utils.AssertEq(t, diag.Message, "Expected newline after statement, got integer")
	utils.AssertEq(t, diag.Span, token.NewSpan(0, 2, 3))
}

func TestIncorrectTokenError(t *testing.T) {
	input := "(1 + 2"

	l := lexer.New(input, "test.lb")
	tokens := l.Tokenise()
	p := parser.New(tokens, l.Diagnostics)
	p.Parse()

	utils.AssertEq(t, len(p.Diagnostics.Diagnostics), 1)
	diag := p.Diagnostics.Diagnostics[0]
	utils.AssertEq(t, diag.Message, "Expected `)`, found <Eof>")
	utils.AssertEq(t, diag.Span, token.NewSpan(0, 6, 6))
}

func getProgram(t *testing.T, input string) *ast.Program {
	t.Helper()

	l := lexer.New(input, "test.lb")
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()
	utils.AssertEq(t, len(p.Diagnostics.Diagnostics), 0,
		fmt.Sprintf("Expected no diagnostics (got %d)", len(p.Diagnostics.Diagnostics)))

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
	}
}
