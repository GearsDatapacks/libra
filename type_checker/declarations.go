package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

func (t *typeChecker) registerDeclaration(statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.FunctionDeclaration:
		t.registerFunction(stmt)
	case *ast.TypeDeclaration:
		t.registerTypeDeclaration(stmt)
	case *ast.StructDeclaration:
		t.registerStructDeclaration(stmt)
	case *ast.InterfaceDeclaration:
		t.registerInterfaceDeclaration(stmt)
	case *ast.UnionDeclaration:
		t.registerUnionDeclaration(stmt)
	case *ast.TagDeclaration:
		t.registerTagDeclaration(stmt)
	}
}

func (t *typeChecker) typeCheckDeclaration(statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.TypeDeclaration:
		t.typeCheckTypeDeclaration(stmt)
	case *ast.StructDeclaration:
		t.typeCheckStructDeclaration(stmt)
	case *ast.InterfaceDeclaration:
		t.typeCheckInterfaceDeclaration(stmt)
	case *ast.UnionDeclaration:
		t.typeCheckUnionDeclaration(stmt)
	case *ast.TagDeclaration:
		t.typeCheckTagDeclaration(stmt)
	}
}

func (t *typeChecker) registerFunction(fn *ast.FunctionDeclaration) {
	fnType := &types.Function{
		Parameters: []types.Type{},
		ReturnType: types.Void,
	}
	symbol := &symbols.Variable{
		Name:       fn.Name.Value,
		IsMut:      false,
		Type:       fnType,
		ConstValue: nil,
	}

	if fn.MethodOf == nil && fn.MemberOf == nil {
		if !t.symbols.Register(symbol, fn.Exported) {
			t.diagnostics.Report(diagnostics.VariableDefined(fn.Name.Location, fn.Name.Value))
		}
	}
}

func (t *typeChecker) registerTypeDeclaration(typeDec *ast.TypeDeclaration) {
	var ty types.Type
	if typeDec.Explicit {
		ty = &types.Explicit{Name: typeDec.Name.Value, Type: types.Void}
	} else {
		ty = &types.Alias{Type: types.Void}
	}
	symbol := &symbols.Type{
		Name: typeDec.Name.Value,
		Type: ty,
	}
	t.symbols.Register(symbol, typeDec.Exported)
}

func (t *typeChecker) registerStructDeclaration(decl *ast.StructDeclaration) {
	var ty types.Type

	if decl.Body == nil {
		ty = &types.UnitStruct{Name: decl.Name.Value}
	} else {
		isTuple := true
		for _, field := range decl.Body.Fields {
			if field.Name != nil && field.Type != nil {
				isTuple = false
				break
			}
		}

		if isTuple {
			ty = &types.TupleStruct{
				Name:  decl.Name.Value,
				Types: []types.Type{},
			}
		} else {
			ty = &types.Struct{
				Name:     decl.Name.Value,
				ModuleId: t.module.Id,
				Fields:   map[string]types.StructField{},
			}
		}
	}

	symbol := &symbols.Type{
		Name: decl.Name.Value,
		Type: ty,
	}

	t.symbols.Register(symbol, decl.Exported)
}

func (t *typeChecker) registerInterfaceDeclaration(decl *ast.InterfaceDeclaration) {
	symbol := &symbols.Type{
		Name: decl.Name.Value,
		Type: &types.Interface{
			Name:    decl.Name.Value,
			Methods: map[string]*types.Function{},
		},
	}

	t.symbols.Register(symbol, decl.Exported)
}

func (t *typeChecker) registerUnionDeclaration(decl *ast.UnionDeclaration) {
	symbol := &symbols.Type{
		Name: decl.Name.Value,
		Type: &types.Union{
			Name:    decl.Name.Value,
			Members: map[string]types.Type{},
		},
	}

	t.symbols.Register(symbol, decl.Exported)
}

func (t *typeChecker) registerTagDeclaration(decl *ast.TagDeclaration) {
	symbol := &symbols.Type{
		Name: decl.Name.Value,
		Type: &types.Tag{
			Name:  decl.Name.Value,
			Types: []types.Type{},
		},
	}

	t.symbols.Register(symbol, decl.Exported)
}

