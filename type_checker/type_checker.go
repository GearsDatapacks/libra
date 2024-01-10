package typechecker

import (
	"fmt"
	"log"

	"github.com/gearsdatapacks/libra/errors"
	"github.com/gearsdatapacks/libra/modules"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func TypeCheck(manager *modules.ModuleManager) error {
	err := registerStatements(manager)
	if err != nil {
		return err
	}
	err = typeCheckImportStatements(manager)
	if err != nil {
		return err
	}
	err = typeCheckGlobalStatements(manager)
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
	IMPORT
	GLOBAL
	FUNCTION
	STATEMENT
)

func typeCheck(manager *modules.ModuleManager) error {
	if manager.TypeCheckStage > STATEMENT {
		return nil
	}
	manager.TypeCheckStage++

	for _, mod := range manager.Imported {
		err := typeCheck(mod)
		if err != nil {
			return err
		}
	}

	for _, file := range manager.Files {
		for _, stmt := range file.Ast.Body {
			nextType := typeCheckStatement(stmt, manager)
			if nextType.String() == "TypeError" {
				return nextType.(*types.TypeError)
			}
		}
	}
	return nil
}

func registerStatements(manager *modules.ModuleManager) error {
	if manager.TypeCheckStage > REGISTER {
		return nil
	}
	manager.TypeCheckStage++

	for _, mod := range manager.Imported {
		err := registerStatements(mod)
		if err != nil {
			return err
		}
	}

	for _, file := range manager.Files {
		for _, stmt := range file.Ast.Body {
			nextType := registerTypeStatement(stmt, manager)
			if nextType.String() == "TypeError" {
				return nextType.(*types.TypeError)
			}
		}
	}
	return nil
}

func typeCheckGlobalStatements(manager *modules.ModuleManager) error {
	if manager.TypeCheckStage > GLOBAL {
		return nil
	}
	manager.TypeCheckStage++

	for _, mod := range manager.Imported {
		err := typeCheckGlobalStatements(mod)
		if err != nil {
			return err
		}
	}

	for _, file := range manager.Files {
		for _, stmt := range file.Ast.Body {
			nextType := typeCheckGlobalStatement(stmt, manager)
			if nextType.String() == "TypeError" {
				return nextType.(*types.TypeError)
			}
		}
	}
	return nil
}

func typeCheckImportStatements(manager *modules.ModuleManager) error {
	if manager.TypeCheckStage > IMPORT {
		return nil
	}
	manager.TypeCheckStage++

	for _, mod := range manager.Imported {
		err := typeCheckImportStatements(mod)
		if err != nil {
			return err
		}
	}

	for _, file := range manager.Files {
		for _, stmt := range file.Ast.Body {
			if importStmt, ok := stmt.(*ast.ImportStatement); ok {
				nextType := typeCheckImportStatement(importStmt, manager)
				if nextType.String() == "TypeError" {
					return nextType.(*types.TypeError)
				}
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

	for _, mod := range manager.Imported {
		err := typeCheckFunctions(mod)
		if err != nil {
			return err
		}
	}

	for _, file := range manager.Files {
		for _, stmt := range file.Ast.Body {
			if funcDec, ok := stmt.(*ast.FunctionDeclaration); ok {
				nextType := typeCheckFunctionParams(funcDec, manager)
				if nextType.String() == "TypeError" {
					return nextType.(*types.TypeError)
				}
				funcDec.SetType(nextType)
			}
		}
	}

	return nil
}

func TypeCheckType(ty ast.TypeExpression, manager *modules.ModuleManager) types.ValidType {
	var dataType types.ValidType
	if member, ok := ty.(*ast.MemberType); ok {
		dataType = typeCheckMemberType(member, manager)
	} else {
		dataType = FromAst(ty, manager.SymbolTable)
	}
	ty.SetType(dataType)
	return dataType
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

	memberType := types.Member(left, member.Member, false, manager.Id)

	if memberType == nil {
		return types.Error(fmt.Sprintf("Type %q is undefined", member.String()), member)
	}
	if ty, isType := memberType.(*types.Type); isType {
		return ty.DataType
	}
	return types.Error(fmt.Sprintf("Cannot use %q as type, it is a value", member.String()), member)
}


func FromAst(node ast.TypeExpression, table *symbols.SymbolTable) types.ValidType {
	switch typeExpr := node.(type) {
	case *ast.TypeName:
		dataType := types.FromString(typeExpr.Name, table)
		if err, isErr := dataType.(*types.TypeError); isErr {
			err.Line = node.GetToken().Line
			err.Column = node.GetToken().Column
		}
		return dataType

	case *ast.Union:
		dataTypes := []types.ValidType{}

		for _, dataType := range typeExpr.ValidTypes {
			nextType := FromAst(dataType, table)
			if nextType.String() == "TypeError" {
				return nextType
			}
		    dataTypes = append(dataTypes, nextType)
		}

		return types.MakeUnion(dataTypes...)

	case *ast.ListType:
		dataType := FromAst(typeExpr.ElementType, table)
		if dataType.String() == "TypeError" {
			return dataType
		}

		return &types.ListLiteral{
			ElemType: dataType,
		}

	case *ast.ArrayType:
		dataType := FromAst(typeExpr.ElementType, table)
		if dataType.String() == "TypeError" {
			return dataType
		}

		return &types.ArrayLiteral{
			ElemType: dataType,
			Length:   typeExpr.Length,
		}

	case *ast.MapType:
		keyType := FromAst(typeExpr.KeyType, table)
		if keyType.String() == "TypeError" {
			return keyType
		}

		valueType := FromAst(typeExpr.ValueType, table)
		if valueType.String() == "TypeError" {
			return valueType
		}

		return &types.MapLiteral{
			KeyType:   keyType,
			ValueType: valueType,
		}

	case *ast.ErrorType:
		resultType := FromAst(typeExpr.ResultType, table)
		if resultType.String() == "TypeError" {
			return resultType
		}

		return &types.ErrorType{ResultType: resultType}

	case *ast.TupleType:
		members := []types.ValidType{}
		for _, member := range typeExpr.Members {
			resultType := FromAst(member, table)
			if resultType.String() == "TypeError" {
				return resultType
			}
			members = append(members, resultType)
		}

		return &types.Tuple{Members: members}

	case *ast.VoidType:
		return &types.Void{}

	case *ast.InferType:
		return &types.Infer{}

	default:
		log.Fatal(errors.DevError("Unexpected type node: " + node.String()))
		return nil
	}
}
