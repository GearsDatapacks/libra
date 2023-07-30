package interpreter

import "github.com/gearsdatapacks/libra/interpreter/values"

func Register() {
	registerOperators()
}

func extractInts(vals ...values.RuntimeValue) []int {
	result := []int{}
	
	for _, value := range vals {
		intValue := value.(*values.IntegerLiteral).Value().(int)
		result = append(result, intValue)
	}

	return result
}

func registerOperators() {
	RegisterOperator("+", "integer", "integer", func(left values.RuntimeValue, right values.RuntimeValue) values.RuntimeValue {
		ints := extractInts(left, right)

		res := values.MakeInteger(ints[0] + ints[1])
		return &res
	})

	RegisterOperator("*", "integer", "integer", func(left values.RuntimeValue, right values.RuntimeValue) values.RuntimeValue {
		ints := extractInts(left, right)

		res := values.MakeInteger(ints[0] * ints[1])
		return &res
	})
}
