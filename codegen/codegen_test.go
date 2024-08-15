package codegen_test

import (
	"testing"

	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestConstantCodegen(t *testing.T) {
	utils.MatchCodegenSnaps(t,
		"let value = 10",
		"let value: u16 = 301",
		"let weight = 61.2",
		"let half_precision: f32 = 1.3",
		"let b = true",
	)
}

func TestBinaryExpressions(t *testing.T) {
	utils.MatchCodegenSnaps(t,
		"mut a = 1; mut b = 2; let res = a + b",
		"mut a: f32 = 4; mut b: f32 = 9.2; let res = a + b",
		"mut a: u8 = 4; mut b: u8 = 23; let res = a + b",
		"mut a: i16 = 31; mut b: i16 = 4; let res = a - b",
		"mut a: f32 = 1.9; mut b: f32 = 4.8; let res = a - b",
		"mut a: u64 = 10283; mut b: u64 = 732; let res = a * b",
		"mut a = 3.2; mut b = 2.1; let res = a * b",
		"mut a = 203; mut b = 41; let res = a & b",
		"mut a: u64 = 32427; mut b: u64 = 23824523; let res = a & b",
		"mut a = 49; mut b = 4163; let res = a ^ b",
		"mut a: i8 = 91; mut b: i8 = 84; let res = a ^ b",
		"mut a = 924; mut b = 91; let res = a | b",
		"mut a: u32 = 2746; mut b: u32 = 1024; let res = a | b",
		"mut a = true; mut b = false; let or = a || b",
		"mut a = false; mut b = true; let and = a && b",
		"mut a = 20; mut b = 30; let eq = a == b",
		"mut a: i8 = 12; mut b: i8 = 31; let neq = a != b",
		"mut a = 82; mut b = 103; let gt = a > b",
		"mut a = 91; mut b = 91; let ge = a >= b",
		"mut a = 12; mut b = 47; let lt = a < b",
		"mut a = 12; mut b = 47; let le = a <= b",
		"mut a = 1; mut b = 30; let shift = a << b",
		"mut a: u8 = 31; mut b: u8 = 3; let shift = a << b",
		"mut a = 72041; mut b = 3; let shift = a >> b",
		"mut a: i64 = 476354293423; mut b: i64 = 40; let shift = a >> b",
		"mut a = 72041; mut b = 3; let shift = a >>> b",
		"mut a: u16 = 60203; mut b: u16 = 5; let shift = a >>> b",
	)
}

func TestAssignment(t *testing.T) {
	utils.MatchCodegenSnaps(t,
		"mut age = 1; age += 1",
		"mut value = 100; value = 200",
		"mut x: f32 = 1; x = 3.1",
		"mut cond = true; cond = false",
	)
}

func TestUnaryExpressions(t *testing.T) {
	utils.MatchCodegenSnaps(t,
		"mut a = 31; let neg = -a",
		"mut f = 4.2; let neg = -f",
		"mut b = true; let not = !b",
		"mut bits = 478134; let not = ~bits",
	)
}

func TestFunctions(t *testing.T) {
	utils.MatchCodegenSnaps(t,
		`fn add(a, b: i32): i32 {
	return a + b
}

let added = add(1, 4)
let added2 = add(added, 1)`,

		`@extern
fn exit(code: i32)

exit(31)`,
	)
}

func TestPointers(t *testing.T) {
	utils.MatchCodegenSnaps(t,
		"let value = 1; let ptr = &value; let value2 = ptr.*",
		"let ptr = &12; let value = ptr.*",
		`mut value = 1
let ptr = &mut value
ptr.* = value + 1
value = ptr.* + 1`,
	)
}
