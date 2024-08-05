package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (t *typeChecker) typeCheckExpression(expression ast.Expression) ir.Expression {
	expr := t.doTypeCheckExpression(expression)
	if varExpr, ok := expr.(*ir.VariableExpression); ok && varExpr.Symbol.Type == types.RuntimeType {
		if unit, ok := varExpr.ConstValue().(values.TypeValue).Type.(*types.UnitStruct); ok {
			expr = &ir.VariableExpression{
				Symbol: symbols.Variable{
					Name:  varExpr.Symbol.Name,
					IsMut: false,
					Type:  unit,
					ConstValue: &values.UnitValue{
						Name: unit.Name,
					},
				},
			}
		}
	}

	return expr
}

func (t *typeChecker) doTypeCheckExpression(expression ast.Expression) ir.Expression {
	switch expr := expression.(type) {
	case *ast.IntegerLiteral:
		return &ir.IntegerLiteral{
			Value: expr.Value,
		}
	case *ast.FloatLiteral:
		return &ir.FloatLiteral{
			Value: expr.Value,
		}
	case *ast.BooleanLiteral:
		return &ir.BooleanLiteral{
			Value: expr.Value,
		}
	case *ast.StringLiteral:
		return &ir.StringLiteral{
			Value: expr.Value,
		}
	case *ast.ListLiteral:
		return t.typeCheckArray(expr)
	case *ast.MapLiteral:
		return t.typeCheckMap(expr)
	case *ast.TupleExpression:
		return t.typeCheckTuple(expr)

	case *ast.Identifier:
		return t.lookupVariable(expr.Name, expr.Location)
	case *ast.BinaryExpression:
		return t.typeCheckBinaryExpression(expr)
	case *ast.ParenthesisedExpression:
		return t.typeCheckExpression(expr.Expression)
	case *ast.PrefixExpression:
		return t.typeCheckPrefixExpression(expr)
	case *ast.PostfixExpression:
		return t.typeCheckPostfixExpression(expr)
	case *ast.CastExpression:
		return t.typeCheckCastExpression(expr)
	case *ast.IndexExpression:
		return t.typeCheckIndexExpression(expr)
	case *ast.AssignmentExpression:
		return t.typeCheckAssignment(expr)
	case *ast.TypeCheckExpression:
		return t.typeCheckTypeCheck(expr)
	case *ast.FunctionCall:
		return t.typeCheckFunctionCall(expr)
	case *ast.StructExpression:
		return t.typeCheckStructExpression(expr)
	case *ast.MemberExpression:
		return t.typeCheckMemberExpression(expr)
	case *ast.RefExpression:
		return t.typeCheckRefExpression(expr)
	case *ast.DerefExpression:
		return t.typeCheckDerefExpression(expr)
	case *ast.PointerType:
		return t.typeCheckPointerType(expr)
	case *ast.OptionType:
		return t.typeCheckOptionType(expr)

	case *ast.Block:
		return t.typeCheckBlock(expr, true)
	case *ast.IfExpression:
		return t.typeCheckIfExpression(expr)
	case *ast.WhileLoop:
		return t.typeCheckWhileLoop(expr)
	case *ast.ForLoop:
		return t.typeCheckForLoop(expr)
	case *ast.FunctionExpression:
		return t.typeCheckFunctionExpression(expr)

	default:
		panic(fmt.Sprintf("TODO: Type-check %T", expr))
	}
}

func (t *typeChecker) lookupVariable(name string, location text.Location) ir.Expression {
	symbol := t.symbols.Lookup(name)
	if symbol == nil {
		t.diagnostics.Report(diagnostics.VariableUndefined(location, name))
		symbol = &symbols.Variable{
			Name:  name,
			IsMut: true,
			Type:  types.Invalid,
		}
	}
	return &ir.VariableExpression{
		Symbol: symbols.Variable{
			Name:       symbol.GetName(),
			IsMut:      symbol.Mutable(),
			Type:       symbol.GetType(),
			ConstValue: symbol.Value(),
		},
	}
}

