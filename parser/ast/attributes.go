package ast

import "github.com/gearsdatapacks/libra/text"

type Attribute interface {
	GetName() string
}

type FlagAttribute struct {
	Location text.Location
	Name     string
}

func (f *FlagAttribute) GetName() string {
	return f.Name
}

// TODO: impl blocks

type TextAttribute struct {
	Location text.Location
	Name     string
	Text     string
}

func (t *TextAttribute) GetName() string {
	return t.Name
}

type ExpressionAttribute struct {
	Location   text.Location
	Name       string
	Expression Expression
}

func (e *ExpressionAttribute) GetName() string {
	return e.Name
}
