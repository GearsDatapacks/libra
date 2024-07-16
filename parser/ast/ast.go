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

	astPrinter := printer.New(&text, false)
	printer.QueueNodeList(astPrinter, p.Statements, true)

	return text.String()
}

func (p *Program) Print() {
	astPrinter := printer.New(os.Stdout, true)
	printer.QueueNodeList(astPrinter, p.Statements, true)
}
