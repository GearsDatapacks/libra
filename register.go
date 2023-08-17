package main

import (
	"github.com/gearsdatapacks/libra/interpreter"
	typechecker "github.com/gearsdatapacks/libra/type_checker"
)

func register() {
	interpreter.Register()
	typechecker.RegisterOperators()
}
