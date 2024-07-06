package ast

import "github.com/gearsdatapacks/libra/lexer/token"

type Attribute interface {
	GetToken() token.Token
	GetName() string
}

type FlagAttribute struct {
	Token token.Token
}

func (f *FlagAttribute) GetToken() token.Token {
	return f.Token
}

func (f *FlagAttribute) GetName() string {
	return f.Token.Value[1:]
}

// TODO: impl blocks

type TextAttribute struct {
	Token token.Token
	Text  string
}

func (t *TextAttribute) GetToken() token.Token {
	return t.Token
}

func (t *TextAttribute) GetName() string {
	return t.Token.Value[1:]
}

type ExpressionAttribute struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionAttribute) GetToken() token.Token {
	return e.Token
}

func (e *ExpressionAttribute) GetName() string {
	return e.Token.Value[1:]
}
