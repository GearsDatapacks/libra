package values

type ConstValue interface {
	constVal()
}
type constValue struct{}

func (constValue) constVal() {}

type IntValue struct {
	constValue
	Value int64
}

type FloatValue struct {
	constValue
	Value float64
}

type BoolValue struct {
	constValue
	Value bool
}

type StringValue struct {
	constValue
	Value string
}

type ArrayValue struct {
	constValue
	Elements []ConstValue
}

func NumericValue(v ConstValue) float64 {
	switch val := v.(type) {
	case FloatValue:
		return val.Value
	case IntValue:
		return float64(val.Value)
	case BoolValue:
		if val.Value {
			return 1
		}
		return 0
	default:
		panic("Not a numeric value")
	}
}