func (t *typeChecker) typeCheckBinaryExpression(binExpr *ast.BinaryExpression) ir.Expression {
	left := t.typeCheckExpression(binExpr.Left)
	right := t.typeCheckExpression(binExpr.Right)

	// Don't even check for operators with invalid types, to prevent cascading errors
	if left.Type() == types.Invalid || right.Type() == types.Invalid {
		return &ir.BinaryExpression{
			Left:     left,
			Operator: 0,
			Right:    right,
		}
	}

	left, right, operator := getBinaryOperator(binExpr.Operator.Kind, left, right)

	if operator == 0 {
		t.diagnostics.Report(diagnostics.BinaryOperatorUndefined(binExpr.Operator.Location, binExpr.Operator.Value, left.Type(), right.Type()))
	}

	return &ir.BinaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func getBinaryOperator(op token.Kind, left, right ir.Expression) (lhs, rhs ir.Expression, binOp ir.BinaryOperator) {
	lhs = left
	rhs = right
	binOp = 0
	lType := left.Type()
	rType := right.Type()

	leftNum, leftNumeric := lType.(types.Numeric)
	rightNum, rightNumeric := rType.(types.Numeric)
	isFloat :=
		(leftNumeric && leftNum.Kind == types.NumFloat) ||
			(rightNumeric && rightNum.Kind == types.NumFloat)

	lUntyped := leftNumeric && leftNum.Untyped()
	rUntyped := rightNumeric && rightNum.Untyped()
	untyped := lUntyped && rUntyped

	switch op {
	case token.DOUBLE_AMPERSAND:
		if types.Assignable(types.Bool, lType) && types.Assignable(types.Bool, rType) {
			binOp = ir.LogicalAnd
		}
	case token.DOUBLE_PIPE:
		if types.Assignable(types.Bool, lType) && types.Assignable(types.Bool, rType) {
			binOp = ir.LogicalOr
		}
	case token.LEFT_ANGLE:
		if leftNumeric && rightNumeric {
			binOp = ir.Less
			if isFloat {
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			}
		}
	case token.RIGHT_ANGLE:
		if leftNumeric && rightNumeric {
			binOp = ir.Greater
			if isFloat {
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			}
		}
	case token.LEFT_ANGLE_EQUALS:
		if leftNumeric && rightNumeric {
			binOp = ir.LessEq
			if isFloat {
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			}
		}
	case token.RIGHT_ANGLE_EQUALS:
		if leftNumeric && rightNumeric {
			binOp = ir.GreaterEq
			if isFloat {
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			}
		}
	case token.DOUBLE_EQUALS:
		if types.Assignable(rType, lType) || types.Assignable(lType, rType) {
			binOp = ir.Equal
		}
	case token.BANG_EQUALS:
		if types.Assignable(rType, lType) || types.Assignable(lType, rType) {
			binOp = ir.NotEqual
		}
	case token.DOUBLE_LEFT_ANGLE:
		if types.Assignable(types.I32, lType) && types.Assignable(types.I32, rType) {
			binOp = ir.LeftShift
		}
	case token.DOUBLE_RIGHT_ANGLE:
		if types.Assignable(types.I32, lType) && types.Assignable(types.I32, rType) {
			binOp = ir.RightShift
		}
	case token.PLUS:
		if types.Assignable(types.String, lType) && types.Assignable(types.String, rType) {
			binOp = ir.Concat
		}

		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.AddFloat
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			} else {
				binOp = ir.AddInt
			}
		}
	case token.MINUS:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.SubtractFloat
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			} else {
				binOp = ir.SubtractInt
			}
		}
	case token.STAR:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.MultiplyFloat
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			} else {
				binOp = ir.MultiplyInt
			}
		}
	case token.SLASH:
		if leftNumeric && rightNumeric {
			binOp = ir.Divide
			if isFloat {
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			}
		}
	case token.PERCENT:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.ModuloFloat
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			} else {
				binOp = ir.ModuloInt
			}
		}
	case token.DOUBLE_STAR:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.PowerFloat
				lhs = convert(lhs, types.F32, operator)
				rhs = convert(rhs, types.F32, operator)
			} else {
				binOp = ir.PowerInt
			}
		}
	case token.PIPE:
		if types.Assignable(types.I32, lType) && types.Assignable(types.I32, rType) {
			binOp = ir.BitwiseOr
		}
		if types.Assignable(types.RuntimeType, lType) && types.Assignable(types.RuntimeType, rType) {
			binOp = ir.Union
		}
	case token.AMPERSAND:
		if types.Assignable(types.I32, lType) && types.Assignable(types.I32, rType) {
			binOp = ir.BitwiseAnd
		}
	}

	if untyped && binOp != 0 {
		binOp = binOp | ir.UntypedBit
	}
	return
}

