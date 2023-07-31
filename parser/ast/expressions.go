package ast

type BaseExpression struct{}

func (exp *BaseExpression) expressionNode() {}

type IntegerLiteral struct {
	*BaseNode
	*BaseExpression
	Value int
}

func (il *IntegerLiteral) Type() NodeType { return "Integer" }

func (il *IntegerLiteral) String() string {
	return il.Token.Value
}

type BinaryOperation struct {
	*BaseNode
	*BaseExpression
	Left     Expression
	Operator string
	Right    Expression
}

func (binOp *BinaryOperation) Type() NodeType { return "BinaryOperation" }

func (binOp *BinaryOperation) String() string {
	result := ""

	result += "("
	result += binOp.Left.String()
	result += " "
	result += binOp.Operator
	result += " "
	result += binOp.Right.String()
	result += ")"

	return result
}
