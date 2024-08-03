package lowerer

import (
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
)

func (l *lowerer) lowerIntegerLiteral(expr *ir.IntegerLiteral, _ *[]ir.Statement) ir.Expression {
	return expr
}

func (l *lowerer) lowerFloatLiteral(expr *ir.FloatLiteral, _ *[]ir.Statement) ir.Expression {
	return expr
}

func (l *lowerer) lowerBooleanLiteral(expr *ir.BooleanLiteral, _ *[]ir.Statement) ir.Expression {
	return expr
}

func (l *lowerer) lowerStringLiteral(expr *ir.StringLiteral, _ *[]ir.Statement) ir.Expression {
	return expr
}

func (l *lowerer) lowerVariableExpression(expr *ir.VariableExpression, _ *[]ir.Statement) ir.Expression {
	return expr
}

func (l *lowerer) lowerBinaryExpression(binExpr *ir.BinaryExpression, statements *[]ir.Statement) ir.Expression {
	left := l.lowerExpression(binExpr.Left, statements)
	right := l.lowerExpression(binExpr.Right, statements)
	if left == binExpr.Left && right == binExpr.Right {
		return binExpr
	}
	return &ir.BinaryExpression{
		Left:     left,
		Operator: binExpr.Operator,
		Right:    right,
	}
}

func (l *lowerer) lowerUnaryExpression(unExpr *ir.UnaryExpression, statements *[]ir.Statement) ir.Expression {
	operand := l.lowerExpression(unExpr.Operand, statements)
	if operand == unExpr.Operand {
		return unExpr
	}
	return &ir.UnaryExpression{
		Operator: unExpr.Operator,
		Operand:  operand,
	}
}

func (l *lowerer) lowerConversion(conversion *ir.Conversion, statements *[]ir.Statement) ir.Expression {
	expr := l.lowerExpression(conversion.Expression, statements)
	if expr == conversion.Expression {
		return conversion
	}
	return &ir.Conversion{
		Expression: expr,
		To:         conversion.To,
	}
}

func (l *lowerer) lowerInvalidExpression(expr *ir.InvalidExpression, _ *[]ir.Statement) ir.Expression {
	return expr
}

