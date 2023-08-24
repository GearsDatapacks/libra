package registry

import "github.com/gearsdatapacks/libra/type_checker/types"

type params = []types.ValidType

type builtin struct {
	Parameters params
	ReturnType types.ValidType
}

var Builtins = map[string]builtin{}

func registerBuiltin(name string, parameters params, returnType types.ValidType) {
	data := builtin{
		Parameters: parameters,
		ReturnType: returnType,
	}

	Builtins[name] = data
}

func registerBuiltins() {
	registerBuiltin("print", params{&types.Any{}}, &types.Void{})
}
