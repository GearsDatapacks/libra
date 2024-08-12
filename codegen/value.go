package codegen

import (
	"tinygo.org/x/go-llvm"
)

// Represents a value to keep track of, usually a wrapper to an
// llvm.Value, but containing additional information such as whether
// it's really a pointer to an allocated stack variable
type value interface {
	toRValue(*compiler) llvm.Value
	toLValue() llvm.Value
}

type llvmValue llvm.Value

func (l llvmValue) toRValue(*compiler) llvm.Value {
	return llvm.Value(l)
}

func (l llvmValue) toLValue() llvm.Value {
	panic("Cannot be used as an l-value")
}

type stackVariable llvm.Value

func (s stackVariable) toRValue(c *compiler) llvm.Value {
	return c.builder.CreateLoad(llvm.Value(s).AllocatedType(), llvm.Value(s), "load_tmp")
}

func (s stackVariable) toLValue() llvm.Value {
	return llvm.Value(s)
}
