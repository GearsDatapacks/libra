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

func Compile(pkg *ir.LoweredPackage) llvm.Module {
	context := llvm.NewContext()
	compiler := &compiler{
		context:    context,
		mainModule: context.NewModule("main"),
		builder:    context.NewBuilder(),
		table:      newTable(),
	}

	for _, mod := range pkg.Modules {
		compiler.currentModule = compiler.context.NewModule(mod.Name)
		// TODO: Codegen globals and types

		for _, fn := range mod.Functions {
			compiler.registerFn(fn)
		}

		for _, fn := range mod.Functions {
			compiler.compileFn(fn)
		}

		err := llvm.LinkModules(compiler.mainModule, compiler.currentModule)
		if err != nil {
			// TODO: Don't crash here
			panic(err)
		}
	}

	return compiler.mainModule
}

func (c *compiler) registerFn(fn *ir.FunctionDeclaration) {
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
	name := fn.Name
	if fn.Extern != nil {
		name = *fn.Extern
	}
	function := llvm.AddFunction(c.currentModule, name, ty)
	// function.SetLinkage(llvm.ExternalLinkage)
	for i, param := range function.Params() {
		param.SetName(fn.Parameters[i])
	}
	c.table.addValue(fn.Name, llvmValue(function))
}

func (c *compiler) compileFn(fn *ir.FunctionDeclaration) {
	if fn.Extern != nil {
		return
	}

	function := c.table.getValue(fn.Name).toRValue(c)
	c.table = childTable(c.table)
	c.table.context = &fnContext{
		blocks: map[string]llvm.BasicBlock{},
	}

	for i, param := range function.Params() {
		c.table.addValue(fn.Parameters[i], llvmValue(param))
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

	c.table = c.table.parent
	// TODO: Don't crash here
	// llvm.VerifyFunction(function, llvm.AbortProcessAction)
}

func (c *compiler) compileStatement(statement ir.Statement) {
	switch stmt := statement.(type) {
	case *ir.VariableDeclaration:
		alloca := c.builder.CreateAlloca(stmt.Symbol.Type.ToLlvm(c.context), stmt.Symbol.Name)
		if stmt.Value != nil {
			value := c.compileExpression(stmt.Value, true).toRValue(c)
			c.builder.CreateStore(value, alloca)
		}

		c.table.addValue(stmt.Symbol.Name, stackVariable(alloca))
	case *ir.ReturnStatement:
		if stmt.Value == nil {
			c.builder.CreateRetVoid()
		} else {
			value := c.compileExpression(stmt.Value, true).toRValue(c)
			c.builder.CreateRet(value)
		}

	case *ir.Label:
		block := c.table.context.blocks[stmt.Name]
		c.builder.SetInsertPointAtEnd(block)
	case *ir.Goto:
		c.builder.CreateBr(c.table.context.blocks[stmt.Label])
	case *ir.Branch:
		cond := c.compileExpression(stmt.Condition, true).toRValue(c)
		c.builder.CreateCondBr(
			cond,
			c.table.context.blocks[stmt.IfLabel],
			c.table.context.blocks[stmt.ElseLabel],
		)
	case ir.Expression:
		c.compileExpression(stmt, false)
	default:
		panic("Unreachable")
	}
}

func (c *compiler) compileExpression(expression ir.Expression, used bool) value {
	switch expr := expression.(type) {
	case *ir.ArrayExpression:
		panic("TODO")
	case *ir.Assignment:
		lValue := c.compileExpression(expr.Assignee, true).toLValue()
		rValue := c.compileExpression(expr.Value, true).toRValue(c)
		c.builder.CreateStore(rValue, lValue)
		return llvmValue(rValue)

	case *ir.BinaryExpression:
		if !used {
			return llvmValue{}
		}
		return c.compileBinaryExpression(expr)
	case *ir.BooleanLiteral:
		if !used {
			return llvmValue{}
		}
		var value uint64 = 0
		if expr.Value {
			value = 1
		}
		return llvmValue(llvm.ConstInt(c.context.Int1Type(), value, false))
	case *ir.Conversion:
		return c.compileExpression(expr.Expression, used)
	case *ir.DerefExpression:
		value := c.compileExpression(expr.Value, true).toRValue(c)
		return deref{
			value: value,
			ty:    expr.Type().ToLlvm(c.context),
		}
	case *ir.FloatLiteral:
		if !used {
			return llvmValue{}
		}
		return llvmValue(llvm.ConstFloat(expr.Type().ToLlvm(c.context), expr.Value))
	case *ir.FunctionCall:
		callee := c.compileExpression(expr.Function, true).toRValue(c)
		args := make([]llvm.Value, 0, len(expr.Arguments))
		for _, arg := range expr.Arguments {
			args = append(args, c.compileExpression(arg, true).toRValue(c))
		}
		var name string
		if expr.ReturnType != types.Void && used {
			name = "call_tmp"
		}
		return llvmValue(c.builder.CreateCall(callee.GlobalValueType(), callee, args, name))
	case *ir.FunctionExpression:
		panic("TODO")
	case *ir.IndexExpression:
		panic("TODO")
	case *ir.IntegerLiteral:
		if !used {
			return llvmValue{}
		}
		return llvmValue(llvm.ConstInt(expr.Type().ToLlvm(c.context), uint64(expr.Value), true))
	case *ir.UintLiteral:
		if !used {
			return llvmValue{}
		}
		return llvmValue(llvm.ConstInt(expr.Type().ToLlvm(c.context), expr.Value, true))
	case *ir.MapExpression:
		panic("TODO")
	case *ir.MemberExpression:
		panic("TODO")
	case *ir.RefExpression:
		value := c.compileExpression(expr.Value, true)
		return llvmValue(value.toRef(c))

	case *ir.StringLiteral:
		if !used {
			return llvmValue{}
		}
		return llvmValue(c.builder.CreateGlobalString(expr.Value, ".str_const"))
	case *ir.StructExpression:
		return c.compileStructExpression(expr)
	case *ir.TupleExpression:
		panic("TODO")
	case *ir.TupleStructExpression:
		panic("TODO")
	case *ir.TypeCheck:
		panic("TODO")
	case *ir.TypeExpression:
		panic("TODO")
	case *ir.UnaryExpression:
		return c.compileUnaryExpression(expr)
	case *ir.VariableExpression:
		if !used {
			return llvmValue{}
		}
		return c.table.getValue(expr.Symbol.Name)
	case *ir.BitCast:
		if !used {
			return llvmValue{}
		}

		toType := expr.To.ToLlvm(c.context)
		left := c.compileExpression(expr.Value, true)
		alloca := c.builder.CreateAlloca(toType, "bitcast")
		c.builder.CreateStore(left.toRValue(c), alloca)
		return stackVariable(alloca)
	default:
		panic(fmt.Sprintf("Unexpected expression type: %T", expression))
	}
}

func (c *compiler) compileBinaryExpression(binExpr *ir.BinaryExpression) value {
	left := c.compileExpression(binExpr.Left, true).toRValue(c)
	right := c.compileExpression(binExpr.Right, true).toRValue(c)

	var v llvm.Value

	switch binExpr.Operator.Id {
	case ir.AddFloat:
		v = c.builder.CreateFAdd(left, right, "fadd_tmp")
	case ir.AddInt:
		v = c.builder.CreateAdd(left, right, "add_tmp")
	case ir.BitwiseAnd:
		v = c.builder.CreateAnd(left, right, "bit_and_tmp")
	case ir.BitwiseXor:
		v = c.builder.CreateXor(left, right, "xor_tmp")
	case ir.BitwiseOr:
		v = c.builder.CreateOr(left, right, "bit_or_tmp")
	case ir.Concat:
		panic("TODO")
	case ir.Divide:
		panic("TODO")
	case ir.Equal:
		// TODO: Non-integer comparisons
		v = c.builder.CreateICmp(llvm.IntEQ, left, right, "eq_tmp")
	case ir.Greater:
		// TODO: Float and unsigned comparisons
		v = c.builder.CreateICmp(llvm.IntSGT, left, right, "gt_tmp")
	case ir.GreaterEq:
		// TODO: Float and unsigned comparisons
		v = c.builder.CreateICmp(llvm.IntSGE, left, right, "ge_tmp")
	case ir.LeftShift:
		v = c.builder.CreateShl(left, right, "shl_tmp")
	case ir.Less:
		// TODO: Float and unsigned comparisons
		v = c.builder.CreateICmp(llvm.IntSLT, left, right, "lt_tmp")
	case ir.LessEq:
		// TODO: Float and unsigned comparisons
		v = c.builder.CreateICmp(llvm.IntSLE, left, right, "le_tmp")
	case ir.LogicalAnd:
		v = c.builder.CreateAnd(left, right, "and_tmp")
	case ir.LogicalOr:
		v = c.builder.CreateOr(left, right, "or_tmp")
	case ir.ModuloFloat:
		panic("TODO")
	case ir.ModuloInt:
		panic("TODO")
	case ir.MultiplyFloat:
		v = c.builder.CreateFMul(left, right, "fmul_tmp")
	case ir.MultiplyInt:
		v = c.builder.CreateMul(left, right, "mul_tmp")
	case ir.NotEqual:
		// TODO: Non-integer comparisons
		v = c.builder.CreateICmp(llvm.IntNE, left, right, "ne_tmp")
	case ir.PowerFloat:
		panic("TODO")
	case ir.PowerInt:
		panic("TODO")
	case ir.ArithmeticRightShift:
		// TODO: Probably use lshr for unsigned types
		v = c.builder.CreateAShr(left, right, "arsh_tmp")
	case ir.LogicalRightShift:
		v = c.builder.CreateLShr(left, right, "lrsh_tmp")
	case ir.SubtractFloat:
		v = c.builder.CreateFSub(left, right, "fsub_tmp")
	case ir.SubtractInt:
		v = c.builder.CreateSub(left, right, "sub_tmp")
	case ir.Union:
		panic("Unreachable")
	default:
		panic("Unreachable")
	}

	return llvmValue(v)
}

func (c *compiler) compileUnaryExpression(unExpr *ir.UnaryExpression) value {
	operand := c.compileExpression(unExpr.Operand, true).toRValue(c)

	var v llvm.Value

	switch unExpr.Operator.Id {
	case ir.BitwiseNot:
		v = c.builder.CreateNot(operand, "bit_not_tmp")
	case ir.CrashError:
		panic("TODO")
	case ir.DecrementInt:
		panic("TODO")
	case ir.DecrementFloat:
		panic("TODO")
	case ir.IncrementInt:
		panic("TODO")
	case ir.IncrementFloat:
		panic("TODO")
	case ir.LogicalNot:
		v = c.builder.CreateNot(operand, "not_tmp")
	case ir.NegateFloat:
		v = c.builder.CreateFNeg(operand, "fneg_tmp")
	case ir.NegateInt:
		v = c.builder.CreateNeg(operand, "neg_tmp")
	case ir.PropagateError:
		panic("TODO")
	default:
		panic("Unreachable")
	}

	return llvmValue(v)
}

func (c *compiler) compileStructExpression(s *ir.StructExpression) value {
	values := make([]llvm.Value, 0, len(s.Fields))
	structType := types.Unwrap(s.Struct).(*types.Struct)
	for _, name := range structType.FieldOrder {
		values = append(values, c.compileExpression(s.Fields[name], true).toRValue(c))
	}
	return llvmValue(llvm.ConstStruct(values, false))
}
