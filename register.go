package main

import (
	"github.com/gearsdatapacks/libra/interpreter"
	"github.com/gearsdatapacks/libra/type_checker/registry"
)

func register() {
	interpreter.Register()
	registry.Register()
}
