package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gearsdatapacks/libra/lexer"
)

func repl() {
	fmt.Println("Libra repl v0.1.0")
	nextLine := ""
	reader := bufio.NewReader(os.Stdin)

	for strings.ToLower(strings.TrimSpace(nextLine)) != "exit" {
		fmt.Print("> ")

		input, err := reader.ReadBytes('\n')
		nextLine = string(input)

		if err != nil {
			log.Fatal(err)
		}

		lexer := lexer.NewLexer(input)
		tokens := lexer.Tokenise()
		fmt.Println(tokens)
	}
}

func run(file string) {
	code, err := os.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(code))
}

func main() {
	if len(os.Args) == 1 {
		repl()
	} else {
		run(os.Args[1])
	}
}
