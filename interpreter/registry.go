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

func registerOperators() {
	RegisterOperator("+", "integer", "integer", func(left values.RuntimeValue, right values.RuntimeValue) values.RuntimeValue {
		ints := extractValues[int](left, right)

		res := values.MakeInteger(ints[0] + ints[1])
		return &res
	})

	RegisterOperator("*", "integer", "integer", func(left values.RuntimeValue, right values.RuntimeValue) values.RuntimeValue {
		ints := extractValues[int](left, right)

		res := values.MakeInteger(ints[0] * ints[1])
		return &res
	})
}
