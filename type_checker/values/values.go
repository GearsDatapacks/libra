package values

import (
	"hash/fnv"
	"math"
)

type ConstValue interface {
	constVal()
	Hash() uint64
	Index(ConstValue) ConstValue
	Member(string) ConstValue
}
type constValue struct{}

func (constValue) constVal()                   {}
func (constValue) Index(ConstValue) ConstValue { return nil }
func (constValue) Member(string) ConstValue    { return nil }
func (constValue) Hash() uint64                { return 0 }

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

func (s StringValue) Index(index ConstValue) ConstValue {
	return StringValue{
		Value: string(s.Value[index.(IntValue).Value]),
	}
}

type ArrayValue struct {
	constValue
	Elements []ConstValue
}

func (a ArrayValue) Index(index ConstValue) ConstValue {
	return a.Elements[index.(IntValue).Value]
}

type TupleValue struct {
	constValue
	Values []ConstValue
}

func (t TupleValue) Index(index ConstValue) ConstValue {
	return t.Values[index.(IntValue).Value]
}

type MapValue struct {
	constValue
	Values map[uint64]ConstValue
}

func (m MapValue) Index(index ConstValue) ConstValue {
	return m.Values[index.Hash()]
}

type hasMembers interface {
	StaticMemberValue(string) ConstValue
}

type TypeValue struct {
	constValue
	Type any // types.Type, but no import cycles :/
}

func (t TypeValue) Member(member string) ConstValue {
	if hasMembers, ok := t.Type.(hasMembers); ok {
		return hasMembers.StaticMemberValue(member)
	}
	return nil
}

type StructValue struct {
	constValue
	Members map[string]ConstValue
}

func (s StructValue) Member(member string) ConstValue {
	return s.Members[member]
}

type ModuleValue struct {
	constValue
	Module interface{LookupExport(string) interface{Value() ConstValue}}
}

func (m ModuleValue) Member(member string) ConstValue {
	export := m.Module.LookupExport(member)
	if export == nil {
		return nil
	}
	return export.Value()
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
