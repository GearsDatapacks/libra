package ast

import (
	"bytes"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type Node interface {
	Tokens() []token.Token
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Tokens() []token.Token {
	tokens := []token.Token{}

	for _, statement := range p.Statements {
		tokens = append(tokens, statement.Tokens()...)
	}

	return tokens
}

func (p *Program) String() string {
	text := bytes.NewBuffer([]byte{})

	for i, statement := range p.Statements {
		if i != 0 {
			text.WriteByte('\n')
		}

		text.WriteString(statement.String())
	}

	return text.String()
}
