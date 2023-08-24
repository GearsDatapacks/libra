package interpreter

import (
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
)

func Register() {
	registerOperators()
	registerBuiltins()
}

func extractValues[T any](vals ...values.RuntimeValue) []T {
	result := []T{}
	
	for _, value := range vals {
		typedValue := value.Value().(T)
		result = append(result, typedValue)
	}

	return result
}

func makeOperator[A, B, C any] (operator string, operation func(A, B)C) {
	leftType := values.TypeToString[A]()
	rightType := values.TypeToString[B]()
	RegisterOperator(operator, leftType, rightType, func(left, right values.RuntimeValue) values.RuntimeValue {
		leftValue := extractValues[A](left)[0]
		rightValue := extractValues[B](right)[0]

		return values.MakeValue(operation(leftValue, rightValue))
	})
}

func makeOverloadedOperator(operator string, operation opFn, validTypes []string) {
	for _, leftType := range validTypes {
		for _, rightType := range validTypes {
			RegisterOperator(operator, leftType, rightType, operation)
		}
	}
}

func modulo(a, b float64) float64 {
	for a > b {
		a -= b
	}

	for a < -b {
		a += b
	}

	return a
}

func registerOperators() {
	makeOverloadedOperator(
		"+",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			if isAInt && isBInt {
				intValues := extractValues[int](a, b)
				return values.MakeInteger(intValues[0] + intValues[1])
			}

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
			} else {
				valueA = extractValues[float64](a)[0]
			}

			if isBInt {
				valueB = float64(extractValues[int](b)[0])
			} else {
				valueB = extractValues[float64](b)[0]
			}

			return values.MakeFloat(valueA + valueB)
		},
		[]string{"integer", "float"},
	)

	makeOverloadedOperator(
		"-",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			if isAInt && isBInt {
				intValues := extractValues[int](a, b)
				return values.MakeInteger(intValues[0] - intValues[1])
			}

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
			} else {
				valueA = extractValues[float64](a)[0]
			}

			if isBInt {
				valueB = float64(extractValues[int](b)[0])
			} else {
				valueB = extractValues[float64](b)[0]
			}

			return values.MakeFloat(valueA - valueB)
		},
		[]string{"integer", "float"},
	)

	makeOverloadedOperator(
		"*",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			if isAInt && isBInt {
				intValues := extractValues[int](a, b)
				return values.MakeInteger(intValues[0] * intValues[1])
			}

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
			} else {
				valueA = extractValues[float64](a)[0]
			}

			if isBInt {
				valueB = float64(extractValues[int](b)[0])
			} else {
				valueB = extractValues[float64](b)[0]
			}

			return values.MakeFloat(valueA * valueB)
		},
		[]string{"integer", "float"},
	)

	makeOverloadedOperator(
		"/",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
			} else {
				valueA = extractValues[float64](a)[0]
			}

			if isBInt {
				valueB = float64(extractValues[int](b)[0])
			} else {
				valueB = extractValues[float64](b)[0]
			}

			return values.MakeFloat(valueA / valueB)
		},
		[]string{"integer", "float"},
	)

	makeOverloadedOperator(
		"%",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			if isAInt && isBInt {
				intValues := extractValues[int](a, b)
				return values.MakeInteger(intValues[0] % intValues[1])
			}

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
			} else {
				valueA = extractValues[float64](a)[0]
			}

			if isBInt {
				valueB = float64(extractValues[int](b)[0])
			} else {
				valueB = extractValues[float64](b)[0]
			}

			return values.MakeFloat(modulo(valueA, valueB))
		},
		[]string{"integer", "float"},
	)
	

	makeOverloadedOperator(
		">",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
			} else {
				valueA = extractValues[float64](a)[0]
			}

			if isBInt {
				valueB = float64(extractValues[int](b)[0])
			} else {
				valueB = extractValues[float64](b)[0]
			}

			return values.MakeBoolean(valueA > valueB)
		},
		[]string{"integer", "float"},
	)

	makeOverloadedOperator(
		">=",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
			} else {
				valueA = extractValues[float64](a)[0]
			}

			if isBInt {
				valueB = float64(extractValues[int](b)[0])
			} else {
				valueB = extractValues[float64](b)[0]
			}

			return values.MakeBoolean(valueA >= valueB)
		},
		[]string{"integer", "float"},
	)

	makeOverloadedOperator(
		"<",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
			} else {
				valueA = extractValues[float64](a)[0]
			}

			if isBInt {
				valueB = float64(extractValues[int](b)[0])
			} else {
				valueB = extractValues[float64](b)[0]
			}

			return values.MakeBoolean(valueA < valueB)
		},
		[]string{"integer", "float"},
	)

	makeOverloadedOperator(
		"<=",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
			} else {
				valueA = extractValues[float64](a)[0]
			}

			if isBInt {
				valueB = float64(extractValues[int](b)[0])
			} else {
				valueB = extractValues[float64](b)[0]
			}

			return values.MakeBoolean(valueA <= valueB)
		},
		[]string{"integer", "float"},
	)

	makeOverloadedOperator(
		"==",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			return values.MakeBoolean(a.EqualTo(b))
		},
		[]string{"integer", "float", "boolean", "string", "null", "function"},
	)

	makeOverloadedOperator(
		"!=",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			return values.MakeBoolean(!a.EqualTo(b))
		},
		[]string{"integer", "float", "boolean", "string", "null", "function"},
	)

	makeOperator("||", func(a, b bool) bool { return a || b })
	makeOperator("&&", func(a, b bool) bool { return a && b })
}

type builtin func([]values.RuntimeValue, *environment.Environment) values.RuntimeValue
var builtins = map[string]builtin{}

func registerBuiltins() {
	builtins["print"] = print
	builtins["printil"] = printil
	builtins["prompt"] = prompt
	builtins["toString"] = toString
	builtins["parseInt"] = parseInt
	builtins["parseFloat"] = parseFloat
}
