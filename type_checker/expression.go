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
		return t.lookupVariable(expr.Token)
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

func (t *typeChecker) lookupVariable(tok token.Token) ir.Expression {
	symbol := t.symbols.Lookup(tok.Value)
	if symbol == nil {
		t.diagnostics.Report(diagnostics.VariableUndefined(tok.Location, tok.Value))
		symbol = &symbols.Variable{
			Name:  tok.Value,
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

	leftNumeric := types.Assignable(types.Int, lType) || types.Assignable(types.Float, lType)
	rightNumeric := types.Assignable(types.Int, rType) || types.Assignable(types.Float, rType)
	isFloat := !types.Assignable(types.Int, lType) || !types.Assignable(types.Int, rType)
	var lUntyped bool
	var rUntyped bool
	if v, ok := lType.(types.VariableType); ok {
		lUntyped = v.Untyped
	}
	if v, ok := rType.(types.VariableType); ok {
		rUntyped = v.Untyped
	}
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
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			}
		}
	case token.RIGHT_ANGLE:
		if leftNumeric && rightNumeric {
			binOp = ir.Greater
			if isFloat {
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			}
		}
	case token.LEFT_ANGLE_EQUALS:
		if leftNumeric && rightNumeric {
			binOp = ir.LessEq
			if isFloat {
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			}
		}
	case token.RIGHT_ANGLE_EQUALS:
		if leftNumeric && rightNumeric {
			binOp = ir.GreaterEq
			if isFloat {
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
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
		if types.Assignable(types.Int, lType) && types.Assignable(types.Int, rType) {
			binOp = ir.LeftShift
		}
	case token.DOUBLE_RIGHT_ANGLE:
		if types.Assignable(types.Int, lType) && types.Assignable(types.Int, rType) {
			binOp = ir.RightShift
		}
	case token.PLUS:
		if types.Assignable(types.String, lType) && types.Assignable(types.String, rType) {
			binOp = ir.Concat
		}

		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.AddFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.AddInt
			}
		}
	case token.MINUS:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.SubtractFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.SubtractInt
			}
		}
	case token.STAR:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.MultiplyFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.MultiplyInt
			}
		}
	case token.SLASH:
		if leftNumeric && rightNumeric {
			binOp = ir.Divide
			if isFloat {
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			}
		}
	case token.PERCENT:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.ModuloFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.ModuloInt
			}
		}
	case token.DOUBLE_STAR:
		if leftNumeric && rightNumeric {
			if isFloat {
				binOp = ir.PowerFloat
				lhs = convert(lhs, types.Float, operator)
				rhs = convert(rhs, types.Float, operator)
			} else {
				binOp = ir.PowerInt
			}
		}
	case token.PIPE:
		if types.Assignable(types.Int, lType) && types.Assignable(types.Int, rType) {
			binOp = ir.BitwiseOr
		}
		if types.Assignable(types.RuntimeType, lType) && types.Assignable(types.RuntimeType, rType) {
			binOp = ir.Union
		}
	case token.AMPERSAND:
		if types.Assignable(types.Int, lType) && types.Assignable(types.Int, rType) {
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
			Operator: 0,
		}
	}

	operator := getPrefixOperator(unExpr.Operator.Kind, operand)

	if operator == 0 {
		t.diagnostics.Report(diagnostics.UnaryOperatorUndefined(unExpr.Operator.Location, unExpr.Operator.Value, operand.Type()))
	}

	// We can safely ignore the identity operator
	if operator & ^ir.UntypedBit == ir.Identity {
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
			Operator: 0,
		}
	}

	operator, diag := t.getPostfixOperator(unExpr.Operator.Kind, operand)

	if diag != nil {
		t.diagnostics.Report(diag.Location(unExpr.Operand.Location()))
	} else if operator == 0 {
		t.diagnostics.Report(diagnostics.UnaryOperatorUndefined(unExpr.Operator.Location, unExpr.Operator.Value, operand.Type()))
	}

	return &ir.UnaryExpression{
		Operator: operator,
		Operand:  operand,
	}
}

