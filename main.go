package main

import (
	"fmt"
	"os"

	"github.com/gearsdatapacks/libra/module"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
)

func main() {
	mod, diags := module.Load(os.Args[1])

	for _, diag := range diags {
		diag.Print()
	}
	if len(diags) != 0 {
		return
	}

	program, diags := typechecker.TypeCheck(mod, diags)

	for _, diag := range diags {
		diag.Print()
	}

	fmt.Println(program.String())
}
