package ir

import "bytes"

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
	Name string
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
	Condition Expression
	Body *Block
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
	Body *Block
}

func (i *WhileLoop) String() string {
	var result bytes.Buffer
	result.WriteString("while ")
	result.WriteString(i.Condition.String())
	result.WriteByte(' ')
	result.WriteString(i.Body.String())

	return result.String()
}

// TODO:
// ForLoop
// FunctionDeclaration
// ReturnStatement
// TypeDeclaration
// StructDeclaration
// InterfaceDeclaration
// ImportStatement
// EnumDeclaration
