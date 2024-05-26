package symbols

import (
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type Table struct {
	Parent  *Table
	symbols map[string]Symbol
	Context any
}

func New() *Table {
	t := &Table{
		symbols: map[string]Symbol{},
		Context: &globalContext{
			methods:         map[string][]*Method{},
			exportedMethods: map[string][]*Method{},
			exports:         map[string]Symbol{},
		},
	}
	t.registerGlobals()
	return t
}

func (t *Table) Child() *Table {
	return &Table{
		Parent:  t,
		symbols: map[string]Symbol{},
	}
}

func (t *Table) ChildWithContext(context any) *Table {
	return &Table{
		Parent:  t,
		symbols: map[string]Symbol{},
		Context: context,
	}
}

func (t *Table) Register(symbol Symbol, exported ...bool) bool {
	if _, exists := t.symbols[symbol.GetName()]; exists {
		return false
	}
	t.symbols[symbol.GetName()] = symbol

	if len(exported) > 0 && exported[0] {
		context := t.Context.(*globalContext)
		context.exports[symbol.GetName()] = symbol
	}

	return true
}

func (t *Table) Lookup(name string) Symbol {
	symbol, ok := t.symbols[name]
	if ok {
		return symbol
	}
	if t.Parent != nil {
		return t.Parent.Lookup(name)
	}
	return nil
}

func (t *Table) LookupExport(name string) Symbol {
	symbol, ok := t.globalScope().Context.(*globalContext).exports[name]
	if ok {
		return symbol
	}
	return nil
}

func (t *Table) LookupExportType(name string) types.Type {
	symbol := t.LookupExport(name)
	if symbol != nil {
		return symbol.GetType()
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
	context := t.globalScope().Context.(*globalContext)
	methods, ok := context.methods[name]
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

func (t *Table) RegisterMethod(name string, method *Method, exported bool) {
	context := t.globalScope().Context.(*globalContext)
	methods, ok := context.methods[name]
	if !ok {
		context.methods[name] = []*Method{method}
	}
	context.methods[name] = append(methods, method)
	if exported {
		t.addExportedMethod(name, method)
	}
}

func (t *Table) addExportedMethod(name string, method *Method) {
	context := t.globalScope().Context.(*globalContext)
	methods, ok := context.exportedMethods[name]
	if !ok {
		context.exportedMethods[name] = []*Method{method}
	}
	context.exportedMethods[name] = append(methods, method)
}

func (t *Table) Extend(other *Table) {
	context := other.globalScope().Context.(*globalContext)
	for _, export := range context.exports {
		t.Register(export)
	}
	for name, methods := range context.exportedMethods {
		for _, method := range methods {
			t.RegisterMethod(name, method, false)
		}
	}
}

func (t *Table) registerGlobals() {
	t.Register(&Type{"i32", types.Int})
	t.Register(&Type{"f32", types.Float})
	t.Register(&Type{"bool", types.Bool})
	t.Register(&Type{"string", types.String})
	t.Register(&Type{"Type", types.RuntimeType})
}
