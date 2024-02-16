package ast

import (
	"bytes"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type statement struct{}

func (statement) statementNode() {}

type ExpressionStatement struct {
	statement
	Expression Expression
}

func (es *ExpressionStatement) Tokens() []token.Token {
	return es.Expression.Tokens()
}
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

type TypeAnnotation struct {
	Colon token.Token
	Type  TypeExpression
}

func (ta *TypeAnnotation) Tokens() []token.Token {
	tokens := []token.Token{ta.Colon}
	tokens = append(tokens, ta.Type.Tokens()...)

	return tokens
}

func (ta *TypeAnnotation) String() string {
	var result bytes.Buffer

	result.WriteString(": ")
	result.WriteString(ta.Type.String())

	return result.String()
}

type VariableDeclaration struct {
	statement
	Keyword    token.Token
	Identifier token.Token
	Type       *TypeAnnotation
	Equals     token.Token
	Value      Expression
}

func (varDec *VariableDeclaration) Tokens() []token.Token {
	tokens := []token.Token{varDec.Keyword, varDec.Identifier}
	if varDec.Type != nil {
		tokens = append(tokens, varDec.Type.Tokens()...)
	}
	tokens = append(tokens, varDec.Equals)
	tokens = append(tokens, varDec.Value.Tokens()...)

	return tokens
}

func (varDec *VariableDeclaration) String() string {
	var result bytes.Buffer

	result.WriteString(varDec.Keyword.Value)
	result.WriteByte(' ')
	result.WriteString(varDec.Identifier.Value)

	if varDec.Type != nil {
		result.WriteString(varDec.Type.String())
	}

	result.WriteString(" = ")
	result.WriteString(varDec.Value.String())

	return result.String()
}

// TODO:
// Parameter
// FunctionDeclaration
// ReturnStatement
// IfStatement
// ElseStatement
// WhileLoop
// ForLoop
// StructField
// StructDeclaration
// TupleStructDeclaration
// UnitStructDeclaration
// InterfaceMember
// InterfaceDeclaration
// TypeDeclaration
// ImportStatement
// EnumDeclaration
// EnumMember
