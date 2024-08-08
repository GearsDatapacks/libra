package codegen

import "tinygo.org/x/go-llvm"

type table struct {
	parent *table
	values map[string]llvm.Value
	context *fnContext
}

type fnContext struct {
	blocks map[string]llvm.BasicBlock
}

func newTable() *table {
	return &table{
		values: map[string]llvm.Value{},
	}
}

func childTable(parent *table) *table {
	return &table{
		parent: parent,
		values: map[string]llvm.Value{},
	}
}

func (t *table) addValue(name string, value llvm.Value) {
	t.values[name] = value
}

func (t *table) getValue(name string) llvm.Value {
	if value, ok := t.values[name]; ok {
		return value
	}
	if t.parent != nil {
		return t.parent.getValue(name)
	}
	panic("Should find value")
}
