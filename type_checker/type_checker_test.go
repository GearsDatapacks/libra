package typechecker_test

import (
	"testing"

	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestIntegerLiteral(t *testing.T) {
	utils.MatchIrSnaps(t, "+1_23_456")
}

func TestFloatLiteral(t *testing.T) {
	utils.MatchIrSnaps(t, "3.14_15_9")
}

func TestBooleanLiteral(t *testing.T) {
	utils.MatchIrSnaps(t, "true")
}

func TestStringLiteral(t *testing.T) {
	utils.MatchIrSnaps(t, `"Hi \"there\\"`)
}

func TestVariables(t *testing.T) {
	utils.MatchIrSnaps(t,
		"let x = 1; x",
		"mut foo: f32 = 1.4; foo",
		`const greeting: string = "Hi!"; greeting`,
		"mut is_awesome = true; is_awesome",
		"const my_float: f32 = 15; my_float",
	)
}

func TestIfStatements(t *testing.T) {
	utils.MatchIrSnaps(t,
		"if 1 + 2 == 3 {1 + 2}",
		"if true {1} else {0}",
		"if 1 == 2 {3} else if 2 != 3 {7} else {13}",
	)
}

func TestWhileLoops(t *testing.T) {
	utils.MatchIrSnaps(t,
		"while true { break 25 }",
		`
mut i = 0
mut sum = 0
let result = while i < 10 {
	i++
	sum += i
	if i == 10 {
		break i
	}
}
`,
	)
}

func TestForLoops(t *testing.T) {
	utils.MatchIrSnaps(t,
		`for i in [1,2,3] {
	if i % 2 == 0 {
		break i
	}
}`,

		`mut result = ""
for s in ["Hello", "world"] {
	result += s
}`,

		`for kv in {true: 1, false: 0} {
	if !kv[0] {
		break kv[1]
	}
}`,
	)
}

func TestBinaryExpression(t *testing.T) {
	utils.MatchIrSnaps(t,
		"true && false",
		"false || false",
		"1.5 < 2",
		"17 <= 17",
		"3.14 > 2.71",
		"42 >= 69",
		"1 == 2",
		"true == true",
		"1.2 != 7.5",
		"1 << 5",
		"8362 >> 3",
		"10101 | 1010",
		"73 & 52",
		"1 + 6",
		"2.3 + 4",
		`"Hello " + "world"`,
		"8 - 12",
		"3 - 1.3",
		"6 * 7",
		"1.3 * 0.4",
		"0.3 / 2",
		"103 % 2",
		"1.4 % 1",
		"2 ** 7",
		"3 ** 0.5",
	)
}

func TestUnaryExpression(t *testing.T) {
	utils.MatchIrSnaps(t,
		"-1",
		"-2.72",
		"!true",
		"~104",
		"mut a = 0; a++",
		"mut f = 1.5; f++",
		"mut value = 24; value--",
		"mut countdown = 12.3; countdown--",
		// TODO:
		// PropagateError
		// CrashError
	)
}

func TestCastExpression(t *testing.T) {
	utils.MatchIrSnaps(t,
		"1 -> i32",
		"1 -> f32",
		"1.6 -> f32",
		"true -> bool",
		"false -> i32",
	)
}

func TestCompileTimeValues(t *testing.T) {
	utils.MatchIrSnaps(t,
		"1",
		"17.5",
		"1.0",
		"1.0 -> i32",
		"5 -> f32",
		"false",
		"true",
		"-1",
		"!false",
		"1 + 2 * 3",
		"1 + 2 / 4",
		`"test" + "123"`,
		"7 == 10",
		"1.5 != 2.3",
		"true || false",
		"true && false",
	)
}

func TestArrays(t *testing.T) {
	utils.MatchIrSnaps(t,
		"[1, 2, 3]",
		"[true, false, true || true]",
		"[1.5 + 2, 6 / 5, 1.2 ** 2]",
	)
}

func TestMaps(t *testing.T) {
	utils.MatchIrSnaps(t,
		"{1: 2, 3: 4}",
		`{"one": 1, "two": 2, "three": 3}`,
		`{true: "true", false: "false"}`,
		`{"1" + "2": 1 + 2, "7" + "4": 7 + 4}`,
	)
}

func TestTuples(t *testing.T) {
	utils.MatchIrSnaps(t,
		"()",
		"(1, 2, 3)",
		"(1.5, true, -1)",
		`("Hi", 2, false)`,
	)
}

func TestIndexExpressions(t *testing.T) {
	utils.MatchIrSnaps(t,
		"[1, 2, 3][1]",
		"[1.2, 3.4, 1][2]",
		"[7 == 2, 31 > 30.5][0.0]",
	)
}

func TestAssignmentExpressions(t *testing.T) {
	utils.MatchIrSnaps(t,
		"mut a = 1; a = 2",
		"mut pi = 3.15; pi = 3.14",
		`mut greeting = "Hell"; greeting += "o"`,
		"mut count = 10; count -= 2",
	)
}

func TestTypeChecks(t *testing.T) {
	utils.MatchIrSnaps(t,
		"1 is i32",
		"true is bool[1]",
		"({1: 1.0, 3: 3.14}) is {i32: f32}",
		`(1, 3.1, "hi") is (i32, f32, string)`,
	)
}

func TestFunctions(t *testing.T) {
	utils.MatchIrSnaps(t,
		`fn add(a, b: i32): i32 {
	let c = a + b
	return c
}
add(1, 2)`,
		"fn one(): i32 { 1 }",
		"fn nop() {}; nop()",
	)
}

func TestStructs(t *testing.T) {
	utils.MatchIrSnaps(t,
		`struct Person {
	name: string,
	age: i32
}
const me = Person { name: "Self", age: 903 }
let baby = Person { name: "Unnamed" }
const my_age = me.age`,

		`struct Foo { bar, baz: f32 }
mut foo = Foo { bar: 10, baz: 13.1 }
let bar = foo.bar
foo.baz = bar`,
	)
}

func TestTupleStructs(t *testing.T) {
	utils.MatchIrSnaps(t,
		`struct Vector2 { f32, f32 }
const UP = Vector2 { 0, 1 }
const RIGHT = Vector2 { 1, 0 }
mut my_vec = UP
my_vec[0] = RIGHT[0]`,
	)
}

func TestBlockExpressions(t *testing.T) {
	utils.MatchIrSnaps(t,
		"{ 1 + 2 }",
		"{ yield 25 }",
		"{ let a = 10; let b = 20; yield a + b }",
	)
}

func TestIfExpressions(t *testing.T) {
	utils.MatchIrSnaps(t,
		"if true { 1 } else { 2 }",

		`mut value = 10
let other_value = if value > 10 {
	let temp = value
	value += 1
	yield temp
} else if value > 5 {
	value -= 1
	yield value
} else {
	value
}`,
	)
}

func TestPointers(t *testing.T) {
	utils.MatchIrSnaps(t,
		`let value1: i32 = 10
let value_ptr: *i32 = &value1
let value2: i32 = value_ptr.*`,

		`mut value = 1
let ptr: *mut i32 = &mut value
while value < 10 {
	ptr.* += 1
}`,

		`mut mutable = 1
let downcasted_ptr: *i32 = &mut mutable`,
	)
}

func TestFunctionExpressions(t *testing.T) {
	utils.MatchIrSnaps(t,
		`type Callback = fn(i32): i32
fn twice(callback: Callback, value: i32): i32 {
	callback(callback(value))
}
twice(fn(value: i32): i32 { value + 1 }, 10)`,
	)
}

// TODO: ImportStatement

func TestTypeExpressions(t *testing.T) {
	utils.MatchIrSnaps(t,
		"i32", "bool", "Type",
		"i32[]", "bool[][]",
		"string[10]", "f32[2][]",
		"{string: string[]}",
		"(string, i32[], i32)",
		"type Name = string",
		"explicit type CustomStr = string",
		"string | f32",
		"*i32[]", "(*i32)[]", "*mut {string: string}",
		"struct Unit",
		// "?string[]",
		// "!i32",
		"type StrToInt = fn(string): i32",
	)
}

func TestInterfaces(t *testing.T) {
	utils.MatchIrSnaps(t,
		`interface Printable { print() }
struct Message { string }
fn (Message) print() {
	// I would put print(this[0]) if print was implemented
}
let my_printable: Printable = Message { "Hello" }
my_printable.print()`,

		`interface Add { add(i32): i32 }
fn (i32) add(other: i32): i32 { this + other }
fn add(a: Add, b: i32): i32 { a.add(b) }
let result: i32 = add(1, 2)`,
	)
}

func TestUnions(t *testing.T) {
	utils.MatchIrSnaps(t,
		`union IntOrString { i32, string }
mut value: IntOrString = "32"
value = 32
let int_value: i32 = value.i32
let string_value: string = value.string`,

		`union Int { int: i32, other: i32 }
let int1 = 10 -> Int.int
let int2 = 92 -> Int.other`,

		`union Shape {
	Circle { cx, cy, r: i32 },
	Rectangle { x, y, w, h: i32 }
}
mut circle = Shape.Circle { cx: 10, cy: 31, r: 5 }
mut rectangle = Shape.Rectangle { x: 0, y: 0, w: 10, h: 5 }
circle = rectangle`,
	)
}

func TestTags(t *testing.T) {
	utils.MatchIrSnaps(t,
		`tag Tag
@tag Tag
struct Foo
@tag Tag
struct Bar
let foo: Tag = Foo
mut bar: Tag = Bar
bar = Foo`,

		`tag Number {i32, f32}
@tag Number
explicit type Int = i32
mut num: Number = 1.31
num = 10 -> Int`,
	)
}

func TestResults(t *testing.T) {
	utils.MatchIrSnaps(t,
		"let result: !i32 = 10",
		`@tag Error
struct MyError { string }
mut result: !i32 = MyError { "Error: uninitialised" }
result = 10
let must_be_int: i32 = result!`,
	)
}

// TODO: Add a way to create fake modules, for the following errors:
// FieldPrivate
// NoExport

func TestTCDiagnostics(t *testing.T) {
	utils.MatchTCErrorSnaps(t,
		"let x: foo = 1",
		"const text: string = false",
		"let result: !i32 = 10; let int: i32 = result",
		"let foo = 1; let foo = 2",
		"let a = b",
		`mut result = 1 + "hi"`,
		"const neg_bool = -true",
		"fn nop() { return 25 }",
		"let truthy: bool = 1 -> bool",
		"let i = 0; i = 1",
		"mut ptr = &10; ptr.* = 9",
		"1 + 2--",
		"[1, 2, true]",
		"mut a = 0; const b = a + 1",
		"mut i = 1; (1, true, 7.3)[i]",
		`let arr: string[1.5] = ["one", "half"]`,
		`[1, 2, 3][3.14]`,
		"{[1, 2]: 3}",
		"1 = 2",
		"[1, 2, 3][8]",
		"if 21 {12}",
		"for i in true {}",
		"return 23",
		"let func = fn(): bool { return\n }",
		`"print"("Hi")`,
		"fn add(a, b: i32): i32 {}; add(10)",
		`fn print(text: string) {}; print("Hello", "world!")`,
		"struct Empty {}; Empty{}.hello",
		"let value = 10.plus_one",
		"i32 { 1 }",
		"struct MyStruct {foo: string}; MyStruct {bar: 13}",
		"break 10", "continue",
		"while true { let my_func = fn() { break\n }; my_func() }",
		"yield 10", "{ for i in [1, 2, 3] { yield i } }",
		"const my_value: 10 = 10",
		"type Function = fn(i32, second: string)",
		"let func = fn(a: i32, i32[]) {}",
		"let deref = 10.*",
		"const value = 10; let ptr = &mut value",
		"struct Rect { w: i32, h }",
		"struct Wrapper {\nfoo: i32, value\n}",
		"struct Values { i32, i32 }; let values = Values { 1, 2, 3 }",
		"struct Values { i32, i32 }; let values = Values {}",
		"struct Number { i32, f32 }; Number {first: 10, second: 2.5}",
		"struct Vector {x, y: i32}; Vector {1, 2}",
		"struct CustomString {pub string}",
		"union Number { int: i32, float: f32 }; type Uint = Number.uint",
		"union IntArray { one: i32[1], two: i32[2] }; let i: IntArray = [1]; let three = i.three",
		`type NotATag = i32
@tag NotATag
struct Tagged`,
		`import "undefined"`,
		"let value_not_type = [1,2,3][]",
	)
}
