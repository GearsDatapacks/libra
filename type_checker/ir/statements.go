package ir

type statement struct{}

func (statement) irStmt() {}

type ExpressionStatement struct {
	statement
	Expression Expression
}

func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}

// TODO:
// VariableDeclaration
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
