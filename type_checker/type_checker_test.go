package typechecker_test

import (
	"testing"

	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestIntegerLiteral(t *testing.T) {
	input := "+1_23_456"

	utils.MatchIrSnap(t, input)
}

func TestFloatLiteral(t *testing.T) {
	input := "3.14_15_9"

	utils.MatchIrSnap(t, input)
}

func TestBooleanLiteral(t *testing.T) {
	input := "true"

	utils.MatchIrSnap(t, input)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hi \"there\\"`

	utils.MatchIrSnap(t, input)
}

func TestVariables(t *testing.T) {
	tests := []string{
		"let x = 1",
		"mut foo: f32 = 1.4",
		`const greeting: string = "Hi!"`,
		"mut is_awesome = true",
		"const my_float: f32 = 15",
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestIfStatements(t *testing.T) {
	tests := []string{
		"if 1 + 2 == 3 {1 + 2}",
		"if true {1} else {0}",
		"if 1 == 2 {3} else if 2 != 3 {7} else {13}",
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestWhileLoops(t *testing.T) {
	tests := []string{
		"while true {false}",
		"while 9 > 12 {9 * 12}",
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestForLoops(t *testing.T) {
	tests := []string{
		"for i in [1,2,3] {i}",
		`for s in ["Hello", "world"] {s}`,
		"for kv in {true: 1, false: 0} {kv}",
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestBinaryExpression(t *testing.T) {
	tests := []string{
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
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestUnaryExpression(t *testing.T) {
	tests := []string{
		"-1",
		"-2.72",
		"!true",
		"~104",
		// TODO:
		// IncrecementInt (Needs a variable to increment)
		// IncrementFloat
		// DecrecementInt
		// DecrementFloat
		// PropagateError
		// CrashError
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestCastExpression(t *testing.T) {
	tests := []string{
		"1 -> i32",
		"1 -> f32",
		"1.6 -> f32",
		"true -> bool",
		"false -> i32",
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestCompileTimeValues(t *testing.T) {
	tests := []string{
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
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestArrays(t *testing.T) {
	tests := []string{
		"[1, 2, 3]",
		"[true, false, true || true]",
		"[1.5 + 2, 6 / 5, 1.2 ** 2]",
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestMaps(t *testing.T) {
	tests := []string{
		"{1: 2, 3: 4}",
		`{"one": 1, "two": 2, "three": 3}`,
		`{true: "true", false: "false"}`,
		`{"1" + "2": 1 + 2, "7" + "4": 7 + 4}`,
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestTuples(t *testing.T) {
	tests := []string{
		"()",
		"(1, 2, 3)",
		"(1.5, true, -1)",
		`("Hi", 2, false)`,
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestIndexExpressions(t *testing.T) {
	tests := []string{
		"[1, 2, 3][1]",
		"[1.2, 3.4, 1][2]",
		"[7 == 2, 31 > 30.5][0.0]",
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestAssignmentExpressions(t *testing.T) {
	tests := []string{
		"mut a = 1; a = 2",
		"mut pi = 3.15; pi = 3.14",
		`mut greeting = "Hell"; greeting += "o"`,
		"mut count = 10; count -= 2",
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestTypeChecks(t *testing.T) {
	tests := []string{
		"1 is i32",
		"true is bool[1]",
		"({1: 1.0, 3: 3.14}) is {i32: f32}",
		`(1, 3.1, "hi") is (i32, f32, string)`,
	}

	for _, test := range tests {
		utils.MatchIrSnap(t, test)
	}
}

func TestTCDiagnostics(t *testing.T) {
	tests := []string{
		"let x: foo = 1",
		"const text: string = false",
		"let foo = 1; let foo = 2",
		"let a = b",
		`mut result = 1 + "hi"`,
		"const neg_bool = -true",
		"let truthy: bool = 1 -> bool",
		"let i = 0; i = 1",
		"1 + 2--",
		"[1, 2, true]",
		"mut a = 0; const b = a + 1",
		"mut i = 1; (1, true, 7.3)[i]",
		`let arr: string[1.5] = ["one", "half"]`,
		`[1, 2, 3][3.14]`,
		"1 = 2",
		"[1, 2, 3][8]",
		"if 21 {12}",
		"struct Rect { w: i32, h }",
		"struct Wrapper {\nfoo: i32, value\n}",
	}

	for _, test := range tests {
		utils.MatchTCErrorSnap(t, test)
	}
}