func (t *typeChecker) typeCheckPrefixExpression(unExpr *ast.PrefixExpression) ir.Expression {
	operand := t.typeCheckExpression(unExpr.Operand)

	// Don't check for operators with invalid types, to prevent cascading errors
	if operand.Type() == types.Invalid {
		return &ir.UnaryExpression{
			Operand:  operand,
			Operator: ir.UnaryOperator{},
		}
	}

	if unExpr.Operator == token.BANG && operand.Type() == types.RuntimeType {
		return &ir.TypeExpression{
			DataType: &types.Result{
				OkType: t.typeFromExpr(operand, unExpr.Operand.GetLocation()),
			},
		}
	}

	operator := getPrefixOperator(unExpr.Operator, operand)

	if operator.Id == 0 {
		t.diagnostics.Report(diagnostics.UnaryOperatorUndefined(unExpr.Location, unExpr.Operator.String(), operand.Type()))
	}

	// We can safely ignore the identity operator
	if operator.Id & ^ir.UntypedBit == ir.Identity {
		return operand
	}

	return &ir.UnaryExpression{
		Operator: operator,
		Operand:  operand,
	}
}

func (t *typeChecker) typeCheckPostfixExpression(unExpr *ast.PostfixExpression) ir.Expression {
	operand := t.typeCheckExpression(unExpr.Operand)

	// Don't check for operators with invalid types, to prevent cascading errors
	if operand.Type() == types.Invalid {
		return &ir.UnaryExpression{
			Operand:  operand,
			Operator: ir.UnaryOperator{},
		}
	}

	operator, diag := t.getPostfixOperator(unExpr.Operator, operand)

	if diag != nil {
		t.diagnostics.Report(diag.Location(unExpr.Operand.GetLocation()))
	} else if operator.Id == 0 {
		t.diagnostics.Report(diagnostics.UnaryOperatorUndefined(unExpr.OperatorLocation, unExpr.Operator.String(), operand.Type()))
	}

	return &ir.UnaryExpression{
		Operator: operator,
		Operand:  operand,
	}
}

func getPrefixOperator(tokKind token.Kind, operand ir.Expression) ir.UnaryOperator {
	var id ir.UnOpId
	opType := operand.Type()

	num, numeric := opType.(types.Numeric)
	isFloat := numeric && num.Kind == types.NumFloat
	var untyped bool
	if v, ok := opType.(types.Numeric); ok {
		untyped = v.Untyped()
	}

	switch tokKind {
	case token.MINUS:
		if numeric {
			if isFloat {
				id = ir.NegateFloat
			} else {
				id = ir.NegateInt
			}
		}
	case token.PLUS:
		if numeric {
			id = ir.Identity
		}
	case token.BANG:
		if types.Assignable(types.Bool, opType) {
			id = ir.LogicalNot
		}
	case token.TILDE:
		if types.Assignable(types.I32, opType) {
			id = ir.BitwiseNot
		}
	}

	if untyped && id != 0 {
		id = id | ir.UntypedBit
	}

	return ir.UnaryOperator{Id: id}
}

