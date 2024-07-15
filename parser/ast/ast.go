package ast

import (
	"bytes"
	"os"

	"github.com/gearsdatapacks/libra/printer"
	"github.com/gearsdatapacks/libra/text"
)

type Statement interface {
	printer.Printable
	GetLocation() text.Location
}

type Expression interface {
	Statement
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// func (p *Program) Tokens() []token.Token {
// 	tokens := []token.Token{}

// 	for _, statement := range p.Statements {
// 		tokens = append(tokens, statement.Tokens()...)
// 	}

// 	return tokens
// }

func (p *Program) String() string {
	var text bytes.Buffer

	context := printer.New(&text, false)
	printer.WriteNodeList(context, p.Statements)

	return text.String()
}

func (p *Program) Print() {
	context := printer.New(os.Stdout, true)
	printer.WriteNodeList(context, p.Statements)
}
