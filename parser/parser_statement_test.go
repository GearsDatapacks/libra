package parser_test

import (
	"testing"

	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestVariableDeclaration(t *testing.T) {
	utils.MatchAstSnaps(t,
		"let x = 1",
		"mut y: f32 = 7",
		`const message: string = "Hi"`,
		"mut isCool = true",
	)
}

func TestFunctionDeclaration(t *testing.T) {
	utils.MatchAstSnaps(t,
		`fn hello() { "Hello, world!" }`,
		"fn (i32) print() {\nthis\n}",
		"fn (i32) add(\nother: i32\n,\n)\n:i32\n{ 7 }",
		"fn u8.zero(): u8 {0}",
		"fn sum(a,b,c:f64) : usize{ 3.14 }",
		"fn inc(mut x: u32): u32 { x }",
		"fn (mut foo) bar(): foo { this }",
		"fn add(a = 1, mut b: i64 = 2): i64 { c }",
	)
}

func TestReturnStatement(t *testing.T) {
	utils.MatchAstSnaps(t,
		"return",
		"return 7",
		"return false",
	)
}

func TestBreakStatement(t *testing.T) {
	utils.MatchAstSnaps(t,
		"break",
		"break true",
		"break [1,2,3]",
	)
}

func TestYieldStatement(t *testing.T) {
	utils.MatchAstSnaps(t,
		"yield 73",
		`yield "foo"`,
	)
}

func TestTypeDeclaration(t *testing.T) {
	utils.MatchAstSnaps(t,
		"type foo = bar",
		"type int=i32",
		"type boolean\n =\n bool",
		"explicit type ID = u64",
	)
}

func TestStructDeclaration(t *testing.T) {
	utils.MatchAstSnaps(t,
		"struct Unit",
		"struct Mything123",
		"struct Wrapper { value }",
		"struct Three{a,b,c,}",
		"struct Empty {}",
		"struct Rect { w, h: i32 }",
		"struct Vec2{x:f32,y:f32,}",
	)
}

func TestInterfaceDeclaration(t *testing.T) {
	utils.MatchAstSnaps(t,
		"interface Any {}",
		"interface Fooer { foo(bar): baz }",
		`interface Order {
			less ( i32 , f64 ) : bool , 
			greater(u32,i32,):f16
		}`,
	)
}

func TestImportStatement(t *testing.T) {
	utils.MatchAstSnaps(t,
		`import "fs"`,
		`import ".././foo/bar"`,
		`import * from "helpers"`,
		`import { read, write } from "io"`,
		`import "42" as life_universe_everything`,
	)
}

func TestEnumDeclaration(t *testing.T) {
	utils.MatchAstSnaps(t,
		"enum Empty {}",
		"enum Colour: u64 { Invalid, red = 100, green = 783, blue = 1.5 }",
	)
}

func TestUnionDeclaration(t *testing.T) {
	utils.MatchAstSnaps(t,
		"union AOrB { a, b }", "AOrB",
		"union Int { i8, i16, i32, i64 ,}",
		"union Property { Age: i32, Height: f32, Weight:f32,string}",
		"union Shape { Square { f32, f32 }, Circle { radius: f32 } }",
	)
}

func TestTagDeclaration(t *testing.T) {
	utils.MatchAstSnaps(t,
		"tag MyTag",
		"tag Test124",
		"tag Number { i32, f32 }",
	)
}

func TestAttributes(t *testing.T) {
	utils.MatchAstSnaps(t,
		"@tag Error\nstruct MyError { string }",
		"@impl LeInterface\nfn (string) to_string(): string { this }",
		"@untagged\nunion IntOrPtr { int: i32, ptr: *i32 }",
		"@todo Implement it\nfn unimplemented(param: i32) {}",
		"@doc Does cool stuff\nfn do_cool_stuff() {}",
		"@deprecated Use `do_other_thing` instead\nfn do_thing() {}",
		"@doc Has three fields\n@todo Add a third field\n@tag Incomplete\nstruct ThreeFields {i32, f32}",
	)
}
