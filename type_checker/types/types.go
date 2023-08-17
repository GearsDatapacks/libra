package types

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
)

type DataType string

const (
	INT = "int"
	BOOL = "boolean"
	NULL = "null"
)

var typeTable = map[string]DataType{
	"int":     INT,
	"boolean": BOOL,
	"null":    NULL,
}

func FromString(typeString string) DataType {
	dataType, ok := typeTable[typeString]
	if !ok {
		errors.TypeError(fmt.Sprintf("Invalid type %q", typeString))
	}

	return dataType
}
