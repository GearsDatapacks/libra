package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gearsdatapacks/libra/interpreter"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
)

func repl() {
	fmt.Println("Libra repl v0.1.0")
	reader := bufio.NewReader(os.Stdin)

	parser := parser.New()
	manager, _ := modules.NewManager("", symbols.New())
	env := environment.New()

	for {
		fmt.Print("> ")

		input, err := reader.ReadBytes('\n')
		nextLine := string(input)

		if err != nil {
			os.Exit(0)
		}

		if strings.ToLower(strings.TrimSpace(nextLine)) == "exit" {
			os.Exit(0)
		}

		lexer := lexer.New(input)
		tokens, err := lexer.Tokenise()
		if err != nil {
			fmt.Println(err)
			continue
		}

		ast, err := parser.Parse(tokens)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = typechecker.TypeCheck(&ast, manager)
		if err != nil {
			fmt.Println(err)
			continue
		}

		result := interpreter.Evaluate(&ast, env)
		fmt.Println(result.ToString())
	}
}

func run(file string) {
	mods, err := modules.NewManager(file, symbols.New())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	env := environment.New()
	ast := &mods.Main.Ast

	err = typechecker.TypeCheck(ast, mods)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	interpreter.Evaluate(ast, env)
	// fmt.Println(result.ToString())
}

func main() {
	register()

	if len(os.Args) == 1 {
		repl()
	} else {
		run(os.Args[1])
	}
}