func getPrefixOperator(tokKind token.Kind, operand ir.Expression) ir.UnaryOperator {
	var unOp ir.UnaryOperator
	opType := operand.Type()

	numeric := types.Assignable(types.Int, opType) || types.Assignable(types.Float, opType)
	isFloat := !types.Assignable(types.Int, opType)
	var untyped bool
	if v, ok := opType.(types.VariableType); ok {
		untyped = v.Untyped
	}

	switch tokKind {
	case token.MINUS:
		if numeric {
			if isFloat {
				unOp = ir.NegateFloat
			} else {
				unOp = ir.NegateInt
			}
		}
	case token.PLUS:
		if numeric {
			unOp = ir.Identity
		}
	case token.BANG:
		if types.Assignable(types.Bool, opType) {
			unOp = ir.LogicalNot
		}
	case token.TILDE:
		if types.Assignable(types.Int, opType) {
			unOp = ir.BitwiseNot
		}
	}

	if untyped && unOp != 0 {
		unOp = unOp | ir.UntypedBit
	}

	return unOp
}

func (t *typeChecker) getPostfixOperator(tokKind token.Kind, operand ir.Expression) (ir.UnaryOperator, *diagnostics.Partial) {
	var unOp ir.UnaryOperator
	opType := operand.Type()

	numeric := types.Assignable(types.Int, opType) || types.Assignable(types.Float, opType)
	isFloat := !types.Assignable(types.Int, opType)
	var untyped bool
	if v, ok := opType.(types.VariableType); ok {
		untyped = v.Untyped
	}

	switch tokKind {
	case token.DOUBLE_PLUS:
		if numeric {
			if !ir.AssignableExpr(operand) {
				return 0, diagnostics.CannotIncDec("increment")
			} else if !ir.MutableExpr(operand) {
				return 0, diagnostics.ValueImmutablePartial
			}
			if isFloat {
				unOp = ir.IncrementFloat
			} else {
				unOp = ir.IncrecementInt
			}
		}
	case token.DOUBLE_MINUS:
		if numeric {
			if !ir.AssignableExpr(operand) {
				return 0, diagnostics.CannotIncDec("decrement")
			} else if !ir.MutableExpr(operand) {
				return 0, diagnostics.ValueImmutablePartial
			}
			if isFloat {
				unOp = ir.DecrementFloat
			} else {
				unOp = ir.DecrecementInt
			}
		}
	case token.QUESTION:
		panic("TODO: '?' unary operator")

	case token.BANG:
		panic("TODO: '!' postfix operator")
	}

	if untyped && unOp != 0 {
		unOp = unOp | ir.UntypedBit
	}
	return unOp, nil
}

