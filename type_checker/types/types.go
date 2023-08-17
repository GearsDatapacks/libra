package types

type DataType int

const (
	INT = iota
	BOOL
	NULL

	INVALID = -1
)

var typeTable = map[string]DataType{
	"int":     INT,
	"boolean": BOOL,
	"null":    NULL,
}

func FromString(typeString string) DataType {
	dataType, ok := typeTable[typeString]
	if !ok {
		return INVALID
	}

	return dataType
}
