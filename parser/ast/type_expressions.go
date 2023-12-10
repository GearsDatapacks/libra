package ast

import (
	"fmt"
	"strings"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type BaseType struct{}

func (bt *BaseType) typeNode() {}

type TypeName struct {
	BaseNode
	BaseType
	Name string
}

func (tn *TypeName) Type() NodeType {
	return "TypeName"
}

func (tn *TypeName) String() string {
	return tn.Name
}

type Union struct {
	BaseNode
	BaseType
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

type ListType struct {
	BaseNode
	BaseType
	ElementType TypeExpression
}

func (lt *ListType) Type() NodeType { return "ListType" }

func (lt *ListType) String() string {
	if lt.ElementType.Type() == "Union" {
		return fmt.Sprintf("(%s)[]", lt.ElementType.String())
	}
	return lt.ElementType.String() + "[]"
}

type ArrayType struct {
	BaseNode
	BaseType
	ElementType TypeExpression
	Length      int
}

func (at *ArrayType) Type() NodeType { return "ArrayType" }

func (at *ArrayType) String() string {
	lengthString := "_"
	if at.Length != -1 {
		lengthString = fmt.Sprint(at.Length)
	}

	if at.ElementType.Type() == "Union" {
		return fmt.Sprintf("(%s)[%s]", at.ElementType.String(), lengthString)
	}
	return fmt.Sprintf("%s[%s]", at.ElementType.String(), lengthString)
}

type MapType struct {
	BaseNode
	BaseType
	KeyType   TypeExpression
	ValueType TypeExpression
}

func (*MapType) Type() NodeType { return "MapType" }

func (mt *MapType) String() string {
	return fmt.Sprintf("{%s: %s}", mt.KeyType.String(), mt.ValueType.String())
}

type InferType struct {
	BaseNode
	BaseType
}

func (i *InferType) Type() NodeType { return "Infer" }
func (i *InferType) String() string { return "" }

type VoidType struct {
	BaseNode
	BaseType
}

func (v *VoidType) Type() NodeType        { return "Void" }
func (v *VoidType) String() string        { return "void" }
func (v *VoidType) GetToken() token.Token { return token.Token{Value: "void"} }

type ErrorType struct {
	BaseNode
	BaseType
	ResultType TypeExpression
}

func (e *ErrorType) Type() NodeType { return "Error" }
func (e *ErrorType) String() string {
	if e.ResultType.Type() == "Union" {
		return fmt.Sprintf("(%s)!", e.ResultType.String())
	}
	return e.ResultType.String() + "!"
}

type TupleType struct {
	BaseNode
	BaseType
	Members []TypeExpression
}

func (*TupleType) Type() NodeType { return "Tuple" }
func (tuple *TupleType) String() string {
	result := "("

	for i, member := range tuple.Members {
		if i != 0 {
			result += ", "
		}
		result += member.String()
	}

	result += ")"

	return result
}

type MemberType struct {
	BaseNode
	BaseType
	Left TypeExpression
	Member string
}

func (*MemberType) Type() NodeType { return "MemberType" }
func (m *MemberType) String() string {
	return m.Left.String() + "." + m.Member
}
