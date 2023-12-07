package ast

import (
	"github.com/gearsdatapacks/libra/lexer/token"
)

type NodeType string

type BaseNode struct {
	Token       token.Token
}

func (n *BaseNode) GetToken() token.Token {
	return n.Token
}

type Node interface {
	GetToken() token.Token
	Type() NodeType
	String() string
}

type Statement interface {
	Node
	statementNode()
	MarkExport()
	IsExport() bool
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