func (t *typeChecker) getPostfixOperator(tokKind token.Kind, operand ir.Expression) (ir.UnaryOperator, *diagnostics.Partial) {
	var id ir.UnOpId
	var ty types.Type
	opType := operand.Type()

	num, numeric := opType.(types.Numeric)
	isFloat := numeric && num.Kind == types.NumFloat
	var untyped bool
	if v, ok := opType.(types.Numeric); ok {
		untyped = v.Untyped()
	}

	switch tokKind {
	case token.DOUBLE_PLUS:
		if numeric {
			if !ir.AssignableExpr(operand) {
				return ir.UnaryOperator{}, diagnostics.CannotIncDec("increment")
			} else if !ir.MutableExpr(operand) {
				return ir.UnaryOperator{}, diagnostics.ValueImmutablePartial
			}
			if isFloat {
				id = ir.IncrementFloat
			} else {
				id = ir.IncrecementInt
			}
		}
	case token.DOUBLE_MINUS:
		if numeric {
			if !ir.AssignableExpr(operand) {
				return ir.UnaryOperator{}, diagnostics.CannotIncDec("decrement")
			} else if !ir.MutableExpr(operand) {
				return ir.UnaryOperator{}, diagnostics.ValueImmutablePartial
			}
			if isFloat {
				id = ir.DecrementFloat
			} else {
				id = ir.DecrecementInt
			}
		}
	case token.QUESTION:
		var expectedType types.Type = nil
		symbolTable := t.symbols
		for symbolTable != nil {
			if fnContext, ok := symbolTable.Context.(symbols.FunctionContext); ok {
				expectedType = fnContext.ReturnType
				break
			}
			symbolTable = symbolTable.Parent
		}

		if expectedType == nil {
			return ir.UnaryOperator{}, diagnostics.NoPropagateOutsideFunction()
		}

		if result, ok := operand.Type().(*types.Result); ok {
			if _, ok := expectedType.(*types.Result); !ok {
				return ir.UnaryOperator{}, diagnostics.PropagateFnMustReturnResult()
			}

			id = ir.PropagateError
			ty = result.OkType
		} else if option, ok := operand.Type().(*types.Option); ok {
			if _, ok := expectedType.(*types.Option); !ok {
				return ir.UnaryOperator{}, diagnostics.PropagateFnMustReturnOption()
			}

			id = ir.PropagateError
			ty = option.SomeType
		}

	case token.BANG:
		if result, ok := operand.Type().(*types.Result); ok {
			id = ir.CrashError
			ty = result.OkType
		} else if option, ok := operand.Type().(*types.Option); ok {
			id = ir.CrashError
			ty = option.SomeType
		}
	}

	if untyped && id != 0 {
		id = id | ir.UntypedBit
	}
	return ir.UnaryOperator{Id: id, DataType: ty}, nil
}

func (t *typeChecker) typeCheckCastExpression(expr *ast.CastExpression) ir.Expression {
	value := t.typeCheckExpression(expr.Left)
	ty := t.typeCheckType(expr.Type)
	conversion := convert(value, ty, explicit)
	if conversion == nil {
		t.diagnostics.Report(diagnostics.CannotCast(expr.Left.GetLocation(), value.Type(), ty))
		return &ir.InvalidExpression{
			Expression: value,
		}
	}
	return conversion
}

func (t *typeChecker) typeCheckArray(arr *ast.ListLiteral) ir.Expression {
	var elemType types.Type = types.Invalid
	values := []ir.Expression{}

	for _, elem := range arr.Values {
		value := t.typeCheckExpression(elem)
		if elemType == types.Invalid {
			elemType = types.ToReal(value.Type())
		}
		converted := convert(value, elemType, operator)
		if converted == nil {
			t.diagnostics.Report(diagnostics.NotAssignable(elem.GetLocation(), elemType, value.Type()))
		} else {
			values = append(values, converted)
		}
	}

	return &ir.ArrayExpression{
		DataType: &types.ArrayType{
			ElemType: elemType,
			Length:   len(values),
			CanInfer: true,
		},
		Elements: values,
	}
}

func (t *typeChecker) typeCheckIndexExpression(indexExpr *ast.IndexExpression) ir.Expression {
	left := t.typeCheckExpression(indexExpr.Left)

	if left.Type() == types.RuntimeType {
		elemType := t.typeFromExpr(left, indexExpr.Left.GetLocation())

		if indexExpr.Index == nil {
			return &ir.TypeExpression{DataType: &types.ListType{ElemType: elemType}}
		}

		length := -1

		if ident, ok := indexExpr.Index.(*ast.Identifier); !ok || ident.Name != "_" {
			expr := convert(t.typeCheckExpression(indexExpr.Index), types.I32, implicit)
			if expr == nil {
				t.diagnostics.Report(diagnostics.CountMustBeInt(indexExpr.Index.GetLocation()))
			} else if expr.IsConst() {
				value := expr.ConstValue().(values.IntValue)
				length = int(value.Value)
			} else {
				t.diagnostics.Report(diagnostics.NotConst(indexExpr.Index.GetLocation()))
			}
		}

		return &ir.TypeExpression{DataType: &types.ArrayType{
			ElemType: elemType,
			Length:   length,
			CanInfer: false,
		}}
	}

	if indexExpr.Index == nil {
		t.diagnostics.Report(diagnostics.ExpressionIndexWithoutIndex(indexExpr.Location))
		return left
	}

	index := t.typeCheckExpression(indexExpr.Index)
	ty, diag := ir.Index(left, index)

	expr := &ir.IndexExpression{
		Left:     left,
		Index:    index,
		DataType: ty,
	}
	if diag != nil {
		t.diagnostics.Report(diag.Location(indexExpr.Index.GetLocation()))
		return &ir.InvalidExpression{Expression: expr}
	}

	return expr
}

