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
	result := bytes.NewBuffer([]byte{})

	result.WriteString(b.Left.String())
	result.WriteByte(' ')
	result.WriteString(b.Operator.Value)
	result.WriteByte(' ')
	result.WriteString(b.Right.String())

	return result.String()
}

func (b *BinaryExpression) PrecedenceString() string {
	result := bytes.NewBuffer([]byte{})

	result.WriteByte('(')

	if prec, ok := b.Left.(hasPrecedence); ok {
		result.WriteString(prec.PrecedenceString())
	} else {
		result.WriteString(b.Left.String())
	}

	result.WriteByte(' ')
	result.WriteString(b.Operator.Value)
	result.WriteByte(' ')

	if prec, ok := b.Right.(hasPrecedence); ok {
		result.WriteString(prec.PrecedenceString())
	} else {
		result.WriteString(b.Right.String())
	}

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
	result := bytes.NewBuffer([]byte{})

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
	result := bytes.NewBuffer([]byte{})

	result.WriteString(p.Operator.Value)
	result.WriteByte('(')
	if prec, ok := p.Operand.(hasPrecedence); ok {
		result.WriteString(prec.PrecedenceString())
	} else {
		result.WriteString(p.Operand.String())
	}
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
	result := bytes.NewBuffer([]byte{})

	result.WriteByte('(')
	if prec, ok := p.Operand.(hasPrecedence); ok {
		result.WriteString(prec.PrecedenceString())
	} else {
		result.WriteString(p.Operand.String())
	}
	result.WriteByte(')')
	result.WriteString(p.Operator.Value)

	return result.String()
}

type hasPrecedence interface {
	PrecedenceString() string
}

// TODO:
// ListLiteral
// MapLiteral
// FunctionCall
// AssignmentExpression
// IndexExpression
// MemberExpression
// StructExpression
// TupleExpression
// TypeCheckExpression
// CastExpression
