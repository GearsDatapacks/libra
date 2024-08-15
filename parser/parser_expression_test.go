package parser_test

import (
	"testing"

	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestIdentifierExpression(t *testing.T) {
	utils.MatchAstSnaps(t, "hello_World123")
}

func TestIntegerExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"156_098",
		"0b011",
		"0o23",
		"0xdead",
	)
}

func TestFloatExpression(t *testing.T) {
	utils.MatchAstSnaps(t, "3.1415")
}

func TestBooleanExpression(t *testing.T) {
	utils.MatchAstSnaps(t, "true")
}

func TestListExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"[]",
		"[1 ,2, 3,]",
		"[true, false, 5]",
		`[a,b,"c!"]`,
	)
}

func TestMapExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"{}",
		"{1: 2, 2: 3, 3:4}",
		`{"foo": "bar", "hello": "world"}`,
		`{hi: "there", "x": computed}`,
	)
}

func TestFunctionCall(t *testing.T) {
	utils.MatchAstSnaps(t,
		"add(1, 2)",
		`print("Hello, world!" ,)`,
	)
}

func TestIndexExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"(arr[7])",
		`({"a": 1}["b"])`,
	)
}

func TestMemberExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"foo.bar",
		"1.to_string",
		"a\n.b",
		".None",
	)
}

func TestStructExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"foo {bar: 1, baz: 2}",
		"rect {width: 9, height: 7.8}",
		`message {greeting: "Hello", name: name,}`,
		".{a:1, b:2}",
		// TODO: Make this parse the expression somehow
		// `struct {field: "value"}`,
	)
}

func TestCastExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"1->f32",
		"foo -> bar",
		`"_" -> u8`,
	)
}

func TestTypeCheckExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"1 is i32",
		`"Hello" is string`,
		"thing is bool",
	)
}

func TestRangeExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"1..10",
		"1.5..78.03",
	)
}

func TestBinaryExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"1 + 2",
		`"Hello" + "world"`,
		"foo - bar",
		"19 / 27",
		"1 << 2",
		"7 &19",
		"15.04* 1_2_3",
		"true||false",
		"[1,2,3]<< 4",
		"21 ^ 35",
	)
}

func TestAssignmentExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"a = b",
		"foo -= 1",
		`msg += "Hello"`,
	)
}

func TestPrefixExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"-2",
		"!true",
		"+foo",
		"~123",
	)
}

func TestPostfixExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"a?",
		"foo++",
		"5!",
	)
}

func TestDerefExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"ptr.*",
		"72.*",
	)
}

func TestRefExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"&13",
		"&mut my_var",
		"&false",
	)
}

func TestParenthesisedExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"(1 + 2)",
		"(true && false)",
	)
}

func TestTupleExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"()",
		"(1, 2, 3)",
		`(1, "hi", false, thing)`,
	)
}

func TestIfExpression(t *testing.T) {
	utils.MatchAstSnaps(t,
		"if a { 10 }",
		"if false { 10 } else { 20 }",
		`if 69
		{"Nice"}
		else if 42 { "UATLTUAE" }else{
			"Boring"
		}`,
	)
}

func TestWhileLoop(t *testing.T) {
	utils.MatchAstSnaps(t,
		"while true { nop }",
		`while thing { "Hi" }`,
	)
}

func TestForLoop(t *testing.T) {
	utils.MatchAstSnaps(t,
		"for i in [1,2,3] { i }",
		"for foo in 93\n{[foo,bar,]}",
	)
}

func TestFunctionExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"let func = fn() {}",
		"let func = fn(a, b: i32) { a + b }",
		`let func = fn(): string {"Hello, world!"}`,
	)
}

func TestTypeExpressions(t *testing.T) {
	utils.MatchAstSnaps(t,
		"void", "f32", "Type",
		"i64[10]", "bool[2]",
		"string[]", "Foo[]",
		"{string: i32}", "{bool: bool}",
		"(string, string, bool)", "(f64[], i8[])",
		"type Func = fn(SomeType, OtherType): bool", "type Func = fn(): string",
		"i8 | i16 | i32 | i64", "string | bool",
		"*string", "*mut i32",
		"?f32", "?(string[])",
		"!u8", "!{u32: string}",

		"?({string: Value}[10])",
		"!Foo[]",
		"(i32 | f32)[] | i32 | f32",
	)
}

func TestBlockMapInference(t *testing.T) {
	utils.MatchAstSnaps(t,
		"{}",
		"{a: b}",
		"{let value = 10}",
		"{1 + 2}",
	)
}

func TestOperatorPrecedence(t *testing.T) {
	utils.MatchAstSnaps(t,
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
		"a ^ b + 1",
	)
}

func TestParserDiagnostics(t *testing.T) {
	utils.MatchParserErrorSnaps(t,
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
	)
}