func (t *typeChecker) typeCheckMap(mapLit *ast.MapLiteral) ir.Expression {
	var keyType types.Type = types.Invalid
	var valueType types.Type = types.Invalid
	keyValues := []ir.KeyValue{}

	for _, kv := range mapLit.KeyValues {
		key := t.typeCheckExpression(kv.Key)
		if keyType == types.Invalid {
			keyType = types.ToReal(key.Type())
		}
		convertedKey := convert(key, keyType, operator)
		if convertedKey == nil {
			t.diagnostics.Report(diagnostics.NotAssignable(kv.Key.GetLocation(), keyType, key.Type()))
			continue
		}

		value := t.typeCheckExpression(kv.Value)
		if valueType == types.Invalid {
			valueType = types.ToReal(value.Type())
		}
		convertedValue := convert(value, valueType, operator)
		if convertedValue == nil {
			t.diagnostics.Report(diagnostics.NotAssignable(kv.Value.GetLocation(), valueType, value.Type()))
			continue
		}

		keyValues = append(keyValues, ir.KeyValue{
			Key:   key,
			Value: value,
		})
	}

	if len(keyValues) == 1 && keyType == types.RuntimeType && valueType == types.RuntimeType {
		return &ir.TypeExpression{
			DataType: &types.MapType{
				KeyType:   t.typeFromExpr(keyValues[0].Key, mapLit.KeyValues[0].Key.GetLocation()),
				ValueType: t.typeFromExpr(keyValues[0].Value, mapLit.KeyValues[0].Value.GetLocation()),
			},
		}
	}

	mapExpr := &ir.MapExpression{
		KeyValues: keyValues,
		DataType: &types.MapType{
			KeyType:   keyType,
			ValueType: valueType,
		},
	}

	if !types.Hashable(keyType) {
		t.diagnostics.Report(diagnostics.NotHashable(mapLit.KeyValues[0].Key.GetLocation(), keyType))
		return &ir.InvalidExpression{
			Expression: mapExpr,
		}
	}

	return mapExpr
}

