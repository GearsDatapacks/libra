package lowerer

import "github.com/gearsdatapacks/libra/type_checker/ir"

func (l *lowerer) lowerIntegerLiteral(expr *ir.IntegerLiteral) ir.Expression {
	return expr
}

func (l *lowerer) lowerFloatLiteral(expr *ir.FloatLiteral) ir.Expression {
	return expr
}

func (l *lowerer) lowerBooleanLiteral(expr *ir.BooleanLiteral) ir.Expression {
	return expr
}

func (l *lowerer) lowerStringLiteral(expr *ir.StringLiteral) ir.Expression {
	return expr
}

func (l *lowerer) lowerVariableExpression(expr *ir.VariableExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerBinaryExpression(binExpr *ir.BinaryExpression) ir.Expression {
	left := l.lowerExpression(binExpr.Left)
	right := l.lowerExpression(binExpr.Left)
	if left == binExpr.Left && right == binExpr.Right {
		return binExpr
	}
	return &ir.BinaryExpression{
		Left:     left,
		Operator: binExpr.Operator,
		Right:    right,
	}
}

func (l *lowerer) lowerUnaryExpression(unExpr *ir.UnaryExpression) ir.Expression {
	operand := l.lowerExpression(unExpr.Operand)
	if operand == unExpr.Operand {
		return unExpr
	}
	return &ir.UnaryExpression{
		Operator: unExpr.Operator,
		Operand:  operand,
	}
}

func (l *lowerer) lowerConversion(conversion *ir.Conversion) ir.Expression {
	expr := l.lowerExpression(conversion.Expression)
	if expr == conversion.Expression {
		return conversion
	}
	return &ir.Conversion{
		Expression: expr,
		To:         conversion.To,
	}
}

func (l *lowerer) lowerInvalidExpression(expr *ir.InvalidExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerArrayExpression(array *ir.ArrayExpression) ir.Expression {
	values := make([]ir.Expression, 0, len(array.Elements))
	changed := false
	for i, elem := range array.Elements {
		lowered := l.lowerExpression(elem)
		if !changed && lowered != elem {
			changed = true
			values = append(values, array.Elements[:i-1]...)
		}
		if changed {
			values = append(values, lowered)
		}
	}
	if changed {
		return &ir.ArrayExpression{
			DataType: array.DataType,
			Elements: values,
		}
	}
	return array
}

func (l *lowerer) lowerIndexExpression(i *ir.IndexExpression) ir.Expression {
	left := l.lowerExpression(i.Left)
	index := l.lowerExpression(i.Index)
	if left == i.Left && index == i.Index {
		return i
	}
	return &ir.IndexExpression{
		Left:     left,
		Index:    i,
		DataType: i.DataType,
	}
}

func (l *lowerer) lowerMapExpression(mapExpr *ir.MapExpression) ir.Expression {
	keyValues := make([]ir.KeyValue, 0, len(mapExpr.KeyValues))
	changed := false
	for i, kv := range mapExpr.KeyValues {
		key := l.lowerExpression(kv.Key)
		value := l.lowerExpression(kv.Value)
		if !changed && (key != kv.Key || value != kv.Value) {
			keyValues = append(keyValues, mapExpr.KeyValues[:i-1]...)
			changed = true
		}
		if changed {
			keyValues = append(keyValues, ir.KeyValue{})
		}
	}
	if changed {
		return &ir.MapExpression{
			KeyValues: keyValues,
			DataType:  mapExpr.DataType,
		}
	}
	return mapExpr
}

func (l *lowerer) lowerTupleExpression(tuple *ir.TupleExpression) ir.Expression {
	values := make([]ir.Expression, 0, len(tuple.Values))
	changed := false
	for i, value := range tuple.Values {
		lowered := l.lowerExpression(value)
		if !changed && lowered != value {
			changed = true
			values = append(values, tuple.Values[:i-1]...)
		}
		if changed {
			values = append(values, lowered)
		}
	}
	if changed {
		return &ir.TupleExpression{
			Values:   values,
			DataType: tuple.DataType,
		}
	}
	return tuple
}

func (l *lowerer) lowerAssignment(assignment *ir.Assignment) ir.Expression {
	assignee := l.lowerExpression(assignment.Assignee)
	value := l.lowerExpression(assignment.Value)
	if assignee == assignment.Assignee && value == assignment.Value {
		return assignment
	}
	return &ir.Assignment{
		Assignee: assignee,
		Value:    value,
	}
}

func (l *lowerer) lowerTypeCheck(tc *ir.TypeCheck) ir.Expression {
	value := l.lowerExpression(tc.Value)
	if value == tc.Value {
		return tc
	}
	return &ir.TypeCheck{
		Value:    value,
		DataType: tc.DataType,
	}
}

func (l *lowerer) lowerFunctionCall(call *ir.FunctionCall) ir.Expression {
	args := make([]ir.Expression, 0, len(call.Arguments))
	changed := false
	for i, arg := range call.Arguments {
		lowered := l.lowerExpression(arg)
		if !changed && lowered != arg {
			changed = true
			args = append(args, call.Arguments[:i-1]...)
		}
		if changed {
			args = append(args, lowered)
		}
	}
	function := l.lowerExpression(call.Function)
	if !changed && function == call.Function {
		return call
	}
	if !changed {
		return &ir.FunctionCall{
			Function:   function,
			Arguments:  call.Arguments,
			ReturnType: call.ReturnType,
		}
	}
	return &ir.FunctionCall{
		Function:   function,
		Arguments:  args,
		ReturnType: call.ReturnType,
	}
}

func (l *lowerer) lowerStructExpression(structExpr *ir.StructExpression) ir.Expression {
	fields := make(map[string]ir.Expression, len(structExpr.Fields))
	changed := false
	for name, value := range structExpr.Fields {
		lowered := l.lowerExpression(value)
		if lowered != value {
			changed = true
		}

		fields[name] = lowered
	}
	if !changed {
		return structExpr
	}
	return &ir.StructExpression{
		Struct: structExpr.Struct,
		Fields: fields,
	}
}

func (l *lowerer) lowerTupleStructExpression(tuple *ir.TupleStructExpression) ir.Expression {
	fields := make([]ir.Expression, 0, len(tuple.Fields))
	changed := false
	for i, arg := range tuple.Fields {
		lowered := l.lowerExpression(arg)
		if !changed && lowered != arg {
			changed = true
			fields = append(fields, tuple.Fields[:i-1]...)
		}
		if changed {
			fields = append(fields, lowered)
		}
	}
	if !changed {
		return tuple
	}
	return &ir.TupleStructExpression{
		Struct: tuple.Struct,
		Fields: fields,
	}
}

func (l *lowerer) lowerMemberExpression(member *ir.MemberExpression) ir.Expression {
	left := l.lowerExpression(member.Left)
	if left == member.Left {
		return member
	}
	return &ir.MemberExpression{
		Left:     left,
		Member:   member.Member,
		DataType: member.DataType,
	}
}

func (l *lowerer) lowerBlock(expr *ir.Block) ir.Expression {
	panic("TODO")
}

func (l *lowerer) lowerIfExpression(expr *ir.IfExpression) ir.Expression {
	panic("TODO")
}

func (l *lowerer) lowerWhileLoop(expr *ir.WhileLoop) ir.Expression {
	panic("TODO")
}

func (l *lowerer) lowerForLoop(expr *ir.ForLoop) ir.Expression {
	panic("TODO")
}

func (l *lowerer) lowerTypeExpression(expr *ir.TypeExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerFunctionExpression(expr *ir.FunctionExpression) ir.Expression {
	panic("TODO")
}

func (l *lowerer) lowerRefExpression(ref *ir.RefExpression) ir.Expression {
	value := l.lowerExpression(ref.Value)
	if value == ref.Value {
		return ref
	}
	return &ir.RefExpression{
		Value:   value,
		Mutable: ref.Mutable,
	}
}

func (l *lowerer) lowerDerefExpression(deref *ir.DerefExpression) ir.Expression {
	value := l.lowerExpression(deref.Value)
	if value == deref.Value {
		return deref
	}
	return &ir.DerefExpression{
		Value: value,
	}
}
