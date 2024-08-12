package lowerer_test

import (
	"testing"

	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestIfStatements(t *testing.T) {
	utils.MatchLoweredSnaps(t,
		`mut condition = true
if condition {
	let a = 10
}`,
		`mut condition = true
if condition {
	let a = 10
} else {
	let b = 20
}`,
		`mut condition = true
if condition {
	let a = 10
} else if !condition {
	let b = 20
} else {
	let c = 30
}`,
	)
}

func TestWhileLoops(t *testing.T) {
	utils.MatchLoweredSnaps(t,
		`mut i = 0
while i < 10 {
	i++
}`,
	)
}

func TestConstantFolding(t *testing.T) {
	utils.MatchLoweredSnaps(t,
		"let value = 1 + 2",
		"let value = true && false",
		`let value = "Hello, " + "world!"`,
		"let value = 1 == 2 || 3 == 3",
		"let value = (4 / 2) + 1",
	)
}

func TestExpressionOptimisation(t *testing.T) {
	utils.MatchLoweredSnaps(t,
		"mut a = 10; let value = !(a == 10)",
		"mut a = 0; let value = !(a < 20)",
		"mut a = true; let value = !!a",
		"mut a = 5; let value = a * 1",
		"mut a = 31; let value = a + 0",
		"mut a = 1; let value = a * 0",
		"mut a = 13; let value = a / 1",
		"mut a = 21; let value = a ** 0",
		"mut a = 20; let value = -(-a)",
		"mut a = 1; let value = ~(~a)",
		"mut a = false; let value = a || true",
		"mut a = true; let value = a && false",
		"mut a = false; let value = a && true",
		"mut a = false; let value = a || false",
	)
}

func TestUncertainReturns(t *testing.T) {
	utils.MatchLowerErrors(t,
		`fn add(a, b: i32): i32 {
	if a == 0 {
		return b
	} else if b == 0 {
		return a
	}
}`,
		`fn foo(a: i32): i32 {
	while a != 0 {
		return a
	}
}`,
	)
}

func TestCertainReturns(t *testing.T) {
	utils.MatchLoweredSnaps(t,
		`fn add(a, b: i32): i32 {
	if true {
		return a + b
	}
}`,
		`fn add(a, b: i32): i32 {
	mut result = a
	mut counter = b
	while true {
		if counter == 0 {
			return result
		}
		result++
		counter--
	}
}`,
	)
}

func TestUnreachableCode(t *testing.T) {
	utils.MatchLoweredSnaps(t,
		`fn foo() {
	if true {
		return
	}
	let bar = 10
	bar + 1
	return
}`,
		`if false {
	let a = 1
} else {
	let b = 1
}`,
		`mut i = 0
while true {
	i++
}
i--`,
		`while false {
	let foo = 1
}
let bar = 2`,
	)
}
