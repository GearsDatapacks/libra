package symbols

import "github.com/gearsdatapacks/libra/type_checker/types"

type FunctionContext struct {
	ReturnType types.Type
}

type LoopContext struct {
	ResultType types.Type
}
type BlockContext struct {
	ResultType types.Type
}

type globalContext struct {
	methods map[string][]*Method
	exportedMethods map[string][]*Method
	exports map[string]Symbol
}
