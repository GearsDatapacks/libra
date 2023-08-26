package types

func isA[T ValidType](v ValidType) bool {
	_, ok := v.(T)
	return ok
}

type IntLiteral struct {}
func (i *IntLiteral) Valid(t ValidType) bool { return isA[*IntLiteral](t) }
func (i *IntLiteral) String() string { return "int" }

type FloatLiteral struct {}
func (f *FloatLiteral) String() string { return "float" }
func (f *FloatLiteral) Valid(t ValidType) bool { return isA[*FloatLiteral](t) }

type BoolLiteral struct {}
func (b *BoolLiteral) String() string { return "boolean" }
func (b *BoolLiteral) Valid(t ValidType) bool { return isA[*BoolLiteral](t) }

type NullLiteral struct {}
func (n *NullLiteral) String() string { return "null" }
func (n *NullLiteral) Valid(t ValidType) bool { return isA[*NullLiteral](t) }

type StringLiteral struct {}
func (s *StringLiteral) String() string { return "string" }
func (s *StringLiteral) Valid(t ValidType) bool { return isA[*StringLiteral](t) }

type Void struct {}
func (v *Void) String() string { return "void" }
func (v *Void) Valid(t ValidType) bool { return isA[*Void](t) }
