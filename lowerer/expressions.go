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

func (l *lowerer) lowerBinaryExpression(expr *ir.BinaryExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerUnaryExpression(expr *ir.UnaryExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerConversion(expr *ir.Conversion) ir.Expression {
	return expr
}

func (l *lowerer) lowerInvalidExpression(expr *ir.InvalidExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerArrayExpression(expr *ir.ArrayExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerIndexExpression(expr *ir.IndexExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerMapExpression(expr *ir.MapExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerAssignment(expr *ir.Assignment) ir.Expression {
	return expr
}

func (l *lowerer) lowerTupleExpression(expr *ir.TupleExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerTypeCheck(expr *ir.TypeCheck) ir.Expression {
	return expr
}

func (l *lowerer) lowerFunctionCall(expr *ir.FunctionCall) ir.Expression {
	return expr
}

func (l *lowerer) lowerStructExpression(expr *ir.StructExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerTupleStructExpression(expr *ir.TupleStructExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerMemberExpression(expr *ir.MemberExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerBlock(expr *ir.Block) ir.Expression {
	return expr
}

func (l *lowerer) lowerIfExpression(expr *ir.IfExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerWhileLoop(expr *ir.WhileLoop) ir.Expression {
	return expr
}

func (l *lowerer) lowerForLoop(expr *ir.ForLoop) ir.Expression {
	return expr
}

func (l *lowerer) lowerTypeExpression(expr *ir.TypeExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerFunctionExpression(expr *ir.FunctionExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerRefExpression(expr *ir.RefExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerDerefExpression(expr *ir.DerefExpression) ir.Expression {
	return expr
}

func (l *lowerer) lowerEnumMember(expr *ir.EnumMember) ir.Expression {
	return expr
}
