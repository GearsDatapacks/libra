package types

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
)

func FromAst(node ast.TypeExpression) ValidType {
	switch typeExpr := node.(type) {
	case *ast.TypeName:
		return MakeLiteral(FromString(typeExpr.Name))
	
	case *ast.Union:
		types := []ValidType{}

		for _, dataType := range typeExpr.ValidTypes {
			types = append(types, FromAst(dataType))
		}

		return MakeUnion(types...)
	
	case *ast.InferType:
		errors.TypeError("Expected type, got nothing", node)
		return &Literal{}

	default:
		errors.DevError("Unexpected type node: " + node.String())
		return &Literal{}
	}
}

var typeTable = map[string]DataType{
	"int":      INT,
	"float":    FLOAT,
	"boolean":  BOOL,
	"null":     NULL,
	"function": FUNCTION,
}

func FromString(typeString string) DataType {
	dataType, ok := typeTable[typeString]
	if !ok {
		errors.TypeError(fmt.Sprintf("Invalid type %q", typeString))
	}

	return dataType
}
