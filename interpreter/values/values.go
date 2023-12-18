package values

import "github.com/gearsdatapacks/libra/type_checker/types"

type ValueType string

type RuntimeValue interface {
	// Type() ValueType
	ToString() string
	Truthy() bool
	EqualTo(RuntimeValue) bool
	Varname() string
	SetVarname(string)
	Index(RuntimeValue) RuntimeValue
	SetIndex(index RuntimeValue, value RuntimeValue) RuntimeValue
	Member(string) RuntimeValue
	SetMember(member string, value RuntimeValue) RuntimeValue
	Type() types.ValidType
}

type BaseValue struct {
	varname  string
	DataType types.ValidType
}

func (b *BaseValue) Varname() string {
	return b.varname
}

func (b *BaseValue) SetVarname(name string) {
	b.varname = name
}

func (b *BaseValue) Index(RuntimeValue) RuntimeValue {
	return MakeNull()
}

func (b *BaseValue) SetIndex(RuntimeValue, RuntimeValue) RuntimeValue {
	return MakeNull()
}

func (b *BaseValue) Member(string) RuntimeValue {
	return nil
}

func (b *BaseValue) SetMember(string, RuntimeValue) RuntimeValue {
	return MakeNull()
}

func (b *BaseValue) Type() types.ValidType {
	return b.DataType
}

type castable interface {
	castTo(types.ValidType) RuntimeValue
}

func Cast(value RuntimeValue, ty types.ValidType) RuntimeValue {
	if cast, ok := value.(castable); ok {
		return cast.castTo(ty)
	}
	return value
}

type AutoCastable interface {
	AutoCast(types.ValidType) RuntimeValue
}

func Expect(value RuntimeValue, ty types.ValidType) RuntimeValue {
	if cast, ok := value.(AutoCastable); ok {
		return cast.AutoCast(ty)
	}

	return value
}

// func MakeValue(v any) RuntimeValue {
// 	switch value := v.(type) {
// 	case int:
// 		return MakeInteger(value)
// 	case float64:
// 		return MakeFloat(value)
// 	case bool:
// 		return MakeBoolean(value)
// 	case string:
// 		return MakeString(value)
// 	case nil:
// 		return MakeNull()
// 	default:
// 		errors.DevError(fmt.Sprintf("Cannot create runtime value of type %T", v))
// 		return MakeInteger(0)
// 	}
// }

func GetZeroValue(dataType string) RuntimeValue {
	switch dataType {
	case "int":
		return MakeInteger(0)
	case "float":
		return MakeFloat(0)
	case "boolean":
		return MakeBoolean(false)
	case "null":
		return MakeNull()
	case "string":
		return MakeString("")
	default:
		// errors.LogError(errors.DevError(fmt.Sprintf("Cannot create runtime value of type %s", dataType)))
		return MakeNull()
	}
}
