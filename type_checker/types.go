package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (t *typeChecker) typeFromAst(node ast.TypeExpression) types.Type {
	switch ty := node.(type) {
	case *ast.TypeName:
		return t.lookupType(ty.Name.Value)
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
				t.Diagnostics.Report(diagnostics.CountMustBeInt(ty.Count.Location()))
			} else if expr.IsConst() {
				value := expr.ConstValue().(values.IntValue)
				length = int(value.Value)
			} else {
				t.Diagnostics.Report(diagnostics.NotConst(ty.Count.Location()))
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
	case *ast.TupleType:
		dataTypes := []types.Type{}
		for _, ty := range ty.Types {
			dataTypes = append(dataTypes, t.typeFromAst(ty))
		}
		return &types.TupleType{
			Types: dataTypes,
		}
	default:
		panic(fmt.Sprintf("TODO: Types from %T", ty))
	}
}

func (t *typeChecker) lookupType(name string) types.Type {
	switch name {
	case "i32":
		return types.Int
	case "f32":
		return types.Float
	case "bool":
		return types.Bool
	case "string":
		return types.String
	case "Type":
		return types.RuntimeType
	default:
		variable := t.symbols.Lookup(name)
		if variable == nil {
			return nil
		}
		if variable.GetType() != types.RuntimeType {
			return nil
		}
		if variable.Value() == nil {
			return nil
		}
		return variable.Value().(values.TypeValue).Type.(types.Type)
	}
}
