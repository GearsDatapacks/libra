package main

import (
	"fmt"
	"os"

	"github.com/gearsdatapacks/libra/lowerer"
	"github.com/gearsdatapacks/libra/module"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
)

type debugKind int

const (
	none debugKind = iota
	ast
	ir
	lowered
)

func main() {
	mod, diags := module.Load(os.Args[1])
	debugKind := none
	if len(os.Args) > 3 && os.Args[2] == "--debug" {
		switch os.Args[3] {
		case "ast":
			debugKind = ast
		case "ir":
			debugKind = ir
		case "lowered":
			debugKind = lowered
		}
	}

	if len(diags) != 0 {
		for _, diag := range diags {
			diag.Print()
		}
		return
	}

	if debugKind == ast {
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

	if debugKind == ir {
		pkg.Print()
		fmt.Println()
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
