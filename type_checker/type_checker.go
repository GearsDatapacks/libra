package typechecker

import (
	"fmt"

	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TypeCheck(manager *modules.ModuleManager) error {
	err := registerStatements(manager)
	if err != nil {
		return err
	}
	err = typeCheckGlobalStatements(manager)
	if err != nil {
		return err
	}
	err = typeCheckImportStatements(manager)
	if err != nil {
		return err
	}
	err = typeCheckFunctions(manager)
	if err != nil {
		return err
	}

	return typeCheck(manager)
}

const (
	REGISTER = iota
	GLOBAL
	IMPORT
	FUNCTION
	STATEMENT
)

func typeCheck(manager *modules.ModuleManager) error {
	if manager.TypeCheckStage > STATEMENT {
		return nil
	}
	manager.TypeCheckStage++

	for _, mod := range manager.Modules {
		typeCheck(mod)
	}

	for _, stmt := range manager.Main.Ast.Body {
		nextType := typeCheckStatement(stmt, manager)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}
	return nil
}

func registerStatements(manager *modules.ModuleManager) error {
	if manager.TypeCheckStage > REGISTER {
		return nil
	}
	manager.TypeCheckStage++

	for _, mod := range manager.Modules {
		registerStatements(mod)
	}

	for _, stmt := range manager.Main.Ast.Body {
		nextType := registerTypeStatement(stmt, manager)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}
	return nil
}

func typeCheckGlobalStatements(manager *modules.ModuleManager) error {
	if manager.TypeCheckStage > GLOBAL {
		return nil
	}
	manager.TypeCheckStage++

	for _, mod := range manager.Modules {
		typeCheckGlobalStatements(mod)
	}

	for _, stmt := range manager.Main.Ast.Body {
		nextType := typeCheckGlobalStatement(stmt, manager)
		if nextType.String() == "TypeError" {
			return nextType.(*types.TypeError)
		}
	}
	return nil
}

func typeCheckImportStatements(manager *modules.ModuleManager) error {
	if manager.TypeCheckStage > IMPORT {
		return nil
	}
	manager.TypeCheckStage++

	for _, mod := range manager.Modules {
		typeCheckImportStatements(mod)
	}

	for _, stmt := range manager.Main.Ast.Body {
		if importStmt, ok := stmt.(*ast.ImportStatement); ok {
			nextType := typeCheckImportStatement(importStmt, manager)
			if nextType.String() == "TypeError" {
				return nextType.(*types.TypeError)
			}
		}
	}
	return nil
}

func typeCheckFunctions(manager *modules.ModuleManager) error {
	if manager.TypeCheckStage > FUNCTION {
		return nil
	}
	manager.TypeCheckStage++

	for _, mod := range manager.Modules {
		typeCheckFunctions(mod)
	}

	for _, stmt := range manager.Main.Ast.Body {
		if funcDec, ok := stmt.(*ast.FunctionDeclaration); ok {
			nextType := typeCheckFunctionParams(funcDec, manager)
			if nextType.String() == "TypeError" {
				return nextType.(*types.TypeError)
			}
		}
	}
	return nil
}

func typeCheckType(ty ast.TypeExpression, manager *modules.ModuleManager) types.ValidType {
	if member, ok := ty.(*ast.MemberType); ok {
		return typeCheckMemberType(member, manager)
	}
	return types.FromAst(ty, manager.SymbolTable)
}

func typeCheckMemberType(member *ast.MemberType, manager *modules.ModuleManager) types.ValidType {
	var left types.ValidType
	if ident, ok := member.Left.(*ast.TypeName); ok {
		left = manager.SymbolTable.GetSymbol(ident.Name)
	} else {
		left = typeCheckMemberType(member.Left.(*ast.MemberType), manager)
	}

	if left.String() == "TypeError" {
		return left
	}

	memberType := types.Member(left, member.Member, false)

	if memberType == nil {
		return types.Error(fmt.Sprintf("Type %q is undefined", member.String()), member)
	}
	if ty, isType := memberType.(*types.Type); isType {
		return ty.DataType
	}
	return types.Error(fmt.Sprintf("Cannot use %q as type, it is a value", member.String()), member)
}
