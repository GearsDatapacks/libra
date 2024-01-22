package ast

import "github.com/gearsdatapacks/libra/lexer/token"

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

// VariableDeclaration
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
