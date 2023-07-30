package ast

type IntegerLiteral struct {
	*BaseNode
	Value int
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) Type() NodeType { return "Integer" }

func (il *IntegerLiteral) String() string {
	return il.Token.Value
}

type BinaryOperation struct {
	*BaseNode
	Left  Expression
	Op    string
	Right Expression
}

func (binOp *BinaryOperation) expressionNode() {}

func (binOp *BinaryOperation) Type() NodeType { return "BinaryOperation" }

func (binOp *BinaryOperation) String() string {
	result := ""

	result += "("
	result += binOp.Left.String()
	result += " "
	result += binOp.Op
	result += " "
	result += binOp.Right.String()
	result += ")"

	return result
}
