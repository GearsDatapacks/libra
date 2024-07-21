package values

import (
	"hash/fnv"
	"math"

	"github.com/gearsdatapacks/libra/colour"
	"github.com/gearsdatapacks/libra/printer"
)

type ConstValue interface {
	constVal()
	printer.Printable
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

func (i IntValue) Print(node *printer.Node) {
	node.Text(
		"%sINT_VALUE %s%d",
		node.Colour(colour.NodeName),
		node.Colour(colour.Literal),
		i.Value,
	)
}

type FloatValue struct {
	constValue
	Value float64
}

func (f FloatValue) Hash() uint64 {
	return math.Float64bits(f.Value)
}

func (f FloatValue) Print(node *printer.Node) {
	node.Text(
		"%sFLOAT_VALUE %s%f",
		node.Colour(colour.NodeName),
		node.Colour(colour.Literal),
		f.Value,
	)
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

func (b BoolValue) Print(node *printer.Node) {
	node.Text(
		"%sBOOL_VALUE %s%t",
		node.Colour(colour.NodeName),
		node.Colour(colour.Literal),
		b.Value,
	)
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

func (s StringValue) Print(node *printer.Node) {
	node.Text(
		"%sSTRING_VALUE %s%q",
		node.Colour(colour.NodeName),
		node.Colour(colour.Literal),
		s.Value,
	)
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

func (a ArrayValue) Print(node *printer.Node) {
	node.Text(
		"%sARRAY_VALUE",
		node.Colour(colour.NodeName),
	)

	printer.Nodes(node, a.Elements)
}

type TupleValue struct {
	constValue
	Values []ConstValue
}

func (t TupleValue) Index(index ConstValue) ConstValue {
	return t.Values[index.(IntValue).Value]
}

func (t TupleValue) Print(node *printer.Node) {
	node.Text(
		"%sTUPLE_VALUE",
		node.Colour(colour.NodeName),
	)

	printer.Nodes(node, t.Values)
}

type KeyValue struct {
	Key, Value ConstValue
}

func (kv KeyValue) Print(node *printer.Node) {
	node.
		Text("%sKEY_VALUE", node.Colour(colour.NodeName)).
		Node(kv.Key).
		Node(kv.Value)
}

type MapValue struct {
	constValue
	Values map[uint64]KeyValue
}

func (m MapValue) Index(index ConstValue) ConstValue {
	kv, ok := m.Values[index.Hash()]
	if ok {
		return kv.Value
	}
	return nil
}

func (m MapValue) Print(node *printer.Node) {
	node.Text(
		"%sMAP_VALUE",
		node.Colour(colour.NodeName),
	)

	for _, kv := range m.Values {
		node.Node(kv)
	}
}

type hasMembers interface {
	StaticMemberValue(string) ConstValue
}

type TypeValue struct {
	constValue
	Type printer.Printable // types.Type, but no import cycles :/
}

func (t TypeValue) Member(member string) ConstValue {
	if hasMembers, ok := t.Type.(hasMembers); ok {
		return hasMembers.StaticMemberValue(member)
	}
	return nil
}

func (t TypeValue) Print(node *printer.Node) {
	node.
		Text(
			"%sTYPE_VALUE",
			node.Colour(colour.NodeName),
		).
		Node(t.Type)
}

type StructValue struct {
	constValue
	Members map[string]ConstValue
}

func (s StructValue) Member(member string) ConstValue {
	return s.Members[member]
}

func (s StructValue) Print(node *printer.Node) {
	node.Text(
		"%sSTRUCT_VALUE",
		node.Colour(colour.NodeName),
	)

	for name, value := range s.Members {
		node.FakeNode(
			"%sSTRUCT_MEMBER %s%s",
			func(n *printer.Node) { n.Node(value) },
			node.Colour(colour.NodeName),
			node.Colour(colour.Name),
			name,
		)
	}
}

type ModuleValue struct {
	constValue
	Module interface {
		LookupExport(string) interface{ Value() ConstValue }
	}
}

func (m ModuleValue) Member(member string) ConstValue {
	export := m.Module.LookupExport(member)
	if export == nil {
		return nil
	}
	return export.Value()
}

func (ModuleValue) Print(node *printer.Node) {
	node.Text(
		"%sMODULE_VALUE",
		node.Colour(colour.NodeName),
	)
}

type UnitValue struct {
	constValue
	Name string
}

func (u UnitValue) Print(node *printer.Node) {
	node.Text(
		"%sUNIT_VALUE %s%s",
		node.Colour(colour.NodeName),
		node.Colour(colour.Name),
		u.Name,
	)
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
