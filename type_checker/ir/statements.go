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

type BreakStatement struct {
	Value Expression
}

func (b *BreakStatement) String() string {
	if b.Value == nil {
		return "break"
	}
	return fmt.Sprintf("break %s", b.Value.String())
}

type YieldStatement struct {
	Value Expression
}

func (y *YieldStatement) String() string {
	return fmt.Sprintf("yield %s", y.Value.String())
}

type ContinueStatement struct{}

func (*ContinueStatement) String() string {
	return "continue"
}

type ImportStatement struct {
	Module    string
	Name      string
	Symbols   []string
	ImportAll bool
}

func (i *ImportStatement) String() string {
	var result bytes.Buffer

	result.WriteString("import ")
	if i.ImportAll {
		result.WriteString("* from ")
	} else if i.Symbols != nil {
		result.WriteByte('{')
		for i, symbol := range i.Symbols {
			if i != 0 {
				result.WriteString(", ")
			}
			result.WriteString(symbol)
		}
		result.WriteString("} ")
	}

	result.WriteByte('"')
	result.WriteString(i.Module)
	result.WriteByte('"')

	if !i.ImportAll && i.Symbols == nil {
		result.WriteString(" as ")
		result.WriteString(i.Name)
	}

	return result.String()
}

// TODO:
// EnumDeclaration
// UnionDeclaration
