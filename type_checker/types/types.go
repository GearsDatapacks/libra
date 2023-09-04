package types

import (
	"fmt"
	"log"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
)

type ValidType interface {
	Valid(ValidType) bool
	String() string
	WasVariable() bool
	MarkVariable()
	IndexBy(ValidType) ValidType
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

func (*BaseType) IndexBy(ValidType) ValidType {
	return nil
}

type PartialType interface {
	ValidType
	Infer(ValidType) (ValidType, bool)
}

func FromAst(node ast.TypeExpression) ValidType {
	switch typeExpr := node.(type) {
	case *ast.TypeName:
		return FromString(typeExpr.Name)

	case *ast.Union:
		types := []ValidType{}

		for _, dataType := range typeExpr.ValidTypes {
			nextType := FromAst(dataType)
			if nextType.String() == "TypeError" {
				return nextType
			}
			types = append(types, nextType)
		}

		return MakeUnion(types...)

	case *ast.ListType:
		dataType := FromAst(typeExpr.ElementType)
		if dataType.String() == "TypeError" {
			return dataType
		}

		return &ListLiteral{
			ElemType: dataType,
		}

	case *ast.ArrayType:
		dataType := FromAst(typeExpr.ElementType)
		if dataType.String() == "TypeError" {
			return dataType
		}

		return &ArrayLiteral{
			ElemType: dataType,
			Length:   typeExpr.Length,
		}

	case *ast.VoidType:
		return &Void{}

	case *ast.InferType:
		return &Infer{}

	default:
		log.Fatal(errors.DevError("Unexpected type node: " + node.String()))
		return nil
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
		return Error(fmt.Sprintf("Invalid type %q", typeString))
	}

	return dataType
}
