package ast

import (
	"bytes"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/text"
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

func (tn *TypeName) Location() text.Location {
	return tn.Name.Location
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

func (u *Union) Location() text.Location {
	return u.Types[0].Location().To(u.Types[len(u.Types)-1].Location())
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
	Type        TypeExpression
	LeftSquare  token.Token
	Count       Expression
	RightSquare token.Token
}

func (a *ArrayType) Tokens() []token.Token {
	tokens := a.Type.Tokens()
	tokens = append(tokens, a.LeftSquare)
	if a.Count != nil {
		tokens = append(tokens, a.Count.Tokens()...)
	}
	tokens = append(tokens, a.RightSquare)

	return tokens
}

func (a *ArrayType) Location() text.Location {
	return a.Type.Location().To(a.RightSquare.Location)
}

func (a *ArrayType) String() string {
	var result bytes.Buffer

	result.WriteString(a.Type.String())

	result.WriteByte('[')
	if a.Count != nil {
		result.WriteString(a.Count.String())
	}
	result.WriteByte(']')

	return result.String()
}

type PointerType struct {
	typeExpression
	Star token.Token
	Mut  *token.Token
	Type TypeExpression
}

func (ptr *PointerType) Tokens() []token.Token {
	tokens := []token.Token{ptr.Star}
	if ptr.Mut != nil {
		tokens = append(tokens, *ptr.Mut)
	}
	tokens = append(tokens, ptr.Type.Tokens()...)

	return tokens
}

func (ptr *PointerType) Location() text.Location {
	return ptr.Star.Location.To(ptr.Type.Location())
}

func (ptr *PointerType) String() string {
	var result bytes.Buffer

	result.WriteByte('*')
	if ptr.Mut != nil {
		result.WriteString("mut ")
	}

	result.WriteString(ptr.Type.String())

	return result.String()
}

type ErrorType struct {
	typeExpression
	Type TypeExpression
	Bang token.Token
}

func (e *ErrorType) Tokens() []token.Token {
	tokens := []token.Token{}
	if e.Type != nil {
		tokens = append(tokens, e.Type.Tokens()...)
	}
	tokens = append(tokens, e.Bang)

	return tokens
}

func (e *ErrorType) Location() text.Location {
	return e.Type.Location().To(e.Bang.Location)
}

func (e *ErrorType) String() string {
	var result bytes.Buffer

	if e.Type != nil {
		result.WriteString(e.Type.String())
	}
	result.WriteByte('!')

	return result.String()
}

type OptionType struct {
	typeExpression
	Type     TypeExpression
	Question token.Token
}

func (o *OptionType) Tokens() []token.Token {
	tokens := o.Type.Tokens()
	tokens = append(tokens, o.Question)

	return tokens
}

func (o *OptionType) Location() text.Location {
	return o.Type.Location().To(o.Question.Location)
}

func (o *OptionType) String() string {
	var result bytes.Buffer

	result.WriteString(o.Type.String())
	result.WriteByte('?')

	return result.String()
}

type ParenthesisedType struct {
	typeExpression
	LeftParen  token.Token
	Type       TypeExpression
	RightParen token.Token
}

func (p *ParenthesisedType) Tokens() []token.Token {
	tokens := []token.Token{p.LeftParen}
	tokens = append(tokens, p.Type.Tokens()...)
	tokens = append(tokens, p.RightParen)
	return tokens
}

func (p *ParenthesisedType) Location() text.Location {
	return p.LeftParen.Location.To(p.RightParen.Location)
}

func (p *ParenthesisedType) String() string {
	var result bytes.Buffer

	result.WriteByte('(')
	result.WriteString(p.Type.String())
	result.WriteByte(')')

	return result.String()
}

type TupleType struct {
	typeExpression
	LeftParen  token.Token
	Types      []TypeExpression
	RightParen token.Token
}

func (t *TupleType) Tokens() []token.Token {
	tokens := []token.Token{t.LeftParen}

	for _, value := range t.Types {
		tokens = append(tokens, value.Tokens()...)
	}

	tokens = append(tokens, t.RightParen)

	return tokens
}

func (t *TupleType) Location() text.Location {
	return t.LeftParen.Location.To(t.RightParen.Location)
}

func (t *TupleType) String() string {
	var result bytes.Buffer

	result.WriteByte('(')

	for i, value := range t.Types {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(value.String())
	}

	result.WriteByte(')')

	return result.String()
}

type MapType struct {
	typeExpression
	LeftBrace  token.Token
	KeyType    TypeExpression
	Colon      token.Token
	ValueType  TypeExpression
	RightBrace token.Token
}

func (m *MapType) Tokens() []token.Token {
	tokens := []token.Token{m.LeftBrace}
	tokens = append(tokens, m.KeyType.Tokens()...)
	tokens = append(tokens, m.Colon)
	tokens = append(tokens, m.ValueType.Tokens()...)
	tokens = append(tokens, m.RightBrace)

	return tokens
}

func (m *MapType) Location() text.Location {
	return m.LeftBrace.Location.To(m.LeftBrace.Location)
}

func (m *MapType) String() string {
	var result bytes.Buffer

	result.WriteByte('{')
	result.WriteString(m.KeyType.String())
	result.WriteString(": ")
	result.WriteString(m.ValueType.String())
	result.WriteByte('}')

	return result.String()
}

type MemberType struct {
	typeExpression
	Left   TypeExpression
	Dot    token.Token
	Member token.Token
}

func (m *MemberType) Tokens() []token.Token {
	return append(m.Left.Tokens(), m.Dot, m.Member)
}

func (m *MemberType) Location() text.Location {
	return m.Left.Location().To(m.Member.Location)
}

func (m *MemberType) String() string {
	var result bytes.Buffer

	result.WriteString(m.Left.String())
	result.WriteByte('.')
	result.WriteString(m.Member.Value)

	return result.String()
}