func (t *typeChecker) typeCheckCastExpression(expr *ast.CastExpression) ir.Expression {
	value := t.typeCheckExpression(expr.Left)
	ty := t.typeCheckType(expr.Type)
	conversion := convert(value, ty, explicit)
	if conversion == nil {
		t.diagnostics.Report(diagnostics.CannotCast(expr.Left.Location(), value.Type(), ty))
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
			t.diagnostics.Report(diagnostics.NotAssignable(elem.Location(), elemType, value.Type()))
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
		elemType := t.typeFromExpr(left, indexExpr.Left.Location())

		if indexExpr.Index == nil {
			return &ir.TypeExpression{DataType: &types.ListType{ElemType: elemType}}
		}

		length := -1

		if ident, ok := indexExpr.Index.(*ast.Identifier); !ok || ident.Name != "_" {
			expr := convert(t.typeCheckExpression(indexExpr.Index), types.Int, implicit)
			if expr == nil {
				t.diagnostics.Report(diagnostics.CountMustBeInt(indexExpr.Index.Location()))
			} else if expr.IsConst() {
				value := expr.ConstValue().(values.IntValue)
				length = int(value.Value)
			} else {
				t.diagnostics.Report(diagnostics.NotConst(indexExpr.Index.Location()))
			}
		}

		return &ir.TypeExpression{DataType: &types.ArrayType{
			ElemType: elemType,
			Length:   length,
			CanInfer: false,
		}}
	}

	index := t.typeCheckExpression(indexExpr.Index)
	ty, diag := ir.Index(left, index)

	expr := &ir.IndexExpression{
		Left:     left,
		Index:    index,
		DataType: ty,
	}
	if diag != nil {
		t.diagnostics.Report(diag.Location(indexExpr.Index.Location()))
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
			t.diagnostics.Report(diagnostics.NotAssignable(kv.Key.Location(), keyType, key.Type()))
			continue
		}

		value := t.typeCheckExpression(kv.Value)
		if valueType == types.Invalid {
			valueType = types.ToReal(value.Type())
		}
		convertedValue := convert(value, valueType, operator)
		if convertedValue == nil {
			t.diagnostics.Report(diagnostics.NotAssignable(kv.Value.Location(), valueType, value.Type()))
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
				KeyType:   t.typeFromExpr(keyValues[0].Key, mapLit.KeyValues[0].Key.Location()),
				ValueType: t.typeFromExpr(keyValues[0].Value, mapLit.KeyValues[0].Value.Location()),
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
		t.diagnostics.Report(diagnostics.NotHashable(mapLit.KeyValues[0].Key.Location(), keyType))
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
		t.diagnostics.Report(diagnostics.CannotAssign(assignment.Assignee.Location()))
	} else if !ir.MutableExpr(assignee) {
		t.diagnostics.Report(diagnostics.ValueImmutable(assignment.Assignee.Location()))
	} else {
		conversion := convert(value, assignee.Type(), implicit)
		if conversion == nil {
			t.diagnostics.Report(diagnostics.NotAssignable(assignment.Assignee.Location(), assignee.Type(), value.Type()))
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
			tupleMembers = append(tupleMembers, t.typeFromExpr(val, tuple.Values[i].Location()))
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
	funcType, ok := fn.Type().(*types.Function)
	if !ok {
		t.diagnostics.Report(diagnostics.NotCallable(call.Callee.Location(), fn.Type()))
		return &ir.InvalidExpression{
			Expression: &ir.FunctionCall{
				Function:   fn,
				Arguments:  []ir.Expression{},
				ReturnType: types.Invalid,
			},
		}
	}

	if len(call.Arguments) != len(funcType.Parameters) {
		t.diagnostics.Report(diagnostics.WrongNumberAgruments(call.Callee.Location(), len(funcType.Parameters), len(call.Arguments)))
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
			t.diagnostics.Report(diagnostics.NotAssignable(arg.Location(), expectedType, value.Type()))
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
				t.diagnostics.Report(diagnostics.NoNameStructMember(member.Value.Location()))
				continue
			}
			field, ok := structTy.Fields[member.Name.Value]
			if !ok {
				t.diagnostics.Report(diagnostics.NoStructMember(member.Name.Location, structTy.Name, member.Name.Value))
				continue
			}
			var value ir.Expression
			if member.Value != nil {
				value = t.typeCheckExpression(member.Value)
			} else {
				value = t.lookupVariable(*member.Name)
			}

			conversion := convert(value, field.Type, implicit)
			if conversion != nil {
				value = conversion
			} else {
				t.diagnostics.Report(diagnostics.NotAssignable(member.Value.Location(), field.Type, value.Type()))
			}

			fields[member.Name.Value] = value
		}

		return &ir.StructExpression{
			Struct: baseTy,
			Fields: fields,
		}
	} else if tupleTy, ok := ty.(*types.TupleStruct); ok {
		fields := []ir.Expression{}

		if len(structExpr.Members) != len(tupleTy.Types) {
			t.diagnostics.Report(diagnostics.WrongNumberTupleValues(structExpr.Struct.Location(), len(tupleTy.Types), len(structExpr.Members)))
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
			var location text.Location
			if member.Value != nil {
				value = t.typeCheckExpression(member.Value)
				location = member.Value.Location()
				if member.Name != nil {
					t.diagnostics.Report(diagnostics.TupleStructWithNames(member.Name.Location))
				}
			} else {
				value = t.lookupVariable(*member.Name)
				location = member.Name.Location
			}

			conversion := convert(value, field, implicit)
			if conversion != nil {
				value = conversion
			} else {
				t.diagnostics.Report(diagnostics.NotAssignable(location, field, value.Type()))
			}

			fields = append(fields, value)
		}

		return &ir.TupleStructExpression{
			Struct: baseTy,
			Fields: fields,
		}
	} else {
		t.diagnostics.Report(diagnostics.CannotConstruct(structExpr.Struct.Location(), ty))
		return &ir.InvalidExpression{
			Expression: &ir.IntegerLiteral{},
		}
	}
}