func (l *lowerer) lowerArrayExpression(array *ir.ArrayExpression, statements *[]ir.Statement) ir.Expression {
	values := make([]ir.Expression, 0, len(array.Elements))
	changed := false
	for i, elem := range array.Elements {
		lowered := l.lowerExpression(elem, statements)
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

func (l *lowerer) lowerIndexExpression(i *ir.IndexExpression, statements *[]ir.Statement) ir.Expression {
	left := l.lowerExpression(i.Left, statements)
	index := l.lowerExpression(i.Index, statements)
	if left == i.Left && index == i.Index {
		return i
	}
	return &ir.IndexExpression{
		Left:     left,
		Index:    i,
		DataType: i.DataType,
	}
}

func (l *lowerer) lowerMapExpression(mapExpr *ir.MapExpression, statements *[]ir.Statement) ir.Expression {
	keyValues := make([]ir.KeyValue, 0, len(mapExpr.KeyValues))
	changed := false
	for i, kv := range mapExpr.KeyValues {
		key := l.lowerExpression(kv.Key, statements)
		value := l.lowerExpression(kv.Value, statements)
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

func (l *lowerer) lowerTupleExpression(tuple *ir.TupleExpression, statements *[]ir.Statement) ir.Expression {
	values := make([]ir.Expression, 0, len(tuple.Values))
	changed := false
	for i, value := range tuple.Values {
		lowered := l.lowerExpression(value, statements)
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

func (l *lowerer) lowerAssignment(assignment *ir.Assignment, statements *[]ir.Statement) ir.Expression {
	assignee := l.lowerExpression(assignment.Assignee, statements)
	value := l.lowerExpression(assignment.Value, statements)
	if assignee == assignment.Assignee && value == assignment.Value {
		return assignment
	}
	return &ir.Assignment{
		Assignee: assignee,
		Value:    value,
	}
}

func (l *lowerer) lowerTypeCheck(tc *ir.TypeCheck, statements *[]ir.Statement) ir.Expression {
	value := l.lowerExpression(tc.Value, statements)
	if value == tc.Value {
		return tc
	}
	return &ir.TypeCheck{
		Value:    value,
		DataType: tc.DataType,
	}
}

func (l *lowerer) lowerFunctionCall(call *ir.FunctionCall, statements *[]ir.Statement) ir.Expression {
	args := make([]ir.Expression, 0, len(call.Arguments))
	changed := false
	for i, arg := range call.Arguments {
		lowered := l.lowerExpression(arg, statements)
		if !changed && lowered != arg {
			changed = true
			args = append(args, call.Arguments[:i-1]...)
		}
		if changed {
			args = append(args, lowered)
		}
	}
	function := l.lowerExpression(call.Function, statements)
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

func (l *lowerer) lowerStructExpression(structExpr *ir.StructExpression, statements *[]ir.Statement) ir.Expression {
	fields := make(map[string]ir.Expression, len(structExpr.Fields))
	changed := false
	for name, value := range structExpr.Fields {
		lowered := l.lowerExpression(value, statements)
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

func (l *lowerer) lowerTupleStructExpression(tuple *ir.TupleStructExpression, statements *[]ir.Statement) ir.Expression {
	fields := make([]ir.Expression, 0, len(tuple.Fields))
	changed := false
	for i, arg := range tuple.Fields {
		lowered := l.lowerExpression(arg, statements)
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

func (l *lowerer) lowerMemberExpression(member *ir.MemberExpression, statements *[]ir.Statement) ir.Expression {
	left := l.lowerExpression(member.Left, statements)
	if left == member.Left {
		return member
	}
	return &ir.MemberExpression{
		Left:     left,
		Member:   member.Member,
		DataType: member.DataType,
	}
}

func (l *lowerer) lowerBlock(block *ir.Block, statements *[]ir.Statement, endLabel ...string) ir.Expression {
	if len(block.Statements) == 1 {
		if expr, ok := block.Statements[0].(ir.Expression); ok {
			return l.lowerExpression(expr, statements)
		}
	}

	yieldLabel := ""
	addLabel := len(endLabel) == 0
	if len(endLabel) > 0 {
		yieldLabel = endLabel[0]
	} else {
		yieldLabel = l.genLabel()
	}

	yieldVar := symbols.Variable{
		Name:       l.genVar(),
		IsMut:      true,
		Type:       block.ResultType,
		ConstValue: nil,
	}
	defer l.endScope(l.beginScope(blockContext{endLabel: yieldLabel, yieldVariable: yieldVar}))

	for _, stmt := range block.Statements {
		l.lower(stmt, statements)
	}

	if addLabel {
		*statements = append(*statements, &ir.Label{Name: yieldLabel})
	}
	return &ir.VariableExpression{Symbol: yieldVar}
}

func negate(condition ir.Expression) ir.Expression {
	return optimiseExpression(&ir.UnaryExpression{
		Operator: ir.UnaryOperator{Id: ir.LogicalNot},
		Operand:  condition,
	})
}

type ifContext struct {
	finalLabel     string
	returnVariable symbols.Variable
}

func (l *lowerer) lowerIfExpression(ifExpr *ir.IfExpression, statements *[]ir.Statement, context *ifContext) ir.Expression {
	isFirst := context == nil
	if context == nil {
		context = &ifContext{
			finalLabel: l.genLabel(),
			returnVariable: symbols.Variable{
				Name:       l.genVar(),
				IsMut:      true,
				Type:       ifExpr.Type(),
				ConstValue: nil,
			},
		}
	}

	endLabel := ""
	if ifExpr.ElseBranch == nil {
		endLabel = context.finalLabel
	} else {
		endLabel = l.genLabel()
	}

	condition := l.lowerExpression(negate(ifExpr.Condition), statements)
	*statements = append(*statements, &ir.GotoIf{
		Label:     endLabel,
		Condition: condition,
	})

	value := l.lowerBlock(ifExpr.Body, statements, context.finalLabel)
	*statements = append(*statements, &ir.Assignment{
		Assignee: &ir.VariableExpression{Symbol: context.returnVariable},
		Value:    value,
	})

	if ifExpr.ElseBranch != nil {
		*statements = append(*statements, &ir.Goto{Label: context.finalLabel})
		*statements = append(*statements, &ir.Label{Name: endLabel})

		switch eb := ifExpr.ElseBranch.(type) {
		case *ir.IfExpression:
			l.lowerIfExpression(eb, statements, context)
		case *ir.Block:
			value := l.lowerBlock(eb, statements)
			*statements = append(*statements, &ir.Assignment{
				Assignee: &ir.VariableExpression{Symbol: context.returnVariable},
				Value:    value,
			})
			*statements = append(*statements, &ir.Label{Name: context.finalLabel})
		}
	} else if !isFirst {
		*statements = append(*statements, &ir.Label{Name: context.finalLabel})
	} else {
		*statements = append(*statements, &ir.Label{Name: endLabel})
	}

	return &ir.VariableExpression{Symbol: context.returnVariable}
}

func (l *lowerer) lowerWhileLoop(loop *ir.WhileLoop, statements *[]ir.Statement) ir.Expression {
	loopStart := l.genLabel()
	*statements = append(*statements, &ir.Label{Name: loopStart})

	condition := l.lowerExpression(negate(loop.Condition), statements)
	loopEnd := l.genLabel()
	*statements = append(*statements, &ir.GotoIf{
		Condition: condition,
		Label:     loopEnd,
	})

	breakVariable := symbols.Variable{
		Name:       l.genVar(),
		IsMut:      true,
		Type:       loop.Type(),
		ConstValue: nil,
	}

	defer l.endScope(l.beginScope(loopContext{
		breakLabel:    loopEnd,
		continueLabel: loopStart,
		breakVariable: breakVariable,
	}))

	for _, stmt := range loop.Body.Statements {
		l.lower(stmt, statements)
	}
	*statements = append(*statements, &ir.Goto{Label: loopStart})
	*statements = append(*statements, &ir.Label{Name: loopEnd})

	return &ir.VariableExpression{Symbol: breakVariable}
}

func (l *lowerer) lowerForLoop(expr *ir.ForLoop, statements *[]ir.Statement) ir.Expression {
	panic("TODO")
}

func (l *lowerer) lowerTypeExpression(expr *ir.TypeExpression, _ *[]ir.Statement) ir.Expression {
	return expr
}

func (l *lowerer) lowerFunctionExpression(funcExpr *ir.FunctionExpression, _ *[]ir.Statement) ir.Expression {
	statements := []ir.Statement{}
	l.lowerBlock(funcExpr.Body, &statements)
	statements = l.cfa(statements, funcExpr.Location)

	return &ir.FunctionExpression{
		Parameters: funcExpr.Parameters,
		Body:       &ir.Block{Statements: statements, ResultType: funcExpr.Body.ResultType},
		DataType:   funcExpr.DataType,
		Location:   funcExpr.Location,
	}
}

func (l *lowerer) lowerRefExpression(ref *ir.RefExpression, statements *[]ir.Statement) ir.Expression {
	value := l.lowerExpression(ref.Value, statements)
	if value == ref.Value {
		return ref
	}
	return &ir.RefExpression{
		Value:   value,
		Mutable: ref.Mutable,
	}
}

func (l *lowerer) lowerDerefExpression(deref *ir.DerefExpression, statements *[]ir.Statement) ir.Expression {
	value := l.lowerExpression(deref.Value, statements)
	if value == deref.Value {
		return deref
	}
	return &ir.DerefExpression{
		Value: value,
	}
}
