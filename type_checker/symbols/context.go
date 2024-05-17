package symbols

import "github.com/gearsdatapacks/libra/type_checker/types"

type FunctionContext struct {
	ReturnType types.Type
}

type LoopContext struct{}

type GlobalContext struct {
	Methods map[string][]types.Method
}
