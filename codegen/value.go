package codegen

import (
	"tinygo.org/x/go-llvm"
)

// Represents a value to keep track of, usually a wrapper to an
// llvm.Value, but containing additional information such as whether
// it's really a pointer to an allocated stack variable
type value interface {
	toLlvm(*compiler) llvm.Value
}

type llvmValue llvm.Value

func (l llvmValue) toLlvm(*compiler) llvm.Value {
	return llvm.Value(l)
}

type stackVariable llvm.Value

func (s stackVariable) toLlvm(c *compiler) llvm.Value {
	return c.builder.CreateLoad(llvm.Value(s).AllocatedType(), llvm.Value(s), "load_tmp")
}