func (t *typeChecker) typeCheckMemberExpression(member *ast.MemberExpression) ir.Expression {
	left := t.typeCheckExpression(member.Left)
	ty, diag := ir.Member(left, member.Member.Value)

	memberExpr := &ir.MemberExpression{
		Left:     left,
		Member:   member.Member.Value,
		DataType: ty,
	}
	if diag != nil {
		t.diagnostics.Report(diag.Location(member.Member.Location))
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
	for _, stmt := range block.Statements {
		stmts = append(stmts, t.typeCheckStatement(stmt))
	}
	if len(stmts) == 1 {
		if expr, ok := stmts[0].(ir.Expression); ok {
			return &ir.Block{
				Statements: stmts,
				ResultType: expr.Type(),
			}
		}
	}

	var resultType types.Type = types.Void
	if createScope {
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
		t.diagnostics.Report(diagnostics.ConditionMustBeBool(ifStmt.Condition.Location()))
	}

	body := t.typeCheckBlock(ifStmt.Body, true)
	var elseBranch ir.Statement
	if ifStmt.ElseBranch != nil {
		elseBranch = t.typeCheckStatement(ifStmt.ElseBranch.Statement)
	}
	return &ir.IfExpression{
		Condition:  condition,
		Body:       body,
		ElseBranch: elseBranch,
	}
}

func (t *typeChecker) typeCheckWhileLoop(loop *ast.WhileLoop) ir.Expression {
	condition := t.typeCheckExpression(loop.Condition)
	if !types.Assignable(types.Bool, condition.Type()) {
		t.diagnostics.Report(diagnostics.ConditionMustBeBool(loop.Condition.Location()))
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
		t.diagnostics.Report(diagnostics.NotIterable(loop.Iterator.Location()))
	}
	variable := symbols.Variable{
		Name:       loop.Variable.Value,
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
		returnType = t.typeCheckType(fn.ReturnType.Type)
	}

	if fn.Body == nil {
		params := []types.Type{}
		for _, param := range fn.Parameters {
			if param.Name != nil {
				if param.Type == nil {
					if param.Mutable != nil {
						t.diagnostics.Report(diagnostics.MutWithoutParamName(param.Mutable.Location))
					}
					params = append(params, t.lookupType(*param.Name))
				} else {
					t.diagnostics.Report(diagnostics.NamedParamInFnType(param.Name.Location))
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
			t.diagnostics.Report(diagnostics.UnnamedParameter(param.Type.Location()))
			continue
		}
		symbol := &symbols.Variable{
			Name:       param.Name.Value,
			IsMut:      param.Mutable != nil,
			Type:       paramType,
			ConstValue: nil,
		}
		t.symbols.Register(symbol)
		params = append(params, param.Name.Value)
	}

	body := t.typeCheckBlock(fn.Body, false)

	return &ir.FunctionExpression{
		Parameters: params,
		Body:       body,
		DataType: &types.Function{
			Parameters: paramTypes,
			ReturnType: returnType,
		},
	}
}

func (t *typeChecker) typeCheckRefExpression(ref *ast.RefExpression) ir.Expression {
	value := t.typeCheckExpression(ref.Operand)
	isMut := ref.Mutable != nil
	if isMut && !ir.MutableExpr(value) {
		t.diagnostics.Report(diagnostics.MutRefOfNotMut(ref.Mutable.Location))
	}

	return &ir.RefExpression{
		Value:   value,
		Mutable: isMut,
	}
}

func (t *typeChecker) typeCheckDerefExpression(deref *ast.DerefExpression) ir.Expression {
	value := t.typeCheckExpression(deref.Operand)

	if value.Type() == types.RuntimeType {
		return &ir.TypeExpression{
			DataType: &types.Pointer{
				Underlying: value.ConstValue().(values.TypeValue).Type.(types.Type),
				Mutable:    deref.Mutable != nil,
			},
		}
	}

	if deref.Mutable != nil {
		t.diagnostics.Report(diagnostics.MutDerefNotAllowed(deref.Mutable.Location))
	}
	if _, ok := value.Type().(*types.Pointer); !ok {
		t.diagnostics.Report(diagnostics.CannotDeref(deref.Operand.Location(), value.Type()))
	}

	return &ir.DerefExpression{
		Value: value,
	}
}
