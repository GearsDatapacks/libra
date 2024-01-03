package interpreter

import (
	"fmt"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/interpreter/values"
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
	// typechecker "github.com/gearsdatapacks/libra/type_checker"
)

const (
	REGISTER = iota
	IMPORT
	EVALUATE
)

func Evaluate(manager *modules.ModuleManager) values.RuntimeValue {

	register(manager)

	resolveImports(manager)

	return evaluateStatements(manager)
}

func register(manager *modules.ModuleManager) {
	if manager.InterpretStage > REGISTER {
		return
	}
	manager.InterpretStage++

	for _, mod := range manager.Imported {
		register(mod)
	}

	for _, file := range manager.Files {
		for _, stmt := range file.Ast.Body {
			if fn, ok := stmt.(*ast.FunctionDeclaration); ok {
				registerFunctionDeclaration(fn, manager)
			}
		}
	}
}

func resolveImports(manager *modules.ModuleManager) {
	if manager.InterpretStage > IMPORT {
		return
	}
	manager.InterpretStage++

	for _, mod := range manager.Imported {
		resolveImports(mod)
	}

	for _, file := range manager.Files {
		for _, stmt := range file.Ast.Body {
			if imp, ok := stmt.(*ast.ImportStatement); ok {
				evaluateImportStatement(imp, manager)
			}
		}
}
}

func evaluateStatements(manager *modules.ModuleManager) values.RuntimeValue {
	if manager.InterpretStage > EVALUATE {
		return nil
	}
	manager.InterpretStage++

	for _, mod := range manager.Imported {
		evaluateStatements(mod)
	}

	var lastValue values.RuntimeValue
	for _, file := range manager.Files {
		for _, stmt := range file.Ast.Body {
			lastValue = evaluate(stmt, manager)
		}
	}
	return lastValue
}

func evaluate(astNode ast.Statement, manager *modules.ModuleManager) values.RuntimeValue {
	switch statement := astNode.(type) {
	case *ast.ExpressionStatement:
		return evaluateExpressionStatement(statement, manager)

	case *ast.VariableDeclaration:
		return evaluateVariableDeclaration(statement, manager)

	case *ast.FunctionDeclaration:
		// return registerFunctionDeclaration(statement, manager)
		return values.MakeNull()

	case *ast.ReturnStatement:
		return evaluateReturnStatement(statement, manager)

	case *ast.IfStatement:
		return evaluateIfStatement(statement, manager)

	case *ast.WhileLoop:
		return evaluateWhileLoop(statement, manager)

	case *ast.ForLoop:
		return evaluateForLoop(statement, manager)

	case *ast.StructDeclaration:
		// return evaluateStructDeclaration(statement, manager)
		return values.MakeNull()

	case *ast.TupleStructDeclaration:
		return values.MakeNull()

	case *ast.InterfaceDeclaration:
		return values.MakeNull()

	case *ast.TypeDeclaration:
		// env.AddType(statement.Name, typechecker.TypeCheckType(statement.DataType, env))
		return values.MakeNull()

	case *ast.ImportStatement:
		return values.MakeNull()

	case *ast.UnitStructDeclaration:
		return evaluateUnitStructDeclaration(statement, manager)

	default:
		errors.LogError(errors.DevError(fmt.Sprintf("(Interpreter) Unreconised AST node: %s", astNode.String()), astNode))
		return nil
	}
}
