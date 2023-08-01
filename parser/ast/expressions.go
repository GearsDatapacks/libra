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

type Identifier struct {
	*BaseNode
	*BaseExpression
	Symbol string
}

func (ident *Identifier) Type() NodeType { return "Identifier" }

func (ident *Identifier) String() string {
	return ident.Symbol
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

type AssignmentExpression struct {
	*BaseNode
	*BaseExpression
	Assignee     Expression
	Value    Expression
}

func (ae *AssignmentExpression) Type() NodeType { return "AssignmentExpression" }

func (ae *AssignmentExpression) String() string {
	result := ""

	result += ae.Assignee.String()
	result += " = "
	result += ae.Value.String()

	return result
}