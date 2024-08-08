package codegen

import (
	"fmt"

	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"tinygo.org/x/go-llvm"
)

type compiler struct {
	context llvm.Context
	mainModule,
	currentModule llvm.Module
	builder llvm.Builder
	table   *table
}

func Compile(pkg *ir.LoweredPackage) llvm.MemoryBuffer {
	context := llvm.NewContext()
	compiler := &compiler{
		context:    context,
		mainModule: context.NewModule("main"),
		builder:    context.NewBuilder(),
		table: newTable(),
	}

	for _, mod := range pkg.Modules {
		compiler.currentModule = compiler.context.NewModule(mod.Name)
		// TODO: Codegen globals and types

		for _, fn := range mod.Functions {
			compiler.compileFn(fn)
		}

		err := llvm.LinkModules(compiler.mainModule, compiler.currentModule)
		if err != nil {
			// TODO: Don't crash here
			panic(err)
		}
	}

	fmt.Println(compiler.mainModule.String())
	panic("NO")
}

// TODO: Forward-declare functions
func (c *compiler) compileFn(fn *ir.FunctionDeclaration) {
	paramTypes := make([]llvm.Type, 0, len(fn.Parameters))
	for _, param := range fn.Type.Parameters {
		paramTypes = append(paramTypes, param.ToLlvm(c.context))
	}
	var retTy llvm.Type
	if fn.Type.ReturnType == types.Void {
		retTy = c.context.VoidType()
	} else {
		retTy = fn.Type.ReturnType.ToLlvm(c.context)
	}
	ty := llvm.FunctionType(retTy, paramTypes, false)
	function := llvm.AddFunction(c.currentModule, fn.Name, ty)
	// function.SetLinkage(llvm.ExternalLinkage)
	c.table = childTable(c.table)
	c.table.context = &fnContext{
		blocks: map[string]llvm.BasicBlock{},
	}
	for i, param := range function.Params() {
		param.SetName(fn.Parameters[i])
		c.table.addValue(fn.Parameters[i], param)
	}

	for _, stmt := range fn.Body.Statements {
		if label, ok := stmt.(*ir.Label); ok {
			block := c.context.AddBasicBlock(function, label.Name)
			c.table.context.blocks[label.Name] = block
		}
	}

	for _, stmt := range fn.Body.Statements {
		c.compileStatement(stmt)
	}
	c.table.addValue(fn.Name, function)
	// TODO: Don't crash here
	llvm.VerifyFunction(function, llvm.AbortProcessAction)
}

func (c *compiler) compileStatement(statement ir.Statement) {
	switch stmt := statement.(type) {
	case *ir.VariableDeclaration:
		value := c.compileExpression(stmt.Value)
		c.table.addValue(stmt.Symbol.Name, value)
	case *ir.ReturnStatement:
		if stmt.Value == nil {
			c.builder.CreateRetVoid()
		} else {
			value := c.compileExpression(stmt.Value)
			c.builder.CreateRet(value)
		}

	case *ir.Label:
		block := c.table.context.blocks[stmt.Name]
		c.builder.SetInsertPointAtEnd(block)
	case *ir.Goto:
		c.builder.CreateBr(c.table.context.blocks[stmt.Label])
	case *ir.Branch:
		cond := c.compileExpression(stmt.Condition)
		c.builder.CreateCondBr(
			cond,
			c.table.context.blocks[stmt.IfLabel],
			c.table.context.blocks[stmt.ElseLabel],
		)
	case ir.Expression:
		c.compileExpression(stmt)
	default:
		panic("Unreachable")
	}
}

func (c *compiler) compileExpression(expression ir.Expression) llvm.Value {
	switch expr := expression.(type) {
	case *ir.ArrayExpression:
		panic("TODO")
	case *ir.Assignment:
		panic("TODO")
	case *ir.BinaryExpression:
		panic("TODO")
	case *ir.BooleanLiteral:
		var value uint64 = 0
		if expr.Value {
			value = 1
		}
		return llvm.ConstInt(c.context.Int1Type(), value, false)
	case *ir.Conversion:
		panic("TODO")
	case *ir.DerefExpression:
		panic("TODO")
	case *ir.FloatLiteral:
		return llvm.ConstFloat(c.context.DoubleType(), expr.Value)
	case *ir.FunctionCall:
		callee := c.compileExpression(expr.Function)
		args := make([]llvm.Value, 0, len(expr.Arguments))
		for _, arg := range expr.Arguments {
			args = append(args, c.compileExpression(arg))
		}
		return c.builder.CreateCall(callee.GlobalValueType(), callee, args, "call_tmp")
	case *ir.FunctionExpression:
		panic("TODO")
	case *ir.IndexExpression:
		panic("TODO")
	case *ir.IntegerLiteral:
		return llvm.ConstInt(c.context.Int32Type(), uint64(expr.Value), true)
	case *ir.MapExpression:
		panic("TODO")
	case *ir.MemberExpression:
		panic("TODO")
	case *ir.RefExpression:
		panic("TODO")
	case *ir.StringLiteral:
		panic("TODO")
	case *ir.StructExpression:
		panic("TODO")
	case *ir.TupleExpression:
		panic("TODO")
	case *ir.TupleStructExpression:
		panic("TODO")
	case *ir.TypeCheck:
		panic("TODO")
	case *ir.TypeExpression:
		panic("TODO")
	case ir.UintLiteral:
		return llvm.ConstInt(c.context.Int32Type(), uint64(expr.Value), true)
	case *ir.UnaryExpression:
		panic("TODO")
	case *ir.VariableExpression:
		return c.table.getValue(expr.Symbol.Name)
	default:
		panic("Unreachable")
	}
}
