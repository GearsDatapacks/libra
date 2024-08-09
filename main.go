package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gearsdatapacks/libra/codegen"
	"github.com/gearsdatapacks/libra/lowerer"
	"github.com/gearsdatapacks/libra/module"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
	"tinygo.org/x/go-llvm"
)

type debugKind int

const (
	none debugKind = iota
	ast
	ir
	lowered
	llir
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
		case "llir":
			debugKind = llir
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

	loweredPkg, diags := lowerer.Lower(pkg, diags)

	if len(diags) != 0 {
		for _, diag := range diags {
			diag.Print()
		}
		return
	}

	if debugKind == lowered {
		loweredPkg.Print()
		fmt.Println()
		return
	}

	module := codegen.Compile(loweredPkg)

	if debugKind == llir {
		fmt.Println(module.String())
	}

	outputCode(module)
}

func outputCode(module llvm.Module) {
	triple := llvm.DefaultTargetTriple()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
	target, err := llvm.GetTargetFromTriple(triple)
	if err != nil {
		log.Fatal(err)
	}
	cpu := "generic"
	features := ""
	machine := target.CreateTargetMachine(
		triple,
		cpu,
		features,
		llvm.CodeGenLevelDefault,
		llvm.RelocPIC,
		llvm.CodeModelDefault,
	)
	module.SetTarget(triple)

	buffer, err := machine.EmitToMemoryBuffer(module, llvm.ObjectFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.WriteFile("out.o", buffer.Bytes(), os.ModePerm)
}
