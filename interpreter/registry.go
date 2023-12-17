package interpreter

import (
	"fmt"
	"math"
	"os"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func Register() {
	registerOperators()
	registerBuiltins()
}

// func extractValues[T any](vals ...values.RuntimeValue) []T {
// 	result := []T{}

// 	for _, value := range vals {
// 		typedValue := value.Value().(T)
// 		result = append(result, typedValue)
// 	}

// 	return result
// }

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

var errorType = types.Interface{
	Name: "error",
	Members: map[string]types.ValidType{
		"error": &types.Function{
			Name:       "error",
			Parameters: []types.ValidType{},
			ReturnType: &types.StringLiteral{},
		},
	},
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

type numericType int

const (
	INT numericType = iota
	FLOAT
	UNTYPED
)

func extractNumericValue(val values.RuntimeValue) float64 {
	if intVal, isInt := val.(*values.IntegerLiteral); isInt {
		return float64(intVal.Value)
	}

	if floatVal, isFloat := val.(*values.FloatLiteral); isFloat {
		return floatVal.Value
	}

	if untypedVal, isUntyped := val.(*values.UntypedNumber); isUntyped {
		return untypedVal.Value
	}

	errors.DevError(fmt.Sprintf("Type %T is not numeric", val))
	return 0.0
}

func is[T values.RuntimeValue](a values.RuntimeValue) bool {
	_, ok := a.(T)
	return ok
}

func numericOperator(a, b values.RuntimeValue, op func(float64, float64) float64) values.RuntimeValue {
	aVal := extractNumericValue(a)
	bVal := extractNumericValue(b)

	result := op(aVal, bVal)

	if is[*values.FloatLiteral](a) || is[*values.FloatLiteral](b) {
		return values.MakeFloat(result)
	}

	if is[*values.UntypedNumber](a) || is[*values.UntypedNumber](b) {
		isFloat := result == float64(int(result))
		return values.MakeUntypedNumber(result, isFloat)
	}

	return values.MakeInteger(int(result))
}

func registerOperators() {
	RegisterBinaryOperator(
		"+",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			aString, isAString := a.(*values.StringLiteral)
			bString, isBString := b.(*values.StringLiteral)

			if isAString && isBString {
				return values.MakeString(aString.Value + bString.Value)
			}

			return numericOperator(a, b, func(a, b float64) float64 { return a + b })
		},
	)

	RegisterBinaryOperator(
		"-",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			return numericOperator(a, b, func(a, b float64) float64 { return a - b })
		},
	)

	RegisterBinaryOperator(
		"*",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			return numericOperator(a, b, func(a, b float64) float64 { return a * b })
		},
	)

	RegisterBinaryOperator(
		"/",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			aInt, isAInt := a.(*values.IntegerLiteral)
			bInt, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(aInt.Value)
			} else {
				valueA = a.(*values.FloatLiteral).Value
			}

			if isBInt {
				valueB = float64(bInt.Value)
			} else {
				valueB = b.(*values.FloatLiteral).Value
			}

			return values.MakeFloat(valueA / valueB)
		},
	)

	RegisterBinaryOperator(
		"**",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			aInt, isAInt := a.(*values.IntegerLiteral)
			var valueA float64
			valueB := float64(b.(*values.IntegerLiteral).Value)

			if isAInt {
				valueA = float64(aInt.Value)
				return values.MakeInteger(int(math.Pow(valueA, valueB)))
			}

			valueA = a.(*values.FloatLiteral).Value

			return values.MakeFloat(math.Pow(valueA, valueB))
		},
	)

	RegisterBinaryOperator(
		"%",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			return numericOperator(a, b, func(a, b float64) float64 { return modulo(a, b) })
		},
	)

	RegisterBinaryOperator(
		">",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			aInt, isAInt := a.(*values.IntegerLiteral)
			bInt, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(aInt.Value)
			} else {
				valueA = a.(*values.FloatLiteral).Value
			}

			if isBInt {
				valueB = float64(bInt.Value)
			} else {
				valueB = b.(*values.FloatLiteral).Value
			}

			return values.MakeBoolean(valueA > valueB)
		},
	)

	RegisterBinaryOperator(
		">=",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			aInt, isAInt := a.(*values.IntegerLiteral)
			bInt, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(aInt.Value)
			} else {
				valueA = a.(*values.FloatLiteral).Value
			}

			if isBInt {
				valueB = float64(bInt.Value)
			} else {
				valueB = b.(*values.FloatLiteral).Value
			}

			return values.MakeBoolean(valueA >= valueB)
		},
	)

	RegisterBinaryOperator(
		"<",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			aInt, isAInt := a.(*values.IntegerLiteral)
			bInt, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(aInt.Value)
			} else {
				valueA = a.(*values.FloatLiteral).Value
			}

			if isBInt {
				valueB = float64(bInt.Value)
			} else {
				valueB = b.(*values.FloatLiteral).Value
			}

			return values.MakeBoolean(valueA < valueB)
		},
	)

	RegisterBinaryOperator(
		"<=",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			aInt, isAInt := a.(*values.IntegerLiteral)
			bInt, isBInt := b.(*values.IntegerLiteral)

			var valueA float64
			var valueB float64

			if isAInt {
				valueA = float64(aInt.Value)
			} else {
				valueA = a.(*values.FloatLiteral).Value
			}

			if isBInt {
				valueB = float64(bInt.Value)
			} else {
				valueB = b.(*values.FloatLiteral).Value
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

	RegisterBinaryOperator(
		"<<",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			list, isList := a.(*values.ListLiteral)
			if isList {
				list.Elements = append(list.Elements, b)
				return list
			}

			aInt := a.(*values.IntegerLiteral).Value
			bInt := b.(*values.IntegerLiteral).Value

			return values.MakeInteger(aInt << bInt)
		},
	)

	RegisterBinaryOperator(
		">>",
		func(a, b values.RuntimeValue) values.RuntimeValue {
			list, isList := b.(*values.ListLiteral)
			if isList {
				newElems := []values.RuntimeValue{a}
				newElems = append(newElems, list.Elements...)
				list.Elements = newElems
				return list
			}

			aInt := a.(*values.IntegerLiteral).Value
			bInt := b.(*values.IntegerLiteral).Value

			return values.MakeInteger(aInt >> bInt)
		},
	)

	RegisterUnaryOperator("-", func(value values.RuntimeValue, _ bool, env *environment.Environment) values.RuntimeValue {
		if intValue, isInt := value.(*values.IntegerLiteral); isInt {
			intVal := intValue.Value
			return values.MakeInteger(-intVal)
		}

		floatVal := value.(*values.FloatLiteral).Value
		return values.MakeFloat(-floatVal)
	})

	RegisterUnaryOperator("++", func(value values.RuntimeValue, _ bool, env *environment.Environment) values.RuntimeValue {
		intValue, isInt := value.(*values.IntegerLiteral)
		if isInt {
			intValue.Value++
			return intValue
			// return env.AssignVariable(value.Varname(), values.MakeInteger(intValue + 1))
		}
		floatValue := value.(*values.FloatLiteral)
		floatValue.Value++
		return floatValue
		// return env.AssignVariable(value.Varname(), values.MakeFloat(floatValue + 1.0))
	})

	RegisterUnaryOperator("--", func(value values.RuntimeValue, _ bool, env *environment.Environment) values.RuntimeValue {
		intValue, isInt := value.(*values.IntegerLiteral)
		if isInt {
			intValue.Value--
			return intValue
			// return env.AssignVariable(value.Varname(), values.MakeInteger(intValue - 1))
		}
		floatValue := value.(*values.FloatLiteral)
		floatValue.Value--
		return floatValue
		// return env.AssignVariable(value.Varname(), values.MakeFloat(floatValue - 1.0))
	})

	RegisterUnaryOperator("!", func(value values.RuntimeValue, postfix bool, env *environment.Environment) values.RuntimeValue {
		if !postfix {
			return values.MakeBoolean(!value.Truthy())
		}
		if isError(value) {
			fmt.Println(value.ToString())
			os.Exit(1)
		}
		return value
	})

	RegisterUnaryOperator("?", func(value values.RuntimeValue, _ bool, env *environment.Environment) values.RuntimeValue {
		if isError(value) {
			functionScope := env.FindFunctionScope()
			functionScope.ReturnValue = value
			return values.MakeNull()
		}
		return value
	})
}

func isError(value values.RuntimeValue) bool {
	if _, isRuntimeErr := value.(*values.Error); isRuntimeErr {
		return true
	}

	return errorType.Valid(value.Type())
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
	builtins["readFile"] = readFile
	builtins["writeFile"] = writeFile
}
