package symbols

import (
	"fmt"
	"log"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/type_checker/registry"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type scopeKind int

const (
	GLOBAL_SCOPE = iota
	GENERIC_SCOPE
	FUNCTION_SCOPE
)

type SymbolTable struct {
	parent     *SymbolTable
	variables  map[string]types.ValidType
	types      map[string]types.ValidType
	kind       scopeKind
	returnType types.ValidType
}

func New() *SymbolTable {
	return &SymbolTable{
		parent:    nil,
		variables: map[string]types.ValidType{},
		types:     map[string]types.ValidType{},
		kind:      GLOBAL_SCOPE,
	}
}

func NewChild(parent *SymbolTable, kind scopeKind) *SymbolTable {
	return &SymbolTable{
		parent:    parent,
		variables: map[string]types.ValidType{},
		types:     map[string]types.ValidType{},
		kind:      kind,
	}
}

func NewFunction(parent *SymbolTable, returnType types.ValidType) *SymbolTable {
	table := NewChild(parent, FUNCTION_SCOPE)
	table.returnType = returnType
	return table
}

func (st *SymbolTable) RegisterSymbol(name string, dataType types.ValidType, constant bool) *types.TypeError {
	if _, ok := st.variables[name]; ok {
		return types.Error(fmt.Sprintf("Cannot redeclare variable %q, it is already defined", name))
	}

	if _, ok := registry.Builtins[name]; ok {
		return types.Error(fmt.Sprintf("Cannot redefine builtin function %q", name))
	}

	if array, isArray := dataType.(*types.ArrayLiteral); isArray {
		array.CanInfer = false
	}

	if constant {
		dataType.MarkConstant()
	}

	dataType.MarkVariable()
	st.variables[name] = dataType
	return nil
}

func (st *SymbolTable) GetSymbol(name string) types.ValidType {
	table := st.resolveVariable(name)

	if table == nil {
		return types.Error(fmt.Sprintf("Variable %q is undefined", name))
	}

	return table.variables[name]
}

func (st *SymbolTable) Exists(name string) bool {
	table := st.resolveVariable(name)

	return table != nil
}

func (st *SymbolTable) resolveVariable(varName string) *SymbolTable {
	if _, ok := st.variables[varName]; ok {
		return st
	}

	if st.parent == nil {
		return nil
	}

	return st.parent.resolveVariable(varName)
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

func (st *SymbolTable) AddType(name string, dataType types.ValidType) *types.TypeError {
	_, hasType := st.types[name]
	if hasType {
		return types.Error(fmt.Sprintf("Cannot redeclare type %q", name))
	}

	st.types[name] = dataType
	return nil
}

func (st *SymbolTable) GetType(name string) types.ValidType {
	table := st.resolveType(name)
	if table == nil {
		return types.Error(fmt.Sprintf("Type %q is undefind", name))
	}

	return table.types[name]
}

func (st *SymbolTable) resolveType(name string) *SymbolTable {
	if _, ok := st.types[name]; ok {
		return st
	}
	if st.parent == nil {
		return nil
	}
	return st.parent.resolveType(name)
}
