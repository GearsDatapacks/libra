package ast

import (
	"bytes"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type expression struct{}

func (expression) expressionNode() {}

type IntegerLiteral struct {
	expression
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) Tokens() []token.Token {
	return []token.Token{il.Token}
}
func (il *IntegerLiteral) String() string {
	return il.Token.Value
}

type FloatLiteral struct {
	expression
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) Tokens() []token.Token {
	return []token.Token{fl.Token}
}
func (fl *FloatLiteral) String() string {
	return fl.Token.Value
}

type BooleanLiteral struct {
	expression
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) Tokens() []token.Token {
	return []token.Token{bl.Token}
}
func (bl *BooleanLiteral) String() string {
	return bl.Token.Value
}

type StringLiteral struct {
	expression
	Token token.Token
	Value string
}

func (sl *StringLiteral) Tokens() []token.Token {
	return []token.Token{sl.Token}
}
func (sl *StringLiteral) String() string {
	return sl.Token.Value
}

type Identifier struct {
	expression
	Token token.Token
	Name  string
}

func (i *Identifier) Tokens() []token.Token {
	return []token.Token{i.Token}
}
func (i *Identifier) String() string {
	return i.Name
}

type ErrorExpression struct {
	expression
}

func (e *ErrorExpression) Tokens() []token.Token {
	return []token.Token{}
}
func (e *ErrorExpression) String() string {
	return ""
}

type BinaryExpression struct {
	expression
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (b *BinaryExpression) Tokens() []token.Token {
	tokens := []token.Token{}
	tokens = append(tokens, b.Left.Tokens()...)
	tokens = append(tokens, b.Operator)
	tokens = append(tokens, b.Right.Tokens()...)
	return tokens
}

func (b *BinaryExpression) String() string {
	var result bytes.Buffer

	result.WriteString(b.Left.String())
	result.WriteByte(' ')
	result.WriteString(b.Operator.Value)
	result.WriteByte(' ')
	result.WriteString(b.Right.String())

	return result.String()
}

func (b *BinaryExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteByte('(')

	result.WriteString(maybePrecedence(b.Left))

	result.WriteByte(' ')
	result.WriteString(b.Operator.Value)
	result.WriteByte(' ')

	result.WriteString(maybePrecedence(b.Right))

	result.WriteByte(')')

	return result.String()
}

type ParenthesisedExpression struct {
	expression
	LeftParen  token.Token
	Expression Expression
	RightParen token.Token
}

func (p *ParenthesisedExpression) Tokens() []token.Token {
	tokens := []token.Token{p.LeftParen}
	tokens = append(tokens, p.Expression.Tokens()...)
	tokens = append(tokens, p.RightParen)
	return tokens
}

func (p *ParenthesisedExpression) String() string {
	var result bytes.Buffer

	result.WriteByte('(')
	result.WriteString(p.Expression.String())
	result.WriteByte(')')

	return result.String()
}

type PrefixExpression struct {
	expression
	Operator token.Token
	Operand  Expression
}

func (p *PrefixExpression) Tokens() []token.Token {
	return append([]token.Token{p.Operator}, p.Operand.Tokens()...)
}

func (p *PrefixExpression) String() string {
	return p.Operator.Value + p.Operand.String()
}

func (p *PrefixExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteString(p.Operator.Value)
	result.WriteByte('(')
	result.WriteString(maybePrecedence(p.Operand))
	result.WriteByte(')')

	return result.String()
}

type PostfixExpression struct {
	expression
	Operand  Expression
	Operator token.Token
}

func (p *PostfixExpression) Tokens() []token.Token {
	return append(p.Operand.Tokens(), p.Operator)
}

func (p *PostfixExpression) String() string {
	return p.Operand.String() + p.Operator.Value
}

func (p *PostfixExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteByte('(')
	result.WriteString(maybePrecedence(p.Operand))
	result.WriteByte(')')
	result.WriteString(p.Operator.Value)

	return result.String()
}

// We don't store the tokens of the commas because they probably won't be needed
type ListLiteral struct {
	expression
	LeftSquare  token.Token
	Values      []Expression
	RightSquare token.Token
}

func (l *ListLiteral) Tokens() []token.Token {
	tokens := []token.Token{l.LeftSquare}
	for _, value := range l.Values {
		tokens = append(tokens, value.Tokens()...)
	}
	tokens = append(tokens, l.RightSquare)
	return tokens
}

func (l *ListLiteral) String() string {
	var result bytes.Buffer

	result.WriteByte('[')
	for i, value := range l.Values {
		if i != 0 {
			result.WriteString(", ")
		}

		result.WriteString(value.String())
	}
	result.WriteByte(']')

	return result.String()
}

type KeyValue struct {
	Key   Expression
	Colon token.Token
	Value Expression
}

func (kv *KeyValue) Tokens() []token.Token {
	tokens := kv.Key.Tokens()
	tokens = append(tokens, kv.Colon)
	tokens = append(tokens, kv.Value.Tokens()...)

	return tokens
}

func (kv *KeyValue) String() string {
	var result bytes.Buffer

	result.WriteString(kv.Key.String())
	result.WriteString(": ")
	result.WriteString(kv.Value.String())

	return result.String()
}

type MapLiteral struct {
	expression
	LeftBrace  token.Token
	KeyValues  []KeyValue
	RightBrace token.Token
}

func (m *MapLiteral) Tokens() []token.Token {
	tokens := []token.Token{m.LeftBrace}

	for _, kv := range m.KeyValues {
		tokens = append(tokens, kv.Tokens()...)
	}

	tokens = append(tokens, m.RightBrace)

	return tokens
}

func (m *MapLiteral) String() string {
	var result bytes.Buffer

	result.WriteByte('{')

	for i, kv := range m.KeyValues {
		if i != 0 {
			result.WriteString(", ")
		}

		result.WriteString(kv.String())
	}

	result.WriteByte('}')

	return result.String()
}

type FunctionCall struct {
	expression
	Callee     Expression
	LeftParen  token.Token
	Arguments  []Expression
	RightParen token.Token
}

func (call *FunctionCall) Tokens() []token.Token {
	tokens := append(call.Callee.Tokens(), call.LeftParen)

	for _, arg := range call.Arguments {
		tokens = append(tokens, arg.Tokens()...)
	}

	tokens = append(tokens, call.RightParen)

	return tokens
}

func (call *FunctionCall) String() string {
	var result bytes.Buffer

	result.WriteString(call.Callee.String())
	result.WriteByte('(')

	for i, arg := range call.Arguments {
		if i != 0 {
			result.WriteString(", ")
		}

		result.WriteString(arg.String())
	}

	result.WriteByte(')')

	return result.String()
}

type IndexExpression struct {
	expression
	Left        Expression
	LeftSquare  token.Token
	Index       Expression
	RightSquare token.Token
}

func (index *IndexExpression) Tokens() []token.Token {
	tokens := append(index.Left.Tokens(), index.LeftSquare)

	tokens = append(tokens, index.Index.Tokens()...)

	tokens = append(tokens, index.RightSquare)

	return tokens
}

func (index *IndexExpression) String() string {
	var result bytes.Buffer

	result.WriteString(index.Left.String())
	result.WriteByte('[')

	result.WriteString(index.Index.String())

	result.WriteByte(']')

	return result.String()
}

type AssignmentExpression struct {
	expression
	Assignee Expression
	Operator token.Token
	Value    Expression
}

func (a *AssignmentExpression) Tokens() []token.Token {
	tokens := append(a.Assignee.Tokens(), a.Operator)
	return append(tokens, a.Value.Tokens()...)
}

func (a *AssignmentExpression) String() string {
	var result bytes.Buffer

	result.WriteString(a.Assignee.String())
	result.WriteByte(' ')
	result.WriteString(a.Operator.Value)
	result.WriteByte(' ')
	result.WriteString(a.Value.String())

	return result.String()
}

func (a *AssignmentExpression) PrecedenceString() string {
	var result bytes.Buffer

	result.WriteByte('(')
	result.WriteString(maybePrecedence(a.Assignee))
	result.WriteByte(' ')
	result.WriteString(a.Operator.Value)
	result.WriteByte(' ')
	result.WriteString(maybePrecedence(a.Value))
	result.WriteByte(')')

	return result.String()
}

type TupleExpression struct {
	expression
	LeftParen  token.Token
	Values     []Expression
	RightParen token.Token
}

func (t *TupleExpression) Tokens() []token.Token {
	tokens := []token.Token{t.LeftParen}

	for _, value := range t.Values {
		tokens = append(tokens, value.Tokens()...)
	}

	tokens = append(tokens, t.RightParen)

	return tokens
}

func (t *TupleExpression) String() string {
	var result bytes.Buffer

	result.WriteByte('(')

	for i, value := range t.Values {
		if i != 0 {
			result.WriteString(", ")
		}
		result.WriteString(value.String())
	}

	result.WriteByte(')')

	return result.String()
}

type MemberExpression struct {
	expression
	Left   Expression
	Dot    token.Token
	Member token.Token
}

func (m *MemberExpression) Tokens() []token.Token {
	return append(m.Left.Tokens(), m.Dot, m.Member)
}

func (m *MemberExpression) String() string {
	var result bytes.Buffer

	result.WriteString(m.Left.String())
	result.WriteByte('.')
	result.WriteString(m.Member.Value)

	return result.String()
}

type StructMember struct {
	Name  token.Token
	Colon token.Token
	Value Expression
}

func (sm *StructMember) Tokens() []token.Token {
	return append([]token.Token{sm.Name, sm.Colon}, sm.Value.Tokens()...)
}

func (sm *StructMember) String() string {
	var result bytes.Buffer

	result.WriteString(sm.Name.Value)
	result.WriteString(": ")
	result.WriteString(sm.Value.String())

	return result.String()
}

type StructExpression struct {
	expression
	Struct     Expression
	LeftBrace  token.Token
	Members    []StructMember
	RightBrace token.Token
}

func (s *StructExpression) Tokens() []token.Token {
	tokens := append(s.Struct.Tokens(), s.LeftBrace)

	for _, member := range s.Members {
		tokens = append(tokens, member.Tokens()...)
	}

	return append(tokens, s.RightBrace)
}

func (s *StructExpression) String() string {
	var result bytes.Buffer

	result.WriteString(s.Struct.String())
	result.WriteString(" { ")

	for i, member := range s.Members {
		if i != 0 {
			result.WriteString(", ")
		}

		result.WriteString(member.String())
	}

	result.WriteString(" }")

	return result.String()
}

// TODO:
// TypeCheckExpression
// CastExpression

type HasPrecedence interface {
	Expression
	PrecedenceString() string
}

func maybePrecedence(expr Expression) string {
	if prec, ok := expr.(HasPrecedence); ok {
		return prec.PrecedenceString()
	}

	return expr.String()
}
