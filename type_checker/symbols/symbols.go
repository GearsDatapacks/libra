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
	CONDITIONAL_SCOPE
	FALLBACK_SCOPE
)

type SymbolTable struct {
	Parent               *SymbolTable
	variables            map[string]types.ValidType
	types                map[string]types.ValidType
	kind                 scopeKind
	returnType           types.ValidType
	hasReturn            bool
	hasConditionalReturn bool
	Exports              map[string]types.ValidType
}

func New() *SymbolTable {
	return &SymbolTable{
		Parent:    nil,
		variables: map[string]types.ValidType{},
		types:     map[string]types.ValidType{},
		kind:      GLOBAL_SCOPE,
		Exports:   map[string]types.ValidType{},
	}
}

func NewChild(parent *SymbolTable, kind scopeKind) *SymbolTable {
	if kind == CONDITIONAL_SCOPE {
		parent.removeConditionalReturn()
	}
	return &SymbolTable{
		Parent:    parent,
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

	if st.Parent == nil {
		return nil
	}

	return st.Parent.resolveVariable(varName)
}

func (st *SymbolTable) ReturnType() types.ValidType {
	scope := st.FindFunctionScope()
	if scope == nil {
		log.Fatal(errors.DevError("Cannot get return type of non-function scope"))
	}

	return scope.returnType
}

func (st *SymbolTable) HasReturn() bool {
	scope := st.FindFunctionScope()
	if scope == nil {
		log.Fatal(errors.DevError("Cannot check return value of non-function scope"))
	}

	return scope.hasReturn
}

func (st *SymbolTable) AddReturn() {
	scope, conditional, fallback := st.findFunctionScope(false, false)
	if scope == nil {
		log.Fatal(errors.DevError("Cannot set return value of non-function scope"))
	}

	if !conditional && !fallback {
		scope.hasReturn = true
		return
	}

	if conditional {
		scope.hasConditionalReturn = true
	}

	if fallback && scope.hasConditionalReturn {
		scope.hasReturn = true
	}
}

func (st *SymbolTable) removeConditionalReturn() {
	scope := st.FindFunctionScope()
	if scope == nil {
		return
	}

	scope.hasConditionalReturn = false
}

func (st *SymbolTable) IsInFunctionScope() bool {
	return st.FindFunctionScope() != nil
}

func (st *SymbolTable) FindFunctionScope() *SymbolTable {
	scope, _, _ := st.findFunctionScope(false, false)
	return scope
}

func (st *SymbolTable) findFunctionScope(conditional, fallback bool) (table *SymbolTable, isConditional bool, isFallback bool) {
	if st.kind == FUNCTION_SCOPE {
		return st, conditional, fallback
	}

	if st.kind == CONDITIONAL_SCOPE {
		conditional = true
	}

	if st.kind == FALLBACK_SCOPE {
		fallback = true
	}

	if st.Parent == nil {
		return nil, false, false
	}

	return st.Parent.findFunctionScope(conditional, fallback)
}

func (st *SymbolTable) AddType(name string, dataType types.ValidType) *types.TypeError {
	_, hasType := st.types[name]
	if hasType {
		return types.Error(fmt.Sprintf("Cannot redeclare type %q", name))
	}

	st.types[name] = dataType
	return nil
}

func (st *SymbolTable) UpdateType(name string, dataType types.ValidType) *types.TypeError {
	_, hasType := st.types[name]
	if !hasType {
		errors.DevError(fmt.Sprintf("Cannot update type %q, it does not exist", name))
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
	if st.Parent == nil {
		return nil
	}
	return st.Parent.resolveType(name)
}

func (st *SymbolTable) GlobalScope() *SymbolTable {
	if st.Parent == nil {
		return st
	}

	return st.Parent.GlobalScope()
}
