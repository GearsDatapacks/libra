package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gearsdatapacks/libra/interpreter"
	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/parser"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
)

func repl() {
	fmt.Println("Libra repl v0.1.0")
	nextLine := ""
	reader := bufio.NewReader(os.Stdin)

	parser := parser.New()
	env := environment.New()

	for strings.ToLower(strings.TrimSpace(nextLine)) != "exit" {
		fmt.Print("> ")

		input, err := reader.ReadBytes('\n')
		nextLine = string(input)

		if err != nil {
			log.Fatal(err)
		}

		lexer := lexer.New(input)
		tokens := lexer.Tokenise()

		ast := parser.Parse(tokens)

		result := interpreter.Evaluate(ast, env)
		fmt.Println(result.ToString())
	}
}

func run(file string) {
	code, err := os.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	lexer := lexer.New(code)
	parser := parser.New()
	env := environment.New()

	tokens := lexer.Tokenise()
	ast := parser.Parse(tokens)

	if !typechecker.TypeCheck(ast) {
		log.Fatal("Invalid types")
	}

	result := interpreter.Evaluate(ast, env)
	fmt.Println(result.ToString())
}

func main() {
	if len(os.Args) == 1 {
		repl()
	} else {
		run(os.Args[1])
	}
}
