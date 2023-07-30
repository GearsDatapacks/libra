package ast

type IntegerLiteral struct {
	*BaseNode
	Value int
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) ToString() string {
	return il.Token.Value
}

type BinaryOperation struct {
	*BaseNode
	Left Expression
	Op string
	Right Expression
}

func (binOp *BinaryOperation) expressionNode() {}

func (binOp *BinaryOperation) ToString() string {
	result := ""

	result += "("
	result += binOp.Left.ToString()
	result += " "
	result += binOp.Op
	result += " "
	result += binOp.Right.ToString()
	result += ")"

	return result
}
