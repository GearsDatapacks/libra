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

// TODO:
// BlockStatement
// IfStatement
// ElseBranch
// WhileLoop
// ForLoop
// FunctionDeclaration
// ReturnStatement
// TypeDeclaration
// StructDeclaration
// InterfaceDeclaration
// ImportStatement
// EnumDeclaration
