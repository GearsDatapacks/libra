package types

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
)

type ValidType interface {
	Valid(ValidType) bool
	String() string
	WasVariable() bool
	MarkVariable()
}

type BaseType struct {
	wasVariable bool
}

func (b *BaseType) WasVariable() bool {
	return b.wasVariable
}

func (b *BaseType) MarkVariable() {
	b.wasVariable = true
}

func FromAst(node ast.TypeExpression) ValidType {
	switch typeExpr := node.(type) {
	case *ast.TypeName:
		return FromString(typeExpr.Name)

	case *ast.Union:
		types := []ValidType{}

		for _, dataType := range typeExpr.ValidTypes {
			types = append(types, FromAst(dataType))
		}

		return MakeUnion(types...)

	case *ast.ListType:
		return &ListLiteral{
			ElemType: FromAst(typeExpr.ElementType),
		}
	
	case *ast.ArrayType:
		return &ArrayLiteral{
			ElemType: FromAst(typeExpr.ElementType),
			Length: typeExpr.Length,
		}

	case *ast.VoidType:
		return &Void{}

	case *ast.InferType:
		errors.TypeError("Expected type, got nothing", node)
		return &IntLiteral{}

	default:
		errors.DevError("Unexpected type node: " + node.String())
		return &IntLiteral{}
	}
}

var typeTable = map[string]ValidType{
	"int":      &IntLiteral{},
	"float":    &FloatLiteral{},
	"boolean":  &BoolLiteral{},
	"null":     &NullLiteral{},
	"function": &Function{},
	"string":   &StringLiteral{},
}

func FromString(typeString string) ValidType {
	dataType, ok := typeTable[typeString]
	if !ok {
		errors.TypeError(fmt.Sprintf("Invalid type %q", typeString))
	}

	return dataType
}
