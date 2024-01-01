package environment

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type scopeKind int

const (
	GLOBAL_SCOPE scopeKind = iota
	GENERIC_SCOPE
	FUNCTION_SCOPE
)

type Environment struct {
	Parent    *Environment
	variables map[string]values.RuntimeValue
	// types       map[string]types.ValidType
	kind        scopeKind
	ReturnValue values.RuntimeValue
	Exports     map[string]values.RuntimeValue
}

func New() *Environment {
	env := Environment{
		variables: map[string]values.RuntimeValue{},
		// types:       map[string]types.ValidType{},
		kind:    GLOBAL_SCOPE,
		Exports: map[string]values.RuntimeValue{},
	}

	return &env
}

func NewChild(parent *Environment, kind scopeKind) *Environment {
	return &Environment{
		Parent:    parent,
		variables: map[string]values.RuntimeValue{},
		// types:       parent.types,
		kind:    kind,
	}
}

func (env *Environment) DeclareVariable(name string, varType types.ValidType, value values.RuntimeValue) values.RuntimeValue {
	return env.setVariable(name, varType, value)
}

func (env *Environment) AssignVariable(name string, varType types.ValidType, value values.RuntimeValue) values.RuntimeValue {
	declaredenvironment := env.resolve(name)

	return declaredenvironment.setVariable(name, varType, value)
}

func (env *Environment) setVariable(name string, varType types.ValidType, value values.RuntimeValue) values.RuntimeValue {
	if castable, ok := value.(values.AutoCastable); ok {
		value = castable.AutoCast(varType)
	} else {
		value = value.Copy()
	}

	env.variables[name] = value
	value.SetVarname(name)
	return value
}

func (env *Environment) GetVariable(name string) values.RuntimeValue {
	declaredEnvironment := env.resolve(name)
	if declaredEnvironment == nil {
		errors.LogError(errors.DevError(fmt.Sprintf("Cannot find variable %q, it does not exist", name)))
	}
	return declaredEnvironment.variables[name]
}

func (env *Environment) resolve(varName string) *Environment {
	if _, ok := env.variables[varName]; ok {
		return env
	}

	if env.Parent == nil {
		return nil
	}

	return env.Parent.resolve(varName)
}

func (env *Environment) isFunctionScope() bool {
	return env.kind == FUNCTION_SCOPE
}

func (env *Environment) FindFunctionScope() *Environment {
	if env.isFunctionScope() {
		return env
	}
	if env.Parent == nil {
		errors.LogError(errors.DevError("Cannot use return statement outside of a function"))
	}
	return env.Parent.FindFunctionScope()
}

/*
	func (env *Environment) AddType(name string, dataType types.ValidType) {
		env.types[name] = dataType
	}

	func (env *Environment) GetType(name string) types.ValidType {
		return env.types[name]
	}
*/
func (env *Environment) Exists(name string) bool {
	return env.resolve(name) != nil
}

var methods = map[string][]*values.FunctionValue{}

func GetMethod(name string, methodOf types.ValidType) *values.FunctionValue {
	overloads, ok := methods[name]
	if !ok {
		return nil
	}

	for _, overload := range overloads {
		if overload.Type().(*types.Function).MethodOf.Valid(methodOf) {
			return overload
		}
	}

	return nil
}

func AddMethod(name string, method *values.FunctionValue) {
	overloads, ok := methods[name]
	if !ok {
		methods[name] = []*values.FunctionValue{method}
	}
	overloads = append(overloads, method)
	methods[name] = overloads
}

func (env *Environment) GlobalScope() *Environment {
	if env.Parent == nil {
		return env
	}

	return env.Parent.GlobalScope()
}
