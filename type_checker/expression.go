package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (t *typeChecker) typeCheckExpression(expression ast.Expression) ir.Expression {
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
		return t.typeCheckIdentifier(expr)
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

	case *ast.Block:
		return t.typeCheckBlock(expr, true)
	case *ast.IfExpression:
		return t.typeCheckIfExpression(expr)
	case *ast.WhileLoop:
		return t.typeCheckWhileLoop(expr)
	case *ast.ForLoop:
		return t.typeCheckForLoop(expr)

	default:
		panic(fmt.Sprintf("TODO: Type-check %T", expr))
	}
}

func (t *typeChecker) typeCheckIdentifier(ident *ast.Identifier) ir.Expression {
	symbol := t.symbols.Lookup(ident.Name)
	if symbol == nil {
		t.Diagnostics.Report(diagnostics.VariableUndefined(ident.Token.Location, ident.Name))
		symbol = &symbols.Variable{
			Name:  ident.Name,
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
		t.Diagnostics.Report(diagnostics.BinaryOperatorUndefined(binExpr.Operator.Location, binExpr.Operator.Value, left.Type(), right.Type()))
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
		t.Diagnostics.Report(diagnostics.UnaryOperatorUndefined(unExpr.Operator.Location, unExpr.Operator.Value, operand.Type()))
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
		t.Diagnostics.Report(diag.Location(unExpr.Operand.Location()))
	} else if operator == 0 {
		t.Diagnostics.Report(diagnostics.UnaryOperatorUndefined(unExpr.Operator.Location, unExpr.Operator.Value, operand.Type()))
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
		t.Diagnostics.Report(diagnostics.CannotCast(expr.Left.Location(), value.Type(), ty))
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
			t.Diagnostics.Report(diagnostics.NotAssignable(elem.Location(), elemType, value.Type()))
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
				t.Diagnostics.Report(diagnostics.CountMustBeInt(indexExpr.Index.Location()))
			} else if expr.IsConst() {
				value := expr.ConstValue().(values.IntValue)
				length = int(value.Value)
			} else {
				t.Diagnostics.Report(diagnostics.NotConst(indexExpr.Index.Location()))
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
		t.Diagnostics.Report(diag.Location(indexExpr.Index.Location()))
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
			t.Diagnostics.Report(diagnostics.NotAssignable(kv.Key.Location(), keyType, key.Type()))
			continue
		}

		value := t.typeCheckExpression(kv.Value)
		if valueType == types.Invalid {
			valueType = types.ToReal(value.Type())
		}
		convertedValue := convert(value, valueType, operator)
		if convertedValue == nil {
			t.Diagnostics.Report(diagnostics.NotAssignable(kv.Value.Location(), valueType, value.Type()))
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
		t.Diagnostics.Report(diagnostics.NotHashable(mapLit.KeyValues[0].Key.Location(), keyType))
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
			t.Diagnostics.Report(diagnostics.BinaryOperatorUndefined(assignment.Operator.Location, assignment.Operator.Value, left.Type(), right.Type()))
		}

		value = &ir.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	if !ir.AssignableExpr(assignee) {
		t.Diagnostics.Report(diagnostics.CannotAssign(assignment.Assignee.Location()))
	} else if !ir.MutableExpr(assignee) {
		t.Diagnostics.Report(diagnostics.ValueImmutable(assignment.Assignee.Location()))
	} else {
		conversion := convert(value, assignee.Type(), implicit)
		if conversion == nil {
			t.Diagnostics.Report(diagnostics.NotAssignable(assignment.Assignee.Location(), assignee.Type(), value.Type()))
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
		t.Diagnostics.Report(diagnostics.NotCallable(call.Callee.Location(), fn.Type()))
		return &ir.InvalidExpression{
			Expression: &ir.FunctionCall{
				Function:   fn,
				Arguments:  []ir.Expression{},
				ReturnType: types.Invalid,
			},
		}
	}

	if len(call.Arguments) != len(funcType.Parameters) {
		t.Diagnostics.Report(diagnostics.WrongNumberAgruments(call.Callee.Location(), len(funcType.Parameters), len(call.Arguments)))
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
			t.Diagnostics.Report(diagnostics.NotAssignable(arg.Location(), expectedType, value.Type()))
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
	struc := t.typeCheckExpression(structExpr.Struct)
	if !types.Assignable(struc.Type(), types.RuntimeType) {
		t.Diagnostics.Report(diagnostics.OnlyConstructTypes(structExpr.Struct.Location()))
		return &ir.InvalidExpression{
			Expression: &ir.IntegerLiteral{},
		}
	}
	if !struc.IsConst() {
		t.Diagnostics.Report(diagnostics.NotConst(structExpr.Struct.Location()))
		return &ir.InvalidExpression{
			Expression: &ir.IntegerLiteral{},
		}
	}
	ty := struc.ConstValue().(values.TypeValue)
	structTy, ok := ty.Type.(*types.Struct)
	if !ok {
		t.Diagnostics.Report(diagnostics.CannotConstruct(structExpr.Struct.Location(), ty.Type.(types.Type)))
		return &ir.InvalidExpression{
			Expression: &ir.IntegerLiteral{},
		}
	}

	fields := map[string]ir.Expression{}

	for _, member := range structExpr.Members {
		ty, ok := structTy.Fields[member.Name.Value]
		if !ok {
			t.Diagnostics.Report(diagnostics.NoStructMember(member.Name.Location, structTy.Name, member.Name.Value))
			continue
		}
		value := t.typeCheckExpression(member.Value)
		conversion := convert(value, ty, implicit)
		if conversion != nil {
			value = conversion
		} else {
			t.Diagnostics.Report(diagnostics.NotAssignable(member.Value.Location(), ty, value.Type()))
		}

		fields[member.Name.Value] = value
	}

	return &ir.StructExpression{
		Struct: structTy,
		Fields: fields,
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
		t.Diagnostics.Report(diag.Location(member.Member.Location))
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
		t.Diagnostics.Report(diagnostics.ConditionMustBeBool(ifStmt.Condition.Location()))
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
		t.Diagnostics.Report(diagnostics.ConditionMustBeBool(loop.Condition.Location()))
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
		t.Diagnostics.Report(diagnostics.NotIterable(loop.Iterator.Location()))
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
