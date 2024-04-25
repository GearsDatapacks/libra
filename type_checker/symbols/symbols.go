package symbols

import (
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

type Variable struct {
	Name string
	Mutable bool
	Type types.Type
	ConstValue values.ConstValue
}
