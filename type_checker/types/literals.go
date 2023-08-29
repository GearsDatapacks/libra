package types

import "fmt"

func isA[T ValidType](v ValidType) bool {
	_, ok := v.(T)
	return ok
}

type IntLiteral struct{}

func (i *IntLiteral) Valid(t ValidType) bool { return isA[*IntLiteral](t) }
func (i *IntLiteral) String() string         { return "int" }

type FloatLiteral struct{}

func (f *FloatLiteral) String() string         { return "float" }
func (f *FloatLiteral) Valid(t ValidType) bool { return isA[*FloatLiteral](t) }

type BoolLiteral struct{}

func (b *BoolLiteral) String() string         { return "boolean" }
func (b *BoolLiteral) Valid(t ValidType) bool { return isA[*BoolLiteral](t) }

type NullLiteral struct{}

func (n *NullLiteral) String() string         { return "null" }
func (n *NullLiteral) Valid(t ValidType) bool { return isA[*NullLiteral](t) }

type StringLiteral struct{}

func (s *StringLiteral) String() string         { return "string" }
func (s *StringLiteral) Valid(t ValidType) bool { return isA[*StringLiteral](t) }

type ListLiteral struct {
	ElemType ValidType
}

func (list *ListLiteral) String() string {
	if isA[*Union](list.ElemType) {
		return fmt.Sprintf("(%s)[]", list.ElemType.String())
	}
	return list.ElemType.String() + "[]"
}
func (list *ListLiteral) Valid(t ValidType) bool {
	return isA[*ListLiteral](t) && list.ElemType.Valid(t.(*ListLiteral).ElemType)
}

type Void struct{}

func (v *Void) String() string         { return "void" }
func (v *Void) Valid(t ValidType) bool { return isA[*Void](t) }
