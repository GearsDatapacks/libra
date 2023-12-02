package ast

import (
	"fmt"
	"strings"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type BaseExpression struct{}

func (exp *BaseExpression) expressionNode() {}

type IntegerLiteral struct {
	*BaseNode
	*BaseExpression
	Value int
}

func (il *IntegerLiteral) Type() NodeType { return "Integer" }

func (il *IntegerLiteral) String() string {
	return il.Token.Value
}

type FloatLiteral struct {
	*BaseNode
	*BaseExpression
	Value float64
}

func (fl *FloatLiteral) Type() NodeType { return "Float" }

func (fl *FloatLiteral) String() string {
	return fl.Token.Value
}

type StringLiteral struct {
	*BaseNode
	*BaseExpression
	Value string
}

func (sl *StringLiteral) Type() NodeType { return "String" }

func (sl *StringLiteral) String() string {
	return "\"" + sl.Value + "\""
}

type BooleanLiteral struct {
	*BaseNode
	*BaseExpression
	Value bool
}

func (bl *BooleanLiteral) Type() NodeType { return "Boolean" }

func (bl *BooleanLiteral) String() string {
	return bl.Token.Value
}

type NullLiteral struct {
	*BaseNode
	*BaseExpression
}

func (nl *NullLiteral) Type() NodeType { return "Null" }

func (nl *NullLiteral) String() string { return "null" }

type VoidValue struct {
	*BaseNode
	*BaseExpression
}

func (nl *VoidValue) Type() NodeType       { return "Void" }
func (nl *VoidValue) String() string       { return "void" }
func (v *VoidValue) GetToken() token.Token { return token.Token{Value: "void"} }

type Identifier struct {
	*BaseNode
	*BaseExpression
	Symbol string
}

func (ident *Identifier) Type() NodeType { return "Identifier" }

func (ident *Identifier) String() string {
	return ident.Symbol
}

type ListLiteral struct {
	*BaseNode
	*BaseExpression
	Elements []Expression
}

func (list *ListLiteral) Type() NodeType { return "List" }

func (list *ListLiteral) String() string {
	result := "["
	valueStrings := []string{}

	for _, value := range list.Elements {
		valueStrings = append(valueStrings, value.String())
	}

	result += strings.Join(valueStrings, ", ")

	result += "]"
	return result
}

type MapLiteral struct {
	*BaseNode
	*BaseExpression
	Elements map[Expression]Expression
}

func (*MapLiteral) Type() NodeType { return "Map" }

func (maplit *MapLiteral) String() string {
	result := "{"
	valueStrings := []string{}

	for key, value := range maplit.Elements {
		valueStrings = append(valueStrings, key.String())
		valueStrings[len(valueStrings)-1] += ": "
		valueStrings[len(valueStrings)-1] += value.String()
	}

	result += strings.Join(valueStrings, ", ")

	result += "}"
	return result
}

type FunctionCall struct {
	*BaseNode
	*BaseExpression
	Left Expression
	Args []Expression
}

func (fn *FunctionCall) Type() NodeType { return "FunctionCall" }

func (fn *FunctionCall) String() string {
	result := fn.Left.String()

	result += "("

	for i, arg := range fn.Args {
		result += arg.String()
		if i != len(fn.Args)-1 {
			result += ", "
		}
	}

	result += ")"

	return result
}

type BinaryOperation struct {
	*BaseNode
	*BaseExpression
	Left     Expression
	Operator string
	Right    Expression
}

func (binOp *BinaryOperation) Type() NodeType { return "BinaryOperation" }

func (binOp *BinaryOperation) String() string {
	result := ""

	result += "("
	result += binOp.Left.String()
	result += " "
	result += binOp.Operator
	result += " "
	result += binOp.Right.String()
	result += ")"

	return result
}

type UnaryOperation struct {
	*BaseNode
	*BaseExpression
	Value    Expression
	Operator string
	Postfix  bool
}

func (unOp *UnaryOperation) Type() NodeType { return "UnaryOperation" }

func (unOp *UnaryOperation) String() string {
	if unOp.Postfix {
		return unOp.Value.String() + unOp.Operator
	}
	return unOp.Operator + unOp.Value.String()
}

type AssignmentExpression struct {
	*BaseNode
	*BaseExpression
	Assignee  Expression
	Value     Expression
	Operation string
}

func (ae *AssignmentExpression) Type() NodeType { return "AssignmentExpression" }

func (ae *AssignmentExpression) String() string {
	result := ""

	result += ae.Assignee.String()
	result += " "
	result += ae.Operation
	result += " "
	result += ae.Value.String()

	return result
}

type IndexExpression struct {
	*BaseNode
	*BaseExpression
	Left Expression
	Index Expression
}

func (index *IndexExpression) Type() NodeType { return "IndexExpression" }

func (index *IndexExpression) String() string {
	return fmt.Sprintf("%s[%s]", index.Left.String(), index.Index.String())
}

type MemberExpression struct {
	*BaseNode
	*BaseExpression
	Left Expression
	Member string
}

func (member *MemberExpression) Type() NodeType { return "MemberExpression" }

func (member *MemberExpression) String() string {
	return fmt.Sprintf("%s.%s", member.Left.String(), member.Member)
}

type StructExpression struct {
	*BaseNode
	*BaseExpression
	Name string
	Members map[string]Expression
}

func (tuple *StructExpression) Type() NodeType { return "StructExpression" }

func (tuple *TupleExpression) String() string {
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

type TupleExpression struct {
	*BaseNode
	*BaseExpression
	Members []Expression
}

func (structExpr *TupleExpression) Type() NodeType { return "TupleExpression" }

func (structExpr *StructExpression) String() string {
	result := structExpr.Name
	result += " {\n"

	for memberName, memberValue := range structExpr.Members {
		result += memberName
		result += ": "
		result += memberValue.String()
		result += ",\n"
	}

	result += "}"

	return result
}
