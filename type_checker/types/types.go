package types

type Type interface {
}

type PrimaryType int

const (
	_ PrimaryType = iota
	Int
)
