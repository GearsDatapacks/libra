package main

import (
	"fmt"
	"os"

	"github.com/gearsdatapacks/libra/lexer"
)

func main() {
  code, err := os.ReadFile(os.Args[1])
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  lexer := lexer.New(string(code))
  tokens := lexer.Tokenise()

  /*if len(lexer.Diagnostics > 0) {
    for _, diag := range lexer.Diagnostics {
      fmt.Println(diag)
    }
  }*/

  fmt.Println(tokens)
}