func (t *typeChecker) typeCheckTypeDeclaration(typeDec *ast.TypeDeclaration) {
	symbol := t.symbols.Lookup(typeDec.Name.Value).(*symbols.Type)
	if typeDec.Explicit {
		symbol.Type.(*types.Explicit).Type = t.typeCheckType(typeDec.Type)
	} else {
		symbol.Type.(*types.Alias).Type = t.typeCheckType(typeDec.Type)
	}

	if typeDec.Tag != nil {
		t.addToTag(typeDec.Tag, symbol.Type)
	}
}

func (t *typeChecker) typeCheckFunctionType(fn *ast.FunctionDeclaration) {
	var fnType *types.Function
	if fn.MethodOf == nil && fn.MemberOf == nil {
		fnType = t.symbols.Lookup(fn.Name.Value).GetType().(*types.Function)
	} else {
		fnType = &types.Function{
			Parameters: []types.Type{},
			ReturnType: types.Void,
		}
	}

	for _, param := range fn.Parameters {
		if param.Type != nil {
			paramType := t.typeCheckType(param.Type)
			for i := len(fnType.Parameters) - 1; i >= 0; i-- {
				if fnType.Parameters[i] == nil {
					fnType.Parameters[i] = paramType
				} else {
					break
				}
			}
			fnType.Parameters = append(fnType.Parameters, paramType)
		} else {
			fnType.Parameters = append(fnType.Parameters, nil)
		}
	}

	if fn.ReturnType != nil {
		fnType.ReturnType = t.typeCheckType(fn.ReturnType.Type)
	}

	if fn.MethodOf != nil {
		methodOf := t.typeCheckType(fn.MethodOf.Type)
		t.symbols.RegisterMethod(fn.Name.Value, &symbols.Method{
			MethodOf: methodOf,
			Static:   false,
			Function: fnType,
		}, fn.Exported)
	} else if fn.MemberOf != nil {
		methodOf := t.lookupType(fn.MemberOf.Name)
		t.symbols.RegisterMethod(fn.Name.Value, &symbols.Method{
			MethodOf: methodOf,
			Static:   true,
			Function: fnType,
		}, fn.Exported)
	}
}

func (t *typeChecker) typeCheckStructDeclaration(decl *ast.StructDeclaration) {
	ty := t.symbols.Lookup(decl.Name.Value).(*symbols.Type).Type

	t.typeCheckStructBody(decl.Name, decl.Body, ty)
	if decl.Tag != nil {
		t.addToTag(decl.Tag, ty)
	}
}

func (t *typeChecker) typeCheckStructBody(name token.Token, body *ast.StructBody, ty types.Type) {
	if structTy, ok := ty.(*types.Struct); ok {
		fields := []types.StructField{}
		for _, field := range body.Fields {
			if field.Name == nil {
				t.diagnostics.Report(diagnostics.MixedNamedUnnamedStructFields(field.Type.Location()))
			}
			structField := types.StructField{Type: nil, Exported: field.Pub != nil}
			if field.Type == nil {
				fields = append(fields, structField)
			} else {
				ty := t.typeCheckType(field.Type)
				for i := len(fields) - 1; i >= 0; i-- {
					if fields[i].Type == nil {
						fields[i].Type = ty
					} else {
						break
					}
				}
				structField.Type = ty
				fields = append(fields, structField)
			}
		}

		if len(fields) > 0 && fields[len(fields)-1].Type == nil {
			lastField := body.Fields[len(fields)-1]
			t.diagnostics.ReportMany(diagnostics.LastStructFieldMustHaveType(lastField.Name.Location, name.Location))
		}

		for i, field := range body.Fields {
			if fields[i].Type == nil {
				fields[i].Type = types.Invalid
			}
			structTy.Fields[field.Name.Value] = fields[i]
		}
	} else if structTy, ok := ty.(*types.TupleStruct); ok {
		for _, field := range body.Fields {
			if field.Pub != nil {
				t.diagnostics.Report(diagnostics.PubUnnamedStructField(field.Pub.Location))
			}
			if field.Type == nil {
				structTy.Types = append(structTy.Types, t.lookupType(*field.Name))
			} else {
				structTy.Types = append(structTy.Types, t.typeCheckType(field.Type))
			}
		}
	}
}

