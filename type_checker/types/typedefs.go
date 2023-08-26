package types

import (
	"strings"
)

type Union struct {
	Types []ValidType
}

func MakeUnion(types ...ValidType) *Union {
	return &Union{Types: types}
}

func (u *Union) Valid(dataType ValidType) bool {
	union, isUnion := dataType.(*Union)

	for _, unionType := range u.Types {
		if !isUnion && dataType.Valid(unionType) {
			return true
		}
		
		if isUnion && !union.Valid(unionType) {
			return false
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
	Parameters []ValidType
	ReturnType ValidType
}

func (fn *Function) Valid(dataType ValidType) bool {
	otherFn, isFn := dataType.(*Function)
	if !isFn { return false }

	return fn == otherFn
}

func (fn *Function) String() string {
	return "function"
}

type Any struct {}

func (a *Any) Valid(dataType ValidType) bool {
	_, isVoid := dataType.(*Void)
	return !isVoid
}

func (a *Any) String() string {
	return "any"
}
