package lowerer

import (
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func optimiseExpression(expression ir.Expression) ir.Expression {
	if expression.IsConst() {
		return foldConstants(expression)
	}

	switch expr := expression.(type) {
	case *ir.UnaryExpression:
		switch expr.Operator.Id {
		case ir.LogicalNot:
			return optimiseLogicalNot(expr.Operand, expr)
		case ir.NegateInt, ir.NegateFloat:
			return optimiseNegation(expr.Operand, expr)
		case ir.BitwiseNot:
			return optimiseComplement(expr.Operand, expr)
		}

	case *ir.BinaryExpression:
		switch expr.Operator {
		case ir.LogicalAnd:
			if boolean, ok := expr.Left.ConstValue().(values.BoolValue); ok {
				if boolean.Value {
					return expr.Right
				} else {
					return &ir.BooleanLiteral{Value: false}
				}
			}
			if boolean, ok := expr.Right.ConstValue().(values.BoolValue); ok {
				if boolean.Value {
					return expr.Left
				} else {
					return &ir.BooleanLiteral{Value: false}
				}
			}
		case ir.LogicalOr:
			if boolean, ok := expr.Left.ConstValue().(values.BoolValue); ok {
				if boolean.Value {
					return &ir.BooleanLiteral{Value: true}
				} else {
					return expr.Right
				}
			}
			if boolean, ok := expr.Right.ConstValue().(values.BoolValue); ok {
				if boolean.Value {
					return &ir.BooleanLiteral{Value: true}
				} else {
					return expr.Left
				}
			}
		case ir.LeftShift:
			if int, ok := expr.Left.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 0}
				}
			}
			if int, ok := expr.Right.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 0}
				}
			}

		case ir.RightShift:
			if int, ok := expr.Left.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 0}
				}
			}
			if int, ok := expr.Right.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 0}
				}
			}
		case ir.BitwiseOr:
			if int, ok := expr.Left.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return expr.Right
				}
			}
			if int, ok := expr.Right.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return expr.Left
				}
			}
		case ir.BitwiseAnd:
			if int, ok := expr.Left.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 0}
				}
			}
			if int, ok := expr.Right.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 0}
				}
			}
		case ir.AddInt:
			if int, ok := expr.Left.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return expr.Right
				}
			}
			if int, ok := expr.Right.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return expr.Left
				}
			}
		case ir.AddFloat:
			if float, ok := expr.Left.ConstValue().(values.FloatValue); ok {
				if float.Value == 0 {
					return expr.Right
				}
			}
			if float, ok := expr.Right.ConstValue().(values.FloatValue); ok {
				if float.Value == 0 {
					return expr.Left
				}
			}
		case ir.Concat:
			if str, ok := expr.Left.ConstValue().(values.StringValue); ok {
				if len(str.Value) == 0 {
					return expr.Right
				}
			}
			if str, ok := expr.Right.ConstValue().(values.StringValue); ok {
				if len(str.Value) == 0 {
					return expr.Left
				}
			}
		case ir.SubtractInt:
			if int, ok := expr.Left.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.UnaryExpression{
						Operator: ir.UnaryOperator{Id: ir.NegateInt},
						Operand:  expr.Right,
					}
				}
			}
			if int, ok := expr.Right.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return expr.Left
				}
			}
		case ir.SubtractFloat:
			if float, ok := expr.Left.ConstValue().(values.FloatValue); ok {
				if float.Value == 0 {
					return &ir.UnaryExpression{
						Operator: ir.UnaryOperator{Id: ir.NegateFloat},
						Operand:  expr.Right,
					}
				}
			}
			if float, ok := expr.Right.ConstValue().(values.FloatValue); ok {
				if float.Value == 0 {
					return expr.Left
				}
			}
		case ir.MultiplyInt:
			if int, ok := expr.Left.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 0}
				} else if int.Value == 1 {
					return expr.Right
				}
			}
			if int, ok := expr.Right.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 0}
				} else if int.Value == 1 {
					return expr.Left
				}
			}
		case ir.MultiplyFloat:
			if float, ok := expr.Left.ConstValue().(values.FloatValue); ok {
				if float.Value == 0 {
					return &ir.FloatLiteral{Value: 0}
				} else if float.Value == 1 {
					return expr.Right
				}
			}
			if float, ok := expr.Right.ConstValue().(values.FloatValue); ok {
				if float.Value == 0 {
					return &ir.FloatLiteral{Value: 0}
				} else if float.Value == 1 {
					return expr.Left
				}
			}
		case ir.Divide:
			if float, ok := expr.Right.ConstValue().(values.FloatValue); ok {
				if float.Value == 1 {
					return expr.Left
				}
			}
			if int, ok := expr.Right.ConstValue().(values.IntValue); ok {
				if int.Value == 1 {
					return expr.Left
				}
			}
		case ir.PowerInt:
			if int, ok := expr.Left.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 0}
				} else if int.Value == 1 {
					return &ir.IntegerLiteral{Value: 1}
				}
			}
			if int, ok := expr.Right.ConstValue().(values.IntValue); ok {
				if int.Value == 0 {
					return &ir.IntegerLiteral{Value: 1}
				} else if int.Value == 1 {
					return expr.Left
				}
			}
		case ir.PowerFloat:
			if float, ok := expr.Left.ConstValue().(values.FloatValue); ok {
				if float.Value == 0 {
					return &ir.FloatLiteral{Value: 0}
				} else if float.Value == 1 {
					return &ir.FloatLiteral{Value: 1}
				}
			}
			if float, ok := expr.Right.ConstValue().(values.FloatValue); ok {
				if float.Value == 0 {
					return &ir.FloatLiteral{Value: 1}
				} else if float.Value == 1 {
					return expr.Left
				}
			}
		}
	}

	return expression
}

