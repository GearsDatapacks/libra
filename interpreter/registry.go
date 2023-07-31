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

func makeOperator[A, B, C any] (operator, leftType, rightType string,  operation func(A, B)C) {
	RegisterOperator(operator, leftType, rightType, func(left values.RuntimeValue, right values.RuntimeValue) values.RuntimeValue {
		leftValue := extractValues[A](left)[0]
		rightValue := extractValues[B](right)[0]

		return values.MakeValue(operation(leftValue, rightValue))
	})
}

func registerOperators() {
	makeOperator("+", "integer", "integer", func(a, b int) int { return a + b })
	makeOperator("-", "integer", "integer", func(a, b int) int { return a - b })
	makeOperator("*", "integer", "integer", func(a, b int) int { return a * b })
	makeOperator("/", "integer", "integer", func(a, b int) int { return a / b })
}
