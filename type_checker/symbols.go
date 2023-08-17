package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type SymbolTable struct {
	parent  *SymbolTable
	symbols map[string]types.DataType
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		parent:  nil,
		symbols: map[string]types.DataType{},
	}
}

func NewChildSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		parent:  parent,
		symbols: map[string]types.DataType{},
	}
}

func (st *SymbolTable) RegisterSymbol(name string, dataType types.DataType) {
	if _, ok := st.symbols[name]; ok {
		errors.RuntimeError(fmt.Sprintf("Cannot redeclare variable %q, it is already defined", name))
	}

	st.symbols[name] = dataType
}

func (st *SymbolTable) GetSymbol(name string) types.DataType {
	table := st.resolve(name)
	return table.symbols[name]
}

func (st *SymbolTable) resolve(varName string) *SymbolTable {
	if _, ok := st.symbols[varName]; ok {
		return st
	}

	if st.parent == nil {
		errors.RuntimeError(fmt.Sprintf("Variable %q is undefined", varName))
	}

	return st.parent.resolve(varName)
}
