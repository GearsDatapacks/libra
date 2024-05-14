package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

func (t *typeChecker) registerDeclaration(statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.FunctionDeclaration:
		t.registerFunction(stmt)
	case *ast.TypeDeclaration:
		t.registerTypeDeclaration(stmt)
	case *ast.StructDeclaration:
		t.registerStructDeclaration(stmt)
	}
}

func (t *typeChecker) typeCheckDeclaration(statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.TypeDeclaration:
		t.typeCheckTypeDeclaration(stmt)
	case *ast.StructDeclaration:
		t.typeCheckStructDeclaration(stmt)
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
		if !t.symbols.Register(symbol) {
			t.Diagnostics.Report(diagnostics.VariableDefined(fn.Name.Location, fn.Name.Value))
		}
	}
}

func (t *typeChecker) registerTypeDeclaration(typeDec *ast.TypeDeclaration) {
	symbol := &symbols.Type{
		Name: typeDec.Name.Value,
		Type: &types.Alias{Type: types.Void},
	}
	t.symbols.Register(symbol)
}

func (t *typeChecker) registerStructDeclaration(decl *ast.StructDeclaration) {
	if decl.StructType == nil {
		panic("TODO: Tuple struct declarations")
	}
	symbol := &symbols.Type{
		Name: decl.Name.Value,
		Type: &types.Struct{
			Name:   decl.Name.Value,
			Fields: map[string]types.Type{},
		},
	}

	t.symbols.Register(symbol)
}

func (t *typeChecker) typeCheckTypeDeclaration(typeDec *ast.TypeDeclaration) {
	symbol := t.symbols.Lookup(typeDec.Name.Value).(*symbols.Type)
	symbol.Type.(*types.Alias).Type = t.typeFromAst(typeDec.Type)
}

func (t *typeChecker) typeCheckFunctionType(fn *ast.FunctionDeclaration) {
	var fnType *types.Function
	if fn.MethodOf == nil && fn.MemberOf == nil {
		fnType = t.symbols.Lookup(fn.Name.Value).GetType().(*types.Function)
	} else {
		panic("TODO: Methods and static methods")
	}

	for _, param := range fn.Parameters {
		if param.Type != nil {
			paramType := t.typeFromAst(param.Type.Type)
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

	fnType.ReturnType = t.typeFromAst(fn.ReturnType.Type)
}

func (t *typeChecker) typeCheckStructDeclaration(decl *ast.StructDeclaration) {
	ty := t.symbols.Lookup(decl.Name.Value).(*symbols.Type).Type.(*types.Struct)

	fields := []types.Type{}
	for _, field := range decl.StructType.Fields {
		if field.Type == nil {
			fields = append(fields, nil)
		} else {
			ty := t.typeFromAst(field.Type.Type)
			for i := len(fields) - 1; i >= 0; i-- {
				if fields[i] == nil {
					fields[i] = ty
				} else {
					break
				}
			}
			fields = append(fields, ty)
		}
	}

	for i, field := range decl.StructType.Fields {
		ty.Fields[field.Name.Value] = fields[i]
	}
}
