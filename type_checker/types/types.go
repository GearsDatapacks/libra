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

func FromAst(node ast.TypeExpression) (ValidType, error) {
	switch typeExpr := node.(type) {
	case *ast.TypeName:
		return FromString(typeExpr.Name)

	case *ast.Union:
		types := []ValidType{}

		for _, dataType := range typeExpr.ValidTypes {
			nextType, err := FromAst(dataType)
			if err != nil {
				return nil, err
			}
			types = append(types, nextType)
		}

		return MakeUnion(types...), nil

	case *ast.ListType:
		dataType, err := FromAst(typeExpr.ElementType)
		if err != nil {
			return nil, err
		}

		return &ListLiteral{
			ElemType: dataType,
		}, nil
	
	case *ast.ArrayType:
		dataType, err := FromAst(typeExpr.ElementType)
		if err != nil {
			return nil, err
		}

		return &ArrayLiteral{
			ElemType: dataType,
			Length: typeExpr.Length,
		}, nil

	case *ast.VoidType:
		return &Void{}, nil

	case *ast.InferType:
		return nil, errors.TypeError("Expected type, got nothing", node)

	default:
		return nil, errors.DevError("Unexpected type node: " + node.String())
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

func FromString(typeString string) (ValidType, error) {
	dataType, ok := typeTable[typeString]
	if !ok {
		return nil, errors.TypeError(fmt.Sprintf("Invalid type %q", typeString))
	}

	return dataType, nil
}
