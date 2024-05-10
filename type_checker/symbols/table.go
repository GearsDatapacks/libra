package symbols

type Table struct {
	Parent    *Table
	variables map[string]Symbol
	Context   any
}

func New() *Table {
	return &Table{
		variables: map[string]Symbol{},
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
