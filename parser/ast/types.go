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

type ArrayType struct {
	typeExpression
	LeftSquare  token.Token
	Count       Expression
	RightSquare token.Token
	Type        TypeExpression
}

func (a *ArrayType) Tokens() []token.Token {
	tokens := []token.Token{a.LeftSquare}
	if a.Count != nil {
		tokens = append(tokens, a.Count.Tokens()...)
	}
	tokens = append(tokens, a.RightSquare)
	tokens = append(tokens, a.Type.Tokens()...)

	return tokens
}

func (a *ArrayType) String() string {
	var result bytes.Buffer

	result.WriteByte('[')
	if a.Count != nil {
		result.WriteString(a.Count.String())
	}

	result.WriteByte(']')
	result.WriteString(a.Type.String())

	return result.String()
}

// TODO:
// MapType
// InferType
// VoidType
// ErrorType
// TupleType
// MemberType
// PointerType
