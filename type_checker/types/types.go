package types

import (
	"log"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
)

type ValidType interface {
	Valid(ValidType) bool
	String() string
	WasVariable() bool
	MarkVariable()
	Constant() bool
	MarkConstant()
	IndexBy(ValidType) ValidType
	Member(string) ValidType
}

type BaseType struct {
	wasVariable bool
	constant bool
}

func (b *BaseType) WasVariable() bool {
	return b.wasVariable
}

func (b *BaseType) MarkVariable() {
	b.wasVariable = true
}

func (b *BaseType) Constant() bool {
	return b.constant
}

func (b *BaseType) MarkConstant() {
	b.constant = true
}

func (*BaseType) IndexBy(ValidType) ValidType {
	return nil
}

func (*BaseType) Member(string) ValidType {
	return nil
}

type PartialType interface {
	ValidType
	Infer(ValidType) (ValidType, bool)
}

func FromAst(node ast.TypeExpression, table TypeTable) ValidType {
	switch typeExpr := node.(type) {
	case *ast.TypeName:
		return FromString(typeExpr.Name, table)

	case *ast.Union:
		types := []ValidType{}

		for _, dataType := range typeExpr.ValidTypes {
			nextType := FromAst(dataType, table)
			if nextType.String() == "TypeError" {
				return nextType
			}
			types = append(types, nextType)
		}

		return MakeUnion(types...)

	case *ast.ListType:
		dataType := FromAst(typeExpr.ElementType, table)
		if dataType.String() == "TypeError" {
			return dataType
		}

		return &ListLiteral{
			ElemType: dataType,
		}

	case *ast.ArrayType:
		dataType := FromAst(typeExpr.ElementType, table)
		if dataType.String() == "TypeError" {
			return dataType
		}

		return &ArrayLiteral{
			ElemType: dataType,
			Length:   typeExpr.Length,
		}

	case *ast.MapType:
		keyType := FromAst(typeExpr.KeyType, table)
		if keyType.String() == "TypeError" {
			return keyType
		}

		valueType := FromAst(typeExpr.ValueType, table)
		if valueType.String() == "TypeError" {
			return valueType
		}

		return &MapLiteral{
			KeyType:   keyType,
			ValueType: valueType,
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

type TypeTable interface {
	GetType(string) ValidType
}

func FromString(typeString string, table TypeTable) ValidType {
	dataType, ok := typeTable[typeString]
	if !ok {
		return table.GetType(typeString)
	}

	return dataType
}
