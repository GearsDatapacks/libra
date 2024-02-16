package ast

import "github.com/gearsdatapacks/libra/lexer/token"

type typeExpression struct{}

func (typeExpression) typeNode() {}

type TypeName struct {
	typeExpression
	Name token.Token
}

func (tn *TypeName) Tokens() []token.Token {
	return []token.Token{tn.Name}
}

func (tn *TypeName) String() string {
	return tn.Name.Value
}

// TODO:
// Union
// ListType
// ArrayType
// MapType
// InferType
// VoidType
// ErrorType
// TupleType
// MemberType
// PointerType
