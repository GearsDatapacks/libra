package ast

type BaseStatment struct {}

func (stmt *BaseStatment) statementNode() {}

type ExpressionStatement struct {
	*BaseNode
	*BaseStatment
	Expression Expression
}

func (es *ExpressionStatement) Type() NodeType { return "ExpressionStatement" }

func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}
