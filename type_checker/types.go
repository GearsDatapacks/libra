package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (t *typeChecker) typeFromAst(node ast.TypeExpression) types.Type {
	switch ty := node.(type) {
	case *ast.TypeName:
		return lookupType(ty.Name.Value)
	case *ast.ArrayType:
		elemType := t.typeFromAst(ty.Type)

		if ty.Count == nil {
			return &types.ListType{
				ElemType: elemType,
			}
		}

		length := -1

		if _, ok := ty.Count.(*ast.InferredExpression); !ok {
			expr := convert(t.typeCheckExpression(ty.Count), types.Int, implicit)
			if expr == nil {
				t.Diagnostics.ReportCountMustBeInt(ty.Count.Location())
			} else if expr.IsConst() {
				value := expr.ConstValue().(values.IntValue)
				length = int(value.Value)
			} else {
				t.Diagnostics.ReportNotConst(ty.Count.Location())
			}
		}

		return &types.ArrayType{
			ElemType: elemType,
			Length:   length,
			CanInfer: false,
		}
	case *ast.MapType:
		keyType := t.typeFromAst(ty.KeyType)
		valueType := t.typeFromAst(ty.ValueType)
		return &types.MapType{
			KeyType:   keyType,
			ValueType: valueType,
		}
	default:
		panic(fmt.Sprintf("TODO: Types from %T", ty))
	}
}

func lookupType(name string) types.Type {
	switch name {
	case "i32":
		return types.Int
	case "f32":
		return types.Float
	case "bool":
		return types.Bool
	case "string":
		return types.String
	default:
		return nil
	}
}
