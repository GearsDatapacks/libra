package interpreter

import (
	"math"

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

// func makeBinaryOperator[A, B, C any](operator string, operation func(A, B) C) {
// 	RegisterBinaryOperator(operator, func(left, right values.RuntimeValue) values.RuntimeValue {
// 		leftValue := extractValues[A](left)[0]
// 		rightValue := extractValues[B](right)[0]

// 		return values.MakeValue(operation(leftValue, rightValue))
// 	})
// }

// func makeUnaryOperator[A, B any](operator string, operation func(A) B) {
// 	RegisterUnaryOperator(operator, func(runtimeValue values.RuntimeValue) values.RuntimeValue {
// 		value := extractValues[A](runtimeValue)[0]

// 		return values.MakeValue(operation(value))
// 	})
// }

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
	RegisterBinaryOperator(
		"+",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			_, isBInt := b.(*values.IntegerLiteral)

			if isAInt && isBInt {
				intValues := extractValues[int](a, b)
				return values.MakeInteger(intValues[0] + intValues[1])
			}

			_, isAString := a.(*values.StringLiteral)
			_, isBString := b.(*values.StringLiteral)

			if isAString && isBString {
				stringValues := extractValues[string](a, b)
				return values.MakeString(stringValues[0] + stringValues[1])
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
	)

	RegisterBinaryOperator(
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
	)

	RegisterBinaryOperator(
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
	)

	RegisterBinaryOperator(
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
	)

	RegisterBinaryOperator(
		"**",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			_, isAInt := a.(*values.IntegerLiteral)
			var valueA float64
			valueB := float64(extractValues[int](b)[0])

			if isAInt {
				valueA = float64(extractValues[int](a)[0])
				return values.MakeInteger(int(math.Pow(valueA, valueB)))
			}

			valueA = extractValues[float64](a)[0]

			return values.MakeFloat(math.Pow(valueA, valueB))
		},
	)

	RegisterBinaryOperator(
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
	)

	RegisterBinaryOperator(
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
	)

	RegisterBinaryOperator(
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
	)

	RegisterBinaryOperator(
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
	)

	RegisterBinaryOperator(
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
	)

	RegisterBinaryOperator(
		"==",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			return values.MakeBoolean(a.EqualTo(b))
		},
	)

	RegisterBinaryOperator(
		"!=",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			return values.MakeBoolean(!a.EqualTo(b))
		},
	)

	RegisterUnaryOperator("!", func(value values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
		return values.MakeBoolean(!value.Truthy())
	})

	RegisterUnaryOperator("++", func(value values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
		_, isInt := value.(*values.IntegerLiteral)
		if isInt {
			intValue := extractValues[int](value)[0]
			return env.AssignVariable(value.Varname(), values.MakeInteger(intValue + 1))
		}
		floatValue := extractValues[float64](value)[0]
		return env.AssignVariable(value.Varname(), values.MakeFloat(floatValue + 1.0))
	})

	RegisterUnaryOperator("--", func(value values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
		_, isInt := value.(*values.IntegerLiteral)
		if isInt {
			intValue := extractValues[int](value)[0]
			return env.AssignVariable(value.Varname(), values.MakeInteger(intValue - 1))
		}
		floatValue := extractValues[float64](value)[0]
		return env.AssignVariable(value.Varname(), values.MakeFloat(floatValue - 1.0))
	})

	RegisterUnaryOperator("!", func(value values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
		return values.MakeBoolean(!value.Truthy())
	})

	RegisterBinaryOperator("||", func(a, b values.RuntimeValue) values.RuntimeValue {
		if a.Truthy() {
			return a
		}
		return b
	})

	RegisterBinaryOperator("&&", func(a, b values.RuntimeValue) values.RuntimeValue {
		if a.Truthy() {
			return b
		}
		return a
	})
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
