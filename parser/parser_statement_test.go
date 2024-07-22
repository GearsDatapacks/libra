package parser_test

import (
	"testing"
	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestVariableDeclaration(t *testing.T) {
	tests := []string{
		"let x = 1",
		"mut y: f32 = 7",
		`const message: string = "Hi"`,
		"mut isCool = true",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestFunctionDeclaration(t *testing.T) {
	tests := []string{
		`fn hello() { "Hello, world!" }`,
		"fn (i32) print() {\nthis\n}",
		"fn (i32) add(\nother: i32\n,\n)\n:i32\n{ 7 }",
		"fn u8.zero(): u8 {0}",
		"fn sum(a,b,c:f64) : usize{ 3.14 }",
		"fn inc(mut x: u32): u32 { x }",
		"fn (mut foo) bar(): foo { this }",
		"fn add(a = 1, mut b: i64 = 2): i64 { c }",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []string{
		"return",
		"return 7",
		"return false",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestBreakStatement(t *testing.T) {
	tests := []string{
		"break",
		"break true",
		"break [1,2,3]",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestYieldStatement(t *testing.T) {
	tests := []string{
		"yield 73",
		`yield "foo"`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestTypeDeclaration(t *testing.T) {
	tests := []string{
		"type foo = bar",
		"type int=i32",
		"type boolean\n =\n bool",
		"explicit type ID = u64",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestStructDeclaration(t *testing.T) {
	tests := []string{
		"struct Unit",
		"struct Mything123",
		"struct Wrapper { value }",
		"struct Three{a,b,c,}",
		"struct Empty {}",
		"struct Rect { w, h: i32 }",
		"struct Vec2{x:f32,y:f32,}",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestInterfaceDeclaration(t *testing.T) {
	tests := []string{
		"interface Any {}",
		"interface Fooer { foo(bar): baz }",
		`interface Order {
			less ( i32 , f64 ) : bool , 
			greater(u32,i32,):f16
		}`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestImportStatement(t *testing.T) {
	tests := []string{
		`import "fs"`,
		`import ".././foo/bar"`,
		`import * from "helpers"`,
		`import { read, write } from "io"`,
		`import "42" as life_universe_everything`,
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestEnumDeclaration(t *testing.T) {
	tests := []string{
		"enum Empty {}",
		"enum Colour: u64 { Invalid, red = 100, green = 783, blue = 1.5 }",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestUnionDeclaration(t *testing.T) {
	tests := []string{
		"union AOrB { a, b }", "AOrB",
		"union Int { i8, i16, i32, i64 ,}",
		"union Property { Age: i32, Height: f32, Weight:f32,string}",
		"union Shape { Square { f32, f32 }, Circle { radius: f32 } }",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestTagDeclaration(t *testing.T) {
	tests := []string{
		"tag MyTag",
		"tag Test124",
		"tag Number { i32, f32 }",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}

func TestAttributes(t *testing.T) {
	tests := []string{
		"@tag Error\nstruct MyError { string }",
		"@impl LeInterface\nfn (string) to_string(): string { this }",
		"@untagged\nunion IntOrPtr { int: i32, ptr: *i32 }",
		"@todo Implement it\nfn unimplemented(param: i32) {}",
		"@doc Does cool stuff\nfn do_cool_stuff() {}",
		"@deprecated Use `do_other_thing` instead\nfn do_thing() {}",
		"@doc Has three fields\n@todo Add a third field\n@tag Incomplete\nstruct ThreeFields {i32, f32}",
	}

	for _, test := range tests {
		utils.MatchAstSnap(t, test)
	}
}
