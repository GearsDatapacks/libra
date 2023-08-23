package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/interpreter/values"
)

func print(args []values.RuntimeValue, env *environment.Environment) values.RuntimeValue {
	fmt.Println(args[0].(*values.StringLiteral).Value())
	return values.MakeNull()
}
