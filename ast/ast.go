package ast

import (
	"github.com/gearsdatapacks/libra/token"
)

type NodeType string

type BaseNode struct {
	Token token.Token
	isStatement bool
}

func (n *BaseNode) Line() int {
	return n.Token.Line
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

type node interface {
	GetToken() token.Token
	// Type() NodeType
	Line() int
	ToString() string
	IsExpression() bool
	IsStatement() bool
	MarkStatement()
	MarkExpression()
}

type Statement interface {
	node
	statementNode()
}

type Expression interface {
	node
	expressionNode()
}

type Program struct {
	Body []Statement
}

func (p *Program) ToString() string {
	result := ""

	for _, statement := range p.Body {
		result += statement.ToString()
		result += "\n"
	}

	return result
}
