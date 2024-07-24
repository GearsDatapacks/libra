package main

import (
	"fmt"
	"os"

	"github.com/gearsdatapacks/libra/module"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
)

func main() {
	mod, diags := module.Load(os.Args[1])
	debugAst := false
	if len(os.Args) > 2 && os.Args[2] == "--ast" {
		debugAst = true
	}

	for _, diag := range diags {
		diag.Print()
	}
	if len(diags) != 0 {
		return
	}

	if debugAst {
		for _, file := range mod.Files {
			fmt.Println(file.Path + ":")
			file.Ast.Print()
			fmt.Println()
		}
		return
	}

	program, diags := typechecker.TypeCheck(mod, diags)

	for _, diag := range diags {
		diag.Print()
	}

	if len(diags) == 0 {
		program.Print()
		fmt.Println()
	}
}
