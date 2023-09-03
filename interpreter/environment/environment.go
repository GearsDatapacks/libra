package environment

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/values"
)

type scopeKind int

const (
	GLOBAL_SCOPE scopeKind = iota
	GENERIC_SCOPE
	FUNCTION_SCOPE
)

type Environment struct {
	parent      *Environment
	variables   map[string]values.RuntimeValue
	kind        scopeKind
	ReturnValue values.RuntimeValue
}

func New() *Environment {
	env := Environment{
		parent:    nil,
		variables: map[string]values.RuntimeValue{},
		kind:      GLOBAL_SCOPE,
	}

	return &env
}

func NewChild(parent *Environment, kind scopeKind) *Environment {
	return &Environment{
		parent:    parent,
		variables: map[string]values.RuntimeValue{},
		kind:      kind,
	}
}

func (env *Environment) DeclareVariable(name string, value values.RuntimeValue) values.RuntimeValue {
	value.SetVarname(name)
	env.variables[name] = value

	return value
}

func (env *Environment) AssignVariable(name string, value values.RuntimeValue) values.RuntimeValue {
	declaredenvironment := env.resolve(name)

	value.SetVarname(name)
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
		errors.LogError(errors.RuntimeError(fmt.Sprintf("Cannot find variable %q, it does not exist", varName)))
	}

	return env.parent.resolve(varName)
}

func (env *Environment) isFunctionScope() bool {
	return env.kind == FUNCTION_SCOPE
}

func (env *Environment) FindFunctionScope() *Environment {
	if env.isFunctionScope() {
		return env
	}
	if env.parent == nil {
		errors.LogError(errors.DevError("Cannot use return statement outside of a function"))
	}
	return env.parent.FindFunctionScope()
}
