package ast

import "strings"

type BaseType struct {}

func (bt *BaseType) typeNode() {}

type TypeName struct {
	*BaseNode
	*BaseType
	Name string
}

func (tn *TypeName) Type() NodeType {
	return "TypeName"
}

func (tn *TypeName) String() string {
	return tn.Name
}

type Union struct {
	*BaseNode
	*BaseType
	ValidTypes []TypeExpression
}

func (u *Union) Type() NodeType {
	return "Union"
}

func (u *Union) String() string {
	stringTypes := []string{}

	for _, typeExpr := range u.ValidTypes {
		stringTypes = append(stringTypes, typeExpr.String())
	}

	return strings.Join(stringTypes, " | ")
}

type InferType struct {
	*BaseNode
	*BaseType
}

func (i *InferType) Type() NodeType { return "Infer" }
func (i *InferType) String() string { return "" }
