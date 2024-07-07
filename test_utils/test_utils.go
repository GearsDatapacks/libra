package testutils

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/text"
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

func MatchAstSnap(t *testing.T, src string) {
	t.Helper()
	program := getAst(t, src)
	matchSnap(t, program.String())
}

func MatchIrSnap(t *testing.T, src string) {
	t.Helper()
	program := getIr(t, src)
	matchSnap(t, program)
}

func MatchErrorSnap(t *testing.T, src string) {
	t.Helper()

	lexer := lexer.New(text.NewFile("test.lb", src))
	tokens := lexer.Tokenise()

	AssertEq(t, len(lexer.Diagnostics), 0, "Expected no lexer diagnostics")

	p := parser.New(tokens, lexer.Diagnostics)
	p.Parse()
	var diags bytes.Buffer
	for _, diag := range p.Diagnostics {
		diag.WriteTo(&diags, false)
	}

	matchSnap(t, diags.String())
}
