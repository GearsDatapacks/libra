package ast

import (
	"bytes"

	"github.com/gearsdatapacks/libra/lexer/token"
)

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

type Union struct {
	typeExpression
	Types []TypeExpression
}

func (u *Union) Tokens() []token.Token {
	tokens := []token.Token{}

	for _, ty := range u.Types {
		tokens = append(tokens, ty.Tokens()...)
	}

	return tokens
}

func (u *Union) String() string {
	var result bytes.Buffer

	for i, ty := range u.Types {
		if i != 0 {
			result.WriteString(" | ")
		}
		result.WriteString(ty.String())
	}

	return result.String()
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
