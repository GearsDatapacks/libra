package environment

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/values"
)

type Environment struct {
	parent *Environment
	variables map[string]values.RuntimeValue
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

func (env *Environment) DeclareVariable(name string, value values.RuntimeValue) values.RuntimeValue {
	env.variables[name] = value

	return value
}

func (env *Environment) AssignVariable(name string, value values.RuntimeValue) values.RuntimeValue {
	declaredenvironment := env.resolve(name)

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
		errors.RuntimeError(fmt.Sprintf("Cannot find variable %q, it does not exist", varName))
	}

	return env.parent.resolve(varName)
}
