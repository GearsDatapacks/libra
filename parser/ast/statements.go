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

type BlockStatement struct {
	statement
	LeftBrace token.Token
	Statements []Statement
	RightBrace token.Token
}

func (b *BlockStatement) Tokens() []token.Token {
	tokens := []token.Token{b.LeftBrace}
	
	for _, stmt := range b.Statements {
		tokens = append(tokens, stmt.Tokens()...)
	}

	tokens = append(tokens, b.RightBrace)
	return tokens
}

func (b *BlockStatement) String() string {
	var result bytes.Buffer

	result.WriteByte('{')
	for _, stmt := range b.Statements {
		result.WriteByte('\n')
		result.WriteString(stmt.String())
	}
	result.WriteString("\n}")

	return result.String()
}

type IfStatement struct {
	statement
	Keyword token.Token
	Condition Expression
	Body *BlockStatement
	ElseBranch *ElseBranch
}

func (is *IfStatement) Tokens() []token.Token {
	tokens := []token.Token{is.Keyword}
	tokens = append(tokens, is.Body.Tokens()...)
	tokens = append(tokens, is.Condition.Tokens()...)

	if is.ElseBranch != nil {
		tokens = append(tokens, is.ElseBranch.ElseKeyword)
		tokens = append(tokens, is.ElseBranch.Statement.Tokens()...)
	}

	return tokens
}

func (is *IfStatement) String() string {
	var result bytes.Buffer

	result.WriteString("if ")
	result.WriteString(is.Condition.String())
	result.WriteByte(' ')
	result.WriteString(is.Body.String())

	if is.ElseBranch != nil {
		result.WriteString(" else ")
		result.WriteString(is.ElseBranch.Statement.String())
	}

	return result.String()
}

type ElseBranch struct {
	ElseKeyword token.Token
	Statement Statement
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
