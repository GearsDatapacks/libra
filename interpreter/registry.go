package interpreter

import "github.com/gearsdatapacks/libra/interpreter/values"

func Register() {
	registerOperators()
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
	RegisterOperator(operator, leftType, rightType, func(left values.RuntimeValue, right values.RuntimeValue) values.RuntimeValue {
		leftValue := extractValues[A](left)[0]
		rightValue := extractValues[B](right)[0]

		return values.MakeValue(operation(leftValue, rightValue))
	})
}

func registerOperators() {
	makeOperator("+", func(a, b int) int { return a + b })
	makeOperator("-", func(a, b int) int { return a - b })
	makeOperator("*", func(a, b int) int { return a * b })
	makeOperator("/", func(a, b int) int { return a / b })
	makeOperator("%", func(a, b int) int { return a % b })
	makeOperator("+", func(a, b float64) float64 { return a + b })
	makeOperator("-", func(a, b float64) float64 { return a - b })
	makeOperator("*", func(a, b float64) float64 { return a * b })
	makeOperator("/", func(a, b float64) float64 { return a / b })

	makeOperator(">", func(a, b int) bool { return a > b })
	makeOperator(">=", func(a, b int) bool { return a >= b })
	makeOperator("<", func(a, b int) bool { return a < b })
	makeOperator("<=", func(a, b int) bool { return a <= b })
	makeOperator(">", func(a, b float64) bool { return a > b })
	makeOperator(">=", func(a, b float64) bool { return a >= b })
	makeOperator("<", func(a, b float64) bool { return a < b })
	makeOperator("<=", func(a, b float64) bool { return a <= b })

	makeOperator("||", func(a, b bool) bool { return a || b })
	makeOperator("&&", func(a, b bool) bool { return a && b })
}
