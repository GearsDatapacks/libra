package ir

import (
	"bytes"

	"github.com/gearsdatapacks/libra/type_checker/types"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	irStmt()
}

type Expression interface {
	Node
	Type() types.Type
	irExpr()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var result bytes.Buffer

	for i, stmt := range p.Statements {
		if i != 0 {
			result.WriteByte('\n')
		}
		result.WriteString(stmt.String())
	}

	return result.String()
}
