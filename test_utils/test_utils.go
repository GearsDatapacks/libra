package testutils

import (
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