func (t *typeChecker) typeCheckAssignment(assignment *ast.AssignmentExpression) ir.Expression {
	assignee := t.typeCheckExpression(assignment.Assignee)
	value := t.typeCheckExpression(assignment.Value)

	if assignment.Operator.Kind != token.EQUALS {
		left, right, operator := getBinaryOperator(assignment.Operator.Kind-token.EQUALS, assignee, value)

		if operator == 0 {
			t.diagnostics.Report(diagnostics.BinaryOperatorUndefined(assignment.Operator.Location, assignment.Operator.Value, left.Type(), right.Type()))
		}

		value = &ir.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	if !ir.AssignableExpr(assignee) {
		t.diagnostics.Report(diagnostics.CannotAssign(assignment.Assignee.GetLocation()))
	} else if !ir.MutableExpr(assignee) {
		t.diagnostics.Report(diagnostics.ValueImmutable(assignment.Assignee.GetLocation()))
	} else {
		conversion := convert(value, assignee.Type(), implicit)
		if conversion == nil {
			t.diagnostics.Report(diagnostics.NotAssignable(assignment.Assignee.GetLocation(), assignee.Type(), value.Type()))
		} else {
			value = conversion
		}
	}

	return &ir.Assignment{
		Assignee: assignee,
		Value:    value,
	}
}

func (t *typeChecker) typeCheckTuple(tuple *ast.TupleExpression) ir.Expression {
	dataTypes := []types.Type{}
	values := []ir.Expression{}
	isType := true

	for _, value := range tuple.Values {
		expr := t.typeCheckExpression(value)
		ty := types.ToReal(expr.Type())
		dataTypes = append(dataTypes, ty)
		if ty != types.RuntimeType {
			isType = false
		}
		values = append(values, convert(expr, ty, implicit))
	}

	if isType && len(dataTypes) != 0 {
		tupleMembers := []types.Type{}
		for i, val := range values {
			tupleMembers = append(tupleMembers, t.typeFromExpr(val, tuple.Values[i].GetLocation()))
		}
		return &ir.TypeExpression{
			DataType: &types.TupleType{
				Types: tupleMembers,
			},
		}
	}

	return &ir.TupleExpression{
		Values: values,
		DataType: &types.TupleType{
			Types: dataTypes,
		},
	}
}

func (t *typeChecker) typeCheckTypeCheck(tc *ast.TypeCheckExpression) ir.Expression {
	value := t.typeCheckExpression(tc.Left)
	ty := t.typeCheckType(tc.Type)

	return &ir.TypeCheck{
		Value:    value,
		DataType: ty,
	}
}

func (t *typeChecker) typeCheckFunctionCall(call *ast.FunctionCall) ir.Expression {
	fn := t.typeCheckExpression(call.Callee)
	funcType, ok := types.Unwrap(fn.Type()).(*types.Function)
	if !ok {
		t.diagnostics.Report(diagnostics.NotCallable(call.Callee.GetLocation(), fn.Type()))
		return &ir.InvalidExpression{
			Expression: &ir.FunctionCall{
				Function:   fn,
				Arguments:  []ir.Expression{},
				ReturnType: types.Invalid,
			},
		}
	}

	if len(call.Arguments) != len(funcType.Parameters) {
		t.diagnostics.Report(diagnostics.WrongNumberAgruments(call.Callee.GetLocation(), len(funcType.Parameters), len(call.Arguments)))
		return &ir.InvalidExpression{
			Expression: &ir.FunctionCall{
				Function:   fn,
				Arguments:  []ir.Expression{},
				ReturnType: types.Invalid,
			},
		}
	}

	args := []ir.Expression{}
	for i, arg := range call.Arguments {
		value := t.typeCheckExpression(arg)
		expectedType := funcType.Parameters[i]
		conversion := convert(value, expectedType, implicit)
		if conversion == nil {
			args = append(args, value)
			t.diagnostics.Report(diagnostics.NotAssignable(arg.GetLocation(), expectedType, value.Type()))
		} else {
			args = append(args, conversion)
		}
	}

	return &ir.FunctionCall{
		Function:   fn,
		Arguments:  args,
		ReturnType: funcType.ReturnType,
	}
}

func (t *typeChecker) typeCheckStructExpression(structExpr *ast.StructExpression) ir.Expression {
	baseTy := t.typeCheckType(structExpr.Struct)
	ty := types.Unwrap(baseTy)
	if structTy, ok := ty.(*types.Struct); ok {
		fields := map[string]ir.Expression{}

		for _, member := range structExpr.Members {
			if member.Name == nil {
				t.diagnostics.Report(diagnostics.NoNameStructMember(member.Value.GetLocation()))
				continue
			}
			field, ok := structTy.Fields[*member.Name]
			if !ok {
				t.diagnostics.Report(diagnostics.NoStructMember(member.Location, structTy.Name, *member.Name))
				continue
			}
			var value ir.Expression
			if member.Value != nil {
				value = t.typeCheckExpression(member.Value)
			} else {
				value = t.lookupVariable(*member.Name, member.Location)
			}

			conversion := convert(value, field.Type, implicit)
			if conversion != nil {
				value = conversion
			} else {
				t.diagnostics.Report(diagnostics.NotAssignable(member.Value.GetLocation(), field.Type, value.Type()))
			}

			fields[*member.Name] = value
		}

		return &ir.StructExpression{
			Struct: baseTy,
			Fields: fields,
		}
	} else if tupleTy, ok := ty.(*types.TupleStruct); ok {
		fields := []ir.Expression{}

		if len(structExpr.Members) != len(tupleTy.Types) {
			t.diagnostics.Report(diagnostics.WrongNumberTupleValues(structExpr.Struct.GetLocation(), len(tupleTy.Types), len(structExpr.Members)))
			return &ir.InvalidExpression{
				Expression: &ir.TupleStructExpression{
					Struct: tupleTy,
					Fields: fields,
				},
			}
		}

		for i, member := range structExpr.Members {
			field := tupleTy.Types[i]
			var value ir.Expression
			if member.Value != nil {
				value = t.typeCheckExpression(member.Value)
				if member.Name != nil {
					t.diagnostics.Report(diagnostics.TupleStructWithNames(member.Location))
				}
			} else {
				value = t.lookupVariable(*member.Name, member.Location)
			}

			conversion := convert(value, field, implicit)
			if conversion != nil {
				value = conversion
			} else {
				t.diagnostics.Report(diagnostics.NotAssignable(member.Location, field, value.Type()))
			}

			fields = append(fields, value)
		}

		return &ir.TupleStructExpression{
			Struct: baseTy,
			Fields: fields,
		}
	} else {
		t.diagnostics.Report(diagnostics.CannotConstruct(structExpr.Struct.GetLocation(), ty))
		return &ir.InvalidExpression{
			Expression: &ir.IntegerLiteral{},
		}
	}
}

func (t *typeChecker) typeCheckMemberExpression(member *ast.MemberExpression) ir.Expression {
	left := t.typeCheckExpression(member.Left)
	ty, diag := ir.Member(left, member.Member)

	memberExpr := &ir.MemberExpression{
		Left:     left,
		Member:   member.Member,
		DataType: ty,
	}
	if diag != nil {
		t.diagnostics.Report(diag.Location(member.MemberLocation))
		return &ir.InvalidExpression{Expression: memberExpr}
	}

	return memberExpr
}

func (t *typeChecker) typeCheckBlock(block *ast.Block, createScope bool) *ir.Block {
	if createScope {
		t.enterScope(&symbols.BlockContext{ResultType: types.Void})
		defer t.exitScope()
	}

	stmts := []ir.Statement{}
	var resultType types.Type = types.Void
	for _, stmt := range block.Statements {
		nextStatement := t.typeCheckStatement(stmt)
		if ir.Diverges(nextStatement, ir.BlockScope) {
			resultType = types.Never
		}
		stmts = append(stmts, nextStatement)
	}
	if len(stmts) == 1 {
		if expr, ok := stmts[0].(ir.Expression); ok {
			return &ir.Block{
				Statements: stmts,
				ResultType: expr.Type(),
			}
		}
	}

	if createScope && resultType != types.Never {
		resultType = t.symbols.Context.(*symbols.BlockContext).ResultType
	}
	return &ir.Block{
		Statements: stmts,
		ResultType: resultType,
	}
}

func (t *typeChecker) typeCheckIfExpression(ifStmt *ast.IfExpression) ir.Expression {
	condition := t.typeCheckExpression(ifStmt.Condition)
	if !types.Assignable(types.Bool, condition.Type()) {
		t.diagnostics.Report(diagnostics.ConditionMustBeBool(ifStmt.Condition.GetLocation()))
	}

	body := t.typeCheckBlock(ifStmt.Body, true)
	resultType := body.ResultType
	var elseBranch ir.Expression

	if ifStmt.ElseBranch != nil {
		elseBranch = t.typeCheckExpression(ifStmt.ElseBranch)

		if resultType == types.Never {
			resultType = elseBranch.Type()
		}
		
		if !types.Assignable(resultType, elseBranch.Type()) {
			t.diagnostics.Report(diagnostics.BranchTypesMustMatch(
				ifStmt.ElseBranch.GetLocation(),
				body.ResultType,
				elseBranch.Type(),
			))
		}
	}

	return &ir.IfExpression{
		Condition:  condition,
		ResultType: resultType,
		Body:       body,
		ElseBranch: elseBranch,
	}
}

func (t *typeChecker) typeCheckWhileLoop(loop *ast.WhileLoop) ir.Expression {
	condition := t.typeCheckExpression(loop.Condition)
	if !types.Assignable(types.Bool, condition.Type()) {
		t.diagnostics.Report(diagnostics.ConditionMustBeBool(loop.Condition.GetLocation()))
	}

	t.enterScope(&symbols.LoopContext{ResultType: types.Void})
	defer t.exitScope()

	body := t.typeCheckBlock(loop.Body, false)
	body.ResultType = t.symbols.Context.(*symbols.LoopContext).ResultType

	return &ir.WhileLoop{
		Condition: condition,
		Body:      body,
	}
}

func (t *typeChecker) typeCheckForLoop(loop *ast.ForLoop) ir.Expression {
	iter := t.typeCheckExpression(loop.Iterator)
	var itemType types.Type = types.Invalid
	if iterator, ok := iter.Type().(types.Iterator); ok {
		itemType = iterator.Item()
	} else {
		t.diagnostics.Report(diagnostics.NotIterable(loop.Iterator.GetLocation()))
	}
	variable := symbols.Variable{
		Name:       loop.Variable,
		IsMut:      false,
		Type:       itemType,
		ConstValue: nil,
	}

	t.enterScope(&symbols.LoopContext{ResultType: types.Void})
	defer t.exitScope()

	t.symbols.Register(&variable)
	body := t.typeCheckBlock(loop.Body, false)
	body.ResultType = t.symbols.Context.(*symbols.LoopContext).ResultType

	return &ir.ForLoop{
		Variable: variable,
		Iterator: iter,
		Body:     body,
	}
}

func (t *typeChecker) typeCheckFunctionExpression(fn *ast.FunctionExpression) ir.Expression {
	// TODO: infer return type
	var returnType types.Type = types.Void
	if fn.ReturnType != nil {
		returnType = t.typeCheckType(fn.ReturnType)
	}

	if fn.Body == nil {
		params := []types.Type{}
		for _, param := range fn.Parameters {
			if param.Name != nil {
				if param.Type == nil {
					if param.Mutable {
						t.diagnostics.Report(diagnostics.MutWithoutParamName(param.Location))
					}
					params = append(params, t.lookupType(*param.Name, param.TypeOrIdent.Location))
				} else {
					t.diagnostics.Report(diagnostics.NamedParamInFnType(param.TypeOrIdent.Location))
				}
			} else if param.Type != nil {
				params = append(params, t.typeCheckType(param.Type))
			}
		}

		return &ir.TypeExpression{DataType: &types.Function{
			Parameters: params,
			ReturnType: returnType,
		}}
	}

	t.enterScope(symbols.FunctionContext{ReturnType: returnType})
	defer t.exitScope()
	params := []string{}
	paramTypes := []types.Type{}

	for _, param := range fn.Parameters {
		if param.Type != nil {
			paramType := t.typeCheckType(param.Type)
			for i := len(paramTypes) - 1; i >= 0; i-- {
				if paramTypes[i] == nil {
					paramTypes[i] = paramType
				} else {
					break
				}
			}
			paramTypes = append(paramTypes, paramType)
		} else {
			paramTypes = append(paramTypes, nil)
		}
	}

	for i, param := range fn.Parameters {
		paramType := paramTypes[i]
		if param.Name == nil {
			t.diagnostics.Report(diagnostics.UnnamedParameter(param.Type.GetLocation()))
			continue
		}
		symbol := &symbols.Variable{
			Name:       *param.Name,
			IsMut:      param.Mutable,
			Type:       paramType,
			ConstValue: nil,
		}
		t.symbols.Register(symbol)
		params = append(params, *param.Name)
	}

	body := t.typeCheckBlock(fn.Body, false)

	return &ir.FunctionExpression{
		Parameters: params,
		Body:       body,
		DataType: &types.Function{
			Parameters: paramTypes,
			ReturnType: returnType,
		},
		Location: fn.Location,
	}
}

func (t *typeChecker) typeCheckRefExpression(ref *ast.RefExpression) ir.Expression {
	value := t.typeCheckExpression(ref.Operand)
	if ref.Mutable && !ir.MutableExpr(value) {
		t.diagnostics.Report(diagnostics.MutRefOfNotMut(ref.Location))
	}

	return &ir.RefExpression{
		Value:   value,
		Mutable: ref.Mutable,
	}
}

func (t *typeChecker) typeCheckDerefExpression(deref *ast.DerefExpression) ir.Expression {
	value := t.typeCheckExpression(deref.Operand)

	if _, ok := value.Type().(*types.Pointer); !ok {
		t.diagnostics.Report(diagnostics.CannotDeref(deref.Operand.GetLocation(), value.Type()))
	}

	return &ir.DerefExpression{
		Value: value,
	}
}

func (t *typeChecker) typeCheckPointerType(ptr *ast.PointerType) ir.Expression {
	ty := t.typeCheckType(ptr.Operand)

	return &ir.TypeExpression{
		DataType: &types.Pointer{
			Underlying: ty,
			Mutable:    ptr.Mutable,
		},
	}
}

func (t *typeChecker) typeCheckOptionType(opt *ast.OptionType) ir.Expression {
	ty := t.typeCheckType(opt.Operand)

	return &ir.TypeExpression{
		DataType: &types.Option{
			SomeType: ty,
		},
	}
}
