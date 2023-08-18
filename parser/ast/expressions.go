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

type FloatLiteral struct {
	*BaseNode
	*BaseExpression
	Value float64
}

func (fl *FloatLiteral) Type() NodeType { return "Float" }

func (fl *FloatLiteral) String() string {
	return fl.Token.Value
}

type BooleanLiteral struct {
	*BaseNode
	*BaseExpression
	Value bool
}

func (bl *BooleanLiteral) Type() NodeType { return "Boolean" }

func (bl *BooleanLiteral) String() string {
	return bl.Token.Value
}

type NullLiteral struct {
	*BaseNode
	*BaseExpression
}

func (nl *NullLiteral) Type() NodeType { return "Null" }

func (nl *NullLiteral) String() string { return "null" }

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