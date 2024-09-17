package lowerer

import (
	"fmt"

	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type abiClass int

const (
	noClass abiClass = iota
	integer
	sse
	sseUp
	// x87
	// x87Up
	// complexX87
	memory
)

func fixAbi(pkg *ir.LoweredPackage) {
	for _, mod := range pkg.Modules {
		for _, fn := range mod.Functions {
			for i, param := range fn.Type.Parameters {
				var low, high abiClass
				classify(param, 0, &low, &high)

				if types.IsStruct(param) {
					if low == integer && high == noClass {
						newType := types.Int(types.BitSize(param))
						fn.Type.Parameters[i] = newType
					} else if low == sse && high == noClass {
						newType := types.Float(types.BitSize(param))
						fn.Type.Parameters[i] = newType
					}
				}
			}
		}

		for _, call := range mod.FunctionCalls {
			for i, arg := range call.Arguments {
				argType := arg.Type()
				var low, high abiClass
				classify(argType, 0, &low, &high)

				if types.IsStruct(argType) {
					if low == integer && high == noClass {
						newType := types.Int(types.BitSize(argType))
						call.Arguments[i] = &ir.BitCast{
							Value: arg,
							To:    newType,
						}
					} else if low == sse && high == noClass {
						newType := types.Float(types.BitSize(argType))
						call.Arguments[i] = &ir.BitCast{
							Value: arg,
							To:    newType,
						}
					}
				}
			}
		}
	}
}

const bits = 1
const bytes = 8 * bits
const eightBytes = 8 * bytes

func classify(
	ty types.Type,
	offset int,
	low, high *abiClass,
) {
	bitWidth := types.BitSize(ty)

	current := low
	if offset > 1*eightBytes {
		current = high
	}

	*low = noClass
	*high = noClass
	*current = memory

	if ty == types.Void {
		*current = noClass
		return
	}

	if types.IsBool(ty) {
		*current = integer
		return
	}

	// TODO: Not CStrings
	if types.IsString(ty) {
		*current = integer
		return
	}

	if types.IsInt(ty) {
		if bitWidth <= 64 {
			*current = integer
		} else if bitWidth <= 128 {
			*low = integer
			*high = integer
		} else {
			panic("TODO: >128 bit-width ints")
		}

		return
	}

	if types.IsFloat(ty) {
		if bitWidth <= 64 {
			*current = sse
		} else if bitWidth == 128 {
			*low = sse
			*high = sseUp
		} else {
			panic("TODO: >128 bit-width floats")
		}

		return
	}

	if types.IsPtr(ty) {
		*low = integer
	}

	if structTy, ok := types.Unwrap(ty).(*types.Struct); ok {
		if bitWidth > 8*eightBytes {
			// current is already set to memory
			return
		}

		fieldOffset := 0
		*current = noClass

		for _, field := range structTy.Fields {
			fieldWidth := types.BitSize(field.Type)
			currentOffset := fieldOffset
			fieldOffset += fieldWidth

			// TODO: Account for unaligned fields (when possible)
			var fieldLow, fieldHigh abiClass
			classify(field.Type, currentOffset, &fieldLow, &fieldHigh)

			*low = merge(*low, fieldLow)
			*high = merge(*high, fieldHigh)

			// As soon as one field is passed in memory, the whole struct is
			if *low == memory || *high == memory {
				break
			}
		}

		postMerge(bitWidth, low, high)
		return
	}
	panic(fmt.Sprintf("TODO: ABI for %T: %s", ty, ty.String()))
}

func merge(main, other abiClass) abiClass {
	if main == noClass {
		return other
	} else if main == other || other == noClass {
		return main
	} else if main == memory || other == memory {
		return memory
	} else if main == integer || other == integer {
		return integer
	} else {
		return sse
	}
}

func postMerge(bitWidth int, low, high *abiClass) {
	if *high == memory {
		*low = memory
	}

	if bitWidth > 2*eightBytes && !(*low == sse && *high == sseUp) {
		*low = memory
	}

	if *high == sseUp && *low != sse {
		*high = sse
	}
}
