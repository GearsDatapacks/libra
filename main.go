package main

import (
	"fmt"
	"os"

	"github.com/gearsdatapacks/libra/lowerer"
	"github.com/gearsdatapacks/libra/module"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
)

func main() {
	mod, diags := module.Load(os.Args[1])
	debugAst := false
	if len(os.Args) > 2 && os.Args[2] == "--ast" {
		debugAst = true
	}

	if len(diags) != 0 {
		for _, diag := range diags {
			diag.Print()
		}
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

	pkg, diags := typechecker.TypeCheck(mod, diags)

	if len(diags) != 0 {
		for _, diag := range diags {
			diag.Print()
		}
		return
	}

	lowered, diags := lowerer.Lower(pkg, diags)

	if len(diags) != 0 {
		for _, diag := range diags {
			diag.Print()
		}
		return
	}

	lowered.Print()
	fmt.Println()
}
