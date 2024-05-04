package ir

import (
	"bytes"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
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
	IsConst() bool
	ConstValue() values.ConstValue
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

func AssignableExpr(expr Expression) bool {
	switch e := expr.(type) {
	case *VariableExpression:
		return true
	case *IndexExpression:
		return AssignableExpr(e.Left)
	case *InvalidExpression:
		return true
	default:
		return false
	}
}

func MutableExpr(expr Expression) bool {
	switch e := expr.(type) {
	case *VariableExpression:
		return e.Symbol.Mutable
	case *IndexExpression:
		return MutableExpr(e.Left)
	case *InvalidExpression:
		return true
	default:
		return false
	}
}

func Index(left, index Expression) (types.Type, *diagnostics.Partial) {
	if index.IsConst() {
		if left.IsConst() {
			return types.Index(left.Type(), index.Type(), index.ConstValue(), left.ConstValue())
		}
		return types.Index(left.Type(), index.Type(), index.ConstValue())
	}
	return types.Index(left.Type(), index.Type())
}
