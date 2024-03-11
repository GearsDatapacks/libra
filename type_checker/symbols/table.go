package symbols

type Table struct {
	parent    *Table
	variables map[string]Variable
}

func New() *Table {
	return &Table{
		variables: map[string]Variable{},
	}
}

func (t *Table) Child() *Table {
	return &Table{
		parent:    t,
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
	if t.parent != nil {
		return t.parent.LookupVariable(name)
	}
	return nil
}
