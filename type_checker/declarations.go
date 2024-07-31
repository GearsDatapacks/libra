package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/ir"
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
	case *ast.EnumDeclaration:
		t.registerEnumDeclaration(stmt)
	case *ast.TagDeclaration:
		t.registerTagDeclaration(stmt)
	}
}

func (t *typeChecker) typeCheckDeclaration(statement ast.Statement) ir.Statement {
	switch stmt := statement.(type) {
	case *ast.TypeDeclaration:
		return t.typeCheckTypeDeclaration(stmt)
	case *ast.StructDeclaration:
		return t.typeCheckStructDeclaration(stmt)
	case *ast.InterfaceDeclaration:
		return t.typeCheckInterfaceDeclaration(stmt)
	case *ast.UnionDeclaration:
		return t.typeCheckUnionDeclaration(stmt)
	case *ast.EnumDeclaration:
		return t.typeCheckEnumDeclaration(stmt)
	case *ast.TagDeclaration:
		return t.typeCheckTagDeclaration(stmt)
	}
	return nil
}

func (t *typeChecker) registerFunction(fn *ast.FunctionDeclaration) {
	fnType := &types.Function{
		Parameters: []types.Type{},
		ReturnType: types.Void,
	}
	symbol := &symbols.Variable{
		Name:       fn.Name,
		IsMut:      false,
		Type:       fnType,
		ConstValue: nil,
	}

	if fn.MethodOf == nil && fn.MemberOf == nil {
		if !t.symbols.Register(symbol, fn.Exported) {
			t.diagnostics.Report(diagnostics.VariableDefined(fn.NameLocation, fn.Name))
		}
	}
}

func (t *typeChecker) registerTypeDeclaration(typeDec *ast.TypeDeclaration) {
	var ty types.Type
	if typeDec.Explicit {
		ty = &types.Explicit{Name: typeDec.Name, Type: types.Void}
	} else {
		ty = &types.Alias{Type: types.Void}
	}
	symbol := &symbols.Type{
		Name: typeDec.Name,
		Type: ty,
	}
	t.symbols.Register(symbol, typeDec.Exported)
}

func (t *typeChecker) registerStructDeclaration(decl *ast.StructDeclaration) {
	var ty types.Type

	if decl.Body == nil {
		ty = &types.UnitStruct{Name: decl.Name}
	} else {
		isTuple := true
		for _, field := range decl.Body {
			if field.Name != nil && field.Type != nil {
				isTuple = false
				break
			}
		}

		if isTuple {
			ty = &types.TupleStruct{
				Name:  decl.Name,
				Types: []types.Type{},
			}
		} else {
			ty = &types.Struct{
				Name:     decl.Name,
				ModuleId: t.module.Id,
				Fields:   map[string]types.StructField{},
			}
		}
	}

	symbol := &symbols.Type{
		Name: decl.Name,
		Type: ty,
	}

	t.symbols.Register(symbol, decl.Exported)
}

func (t *typeChecker) registerInterfaceDeclaration(decl *ast.InterfaceDeclaration) {
	symbol := &symbols.Type{
		Name: decl.Name,
		Type: &types.Interface{
			Name:    decl.Name,
			Methods: map[string]*types.Function{},
		},
	}

	t.symbols.Register(symbol, decl.Exported)
}

func (t *typeChecker) registerUnionDeclaration(decl *ast.UnionDeclaration) {
	symbol := &symbols.Type{
		Name: decl.Name,
		Type: &types.Union{
			Name:    decl.Name,
			Members: map[string]types.Type{},
		},
	}

	t.symbols.Register(symbol, decl.Exported)
}

func (t *typeChecker) registerEnumDeclaration(decl *ast.EnumDeclaration) {
	symbol := &symbols.Type{
		Name: decl.Name,
		Type: &types.Enum{
			Name:       decl.Name,
			Underlying: types.Invalid,
			Members:    map[string]values.ConstValue{},
		},
	}

	t.symbols.Register(symbol, decl.Exported)
}

func (t *typeChecker) registerTagDeclaration(decl *ast.TagDeclaration) {
	symbol := &symbols.Type{
		Name: decl.Name,
		Type: &types.Tag{
			Name:  decl.Name,
			Types: []types.Type{},
		},
	}

	t.symbols.Register(symbol, decl.Exported)
}

