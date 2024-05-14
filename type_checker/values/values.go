package values

import (
	"hash/fnv"
	"math"
)

type ConstValue interface {
	constVal()
	Hash() uint64
}
type constValue struct{}

func (constValue) constVal() {}

type IntValue struct {
	constValue
	Value int64
}

func (i IntValue) Hash() uint64 {
	return uint64(i.Value)
}

type FloatValue struct {
	constValue
	Value float64
}

func (f FloatValue) Hash() uint64 {
	return math.Float64bits(f.Value)
}

type BoolValue struct {
	constValue
	Value bool
}

func (b BoolValue) Hash() uint64 {
	if b.Value {
		return 1
	}
	return 0
}

type StringValue struct {
	constValue
	Value string
}

func (s StringValue) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return h.Sum64()
}

type ArrayValue struct {
	constValue
	Elements []ConstValue
}

func (ArrayValue) Hash() uint64 {
	return 0
}

type TupleValue struct {
	constValue
	Values []ConstValue
}

func (TupleValue) Hash() uint64 {
	return 0
}

type MapValue struct {
	constValue
	Values map[uint64]ConstValue
}

func (MapValue) Hash() uint64 {
	return 0
}

type TypeValue struct {
	constValue
	Type any // types.Type, but no import cycles :/
}

func (TypeValue) Hash() uint64 {
	return 0
}

type StructValue struct {
	constValue
	Members map[string]ConstValue
}

func (StructValue) Hash() uint64 {
	return 0
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
