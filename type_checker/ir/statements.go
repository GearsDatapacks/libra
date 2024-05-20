package ir

import (
	"bytes"
	"fmt"
)

type VariableDeclaration struct {
	Name  string
	Value Expression
}

func (v *VariableDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString("let ")
	result.WriteString(v.Name)
	result.WriteString(" = ")
	result.WriteString(v.Value.String())

	return result.String()
}

type FunctionDeclaration struct {
	Name       string
	Parameters []string
	Body       *Block
}

func (f *FunctionDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString("fn ")
	result.WriteString(f.Name)
	result.WriteByte('(')

	for i, param := range f.Parameters {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(param)
	}

	result.WriteString(") ")
	result.WriteString(f.Body.String())

	return result.String()
}

type ReturnStatement struct {
	Value Expression
}

func (r *ReturnStatement) String() string {
	if r.Value != nil {
		return fmt.Sprintf("return %s", r.Value.String())
	}
	return "return"
}

type BreakStatement struct{}

func (*BreakStatement) String() string {
	return "break"
}

type ContinueStatement struct{}

func (*ContinueStatement) String() string {
	return "continue"
}

// TODO:
// ImportStatement
// EnumDeclaration