func (t *typeChecker) typeCheckTypeDeclaration(typeDec *ast.TypeDeclaration) ir.Statement {
	symbol := t.symbols.Lookup(typeDec.Name).(*symbols.Type)
	if typeDec.Explicit {
		symbol.Type.(*types.Explicit).Type = t.typeCheckType(typeDec.Type)
	} else {
		symbol.Type.(*types.Alias).Type = t.typeCheckType(typeDec.Type)
	}

	if typeDec.Tag != nil {
		t.addToTag(typeDec.Tag, symbol.Type)
	}

	return &ir.TypeDeclaration{
		Name:     typeDec.Name,
		Exported: typeDec.Exported,
		Type:     symbol.Type,
	}
}

func (t *typeChecker) typeCheckFunctionType(fn *ast.FunctionDeclaration) {
	var fnType *types.Function
	if fn.MethodOf == nil && fn.MemberOf == nil {
		fnType = t.symbols.Lookup(fn.Name).GetType().(*types.Function)
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
		fnType.ReturnType = t.typeCheckType(fn.ReturnType)
	}

	if fn.MethodOf != nil {
		methodOf := t.typeCheckType(fn.MethodOf.Type)
		t.symbols.RegisterMethod(fn.Name, &symbols.Method{
			MethodOf: methodOf,
			Static:   false,
			Function: fnType,
		}, fn.Exported)
	} else if fn.MemberOf != nil {
		methodOf := t.lookupType(fn.MemberOf.Name, fn.MemberOf.Location)
		t.symbols.RegisterMethod(fn.Name, &symbols.Method{
			MethodOf: methodOf,
			Static:   true,
			Function: fnType,
		}, fn.Exported)
	}
}

func (t *typeChecker) typeCheckStructDeclaration(decl *ast.StructDeclaration) ir.Statement {
	ty := t.symbols.Lookup(decl.Name).(*symbols.Type).Type

	t.typeCheckStructBody(decl.NameLocation, decl.Body, ty)
	if decl.Tag != nil {
		t.addToTag(decl.Tag, ty)
	}

	return &ir.TypeDeclaration{
		Name:     decl.Name,
		Exported: decl.Exported,
		Type:     ty,
	}
}

func (t *typeChecker) typeCheckStructBody(
	nameLocation text.Location,
	body []ast.StructField,
	ty types.Type,
) {
	if structTy, ok := ty.(*types.Struct); ok {
		fields := []types.StructField{}
		for _, field := range body {
			var name string
			if field.Name == nil {
				t.diagnostics.Report(diagnostics.MixedNamedUnnamedStructFields(field.Type.GetLocation()))
			} else {
				name = *field.Name
			}
			structField := types.StructField{
				Name:     name,
				Type:     nil,
				Exported: field.Pub,
			}
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
			lastField := body[len(fields)-1]
			t.diagnostics.ReportMany(diagnostics.LastStructFieldMustHaveType(lastField.Location, nameLocation))
		}

		for i, field := range body {
			if fields[i].Type == nil {
				fields[i].Type = types.Invalid
			}
			if field.Name != nil {
				structTy.Fields[*field.Name] = fields[i]
			}
		}
	} else if structTy, ok := ty.(*types.TupleStruct); ok {
		for _, field := range body {
			if field.Pub {
				t.diagnostics.Report(diagnostics.PubUnnamedStructField(field.Location))
			}
			if field.Type == nil {
				structTy.Types = append(structTy.Types, t.lookupType(*field.Name, field.TypeOrIdent.Location))
			} else {
				structTy.Types = append(structTy.Types, t.typeCheckType(field.Type))
			}
		}
	}
}

func (t *typeChecker) typeCheckInterfaceDeclaration(decl *ast.InterfaceDeclaration) ir.Statement {
	ty := t.symbols.Lookup(decl.Name).(*symbols.Type).Type.(*types.Interface)

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
			fnType.ReturnType = t.typeCheckType(member.ReturnType)
		}
		ty.Methods[member.Name] = fnType
	}

	return &ir.TypeDeclaration{
		Name:     decl.Name,
		Exported: decl.Exported,
		Type:     ty,
	}
}

func (t *typeChecker) typeCheckUnionDeclaration(decl *ast.UnionDeclaration) ir.Statement {
	ty := t.symbols.Lookup(decl.Name).(*symbols.Type).Type.(*types.Union)

	for _, member := range decl.Members {
		var memberType types.Type
		if member.Type != nil {
			memberType = t.typeCheckType(member.Type)
		} else if member.Compound != nil {
			structTy := &types.Struct{
				Name:     member.Name,
				ModuleId: t.module.Id,
				Fields:   map[string]types.StructField{},
			}
			t.typeCheckStructBody(member.NameLocation, member.Compound, structTy)
			memberType = structTy
		} else {
			memberType = t.lookupType(member.Name, member.NameLocation)
		}

		if memberType == nil {
			continue
		}

		if member.Type == nil && member.Compound == nil {
			ty.Members[member.Name] = memberType
		} else {
			ty.Members[member.Name] = &types.UnionVariant{
				Union: ty,
				Name:  member.Name,
				Type:  memberType,
			}
		}
	}

	if decl.Tag != nil {
		t.addToTag(decl.Tag, ty)
	}

	return &ir.TypeDeclaration{
		Name:     decl.Name,
		Exported: decl.Exported,
		Type:     ty,
	}
}

