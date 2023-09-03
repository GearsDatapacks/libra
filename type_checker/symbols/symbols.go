package symbols

import (
	"fmt"
	"log"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/utils"
)

type scopeKind int

const (
	GLOBAL_SCOPE = iota
	GENERIC_SCOPE
	FUNCTION_SCOPE
)

type SymbolTable struct {
	parent    *SymbolTable
	symbols   map[string]types.ValidType
	constants []string
	kind scopeKind
	returnType types.ValidType
}

func New() *SymbolTable {
	return &SymbolTable{
		parent:  nil,
		symbols: map[string]types.ValidType{},
		kind: GLOBAL_SCOPE,
	}
}

func NewChild(parent *SymbolTable, kind scopeKind) *SymbolTable {
	return &SymbolTable{
		parent:  parent,
		symbols: map[string]types.ValidType{},
		kind: kind,
	}
}

func NewFunction(parent *SymbolTable, returnType types.ValidType) *SymbolTable {
	table := NewChild(parent, FUNCTION_SCOPE)
	table.returnType = returnType
	return table
}

func (st *SymbolTable) RegisterSymbol(name string, dataType types.ValidType, constant bool) error {
	if _, ok := st.symbols[name]; ok {
		return errors.TypeError(fmt.Sprintf("Cannot redeclare variable %q, it is already defined", name))
	}

	if _, ok := registry.Builtins[name]; ok {
		return errors.TypeError(fmt.Sprintf("Cannot redifne builtin function %q", name))
	}

	if constant {
		st.constants = append(st.constants, name)
	}

	dataType.MarkVariable()
	st.symbols[name] = dataType
	return nil
}

func (st *SymbolTable) GetSymbol(name string) (types.ValidType, error) {
	table, err := st.resolve(name)

	if err != nil {
		return nil, err
	}

	return table.symbols[name], nil
}

func (st *SymbolTable) IsConstant(name string) bool {
	return utils.Contains(st.constants, name)
}

func (st *SymbolTable) Exists(name string) bool {
	_, err := st.resolve(name)

	return err == nil
}

func (st *SymbolTable) resolve(varName string) (*SymbolTable, error) {
	if _, ok := st.symbols[varName]; ok {
		return st, nil
	}

	if st.parent == nil {
		return nil, errors.TypeError(fmt.Sprintf("Variable %q is undefined", varName))
	}

	return st.parent.resolve(varName)
}

func (st *SymbolTable) isFunctionScope() bool {
	return st.kind == FUNCTION_SCOPE
}

func (st *SymbolTable) ReturnType() types.ValidType {
	if !st.isFunctionScope() {
		log.Fatal(errors.DevError("Cannot get return type of non-function scope"))
	}

	return st.returnType
}

func (st *SymbolTable) FindFunctionScope() *SymbolTable {
	if st.isFunctionScope() {
		return st
	}

	if st.parent == nil {
		return nil
	}

	return st.parent.FindFunctionScope()
}
