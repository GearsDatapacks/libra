package ast

type ExpressionStatement struct {
	*BaseNode
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) ToString() string {
	return es.Expression.ToString()
}