func (t *typeChecker) typeCheckEnumDeclaration(decl *ast.EnumDeclaration) ir.Statement {
	ty := t.symbols.Lookup(decl.Name).(*symbols.Type).Type.(*types.Enum)
	if decl.ValueType != nil {
		ty.Underlying = t.typeCheckType(decl.ValueType)
	} else {
		ty.Underlying = types.Int
	}
	toEnum, canGenValues := ty.Underlying.(types.HasEnumValue)

	enumValues := []values.ConstValue{}
	for _, member := range decl.Members {
		var value values.ConstValue

		if member.Value != nil {
			expression := t.typeCheckExpression(member.Value)
			conversion := convert(expression, ty.Underlying, implicit)
			if conversion == nil {
				t.diagnostics.Report(diagnostics.NotAssignable(
					member.Value.GetLocation(),
					ty.Underlying,
					expression.Type(),
				))
			} else if conversion.IsConst() {
				value = expression.ConstValue()
			} else {
				t.diagnostics.Report(diagnostics.NotConst(member.Value.GetLocation()))
			}
		}

		if value == nil {
			if canGenValues {
				enumValue, err := toEnum.GetEnumValue(enumValues, member.Name)
				if err != nil {
					t.diagnostics.Report(err.Location(member.Location))
				} else {
					value = enumValue
				}
			} else {
				t.diagnostics.Report(diagnostics.CannotEnum(member.Location, ty.Underlying))
			}
		}

		if value != nil {
			enumValues = append(enumValues, value)
			ty.Members[member.Name] = value
		}
	}

	return &ir.TypeDeclaration{
		Name:     decl.Name,
		Exported: decl.Exported,
		Type:     ty,
	}
}

func (t *typeChecker) typeCheckTagDeclaration(decl *ast.TagDeclaration) ir.Statement {
	ty := t.symbols.Lookup(decl.Name).(*symbols.Type).Type.(*types.Tag)

	if decl.Body != nil {
		for _, member := range decl.Body {
			ty.Types = append(ty.Types, t.typeCheckType(member))
		}
	}

	return &ir.TypeDeclaration{
		Name:     decl.Name,
		Exported: decl.Exported,
		Type:     ty,
	}
}

func (t *typeChecker) addToTag(tag ast.Expression, ty types.Type) {
	typeChecked := t.typeCheckType(tag)
	if typeChecked == types.Invalid {
		return
	}

	tagType, ok := typeChecked.(*types.Tag)
	if !ok {
		t.diagnostics.Report(diagnostics.NotATag(tag.GetLocation(), typeChecked))
		return
	}
	tagType.Types = append(tagType.Types, ty)
}

type moduleWrapper struct{ t *symbols.Table }

func (mod moduleWrapper) LookupExport(name string) interface{ Value() values.ConstValue } {
	return mod.t.LookupExport(name)
}

func (t *typeChecker) typeCheckImport(importStmt *ast.ImportStatement) ir.Statement {
	module, ok := t.subModules[importStmt.Module.ExtraValue]

	if !ok {
		t.diagnostics.Report(diagnostics.ModuleUndefined(importStmt.Module.Location, importStmt.Module.Value))
		return nil
	}

	name := module.module.Name
	if importStmt.Alias != nil {
		name = *importStmt.Alias
	}
	importedSymbols := []string{}

	if importStmt.All {
		t.symbols.Extend(module.symbols)
	} else if importStmt.Symbols != nil {
		for _, symbol := range importStmt.Symbols {
			export := module.symbols.LookupExport(symbol.Name)
			importedSymbols = append(importedSymbols, symbol.Name)
			if export != nil {
				t.symbols.Register(export)
			} else {
				t.diagnostics.Report(diagnostics.NoExport(symbol.Location, module.module.Name, symbol.Name))
			}
		}
	} else {
		moduleType := &types.Module{
			Name:   module.module.Name,
			Module: module.symbols,
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

	return &ir.ImportStatement{
		Module:    importStmt.Module.ExtraValue,
		Name:      name,
		Symbols:   importedSymbols,
		ImportAll: importStmt.All,
	}
}
