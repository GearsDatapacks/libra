package environment

import (
	"log"

	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/utils"
)

type Environment struct {
	parent *Environment
	variables map[string]values.RuntimeValue
	constants []string
}

func New() *Environment {
	env := Environment{
		parent: nil,
		variables: map[string]values.RuntimeValue{},
	}

	return &env
}

func NewChild(parent *Environment) *Environment {
	return &Environment{
		parent: parent,
		variables: map[string]values.RuntimeValue{},
	}
}

func (env *Environment) DeclareVariable(name string, value values.RuntimeValue, constant bool) values.RuntimeValue {
	if _, ok := env.variables[name]; ok {
		log.Fatalf("Cannot redeclare variable %q, it is already defined", name)
	}

	env.variables[name] = value

	if constant {
		env.constants = append(env.constants, name)
	}

	return value
}

func (env *Environment) AssignVariable(name string, value values.RuntimeValue) values.RuntimeValue {
	declaredenvironment := env.resolve(name)

	if utils.Contains(declaredenvironment.constants, name) {
		log.Fatalf("Cannot reassign constant %q", name)
	}

	declaredenvironment.variables[name] = value
	return value
}

func (env *Environment) GetVariable(name string) values.RuntimeValue {
	declaredenvironment := env.resolve(name)
	return declaredenvironment.variables[name]
}

func (env *Environment) resolve(varName string) *Environment {
	if _, ok := env.variables[varName]; ok {
		return env
	}

	if env.parent == nil {
		log.Fatalf("Cannot find variable %q, it does not exist", varName)
	}

	return env.parent.resolve(varName)
}
