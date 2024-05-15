package symbols

import "github.com/gearsdatapacks/libra/type_checker/types"

type Table struct {
	Parent    *Table
	variables map[string]Symbol
	Context   any
}

func New() *Table {
	return &Table{
		variables: map[string]Symbol{},
		Context: &GlobalContext{
			Methods: map[string][]types.Method{},
		},
	}
}

func (t *Table) Child() *Table {
	return &Table{
		Parent:    t,
		variables: map[string]Symbol{},
	}
}

func (t *Table) ChildWithContext(context any) *Table {
	return &Table{
		Parent:    t,
		variables: map[string]Symbol{},
		Context:   context,
	}
}

func (t *Table) Register(symbol Symbol) bool {
	if _, exists := t.variables[symbol.GetName()]; exists {
		return false
	}
	t.variables[symbol.GetName()] = symbol
	return true
}

func (t *Table) Lookup(name string) Symbol {
	symbol, ok := t.variables[name]
	if ok {
		return symbol
	}
	if t.Parent != nil {
		return t.Parent.Lookup(name)
	}
	return nil
}

func (t *Table) globalScope() *Table {
	if t.Parent == nil {
		return t
	}
	return t.Parent.globalScope()
}

func (t *Table) LookupMethod(name string, methodOf types.Type, static bool) *types.Function {
	context := t.globalScope().Context.(*GlobalContext)
	methods, ok := context.Methods[name]
	if !ok {
		return nil
	}
	for _, method := range methods {
		if method.Static == static && types.Assignable(method.MethodOf, methodOf) {
			return method.Function
		}
	}
	return nil
}

func (t *Table) RegisterMethod(name string, method types.Method) {
	context := t.globalScope().Context.(*GlobalContext)
	methods, ok := context.Methods[name]
	if !ok {
		context.Methods[name] = []types.Method{method}
	}
	context.Methods[name] = append(methods, method)
}
