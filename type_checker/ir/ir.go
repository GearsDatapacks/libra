package ir

import (
	"bytes"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

type Statement interface {
	String() string
}

type Expression interface {
	Statement
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
	case *MemberExpression:
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
		return e.Symbol.IsMut
	case *IndexExpression:
		return MutableExpr(e.Left)
	case *MemberExpression:
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

func Member(left Expression, member string) (types.Type, *diagnostics.Partial) {
	if left.IsConst() {
		return types.Member(left.Type(), member, left.ConstValue())
	}
	return types.Member(left.Type(), member)
}