func optimiseLogicalNot(expression ir.Expression, original ir.Expression) ir.Expression {
	switch expr := expression.(type) {
	case *ir.UnaryExpression:
		switch expr.Operator.Id {
		case ir.LogicalNot:
			return expr.Operand
		}
	case *ir.BinaryExpression:
		switch expr.Operator {
		case ir.Less:
			return &ir.BinaryExpression{
				Left:     expr.Left,
				Operator: ir.GreaterEq,
				Right:    expr.Right,
			}
		case ir.LessEq:
			return &ir.BinaryExpression{
				Left:     expr.Left,
				Operator: ir.Greater,
				Right:    expr.Right,
			}
		case ir.Greater:
			return &ir.BinaryExpression{
				Left:     expr.Left,
				Operator: ir.LessEq,
				Right:    expr.Right,
			}
		case ir.GreaterEq:
			return &ir.BinaryExpression{
				Left:     expr.Left,
				Operator: ir.Less,
				Right:    expr.Right,
			}
		case ir.Equal:
			return &ir.BinaryExpression{
				Left:     expr.Left,
				Operator: ir.NotEqual,
				Right:    expr.Right,
			}
		case ir.NotEqual:
			return &ir.BinaryExpression{
				Left:     expr.Left,
				Operator: ir.Equal,
				Right:    expr.Right,
			}
		}
	}
	return original
}

func optimiseNegation(expression ir.Expression, original ir.Expression) ir.Expression {
	switch expr := expression.(type) {
	case *ir.UnaryExpression:
		switch expr.Operator.Id {
		case ir.NegateInt, ir.NegateFloat:
			return expr.Operand
		}
	}
	return original
}

func optimiseComplement(expression ir.Expression, original ir.Expression) ir.Expression {
	switch expr := expression.(type) {
	case *ir.UnaryExpression:
		switch expr.Operator.Id {
		case ir.BitwiseNot:
			return expr.Operand
		}
	}
	return original
}

func foldConstants(expression ir.Expression) ir.Expression {
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
