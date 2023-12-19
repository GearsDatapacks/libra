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

func err(ty types.ValidType) types.ValidType {
	return &types.ErrorType{ResultType: ty}
}

func registerBuiltins() {
	registerBuiltin("print", params{&types.Any{}}, &types.Void{})
	registerBuiltin("printil", params{&types.Any{}}, &types.Void{})
	registerBuiltin("prompt", params{stringType}, stringType)
	registerBuiltin("to_string", params{&types.Any{}}, stringType)
	registerBuiltin("parse_int", params{stringType}, err(intType))
	registerBuiltin("parse_float", params{stringType}, err(floatType))
	registerBuiltin("read_file", params{stringType}, err(stringType))
	registerBuiltin("write_file", params{stringType, stringType}, err(&types.Void{}))
}
