package registry

import "github.com/gearsdatapacks/libra/type_checker/types"

var boolType = &types.BoolLiteral{}
var floatType = &types.FloatLiteral{}
var untypedNumberType = &types.UntypedNumber{}
var intType = &types.IntLiteral{}
var numberType = types.MakeUnion(intType, floatType, untypedNumberType)
var stringType = &types.StringLiteral{}

func Register() {
	registerOperators()
	registerBuiltins()
}
