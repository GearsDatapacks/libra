package types

import (
	"strings"
)

type Union struct {
	BaseType
	Types []ValidType
}

func MakeUnion(types ...ValidType) ValidType {
	if len(types) == 0 {
		return &Void{}
	}
	if len(types) == 1 {
		return types[0]
	}

	return &Union{Types: types}
}

func (u *Union) Valid(dataType ValidType) bool {
	union, isUnion := dataType.(*Union)

	// If it's a union, we want to make sure all possible values it could be are contained within this one
	if isUnion {
		for _, unionType := range union.Types {
			if !u.Valid(unionType) {
				return false
			}
		}

		return true
	}

	// Otherwise, we just make sure the value is contained within this one
	for _, unionType := range u.Types {
		if dataType.Valid(unionType) {
			return true
		}
	}

	return false
}

func (u *Union) String() string {
	typeStrings := []string{}

	for _, dataType := range u.Types {
		typeStrings = append(typeStrings, dataType.String())
	}

	return strings.Join(typeStrings, " | ")
}

type Function struct {
	BaseType
	Parameters []ValidType
	ReturnType ValidType
}

func (fn *Function) Valid(dataType ValidType) bool {
	otherFn, isFn := dataType.(*Function)
	if !isFn {
		return false
	}

	return fn == otherFn
}

func (fn *Function) String() string {
	return "function"
}

type Any struct{ BaseType }

func (a *Any) Valid(dataType ValidType) bool {
	_, isVoid := dataType.(*Void)
	return !isVoid
}

func (a *Any) String() string {
	return "any"
}
