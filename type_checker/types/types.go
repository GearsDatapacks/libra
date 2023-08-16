package types

type Type int

const (
	INT = iota
	BOOL
	NULL

	INVALID = -1
)

var typeTable = map[string]Type{
	"int": INT,
	"boolean": BOOL,
	"null": NULL,
}

func FromString(typeString string) Type {
	dataType, ok := typeTable[typeString]
	if !ok {
		return INVALID
	}

	return dataType
}