func (t *typeChecker) typeCheckInterfaceDeclaration(decl *ast.InterfaceDeclaration) {
	ty := t.symbols.Lookup(decl.Name.Value).(*symbols.Type).Type.(*types.Interface)

	for _, member := range decl.Members {
		params := []types.Type{}
		for _, param := range member.Parameters {
			params = append(params, t.typeCheckType(param))
		}
		fnType := &types.Function{
			Parameters: params,
			ReturnType: types.Void,
		}
		if member.ReturnType != nil {
			fnType.ReturnType = t.typeCheckType(member.ReturnType.Type)
		}
		ty.Methods[member.Name.Value] = fnType
	}
}

func (t *typeChecker) typeCheckUnionDeclaration(decl *ast.UnionDeclaration) {
	ty := t.symbols.Lookup(decl.Name.Value).(*symbols.Type).Type.(*types.Union)

	for _, member := range decl.Members {
		var memberType types.Type
		if member.Type != nil {
			memberType = t.typeCheckType(member.Type.Type)
		} else if member.Compound != nil {
			structTy := &types.Struct{
				Name:     member.Name.Value,
				ModuleId: t.module.Id,
				Fields:   map[string]types.StructField{},
			}
			t.typeCheckStructBody(member.Name, member.Compound, structTy)
			memberType = structTy
		} else {
			memberType = t.lookupType(member.Name)
		}

		if memberType == nil {
			continue
		}

		if member.Type == nil && member.Compound == nil {
			ty.Members[member.Name.Value] = memberType
		} else {
			ty.Members[member.Name.Value] = &types.UnionVariant{
				Union: ty,
				Name:  member.Name.Value,
				Type:  memberType,
			}
		}
	}

	if decl.Tag != nil {
		t.addToTag(decl.Tag, ty)
	}
}

func (t *typeChecker) typeCheckTagDeclaration(decl *ast.TagDeclaration) {
	ty := t.symbols.Lookup(decl.Name.Value).(*symbols.Type).Type.(*types.Tag)

	if decl.Body != nil {
		for _, member := range decl.Body.Types {
			ty.Types = append(ty.Types, t.typeCheckType(member))
		}
	}
}

func (t *typeChecker) addToTag(tag ast.Expression, ty types.Type) {
	typeChecked := t.typeCheckType(tag)
	if typeChecked == types.Invalid {
		return
	}
	
	tagType, ok := typeChecked.(*types.Tag)
	if !ok {
		t.diagnostics.Report(diagnostics.NotATag(tag.Location(), tag))
		return
	}
	tagType.Types = append(tagType.Types, ty)
}

type moduleWrapper struct{ t *symbols.Table }

func (mod moduleWrapper) LookupExport(name string) interface{ Value() values.ConstValue } {
	return mod.t.LookupExport(name)
}

func (t *typeChecker) typeCheckImport(importStmt *ast.ImportStatement) {
	module := t.subModules[importStmt.Module.Value]
	if importStmt.All != nil {
		t.symbols.Extend(module.symbols)
	} else if importStmt.Symbols != nil {
		for _, symbol := range importStmt.Symbols.Symbols {
			export := module.symbols.LookupExport(symbol.Value)
			if export != nil {
				t.symbols.Register(export)
			} else {
				t.diagnostics.Report(diagnostics.NoExport(symbol.Location, module.module.Name, symbol.Value))
			}
		}
	} else {
		moduleType := &types.Module{
			Name:   module.module.Name,
			Module: module.symbols,
		}
		name := module.module.Name
		if importStmt.Alias != nil {
			name = importStmt.Alias.Alias.Value
		}
		symbol := &symbols.Variable{
			Name:  name,
			IsMut: false,
			Type:  moduleType,
			ConstValue: values.ModuleValue{
				Module: moduleWrapper{module.symbols},
			},
		}
		t.symbols.Register(symbol)
	}
}
