package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
)

func print(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	printStr := args[0].ToString()
	
	if _, ok := args[0].(*values.StringLiteral); ok {
		printStr = printStr[1:len(printStr)-1]
	}

	fmt.Println(printStr)

	return values.MakeNull()
}
