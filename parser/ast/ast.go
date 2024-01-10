package ast

import (
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type NodeType string

type BaseNode struct {
	Token    token.Token
	DataType types.ValidType
}

func (n *BaseNode) GetToken() token.Token {
	return n.Token
}

func (n *BaseNode) SetType(dataType types.ValidType) {
	n.DataType = dataType
}

func (n *BaseNode) GetType() types.ValidType {
	return n.DataType
}

type Node interface {
	GetToken() token.Token
	Type() NodeType
	SetType(types.ValidType)
	GetType() types.ValidType
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
