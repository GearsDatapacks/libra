package testutils

import (
	"fmt"
	"testing"
)

func Assert(t *testing.T, condition bool, msg ...string) {
  if condition {
    return
  }

  if len(msg) > 0 {
    t.Error(msg[0])
    return
  }

  t.Error("Condition is not true")
}

func AssertEq[T comparable](t *testing.T, actual, expected T, msg ...string) {
  var defaultMsg = fmt.Sprintf("Expected %v, got %v", expected, actual)
  Assert(t, actual == expected, append(msg, defaultMsg)...)
}

