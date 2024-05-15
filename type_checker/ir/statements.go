package ir

import (
	"bytes"
	"fmt"

	"github.com/gearsdatapacks/libra/type_checker/symbols"
)

type statement struct{}

func (statement) irStmt() {}

type ExpressionStatement struct {
	statement
	Expression Expression
}

func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}

type VariableDeclaration struct {
	statement
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

type Block struct {
	statement
	Statements []Statement
}

func (b *Block) String() string {
	var result bytes.Buffer

	result.WriteByte('{')
	if len(b.Statements) > 0 {
		result.WriteByte('\n')
	}
	for _, stmt := range b.Statements {
		result.WriteString(stmt.String())
		result.WriteByte('\n')
	}
	result.WriteByte('}')

	return result.String()
}

type IfStatement struct {
	statement
	Condition  Expression
	Body       *Block
	ElseBranch Statement
}

func (i *IfStatement) String() string {
	var result bytes.Buffer
	result.WriteString("if ")
	result.WriteString(i.Condition.String())
	result.WriteByte(' ')
	result.WriteString(i.Body.String())

	if i.ElseBranch != nil {
		result.WriteString("\nelse ")
		result.WriteString(i.ElseBranch.String())
	}

	return result.String()
}

type WhileLoop struct {
	statement
	Condition Expression
	Body      *Block
}

func (w *WhileLoop) String() string {
	var result bytes.Buffer
	result.WriteString("while ")
	result.WriteString(w.Condition.String())
	result.WriteByte(' ')
	result.WriteString(w.Body.String())

	return result.String()
}

type ForLoop struct {
	statement
	Variable symbols.Variable
	Iterator Expression
	Body     *Block
}

func (f *ForLoop) String() string {
	var result bytes.Buffer

	result.WriteString("for ")
	result.WriteString(f.Variable.Name)
	result.WriteString(" in ")
	result.WriteString(f.Iterator.String())
	result.WriteByte(' ')
	result.WriteString(f.Body.String())

	return result.String()
}

type FunctionDeclaration struct {
	statement
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
	statement
	Value Expression
}

func (r *ReturnStatement) String() string {
	if r.Value != nil {
		return fmt.Sprintf("return %s", r.Value.String())
	}
	return "return"
}

// TODO:
// ImportStatement
// EnumDeclaration
