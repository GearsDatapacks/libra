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

  tc := typechecker.New(diags)
	program := tc.TypeCheck(mod.Files[0].Ast)
	fmt.Println(program.String())
}
