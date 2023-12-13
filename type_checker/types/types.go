package types

import (
	"log"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/parser/ast"
)

type ValidType interface {
	Valid(ValidType) bool
	String() string
	WasVariable() bool
	MarkVariable()
	Constant() bool
	MarkConstant()
	IndexBy(ValidType) ValidType
}

type hasMembers interface {
	member(string) ValidType
}

type hasNumberMembers interface {
	numberMember(string) ValidType
}

type BaseType struct {
	wasVariable bool
	constant    bool
}

func (b *BaseType) WasVariable() bool {
	return b.wasVariable
}

func (b *BaseType) MarkVariable() {
	b.wasVariable = true
}

func (b *BaseType) Constant() bool {
	return b.constant
}

func (b *BaseType) MarkConstant() {
	b.constant = true
}

func (*BaseType) IndexBy(ValidType) ValidType {
	return nil
}

type CustomCastable interface {
	CanCast(ValidType) bool
}

func CanCast(from, to ValidType) bool {
	if castable, ok := from.(CustomCastable); ok {
		return castable.CanCast(to)
	}

	return from.Valid(to)
}

type PartialType interface {
	ValidType
	Infer(ValidType) (ValidType, bool)
}

func FromAst(node ast.TypeExpression, table TypeTable) ValidType {
	switch typeExpr := node.(type) {
	case *ast.TypeName:
		dataType := FromString(typeExpr.Name, table)
		if err, isErr := dataType.(*TypeError); isErr {
			err.Line = node.GetToken().Line
			err.Column = node.GetToken().Column
		}
		return dataType

	case *ast.Union:
		types := []ValidType{}

		for _, dataType := range typeExpr.ValidTypes {
			nextType := FromAst(dataType, table)
			if nextType.String() == "TypeError" {
				return nextType
			}
			types = append(types, nextType)
		}

		return MakeUnion(types...)

	case *ast.ListType:
		dataType := FromAst(typeExpr.ElementType, table)
		if dataType.String() == "TypeError" {
			return dataType
		}

		return &ListLiteral{
			ElemType: dataType,
		}

	case *ast.ArrayType:
		dataType := FromAst(typeExpr.ElementType, table)
		if dataType.String() == "TypeError" {
			return dataType
		}

		return &ArrayLiteral{
			ElemType: dataType,
			Length:   typeExpr.Length,
		}

	case *ast.MapType:
		keyType := FromAst(typeExpr.KeyType, table)
		if keyType.String() == "TypeError" {
			return keyType
		}

		valueType := FromAst(typeExpr.ValueType, table)
		if valueType.String() == "TypeError" {
			return valueType
		}

		return &MapLiteral{
			KeyType:   keyType,
			ValueType: valueType,
		}

	case *ast.ErrorType:
		resultType := FromAst(typeExpr.ResultType, table)
		if resultType.String() == "TypeError" {
			return resultType
		}

		return &ErrorType{ResultType: resultType}
	
	case *ast.TupleType:
		members := []ValidType{}
		for _, member := range typeExpr.Members {
			resultType := FromAst(member, table)
			if resultType.String() == "TypeError" {
				return resultType
			}
			members = append(members, resultType)
		}

		return &Tuple{Members: members}

	case *ast.VoidType:
		return &Void{}

	case *ast.InferType:
		return &Infer{}

	default:
		log.Fatal(errors.DevError("Unexpected type node: " + node.String()))
		return nil
	}
}

var typeTable = map[string]ValidType{
	"int":      &IntLiteral{},
	"float":    &FloatLiteral{},
	"boolean":  &BoolLiteral{},
	"null":     &NullLiteral{},
	"function": &Function{},
	"string":   &StringLiteral{},
}

type TypeTable interface {
	GetType(string) ValidType
}

func FromString(typeString string, table TypeTable) ValidType {
	dataType, ok := typeTable[typeString]
	if !ok {
		return table.GetType(typeString)
	}

	return dataType
}

func Member(memberOf ValidType, name string, isNumberMember bool) ValidType {
	method := getMethod(memberOf, name)
	if method != nil {
		return method
	}

	if !isNumberMember {
		hasMembers, ok := memberOf.(hasMembers)
		if ok {
			return hasMembers.member(name)
		}
	} else {
		hasNumberMembers, ok := memberOf.(hasNumberMembers)
		if ok {
			return hasNumberMembers.numberMember(name)
		}
	}

	return nil
}

var methods = map[string][]*Function{}

func AddMethod(name string, method *Function) {
	overloads, ok := methods[name]
	if !ok {
		methods[name] = []*Function{method}
	}
	overloads = append(overloads, method)
	methods[name] = overloads
}

func getMethod(methodOf ValidType, name string) *Function {
	overloads, ok := methods[name]
	if !ok {
		return nil
	}

	for _, overload := range overloads {
		if overload.MethodOf.Valid(methodOf) {
			return overload
		}
	}

	return nil
}

type PseudoType interface {
	ToReal() ValidType
}
