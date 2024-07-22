package parser_test

import (
	"testing"

	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestIdentifierExpression(t *testing.T) {
	utils.MatchAstSnap(t, "hello_World123")
}

func TestIntegerExpression(t *testing.T) {
	utils.MatchAstSnap(t, "156_098")
}

func TestFloatExpression(t *testing.T) {
	utils.MatchAstSnap(t, "3.1415")
}

func TestBooleanExpression(t *testing.T) {
	utils.MatchAstSnap(t, "true")
}

func TestListExpression(t *testing.T) {
	tests := []string{
		"[]",
		"[1 ,2, 3,]",
		"[true, false, 5]",
		`[a,b,"c!"]`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestMapExpression(t *testing.T) {
	tests := []string{
		"{}",
		"{1: 2, 2: 3, 3:4}",
		`{"foo": "bar", "hello": "world"}`,
		`{hi: "there", "x": computed}`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []string{
		"add(1, 2)",
		`print("Hello, world!" ,)`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestIndexExpression(t *testing.T) {
	tests := []string{
		"(arr[7])",
		`({"a": 1}["b"])`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestMemberExpression(t *testing.T) {
	tests := []string{
		"foo.bar",
		"1.to_string",
		"a\n.b",
		".None",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestStructExpression(t *testing.T) {
	tests := []string{
		"foo {bar: 1, baz: 2}",
		"rect {width: 9, height: 7.8}",
		`message {greeting: "Hello", name: name,}`,
		".{a:1, b:2}",
		// FIXME: Make this parse the expression somehow
		// `struct {field: "value"}`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestCastExpression(t *testing.T) {
	tests := []string{
		"1->f32",
		"foo -> bar",
		`"_" -> u8`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestTypeCheckExpression(t *testing.T) {
	tests := []string{
		"1 is i32",
		`"Hello" is string`,
		"thing is bool",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestRangeExpression(t *testing.T) {
	tests := []string{
		"1..10",
		"1.5..78.03",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestBinaryExpressions(t *testing.T) {
	tests := []string{
		"1 + 2",
		`"Hello" + "world"`,
		"foo - bar",
		"19 / 27",
		"1 << 2",
		"7 &19",
		"15.04* 1_2_3",
		"true||false",
		"[1,2,3]<< 4",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestAssignmentExpressions(t *testing.T) {
	tests := []string{
		"a = b",
		"foo -= 1",
		`msg += "Hello"`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []string{
		"-2",
		"!true",
		"+foo",
		"~123",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestPostfixExpressions(t *testing.T) {
	tests := []string{
		"a?",
		"foo++",
		"5!",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestDerefExpressions(t *testing.T) {
	tests := []string{
		"ptr.*",
		"72.*",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestRefExpressions(t *testing.T) {
	tests := []string{
		"&13",
		"&mut my_var",
		"&false",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestParenthesisedExpressions(t *testing.T) {
	tests := []string{
		"(1 + 2)",
		"(true && false)",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestTupleExpressions(t *testing.T) {
	tests := []string{
		"()",
		"(1, 2, 3)",
		`(1, "hi", false, thing)`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestIfExpression(t *testing.T) {
	tests := []string{
		"if a { 10 }",
		"if false { 10 } else { 20 }",
		`if 69
		{"Nice"}
		else if 42 { "UATLTUAE" }else{
			"Boring"
		}`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestWhileLoop(t *testing.T) {
	tests := []string{
		"while true { nop }",
		`while thing { "Hi" }`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestForLoop(t *testing.T) {
	tests := []string{
		"for i in [1,2,3] { i }",
		"for foo in 93\n{[foo,bar,]}",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestFunctionExpressions(t *testing.T) {
	tests := []string{
		"fn() {}",
		"fn(a, b: i32) { a + b }",
		`fn(): string {"Hello, world!"}`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, "let func = "+test)
	}
}

func TestTypeExpressions(t *testing.T) {
	tests := []string{
		"void", "f32", "Type",
		"i64[10]", "bool[2]",
		"string[]", "Foo[]",
		"{string: i32}", "{bool: bool}",
		"(string, string, bool)", "(f64[], i8[])",
		"fn(SomeType, OtherType): bool", "fn(): string",
		"i8 | i16 | i32 | i64", "string | bool",
		"*string", "*mut i32",
		"?f32", "?(string[])",
		"!u8", "!{u32: string}",

		"?({string: Value}[10])",
		"!Foo[]",
		"(i32 | f32)[] | i32 | f32",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, "type MyType = "+test)
	}
}

func TestBlockMapInference(t *testing.T) {
	tests := []string{
		"{}",
		"{a: b}",
		"{let value = 10}",
		"{1 + 2}",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []string{
		"1 + 2",
		"1 + 2 + 3",
		"1 + 2 * 3",
		"1 * 2 + 3",
		"foo + bar * baz ** qux",
		"a **b** c",
		"1 << 2 & 3",
		"true || false == true",
		"1 + (2 + 3)",
		"( 2**2 ) ** 2",
		"-1 + 2",
		"foo + -(bar * baz)",
		"1 - foo++",
		"hi + (a || b)!",
		"foo++-- + 1",
		"-a! / 4",
		"!foo() / 79",
		"-a[b] + 4",
		"fns[1]() * 3",
		"a = 1 + 2",
		"foo = bar = baz",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestParserDiagnostics(t *testing.T) {
	tests := []string{
		")",
		"let a = ;",
		"1 2",
		"(1 + 2",
		"else\n {}",
		"for i 42 {}",
		"let in = 1\nfor i in 20 {}",
		"fn add(a: i32, b, c): f32 {}",
		"fn foo(\nbar\n): baz {}",
		"fn func_type(mut i32[]) {}",
		"fn (string) bool.maybe() {}",
		`import * from "io" as in_out`,
		`import {read, write} from * from "io"`,
		`if true { fn a() {} }`,
		`type T = ;`,
		`let value = .`,
		"pub return 10",
		"explicit fn func() {}",
		"@nonexistent\nfn attributed() {}",
		"@tag FunctionTag\nfn tagged() {}",
	}

	for _, test := range tests {
		utils.MatchParserErrorSnap(t, test)
	}
}
