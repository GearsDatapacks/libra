package main

import (
	"fmt"
	"os"

	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/parser"
)

func main() {
  code, err := os.ReadFile(os.Args[1])
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  lexer := lexer.New(string(code), os.Args[1])
  tokens := lexer.Tokenise()
  parser := parser.New(tokens, lexer.Diagnostics)
  program := parser.Parse()

  if len(parser.Diagnostics.Diagnostics) > 0 {
    for _, diag := range parser.Diagnostics.Diagnostics {
      diag.Print()
    }
  }

  fmt.Println(program.String())
}