package typechecker

import (
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func convert(from ir.Expression, to types.Type, maxKind types.CastKind) ir.Expression {
	kind := types.Cast(from.Type(), to, maxKind)

	if kind == types.IdentityCast {
		return from
	}
	if kind == types.NoCast {
		return nil
	}

	return &ir.Conversion{
		Location:   from.GetLocation(),
		Expression: from,
		To:         to,
	}
}

// Converts two numeric types into the same type, following these rules:
//  1. Identical types get preserved
//  2. Similar types get upcasted to the higher number of bits
//  3. If one type can represent all possible values of another, both types get upcasted
//     to that.
//  4. If all previous checks fail, the user must explicitly specify the types.
//
// Examples:
//
//	i16 + i32 -> i32
//	f32 + f64 -> f64
//	i32 + f32 -> f32
//	u32 + i32 -> error (u32 cannot represent negative numbers and i32 cannot represent u32.MAX)
//	u32 + i64 -> i64 (i64 can represent u32.MAX)
//	u32 + f16 -> error (u32 cannot represent non-integer numbers and f16 cannot represent u32.MAX)
func upcastNumbers(a, b types.Numeric) types.Type {
	// TODO: Add support for explicit types
	if types.Cast(a, b, types.OperatorCast) != types.NoCast {
		return combineNumTypes(b, a)
	}
	if types.Cast(b, a, types.OperatorCast) != types.NoCast {
		return combineNumTypes(a, b)
	}

	return nil
}

func combineNumTypes(main, other types.Numeric) types.Numeric {
	if main.Untyped() && other.Untyped() {
		*main.Downcastable = types.Downcastable{}
	} else {
		main.Downcastable = nil
	}
	return main
}
