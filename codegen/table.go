package codegen

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

type table struct {
	parent  *table
	values  map[string]value
	context *fnContext
}

type fnContext struct {
	blocks map[string]llvm.BasicBlock
}

func newTable() *table {
	return &table{
		values: map[string]value{},
	}
}

func childTable(parent *table) *table {
	return &table{
		parent: parent,
		values: map[string]value{},
	}
}

func (t *table) addValue(name string, value value) {
	t.values[name] = value
}

func (t *table) getValue(name string) value {
	if value, ok := t.values[name]; ok {
		return value
	}
	if t.parent != nil {
		return t.parent.getValue(name)
	}
	panic(fmt.Sprintf("Should find value %s", name))
}
