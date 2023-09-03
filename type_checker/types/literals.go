package types

import "fmt"

func isA[T ValidType](v ValidType) bool {
	_, ok := v.(T)
	return ok
}

type IntLiteral struct{ BaseType }

func (i *IntLiteral) Valid(t ValidType) bool { return isA[*IntLiteral](t) }
func (i *IntLiteral) String() string         { return "int" }

type FloatLiteral struct{ BaseType }

func (f *FloatLiteral) String() string         { return "float" }
func (f *FloatLiteral) Valid(t ValidType) bool { return isA[*FloatLiteral](t) }

type BoolLiteral struct{ BaseType }

func (b *BoolLiteral) String() string         { return "boolean" }
func (b *BoolLiteral) Valid(t ValidType) bool { return isA[*BoolLiteral](t) }

type NullLiteral struct{ BaseType }

func (n *NullLiteral) String() string         { return "null" }
func (n *NullLiteral) Valid(t ValidType) bool { return isA[*NullLiteral](t) }

type StringLiteral struct{ BaseType }

func (s *StringLiteral) String() string         { return "string" }
func (s *StringLiteral) Valid(t ValidType) bool { return isA[*StringLiteral](t) }

type ListLiteral struct {
	BaseType
	ElemType ValidType
}

func (list *ListLiteral) String() string {
	if isA[*Union](list.ElemType) {
		return fmt.Sprintf("(%s)[]", list.ElemType.String())
	}
	return list.ElemType.String() + "[]"
}
func (list *ListLiteral) Valid(t ValidType) bool {
	if l, isList := t.(*ListLiteral); isList {
		return list.ElemType.Valid(l.ElemType)
	}
	if array, isArray := t.(*ArrayLiteral); isArray {
		return array.CanInfer && list.ElemType.Valid(array.ElemType)
	}
	return false
}

func (list *ListLiteral) Infer(dataType ValidType) (ValidType, bool) {
	if !list.Valid(dataType) {
		return list, false
	}

	if list.ElemType.String() != "Infer" {
		return list, true
	}
	
	if array, ok := dataType.(*ArrayLiteral); ok {
		return &ListLiteral{
			ElemType: array.ElemType,
		}, true
	}

	return dataType, true
}

type ArrayLiteral struct {
	BaseType
	ElemType ValidType
	Length   int
	CanInfer bool // For array literals to be type inferred
}

func (array *ArrayLiteral) String() string {
	length := "_"
	if array.Length != -1 {
		length = fmt.Sprint(array.Length)
	}
	if isA[*Union](array.ElemType) {
		return fmt.Sprintf("(%s)[%s]", array.ElemType.String(), length)
	}
	return fmt.Sprintf("%s[%s]", array.ElemType.String(), length)
}

func (array *ArrayLiteral) Valid(t ValidType) bool {
	if !isA[*ArrayLiteral](t) {
		return false
	}
	lengthsMatch := array.Length == -1 || array.Length == t.(*ArrayLiteral).Length
	return lengthsMatch && array.ElemType.Valid(t.(*ArrayLiteral).ElemType)
}

func (array *ArrayLiteral) Infer(dataType ValidType) (ValidType, bool) {
	if !array.Valid(dataType) {
		return array, false
	}

	if array.Length != -1 && array.ElemType.String() != "Infer" {
		return array, true
	}
 
	other := dataType.(*ArrayLiteral)
	var length int
	var elemType ValidType
	if array.Length != -1 {
		length = array.Length
	} else {
		length = other.Length
	}

	if array.ElemType.String() != "Infer" {
		elemType = array.ElemType
	} else {
		elemType = other.ElemType
	}

	return &ArrayLiteral{
		ElemType: elemType,
		Length:   length,
		CanInfer: false,
	}, true
}

type Void struct{ BaseType }

func (v *Void) String() string         { return "void" }
func (v *Void) Valid(t ValidType) bool { return isA[*Void](t) }

type Infer struct{ BaseType }

func (i *Infer) String() string         { return "Infer" }
func (i *Infer) Valid(t ValidType) bool { return true }
