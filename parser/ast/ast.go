package ast

import (
	"github.com/gearsdatapacks/libra/lexer/token"
)

type NodeType string

type BaseNode struct {
	Token       token.Token
	isStatement bool
}

func (n *BaseNode) GetToken() token.Token {
	return n.Token
}

func (n *BaseNode) IsStatement() bool {
	return n.isStatement
}

func (n *BaseNode) IsExpression() bool {
	return !n.isStatement
}

func (n *BaseNode) MarkStatement() {
	n.isStatement = true
}

func (n *BaseNode) MarkExpression() {
	n.isStatement = false
}

type Node interface {
	GetToken() token.Token
	Type() NodeType
	String() string
	IsExpression() bool
	IsStatement() bool
	MarkStatement()
	MarkExpression()
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type TypeExpression interface {
	Node
	typeNode()
}

type Program struct {
	Body []Statement
}

func (p *Program) String() string {
	result := ""

	for _, statement := range p.Body {
		result += statement.String()
		result += "\n"
	}

	return result
}
