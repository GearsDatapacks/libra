package main

import (
	"fmt"
	"os"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/module"
)

func main() {
	mod, diags := module.Load(os.Args[1])

	for _, diag := range diags {
		diag.Print()
	}

  printModule(mod)
}

func printModule(mod *module.Module) {
  for _, file := range mod.Files {
    diagnostics.SetColour(diagnostics.Blue)
    fmt.Println(file.Path)
    diagnostics.ResetColour()

    fmt.Println(file.Ast.String())
    fmt.Println()
  }

  for _, imported := range mod.Imported {
    printModule(imported)
  }
}
