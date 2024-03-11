package symbols

import "github.com/gearsdatapacks/libra/type_checker/types"

type Variable struct {
	Name string
	Mutable bool
	Type types.Type
}
