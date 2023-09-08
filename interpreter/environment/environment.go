package environment

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/parser/ast"
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
	structs     map[string]ast.StructDeclaration
	kind        scopeKind
	ReturnValue values.RuntimeValue
}

func New() *Environment {
	env := Environment{
		parent:      nil,
		variables:   map[string]values.RuntimeValue{},
		structs:     map[string]ast.StructDeclaration{},
		kind:        GLOBAL_SCOPE,
	}

	return &env
}

func NewChild(parent *Environment, kind scopeKind) *Environment {
	return &Environment{
		parent:    parent,
		variables: map[string]values.RuntimeValue{},
		structs:     map[string]ast.StructDeclaration{},
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
		errors.LogError(errors.DevError(fmt.Sprintf("Cannot find variable %q, it does not exist", varName)))
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

func (env *Environment) DeclareStruct(structDecl *ast.StructDeclaration) {
	env.structs[structDecl.Name] = *structDecl
}

func (env *Environment) GetStruct(name string) ast.StructDeclaration {
	return env.structs[name]
}