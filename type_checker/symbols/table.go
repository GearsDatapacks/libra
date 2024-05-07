package symbols

type Table struct {
	Parent    *Table
	variables map[string]Variable
}

func New() *Table {
	return &Table{
		variables: map[string]Variable{},
	}
}

func (t *Table) Child() *Table {
	return &Table{
		Parent:    t,
		variables: map[string]Variable{},
	}
}

func (t *Table) DeclareVariable(variable Variable) bool {
	if _, exists := t.variables[variable.Name]; exists {
		return false
	}
	t.variables[variable.Name] = variable
	return true
}

func (t *Table) LookupVariable(name string) *Variable {
	variable, ok := t.variables[name]
	if ok {
		return &variable
	}
	if t.Parent != nil {
		return t.Parent.LookupVariable(name)
	}
	return nil
}
