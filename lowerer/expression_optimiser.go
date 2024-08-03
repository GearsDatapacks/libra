package lowerer

import (
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func optimiseExpression(expr ir.Expression) ir.Expression {
	return foldConstants(expr)
}

func foldConstants(expression ir.Expression) ir.Expression {
	if !expression.IsConst() {
		return expression
	}
	return constValueToExpr(expression.ConstValue(), expression.Type())
}

func constValueToExpr(constValue values.ConstValue, ty types.Type) ir.Expression {
	switch value := constValue.(type) {
	case values.IntValue:
		return &ir.IntegerLiteral{
			Value: value.Value,
		}
	case values.UintValue:
		return &ir.UintLiteral{
			Value: value.Value,
		}
	case values.FloatValue:
		return &ir.FloatLiteral{
			Value: value.Value,
		}
	case values.BoolValue:
		return &ir.BooleanLiteral{
			Value: value.Value,
		}
	case values.StringValue:
		return &ir.StringLiteral{
			Value: value.Value,
		}
	case values.ArrayValue:
		ty := ty.(*types.ArrayType)
		elements := make([]ir.Expression, 0, len(value.Elements))
		for _, elem := range value.Elements {
			elements = append(elements, constValueToExpr(elem, ty.ElemType))
		}
		return &ir.ArrayExpression{
			DataType: ty,
			Elements: elements,
		}

	case values.TupleValue:
		ty := ty.(*types.TupleType)
		values := make([]ir.Expression, 0, len(value.Values))
		for i, value := range value.Values {
			values = append(values, constValueToExpr(value, ty.Types[i]))
		}
		return &ir.TupleExpression{
			Values:   values,
			DataType: ty,
		}
	case values.MapValue:
		ty := ty.(*types.MapType)
		keyValues := make([]ir.KeyValue, 0, len(value.Values))
		for _, kv := range value.Values {
			key := constValueToExpr(kv.Key, ty.KeyType)
			value := constValueToExpr(kv.Key, ty.ValueType)
			keyValues = append(keyValues, ir.KeyValue{
				Key:   key,
				Value: value,
			})
		}
		return &ir.MapExpression{
			KeyValues: keyValues,
			DataType:  ty,
		}

	case values.TypeValue:
		return &ir.TypeExpression{
			DataType: value.Type.(types.Type),
		}

	case values.StructValue:
		ty := ty.(*types.Struct)
		fields := make(map[string]ir.Expression, len(value.Members))
		for name, field := range value.Members {
			fields[name] = constValueToExpr(field, ty.Fields[name].Type)
		}
		return &ir.StructExpression{
			Struct: ty,
			Fields: fields,
		}

	case values.ModuleValue:
		ty := ty.(*types.Module)
		return &ir.VariableExpression{
			Symbol: symbols.Variable{
				Name:       ty.Name,
				IsMut:      false,
				Type:       ty,
				ConstValue: constValue,
			},
		}

	case values.UnitValue:
		return &ir.VariableExpression{
			Symbol: symbols.Variable{
				Name:       value.Name,
				IsMut:      false,
				Type:       ty,
				ConstValue: constValue,
			},
		}
	default:
		panic("Unexpected value")
	}
}
