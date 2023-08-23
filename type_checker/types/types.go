package types

type DataType = string

const (
	INT      = "int"
	FLOAT    = "float"
	BOOL     = "boolean"
	NULL     = "null"
	FUNCTION = "function"
)

type ValidType interface {
	Valid(ValidType) bool
	valid(DataType) bool
	String() string
}
