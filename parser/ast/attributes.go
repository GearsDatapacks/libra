package ast

import "github.com/gearsdatapacks/libra/lexer/token"

type Attribute interface {
	GetToken() token.Token
}

type TagAttribute struct {
	Token token.Token
	Name string
}

func (t *TagAttribute) GetToken() token.Token {
	return t.Token
}

// TODO: impl blocks
type ImplAttribute struct {
	Token token.Token
	Name string
}

func (i *ImplAttribute) GetToken() token.Token {
	return i.Token
}

type UntaggedAttribute struct {
	Token token.Token
}

func (u *UntaggedAttribute) GetToken() token.Token {
	return u.Token
}

type TodoAttribute struct {
	Token token.Token
	Message string
}

func (t *TodoAttribute) GetToken() token.Token {
	return t.Token
}

type DocAttribute struct {
	Token token.Token
	Message string
}

func (d *DocAttribute) GetToken() token.Token {
	return d.Token
}

type DeprecatedAttribute struct {
	Token token.Token
	Message string
}

func (d *DeprecatedAttribute) GetToken() token.Token {
	return d.Token
}
