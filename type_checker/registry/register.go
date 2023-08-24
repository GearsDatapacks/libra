package registry

import "github.com/gearsdatapacks/libra/type_checker/types"

var boolType = types.MakeLiteral(types.BOOL)
var floatType = types.MakeLiteral(types.FLOAT)
var intType = types.MakeLiteral(types.INT)
var numberType = types.MakeUnion(intType, floatType)
var stringType = types.MakeLiteral(types.STRING)

func Register() {
	registerOperators()
	registerBuiltins()
}
