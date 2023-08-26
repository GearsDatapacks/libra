package registry

import "github.com/gearsdatapacks/libra/type_checker/types"

var boolType = &types.BoolLiteral{}
var floatType = &types.FloatLiteral{}
var intType = &types.IntLiteral{}
var numberType = types.MakeUnion(intType, floatType)
var stringType = &types.StringLiteral{}

func Register() {
	registerOperators()
	registerBuiltins()
}
