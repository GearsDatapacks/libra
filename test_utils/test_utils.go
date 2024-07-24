package testutils

import (
	"bytes"
	"fmt"
	"testing"
)

func Assert(t *testing.T, condition bool, msg ...string) {
	t.Helper()
	if condition {
		return
	}

	if len(msg) > 0 {
		t.Fatal(msg[0])
		return
	}

	t.Fatal("Assertion failed")
}

func AssertEq[T comparable](t *testing.T, actual, expected T, msg ...string) {
	t.Helper()
	var defaultMsg = fmt.Sprintf("Expected %v, got %v", expected, actual)
	Assert(t, actual == expected, append(msg, defaultMsg)...)
}

func AssertSingle[T any](t *testing.T, list []T) T {
	t.Helper()

	AssertEq(t, len(list), 1, fmt.Sprintf("Expected a single list item, got %d", len(list)))
	return list[0]
}

func MatchAstSnaps(t *testing.T, tests ...string) {
	t.Helper()

	for _, src := range tests {
		program, diags := getAst(t, src)
		for _, diag := range diags {
			diag.Print()
		}
		AssertEq(t, len(diags), 0,
			fmt.Sprintf("Expected no diagnostics (got %d)", len(diags)))

		matchSnap(t, src, program.String())
	}
}

func MatchIrSnaps(t *testing.T, tests ...string) {
	t.Helper()

	for _, src := range tests {
		program, diags := getIr(t, src)
		for _, diag := range diags {
			diag.Print()
		}
		AssertEq(t, len(diags), 0,
			fmt.Sprintf("Expected no diagnostics (got %d)", len(diags)))
		matchSnap(t, src, program.String())
	}
}

func MatchParserErrorSnaps(t *testing.T, tests ...string) {
	t.Helper()

	for _, src := range tests {
		_, diagnostics := getAst(t, src)
		var diags bytes.Buffer
		for _, diag := range diagnostics {
			diag.WriteTo(&diags, false)
		}

		matchSnap(t, src, diags.String())
	}
}

func MatchTCErrorSnaps(t *testing.T, tests ...string) {
	t.Helper()

	for _, src := range tests {
		_, diagnostics := getIr(t, src)
		var diags bytes.Buffer
		for _, diag := range diagnostics {
			diag.WriteTo(&diags, false)
		}

		matchSnap(t, src, diags.String())
	}
}
